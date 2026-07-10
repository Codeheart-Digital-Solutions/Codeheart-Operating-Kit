package reconcile

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/kitfs"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

const (
	BeginMarker = "<!-- BEGIN CODEHEART OPERATING KIT MANAGED BLOCK -->"
	EndMarker   = "<!-- END CODEHEART OPERATING KIT MANAGED BLOCK -->"
)

type Action struct {
	Kind    string `json:"kind"`
	Target  string `json:"target"`
	Owner   string `json:"owner"`
	Content []byte `json:"-"`
	Mode    uint32 `json:"mode,omitempty"`
}

type Request struct {
	Command       string
	Root          string
	Observed      state.Observed
	Graph         state.Graph
	DesiredLock   map[string]any
	DesiredConfig map[string]any
	ExtraFiles    map[string][]byte
	EnsureIgnore  bool
	ExpectedAfter []state.Classification
}

type Plan struct {
	SchemaVersion int                    `json:"schema_version"`
	ID            string                 `json:"id"`
	Command       string                 `json:"command"`
	Root          string                 `json:"root"`
	StateBefore   string                 `json:"state_before"`
	Actions       []Action               `json:"actions"`
	Blockers      []Blocker              `json:"blockers"`
	DesiredLock   map[string]any         `json:"-"`
	ExpectedAfter []state.Classification `json:"expected_after,omitempty"`
}

// BuildPlan compares desired state with observed bytes without mutating the target.
func BuildPlan(request Request) (Plan, error) {
	root, err := filepath.Abs(request.Root)
	if err != nil {
		return Plan{}, err
	}
	plan := Plan{
		SchemaVersion: 1,
		Command:       request.Command,
		Root:          root,
		StateBefore:   string(request.Observed.Classification),
		Actions:       []Action{},
		Blockers:      []Blocker{},
		DesiredLock:   request.DesiredLock,
		ExpectedAfter: append([]state.Classification{}, request.ExpectedAfter...),
	}
	if len(plan.ExpectedAfter) == 0 {
		plan.ExpectedAfter = []state.Classification{state.StateCurrent}
	}
	for _, node := range request.Graph.Nodes {
		if node.DirectoryTarget || node.Ownership == state.OwnershipLocalMachine {
			continue
		}
		if node.Target == state.LockPath || node.Target == state.ConfigPath {
			continue
		}
		switch node.Update {
		case state.UpdateReplace:
			content, err := kitfs.ReadFile(node.Source)
			if err != nil {
				return Plan{}, fmt.Errorf("read desired source %s: %w", node.Source, err)
			}
			plan.addFileAction(root, node.Target, string(node.Ownership), content, node.Presence == state.PresenceCreateWhenAbsent)
		case state.UpdateManagedSection:
			template, err := kitfs.ReadFile(node.Source)
			if err != nil {
				return Plan{}, err
			}
			desired, blocker, err := desiredManagedSection(filepath.Join(root, filepath.FromSlash(node.Target)), template)
			if err != nil {
				return Plan{}, err
			}
			if blocker != nil {
				plan.Blockers = append(plan.Blockers, *blocker)
				continue
			}
			plan.addFileAction(root, node.Target, string(node.Ownership), desired, false)
		case state.UpdatePreserve:
			if node.Presence == state.PresenceCreateWhenAbsent && node.Source != "" && kitfs.Exists(node.Source) {
				content, err := kitfs.ReadFile(node.Source)
				if err != nil {
					return Plan{}, err
				}
				plan.addFileAction(root, node.Target, string(node.Ownership), content, true)
			}
		case state.UpdateAppendOnly, state.UpdateReportOnly:
			continue
		}
	}
	for target, content := range request.ExtraFiles {
		plan.addFileAction(root, target, "generated-surface", content, true)
	}
	if request.EnsureIgnore {
		content, changed, err := desiredGitignore(filepath.Join(root, ".gitignore"))
		if err != nil {
			return Plan{}, err
		}
		if changed {
			plan.addFileAction(root, ".gitignore", "template", content, false)
		}
	}
	if request.DesiredConfig != nil {
		data, err := state.EncodeYAML(request.DesiredConfig)
		if err != nil {
			return Plan{}, err
		}
		plan.addFileAction(root, state.ConfigPath, "generated-surface", data, false)
	}
	if len(request.Graph.Nodes) > 0 {
		plan.addSafeRemovals(root, request.Observed.Lock, request.Graph)
	}
	sort.Slice(plan.Actions, func(i, j int) bool {
		if plan.Actions[i].Target == state.LockPath {
			return false
		}
		if plan.Actions[j].Target == state.LockPath {
			return true
		}
		return plan.Actions[i].Target < plan.Actions[j].Target
	})
	writeLock := request.DesiredLock != nil
	if writeLock && len(plan.Actions) == 0 && lockMateriallyEqual(request.Observed.Lock, request.DesiredLock) {
		writeLock = false
		plan.DesiredLock = nil
	}
	if writeLock {
		previousGeneration := state.AsInt(request.Observed.Lock["state_generation"])
		if state.AsInt(request.Observed.Lock["schema_version"]) != 2 {
			previousGeneration = 0
		}
		request.DesiredLock["state_generation"] = previousGeneration + 1
		operation := state.Map(request.DesiredLock["last_operation"])
		if operation != nil {
			operation["previous_generation"] = previousGeneration
		}
		plan.DesiredLock = request.DesiredLock
	}
	plan.ID = planDigest(plan)
	if writeLock {
		operation := state.Map(request.DesiredLock["last_operation"])
		if operation != nil {
			operation["transaction_id"] = plan.ID
		}
		data, err := state.EncodeYAML(request.DesiredLock)
		if err != nil {
			return Plan{}, err
		}
		plan.addFileAction(root, state.LockPath, "generated-surface", data, false)
	}
	return plan, nil
}

func (plan *Plan) addFileAction(root, target, owner string, desired []byte, createOnly bool) {
	path := filepath.Join(root, filepath.FromSlash(target))
	existing, err := os.ReadFile(path)
	if err == nil {
		if createOnly || string(existing) == string(desired) {
			return
		}
		plan.Actions = append(plan.Actions, Action{Kind: "replace", Target: target, Owner: owner, Content: desired, Mode: 0o644})
		return
	}
	if os.IsNotExist(err) {
		plan.Actions = append(plan.Actions, Action{Kind: "create", Target: target, Owner: owner, Content: desired, Mode: 0o644})
		return
	}
	plan.Blockers = append(plan.Blockers, Blocker{Code: "target_unreadable", Message: err.Error(), Path: target})
}

func (plan *Plan) addSafeRemovals(root string, lock map[string]any, graph state.Graph) {
	if lock == nil {
		return
	}
	desired := map[string]bool{}
	for _, node := range graph.Nodes {
		if node.Ownership == state.OwnershipManaged {
			desired[node.Target] = true
		}
	}
	for _, item := range state.AnySlice(lock["managed_paths"]) {
		record := state.Map(item)
		if record == nil || state.AsString(record["ownership"]) != "managed" {
			continue
		}
		target := state.AsString(record["path"])
		if target == "" || desired[target] {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(target)))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			plan.Blockers = append(plan.Blockers, Blocker{Code: "managed_path_unreadable", Message: err.Error(), Path: target})
			continue
		}
		digest := sha256.Sum256(data)
		if hex.EncodeToString(digest[:]) != state.AsString(record["checksum_sha256"]) {
			plan.Blockers = append(plan.Blockers, Blocker{
				Code:         "managed_path_modified",
				Message:      "retired managed path contains changes and was preserved",
				Path:         target,
				Remediation:  "move or preserve the changed file, then retry the lifecycle command",
				RetryCommand: plan.Command,
			})
			continue
		}
		plan.Actions = append(plan.Actions, Action{Kind: "remove", Target: target, Owner: "managed"})
	}
}

func desiredManagedSection(path string, template []byte) ([]byte, *Blocker, error) {
	desiredBlock, err := managedBlock(string(template))
	if err != nil {
		return nil, nil, err
	}
	existing, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return template, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	text := string(existing)
	hasBegin := strings.Contains(text, BeginMarker)
	hasEnd := strings.Contains(text, EndMarker)
	if hasBegin != hasEnd {
		return nil, &Blocker{
			Code:        "managed_section_invalid",
			Message:     "shared file contains only one Operating Kit managed marker",
			Path:        filepath.Base(path),
			Remediation: "repair the marker boundary before retrying",
		}, nil
	}
	if hasBegin {
		before, rest, _ := strings.Cut(text, BeginMarker)
		_, after, _ := strings.Cut(rest, EndMarker)
		return []byte(before + desiredBlock + after), nil, nil
	}
	return []byte(string(template) + "\n\n# Existing Instructions Preserved\n\n" + text), nil, nil
}

func managedBlock(text string) (string, error) {
	before, rest, ok := strings.Cut(text, BeginMarker)
	_ = before
	if !ok {
		return "", fmt.Errorf("managed template has no begin marker")
	}
	middle, _, ok := strings.Cut(rest, EndMarker)
	if !ok {
		return "", fmt.Errorf("managed template has no end marker")
	}
	return BeginMarker + middle + EndMarker, nil
}

func desiredGitignore(path string) ([]byte, bool, error) {
	lines := []string{}
	data, err := os.ReadFile(path)
	if err == nil {
		lines = strings.Split(strings.TrimSuffix(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n"), "\n")
	} else if !os.IsNotExist(err) {
		return nil, false, err
	}
	required := []string{
		"# Codeheart Operating Kit local user layer",
		".codeheart/user/preferences.yaml",
		".codeheart/user/*.local.yaml",
		".codeheart/user/feedback/",
		"# Codeheart Operating Kit local machine layer",
		".codeheart/local/",
		".codeheart/kit.transaction.json",
	}
	present := map[string]bool{}
	for _, line := range lines {
		present[line] = true
	}
	changed := false
	for _, line := range required {
		if !present[line] {
			if len(lines) > 0 && !changed {
				lines = append(lines, "")
			}
			lines = append(lines, line)
			present[line] = true
			changed = true
		}
	}
	return []byte(strings.TrimLeft(strings.Join(lines, "\n"), "\n") + "\n"), changed, nil
}

func planDigest(plan Plan) string {
	type digestAction struct {
		Kind   string `json:"kind"`
		Target string `json:"target"`
		Digest string `json:"digest,omitempty"`
	}
	actions := make([]digestAction, len(plan.Actions))
	for index, action := range plan.Actions {
		digest := sha256.Sum256(action.Content)
		actions[index] = digestAction{Kind: action.Kind, Target: action.Target, Digest: hex.EncodeToString(digest[:])}
	}
	desiredLock := state.DeepCopy(plan.DesiredLock)
	if operation := state.Map(desiredLock["last_operation"]); operation != nil {
		delete(operation, "transaction_id")
	}
	data, _ := json.Marshal(map[string]any{"command": plan.Command, "root": plan.Root, "state": plan.StateBefore, "actions": actions, "desired_lock": desiredLock})
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:16])
}

// ManagedPathRecords returns deterministic lock records for whole-file managed graph nodes.
func ManagedPathRecords(graph state.Graph) []any {
	records := []any{}
	for _, node := range graph.Nodes {
		if node.Ownership != state.OwnershipManaged || node.Update != state.UpdateReplace || node.DirectoryTarget || node.ExpectedSHA256 == "" {
			continue
		}
		records = append(records, map[string]any{
			"path":            node.Target,
			"component":       node.Component,
			"source":          node.Source,
			"ownership":       "managed",
			"checksum_sha256": node.ExpectedSHA256,
		})
	}
	return records
}

// ManagedSectionRecords returns deterministic lock records for bounded shared-file sections.
func ManagedSectionRecords(graph state.Graph) []any {
	records := []any{}
	for _, node := range graph.Nodes {
		if node.Update != state.UpdateManagedSection || node.Source == "" {
			continue
		}
		data, err := kitfs.ReadFile(node.Source)
		if err != nil {
			continue
		}
		block, err := managedBlock(string(data))
		if err != nil {
			continue
		}
		digest := sha256.Sum256([]byte(block))
		records = append(records, map[string]any{
			"path":            node.Target,
			"begin_marker":    BeginMarker,
			"end_marker":      EndMarker,
			"checksum_sha256": hex.EncodeToString(digest[:]),
		})
	}
	return records
}

// GeneratedSurfaceRecords returns deterministic lock records for non-managed graph nodes.
func GeneratedSurfaceRecords(graph state.Graph) []any {
	records := []any{}
	for _, node := range graph.Nodes {
		if node.Ownership == state.OwnershipManaged {
			continue
		}
		record := map[string]any{
			"path":      node.Target,
			"ownership": string(node.Ownership),
		}
		if node.Component != "" {
			record["component"] = node.Component
		}
		if node.Source != "" {
			record["source"] = node.Source
		}
		if node.ExpectedSHA256 != "" && !node.DirectoryTarget {
			record["checksum_sha256"] = node.ExpectedSHA256
		}
		records = append(records, record)
	}
	return records
}

func lockMateriallyEqual(existing, desired map[string]any) bool {
	if existing == nil || desired == nil || state.AsInt(existing["schema_version"]) != state.AsInt(desired["schema_version"]) {
		return false
	}
	left := state.DeepCopy(existing)
	right := state.DeepCopy(desired)
	for _, value := range []map[string]any{left, right} {
		delete(value, "state_generation")
		delete(value, "last_operation")
	}
	leftData, _ := json.Marshal(left)
	rightData, _ := json.Marshal(right)
	return string(leftData) == string(rightData)
}
