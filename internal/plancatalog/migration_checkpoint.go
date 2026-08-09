package plancatalog

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

type ledgerCheckpoint struct {
	EvidenceRevision       string
	ActivationBaseRevision string
	LedgerPath             string
	LedgerSHA256           string
}

func LedgerActivationBaseRevision(root string, ledger MigrationLedger) (string, error) {
	checkpoint, problems := verifyLedgerCheckpoint(root, ledger.EvidenceRevision, ledger.ArtifactPath, ledger.ArtifactSHA256)
	if HasErrors(problems) || checkpoint.ActivationBaseRevision == "" {
		if len(problems) > 0 {
			return "", fmt.Errorf("%s: %s", problems[0].Code, problems[0].Message)
		}
		return "", fmt.Errorf("ledger_checkpoint_unavailable: activation base could not be resolved")
	}
	return checkpoint.ActivationBaseRevision, nil
}

func ActivationCheckpointRevision(root, activationBase string) (string, error) {
	history, err := gitText(root, "rev-list", "--first-parent", "--reverse", activationBase+"..HEAD")
	if err != nil || strings.TrimSpace(history) == "" {
		return "", fmt.Errorf("activation_checkpoint_missing: no first-parent activation checkpoint follows the ledger checkpoint")
	}
	activation := strings.Fields(history)[0]
	parent, err := gitText(root, "rev-parse", "--verify", activation+"^1")
	if err != nil || parent != activationBase {
		return "", fmt.Errorf("activation_checkpoint_parent_mismatch: activation checkpoint is not the immediate first-parent child of the ledger checkpoint")
	}
	return activation, nil
}

func bindLedgerArtifact(root, ledgerPath string, data []byte) (string, string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", "", err
	}
	ledgerAbs, err := filepath.Abs(ledgerPath)
	if err != nil {
		return "", "", err
	}
	relative, err := filepath.Rel(rootAbs, ledgerAbs)
	if err != nil {
		return "", "", err
	}
	relative = filepath.ToSlash(relative)
	if problem := validateGitPath(relative); problem != nil || strings.HasPrefix(relative, "../") || relative == ".." {
		return "", "", fmt.Errorf("ledger_path_invalid: ledger must be a normalized repository-relative path")
	}
	if protectedMigrationPath(relative) {
		return "", "", fmt.Errorf("ledger_path_invalid: ledger must be outside protected catalog paths")
	}
	digest := sha256.Sum256(data)
	return relative, hex.EncodeToString(digest[:]), nil
}

func verifyLedgerCheckpoint(root, evidenceRevision, ledgerPath, ledgerSHA string) (ledgerCheckpoint, []Problem) {
	checkpoint := ledgerCheckpoint{EvidenceRevision: evidenceRevision, LedgerPath: ledgerPath, LedgerSHA256: ledgerSHA}
	problems := []Problem{}
	head, err := gitText(root, "rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		return checkpoint, []Problem{{Code: "ledger_checkpoint_unavailable", Message: err.Error(), Severity: SeverityError}}
	}
	checkpoint.ActivationBaseRevision = head
	parent, err := gitText(root, "rev-parse", "--verify", head+"^1")
	if err != nil || parent != evidenceRevision {
		message := "ledger checkpoint must have the reviewed evidence revision as its first parent"
		if err != nil {
			message = err.Error()
		}
		problems = append(problems, Problem{Code: "ledger_checkpoint_parent_mismatch", Message: message, Path: ledgerPath, Severity: SeverityError, Remediation: "reinventory and review from a new evidence revision"})
	}
	changed, diffErr := gitNULFields(root, "diff-tree", "--no-commit-id", "--name-only", "-z", "-r", evidenceRevision, head)
	if diffErr != nil || len(changed) != 1 || changed[0] != ledgerPath {
		message := fmt.Sprintf("ledger checkpoint delta is %v; expected only %s", changed, ledgerPath)
		if diffErr != nil {
			message = diffErr.Error()
		}
		problems = append(problems, Problem{Code: "ledger_checkpoint_delta_mismatch", Message: message, Path: ledgerPath, Severity: SeverityError, Remediation: "commit the reviewed ledger as the only tree change directly after the evidence revision"})
	}
	if !gitRevisionHasRegularFile(root, head, ledgerPath) {
		problems = append(problems, Problem{Code: "ledger_checkpoint_blob_invalid", Message: "ledger checkpoint does not contain a regular ledger blob", Path: ledgerPath, Severity: SeverityError})
	} else if committed, readErr := gitBytes(root, "show", head+":"+ledgerPath); readErr != nil || sha256Text(committed) != ledgerSHA {
		message := "committed ledger bytes do not match the reviewed ledger hash"
		if readErr != nil {
			message = readErr.Error()
		}
		problems = append(problems, Problem{Code: "ledger_checkpoint_blob_mismatch", Message: message, Path: ledgerPath, Severity: SeverityError})
	}
	if clean, cleanErr := gitIndexMatchesHead(root); cleanErr != nil || !clean {
		message := "index must exactly match the ledger checkpoint"
		if cleanErr != nil {
			message = cleanErr.Error()
		}
		problems = append(problems, Problem{Code: "ledger_checkpoint_index_dirty", Message: message, Severity: SeverityError})
	}
	return checkpoint, problems
}

func gitNULFields(root string, args ...string) ([]string, error) {
	data, err := gitBytes(root, args...)
	if err != nil {
		return nil, err
	}
	fields := []string{}
	for _, raw := range bytes.Split(data, []byte{0}) {
		if len(raw) > 0 {
			fields = append(fields, string(raw))
		}
	}
	return fields, nil
}

func gitIndexMatchesHead(root string) (bool, error) {
	command := exec.Command("git", "-C", root, "diff", "--cached", "--quiet", "--exit-code", "HEAD", "--")
	command.Env = append(command.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_NO_LAZY_FETCH=1", "GIT_NO_REPLACE_OBJECTS=1", "LC_ALL=C")
	err := command.Run()
	if err == nil {
		return true, nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, fmt.Errorf("git_index_unavailable: %w", err)
}

func gitWorktreeClean(root string) (bool, error) {
	fields, err := gitNULFields(root, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	return len(fields) == 0, err
}

type migrationActionEvidence struct {
	Kind           string `json:"kind"`
	Target         string `json:"target"`
	Mode           uint32 `json:"mode,omitempty"`
	ContentSHA256  string `json:"content_sha256,omitempty"`
	ExpectedSHA256 string `json:"expected_sha256,omitempty"`
}

func migrationActionDigest(actions []reconcile.Action) string {
	items := make([]migrationActionEvidence, 0, len(actions))
	for _, action := range actions {
		item := migrationActionEvidence{Kind: action.Kind, Target: action.Target, Mode: action.Mode, ExpectedSHA256: action.ExpectedSHA256}
		if action.Kind != "remove" {
			item.ContentSHA256 = sha256Text(action.Content)
		}
		items = append(items, item)
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Target != items[j].Target {
			return items[i].Target < items[j].Target
		}
		return items[i].Kind < items[j].Kind
	})
	data, _ := json.Marshal(items)
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

// MigrationActionDigestForCheckpoint reconstructs the reviewed plan actions
// from the immutable L..A commit delta. The guarded config write is excluded
// because the persisted binding stores the digest of plan actions only.
func MigrationActionDigestForCheckpoint(root, activationBase, activationRevision string) (string, []string, error) {
	changed, err := gitNULFields(root, "diff-tree", "--no-commit-id", "--name-only", "-z", "-r", activationBase, activationRevision)
	if err != nil {
		return "", nil, err
	}
	changed = uniqueSortedStrings(changed)
	actions := make([]reconcile.Action, 0, len(changed))
	planPaths := make([]string, 0, len(changed))
	for _, target := range changed {
		if target == state.ConfigPath {
			continue
		}
		before, beforePresent, beforeMode, beforeErr := checkpointFileState(root, activationBase, target)
		after, afterPresent, afterMode, afterErr := checkpointFileState(root, activationRevision, target)
		if beforeErr != nil || afterErr != nil || (!beforePresent && !afterPresent) {
			return "", nil, fmt.Errorf("activation_checkpoint_blob_mismatch: cannot reconstruct %s", target)
		}
		action := reconcile.Action{Target: target, Owner: "repo-plan"}
		switch {
		case beforePresent && !afterPresent:
			action.Kind = "remove"
			action.ExpectedSHA256 = sha256Text(before)
		case !beforePresent && afterPresent:
			action.Kind = "create"
			action.Content = after
			action.Mode = afterMode
		case beforePresent && afterPresent:
			action.Kind = "replace"
			action.Content = after
			action.Mode = afterMode
			action.ExpectedSHA256 = sha256Text(before)
		}
		if beforePresent && beforeMode != 0 && action.Kind == "remove" {
			// Removal actions intentionally do not carry a mode.
			action.Mode = 0
		}
		actions = append(actions, action)
		planPaths = append(planPaths, target)
	}
	return migrationActionDigest(actions), planPaths, nil
}

func checkpointFileState(root, revision, target string) ([]byte, bool, uint32, error) {
	output, err := gitText(root, "ls-tree", revision, "--", target)
	if err != nil {
		return nil, false, 0, err
	}
	if output == "" {
		return nil, false, 0, nil
	}
	fields := strings.Fields(output)
	if len(fields) < 3 || fields[1] != "blob" || (fields[0] != "100644" && fields[0] != "100755") {
		return nil, false, 0, fmt.Errorf("unsupported checkpoint path state")
	}
	data, err := gitBytes(root, "show", revision+":"+target)
	if err != nil {
		return nil, false, 0, err
	}
	mode := uint32(0o644)
	if fields[0] == "100755" {
		mode = 0o755
	}
	return data, true, mode, nil
}

func verifyMigrationWorktree(root string, actions []reconcile.Action) []Problem {
	problems := []Problem{}
	clean, err := gitIndexMatchesHead(root)
	if err != nil || !clean {
		message := "index must exactly match the activation base revision"
		if err != nil {
			message = err.Error()
		}
		problems = append(problems, Problem{Code: "catalog_activation_index_dirty", Message: message, Severity: SeverityError})
	}
	tracked, trackedErr := gitNULFields(root, "diff", "--name-only", "-z", "HEAD", "--")
	untracked, untrackedErr := gitNULFields(root, "ls-files", "--others", "--exclude-standard", "-z")
	if trackedErr != nil || untrackedErr != nil {
		message := "worktree delta could not be enumerated"
		if trackedErr != nil {
			message = trackedErr.Error()
		} else if untrackedErr != nil {
			message = untrackedErr.Error()
		}
		return append(problems, Problem{Code: "catalog_activation_worktree_unavailable", Message: message, Severity: SeverityError})
	}
	actualPaths := append(tracked, untracked...)
	sort.Strings(actualPaths)
	expectedPaths := make([]string, 0, len(actions))
	for _, action := range actions {
		expectedPaths = append(expectedPaths, action.Target)
		absolute := filepath.Join(root, filepath.FromSlash(action.Target))
		if action.Kind == "remove" {
			if _, statErr := os.Lstat(absolute); !os.IsNotExist(statErr) {
				problems = append(problems, Problem{Code: "catalog_activation_action_mismatch", Message: "reviewed removal is not present in the worktree", Path: action.Target, Severity: SeverityError})
			}
			continue
		}
		data, readErr := readRegularSource(root, action.Target)
		matches := false
		var matchErr error
		if readErr == nil {
			matches, matchErr = checkoutLineEndingEquivalentForPath(root, action.Target, data, action.Content)
		}
		if readErr != nil || matchErr != nil || !matches {
			message := "worktree bytes do not match the reviewed migration output"
			if readErr != nil {
				message = readErr.Error()
			} else if matchErr != nil {
				message = matchErr.Error()
			}
			problems = append(problems, Problem{Code: "catalog_activation_action_mismatch", Message: message, Path: action.Target, Severity: SeverityError})
		}
	}
	sort.Strings(expectedPaths)
	if !equalStrings(actualPaths, expectedPaths) {
		problems = append(problems, Problem{Code: "catalog_activation_worktree_mismatch", Message: fmt.Sprintf("worktree paths %v do not match reviewed migration paths %v", actualPaths, expectedPaths), Severity: SeverityError, Remediation: "preserve unrelated work and activate from a worktree containing only the exact migration projection"})
	}
	SortProblems(problems)
	return problems
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func configPrecondition(root string) (string, error) {
	data, err := readRegularSource(root, state.ConfigPath)
	if err != nil {
		return "", err
	}
	return sha256Text(data), nil
}
