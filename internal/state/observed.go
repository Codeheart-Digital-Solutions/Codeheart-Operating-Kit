package state

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/kitfs"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

type Classification string

const (
	StateAbsent                   Classification = "absent"
	StateAdoptable                Classification = "adoptable"
	StateCurrent                  Classification = "current"
	StateDrifted                  Classification = "drifted"
	StateStaleCLI                 Classification = "stale-cli"
	StatePartial                  Classification = "partial"
	StateSchemaInvalid            Classification = "schema-invalid"
	StateLegacyV1Compatible       Classification = "legacy-v1-compatible"
	StateTransactionInProgress    Classification = "transaction-in-progress"
	StateRecoveryRequired         Classification = "recovery-required"
	StateUnsupportedFutureVersion Classification = "unsupported-future-version"
)

const (
	LockPath        = ".codeheart/kit.lock.yaml"
	ConfigPath      = ".codeheart/kit.config.yaml"
	TransactionPath = ".codeheart/kit.transaction.json"
)

type Observed struct {
	Root              string         `json:"root"`
	CanonicalRoot     string         `json:"canonical_root,omitempty"`
	Classification    Classification `json:"classification"`
	Traits            []string       `json:"traits,omitempty"`
	LockSchemaVersion int            `json:"lock_schema_version,omitempty"`
	Lock              map[string]any `json:"lock,omitempty"`
	Config            map[string]any `json:"config,omitempty"`
	Graph             Graph          `json:"graph,omitempty"`
	MissingPaths      []string       `json:"missing_paths,omitempty"`
	DriftedPaths      []string       `json:"drifted_paths,omitempty"`
	Errors            []string       `json:"errors,omitempty"`
	LegacyAnomalies   []string       `json:"legacy_anomalies,omitempty"`
	Transaction       map[string]any `json:"transaction,omitempty"`
}

// Inspect classifies an installation without mutating it.
func Inspect(root string) (Observed, error) {
	return inspect(root, false)
}

// InspectIgnoringTransaction validates staged or just-committed state while the reconciler holds
// the transaction marker. Normal callers must use Inspect.
func InspectIgnoringTransaction(root string) (Observed, error) {
	return inspect(root, true)
}

func inspect(root string, ignoreTransaction bool) (Observed, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return Observed{}, err
	}
	result := Observed{Root: abs}
	info, err := os.Stat(abs)
	if os.IsNotExist(err) {
		result.Classification = StateAbsent
		return result, nil
	}
	if err != nil {
		return result, err
	}
	if !info.IsDir() {
		result.Classification = StateSchemaInvalid
		result.Errors = []string{"target is not a directory"}
		return result, nil
	}
	canonical, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return result, err
	}
	result.CanonicalRoot = canonical

	if transactionData, err := os.ReadFile(filepath.Join(canonical, filepath.FromSlash(TransactionPath))); err == nil && !ignoreTransaction {
		var marker map[string]any
		decoder := json.NewDecoder(bytes.NewReader(transactionData))
		decoder.UseNumber()
		if err := decoder.Decode(&marker); err != nil {
			result.Classification = StateRecoveryRequired
			result.Errors = []string{"transaction marker is invalid: " + err.Error()}
			return result, nil
		}
		result.Transaction = marker
		phase := AsString(marker["phase"])
		if phase == "failed" || phase == "recovery-required" {
			result.Classification = StateRecoveryRequired
		} else {
			result.Classification = StateTransactionInProgress
		}
		return result, nil
	} else if err != nil && !os.IsNotExist(err) {
		return result, err
	}

	lockExists := pathExists(filepath.Join(canonical, filepath.FromSlash(LockPath)))
	configExists := pathExists(filepath.Join(canonical, filepath.FromSlash(ConfigPath)))
	kitExists := pathExists(filepath.Join(canonical, ".codeheart", "kit"))
	agentsExists := pathExists(filepath.Join(canonical, "AGENTS.md"))
	if !lockExists && !configExists && !kitExists {
		result.Classification = StateAdoptable
		return result, nil
	}
	if !lockExists {
		result.Classification = StatePartial
		for path, exists := range map[string]bool{LockPath: lockExists, ConfigPath: configExists, ".codeheart/kit/": kitExists, "AGENTS.md": agentsExists} {
			if !exists {
				result.MissingPaths = append(result.MissingPaths, path)
			}
		}
		sort.Strings(result.MissingPaths)
		return result, nil
	}

	lockBytes, err := os.ReadFile(filepath.Join(canonical, filepath.FromSlash(LockPath)))
	if err != nil {
		return result, err
	}
	lock, err := DecodeYAMLMap(lockBytes)
	if err != nil {
		result.Classification = StateSchemaInvalid
		result.Errors = []string{err.Error()}
		return result, nil
	}
	result.LockSchemaVersion = AsInt(lock["schema_version"])
	if result.LockSchemaVersion > 2 {
		result.Classification = StateUnsupportedFutureVersion
		result.Lock = lock
		return result, nil
	}
	validatedLock := lock
	if result.LockSchemaVersion == 1 {
		validatedLock, result.LegacyAnomalies = NormalizeLegacyV1(lock)
		result.Traits = append(result.Traits, "lock-v1")
	}
	schemaPath, err := SchemaForLockVersion(result.LockSchemaVersion)
	if err != nil {
		result.Classification = StateSchemaInvalid
		result.Errors = []string{err.Error()}
		return result, nil
	}
	if err := Validate(schemaPath, validatedLock); err != nil {
		result.Classification = StateSchemaInvalid
		result.Errors = []string{err.Error()}
		result.Lock = lock
		return result, nil
	}
	result.Lock = validatedLock
	if !configExists || !kitExists || !agentsExists {
		result.Classification = StatePartial
		for path, exists := range map[string]bool{ConfigPath: configExists, ".codeheart/kit/": kitExists, "AGENTS.md": agentsExists} {
			if !exists {
				result.MissingPaths = append(result.MissingPaths, path)
			}
		}
		sort.Strings(result.MissingPaths)
		return result, nil
	}

	configBytes, err := os.ReadFile(filepath.Join(canonical, filepath.FromSlash(ConfigPath)))
	if err != nil {
		return result, err
	}
	config, err := DecodeAndValidateYAML(ConfigV1Schema, configBytes)
	if err != nil {
		result.Classification = StateSchemaInvalid
		result.Errors = []string{err.Error()}
		return result, nil
	}
	result.Config = config
	profileID := AsString(validatedLock["selected_profile"])
	if profileID == "" {
		profileID = "standard"
	}
	graph, err := CompileGraph(profileID)
	if err != nil {
		result.Classification = StateSchemaInvalid
		result.Errors = []string{err.Error()}
		return result, nil
	}
	result.Graph = graph
	for _, node := range graph.Nodes {
		if node.Presence != PresenceRequired {
			continue
		}
		// The lock and config are validated structurally above. Their graph sources are
		// schemas rather than desired file bytes, so content-digest comparison would
		// always report a false drift.
		if node.Target == LockPath || node.Target == ConfigPath {
			continue
		}
		target := filepath.Join(canonical, filepath.FromSlash(strings.TrimSuffix(node.Target, "/")))
		if !pathExists(target) {
			result.MissingPaths = append(result.MissingPaths, node.Target)
			continue
		}
		if node.Update == UpdateManagedSection {
			actual, err := os.ReadFile(target)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("read %s: %v", node.Target, err))
				continue
			}
			expected, err := kitfs.ReadFile(node.Source)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("read desired %s: %v", node.Source, err))
				continue
			}
			actualBlock, actualOK := boundedManagedBlock(string(actual))
			expectedBlock, expectedOK := boundedManagedBlock(string(expected))
			if !actualOK || !expectedOK || actualBlock != expectedBlock {
				result.DriftedPaths = append(result.DriftedPaths, node.Target)
			}
			continue
		}
		if node.ExpectedSHA256 != "" && node.Update == UpdateReplace {
			data, err := os.ReadFile(target)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("read %s: %v", node.Target, err))
				continue
			}
			digest := sha256.Sum256(data)
			if hex.EncodeToString(digest[:]) != node.ExpectedSHA256 {
				result.DriftedPaths = append(result.DriftedPaths, node.Target)
			}
		}
	}
	if result.LockSchemaVersion == 2 {
		validateGraphRecords(&result, graph, validatedLock)
	}
	sort.Strings(result.MissingPaths)
	sort.Strings(result.DriftedPaths)
	if len(result.Errors) > 0 {
		result.Classification = StateSchemaInvalid
	} else if len(result.MissingPaths) > 0 {
		result.Classification = StatePartial
	} else if len(result.LegacyAnomalies) > 0 {
		result.Classification = StateLegacyV1Compatible
	} else if len(result.DriftedPaths) > 0 {
		result.Classification = StateDrifted
	} else if AsString(validatedLock["kit_version"]) != "" && AsString(validatedLock["kit_version"]) != version.Version {
		result.Classification = StateStaleCLI
	} else {
		result.Classification = StateCurrent
	}
	return result, nil
}

const (
	managedBeginMarker = "<!-- BEGIN CODEHEART OPERATING KIT MANAGED BLOCK -->"
	managedEndMarker   = "<!-- END CODEHEART OPERATING KIT MANAGED BLOCK -->"
)

func boundedManagedBlock(text string) (string, bool) {
	_, rest, begin := strings.Cut(text, managedBeginMarker)
	if !begin {
		return "", false
	}
	middle, _, end := strings.Cut(rest, managedEndMarker)
	if !end {
		return "", false
	}
	return managedBeginMarker + middle + managedEndMarker, true
}

func validateGraphRecords(result *Observed, graph Graph, lock map[string]any) {
	managed := recordIndex(lock["managed_paths"])
	sections := recordIndex(lock["managed_sections"])
	generated := recordIndex(lock["generated_surfaces"])
	for _, node := range graph.Nodes {
		switch {
		case node.Ownership == OwnershipManaged && node.Update == UpdateReplace && !node.DirectoryTarget:
			record := managed[node.Target]
			if record == nil || AsString(record["checksum_sha256"]) != node.ExpectedSHA256 {
				result.DriftedPaths = append(result.DriftedPaths, LockPath+":managed_paths:"+node.Target)
			}
		case node.Update == UpdateManagedSection:
			record := sections[node.Target]
			expected, err := kitfs.ReadFile(node.Source)
			if err != nil {
				continue
			}
			block, ok := boundedManagedBlock(string(expected))
			digest := sha256.Sum256([]byte(block))
			if !ok || record == nil || AsString(record["checksum_sha256"]) != hex.EncodeToString(digest[:]) {
				result.DriftedPaths = append(result.DriftedPaths, LockPath+":managed_sections:"+node.Target)
			}
		case node.Ownership != OwnershipManaged:
			record := generated[node.Target]
			if record == nil || AsString(record["ownership"]) != string(node.Ownership) {
				result.DriftedPaths = append(result.DriftedPaths, LockPath+":generated_surfaces:"+node.Target)
			}
		}
	}
}

func recordIndex(value any) map[string]map[string]any {
	result := map[string]map[string]any{}
	for _, item := range AnySlice(value) {
		record := Map(item)
		if path := AsString(record["path"]); path != "" {
			result[path] = record
		}
	}
	return result
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
