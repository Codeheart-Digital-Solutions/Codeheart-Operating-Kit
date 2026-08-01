package commands

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/plancatalog"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/portfolio"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

type plansValidationOutput struct {
	SchemaVersion  int                     `json:"schema_version"`
	Mode           plancatalog.CatalogMode `json:"mode"`
	RepositoryID   string                  `json:"repository_id,omitempty"`
	Valid          bool                    `json:"valid"`
	RecordCount    int                     `json:"record_count"`
	Problems       []plancatalog.Problem   `json:"problems"`
	RemoteOverlays *portfolio.ScanResult   `json:"remote_overlays,omitempty"`
}

func RunPlans(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		return writeArgError(stderr, "plans", fmt.Errorf("the following arguments are required: subcommand"))
	}
	switch args[0] {
	case "validate":
		return runPlansValidate(args[1:], stdout, stderr)
	case "list":
		return runPlansList(args[1:], stdout, stderr)
	case "inventory":
		return runPlansInventory(args[1:], stdout, stderr)
	case "migrate":
		return runPlansMigrate(args[1:], stdout, stderr)
	default:
		return writeArgError(stderr, "plans", fmt.Errorf("invalid subcommand %q (choose from validate, list, inventory, migrate)", args[0]))
	}
}

func runPlansValidate(args []string, stdout io.Writer, stderr io.Writer) int {
	_, bools, positionals, err := parseValueArgs(args, map[string]bool{"--json": false, "--remote-overlays": false})
	if err != nil {
		return writeArgError(stderr, "plans validate", err)
	}
	root, err := singleRoot(positionals)
	if err != nil {
		return writeArgError(stderr, "plans validate", err)
	}
	snapshot, err := plancatalog.LoadRepositorySnapshot(root)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans validate: error: %v\n", err)
		return 1
	}
	output := plansValidationOutput{SchemaVersion: 1, Mode: snapshot.Settings.Mode, RepositoryID: snapshot.Settings.RepositoryID, Valid: !plancatalog.HasErrors(snapshot.Problems), RecordCount: len(snapshot.Records), Problems: snapshot.Problems}
	if bools["--remote-overlays"] {
		remote, scanErr := scanRemoteOverlays(root)
		if scanErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans validate: error: %v\n", scanErr)
			return 1
		}
		output.RemoteOverlays = &remote
		if !remote.Catalog.Complete || !catalogContainsMember(remote.Catalog, snapshot.Settings.RepositoryID) {
			output.Valid = false
		}
	}
	if bools["--json"] {
		if err := writeJSON(stdout, output); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans validate: error: %v\n", err)
			return 1
		}
	} else {
		fmt.Fprintf(stdout, "Plan validation (%s mode): %d record(s); valid=%t.\n", output.Mode, output.RecordCount, output.Valid)
		for _, problem := range output.Problems {
			fmt.Fprintf(stdout, "- %s %s: %s", problem.Severity, problem.Code, problem.Message)
			if problem.Path != "" {
				fmt.Fprintf(stdout, " (%s)", problem.Path)
			}
			fmt.Fprintln(stdout)
		}
		if output.RemoteOverlays != nil {
			fmt.Fprintf(stdout, "Remote overlays: complete=%t; %d pushed observation(s) across %d member(s). Local heads, worktree changes, and unpushed commits are omitted.\n", output.RemoteOverlays.Catalog.Complete, len(output.RemoteOverlays.Catalog.Observations), len(output.RemoteOverlays.Catalog.Members))
		}
	}
	if output.Valid {
		return 0
	}
	return 1
}

func runPlansList(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{"--format": true, "--json": false})
	if err != nil {
		return writeArgError(stderr, "plans list", err)
	}
	format := values["--format"]
	switch format {
	case "", "text", "json":
	default:
		return writeArgError(stderr, "plans list", fmt.Errorf("invalid --format %q (choose from text, json)", format))
	}
	if bools["--json"] {
		if format == "text" {
			return writeArgError(stderr, "plans list", fmt.Errorf("--json conflicts with --format text"))
		}
		format = "json"
	}
	if format == "" {
		format = "text"
	}
	root, err := singleRoot(positionals)
	if err != nil {
		return writeArgError(stderr, "plans list", err)
	}
	snapshot, err := plancatalog.LoadRepositorySnapshot(root)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans list: error: %v\n", err)
		return 1
	}
	view := plancatalog.BuildView(snapshot)
	if format == "json" {
		if err := plancatalog.WriteViewJSON(stdout, view); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans list: error: %v\n", err)
			return 1
		}
	} else {
		plancatalog.WriteViewText(stdout, view)
	}
	if plancatalog.HasErrors(view.Problems) {
		return 1
	}
	return 0
}

func runPlansInventory(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{"--output": true, "--json": false, "--remote-overlays": false})
	if err != nil {
		return writeArgError(stderr, "plans inventory", err)
	}
	if values["--output"] == "" {
		return writeArgError(stderr, "plans inventory", fmt.Errorf("option --output requires a value"))
	}
	root, err := singleRoot(positionals)
	if err != nil {
		return writeArgError(stderr, "plans inventory", err)
	}
	inventory, err := plancatalog.BuildInventory(root, time.Now())
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans inventory: error: %v\n", err)
		return 1
	}
	if bools["--remote-overlays"] {
		remote, scanErr := scanRemoteOverlays(root)
		if scanErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans inventory: error: %v\n", scanErr)
			return 1
		}
		inventory.RemoteOverlays = remote
		if !remote.Catalog.Complete || !catalogContainsMember(remote.Catalog, inventory.RepositoryID) {
			inventory.Problems = append(inventory.Problems, plancatalog.Problem{Code: "remote_overlay_coverage_incomplete", Message: "pushed default and unmerged-branch evidence is incomplete for this repository", Severity: plancatalog.SeverityError})
		}
		plancatalog.SortProblems(inventory.Problems)
	}
	if err := writeInventoryArtifact(root, values["--output"], inventory); err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans inventory: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, inventory); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans inventory: error: %v\n", err)
			return 1
		}
	} else {
		fmt.Fprintf(stdout, "Plan inventory: %d formal record(s), %d canonical, %d legacy, %d invalid, %d unreadable, %d unsafe, %d unpaired legacy; revision %s.\n", inventory.Coverage.FormalRecords, inventory.Coverage.CanonicalMetadata, inventory.Coverage.LegacyRecords, inventory.Coverage.InvalidRecords, inventory.Coverage.UnreadableRecords, inventory.Coverage.UnsafeRecords, inventory.Coverage.UnpairedLegacyEvidence, inventory.SourceRevision)
		fmt.Fprintf(stdout, "Wrote inventory: %s\n", values["--output"])
	}
	if plancatalog.HasErrors(inventory.Problems) {
		return 1
	}
	return 0
}

func scanRemoteOverlays(root string) (portfolio.ScanResult, error) {
	config, err := portfolio.LoadConfig(root)
	if err != nil {
		return portfolio.ScanResult{}, err
	}
	runner := portfolio.ExecRunner{}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	result, err := portfolio.Scan(ctx, portfolio.ScanOptions{Root: root, Config: config, Sources: portfolio.BuildSources(config, runner), Runner: runner, WriteCache: false})
	if err != nil {
		return portfolio.ScanResult{}, err
	}
	return publicRemoteOverlayEvidence(result), nil
}

// publicRemoteOverlayEvidence removes machine-local discovery paths from plan
// validation and migration artifacts. The portfolio cache remains ignored local
// evidence and retains the scanner's exact locator for operational recovery.
func publicRemoteOverlayEvidence(result portfolio.ScanResult) portfolio.ScanResult {
	for index := range result.Catalog.Members {
		member := &result.Catalog.Members[index]
		member.SourceLocator = publicSourceLocator(member.SourceLocator, member.RepositoryID)
	}
	for index := range result.Catalog.Candidates {
		candidate := &result.Catalog.Candidates[index]
		candidate.SourceLocator = publicSourceLocator(candidate.SourceLocator, candidate.RepositoryID)
	}
	for index := range result.Catalog.Errors {
		scanError := &result.Catalog.Errors[index]
		scanError.SourceLocator = publicSourceLocator(scanError.SourceLocator, scanError.RepositoryID)
	}
	return result
}

func publicSourceLocator(locator, repositoryID string) string {
	trimmed := strings.TrimSpace(locator)
	if trimmed == "" || (!filepath.IsAbs(trimmed) && !strings.HasPrefix(strings.ToLower(trimmed), "file:")) {
		return locator
	}
	if repositoryID != "" {
		return "local-git:" + repositoryID
	}
	return "local-git"
}

func catalogContainsMember(catalog portfolio.Catalog, repositoryID string) bool {
	for _, member := range catalog.Members {
		if member.RepositoryID == repositoryID {
			return true
		}
	}
	return false
}

func runPlansMigrate(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{"--ledger": true, "--dry-run": false, "--yes": false, "--json": false})
	if err != nil {
		return writeArgError(stderr, "plans migrate", err)
	}
	if values["--ledger"] == "" {
		return writeArgError(stderr, "plans migrate", fmt.Errorf("option --ledger requires a value"))
	}
	if bools["--dry-run"] == bools["--yes"] {
		return writeArgError(stderr, "plans migrate", fmt.Errorf("choose exactly one of --dry-run or --yes"))
	}
	root, err := singleRoot(positionals)
	if err != nil {
		return writeArgError(stderr, "plans migrate", err)
	}
	ledgerData, err := os.ReadFile(values["--ledger"])
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans migrate: error: invalid_ledger: %v\n", err)
		return 1
	}
	ledger, err := plancatalog.LoadMigrationLedger(ledgerData)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans migrate: error: %v\n", err)
		return 1
	}
	plan, err := plancatalog.BuildMigrationPlan(root, ledger)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans migrate: error: %v\n", err)
		return 1
	}
	outcome, err := plancatalog.ExecuteMigration(plan, plancatalog.MigrationApplyOptions{DryRun: bools["--dry-run"], Now: time.Now()})
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans migrate: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, outcome); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans migrate: error: %v\n", err)
			return 1
		}
	} else {
		reconcile.WriteText(stdout, outcome.Result)
		for _, skip := range outcome.Skips {
			fmt.Fprintf(stdout, "- skipped %s: %s (%s)\n", skip.Code, skip.Message, skip.Path)
		}
	}
	if !outcome.Result.OK() || plancatalog.HasMaterialMigrationSkips(outcome.Skips) {
		return 1
	}
	return 0
}

func singleRoot(positionals []string) (string, error) {
	root := "."
	if len(positionals) > 0 {
		root = expandPath(positionals[0])
	}
	if len(positionals) > 1 {
		return "", fmt.Errorf("unexpected argument %q", positionals[1])
	}
	return root, nil
}

func writeInventoryArtifact(root, destination string, inventory plancatalog.Inventory) error {
	return writeInventoryArtifactWithHook(root, destination, inventory, nil)
}

func writeInventoryArtifactWithHook(root, destination string, inventory plancatalog.Inventory, hook func(string) error) error {
	destination = expandPath(destination)
	absDestination, err := filepath.Abs(destination)
	if err != nil {
		return err
	}
	rootAbsolute, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	parentCandidate := filepath.Dir(absDestination)
	parentCandidateInfo, err := os.Lstat(parentCandidate)
	if os.IsNotExist(err) {
		return fmt.Errorf("inventory_target_parent_missing: create the explicit output directory before writing inventory")
	}
	if err != nil {
		return err
	}
	if !parentCandidateInfo.IsDir() || parentCandidateInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("inventory_target_unsafe: output parent is not a regular directory")
	}
	parentPath, err := filepath.EvalSymlinks(parentCandidate)
	if err != nil {
		return err
	}
	parentRoot, err := os.OpenRoot(parentPath)
	if err != nil {
		return err
	}
	defer parentRoot.Close()
	parentInfo, err := parentRoot.Stat(".")
	if err != nil {
		return err
	}
	if err := validateBoundInventoryParent(parentPath, parentRoot, parentInfo); err != nil {
		return err
	}
	if err := invokeInventoryHook(hook, "inventory-parent-bound-before-protection"); err != nil {
		return err
	}
	if err := validateBoundInventoryParent(parentPath, parentRoot, parentInfo); err != nil {
		return err
	}
	absDestination = filepath.Join(parentPath, filepath.Base(absDestination))
	if err := rejectRepoLocalSymlinkTraversal(rootAbsolute, absDestination); err != nil {
		return err
	}
	if info, err := os.Lstat(absDestination); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("inventory_target_unsafe: output target is a symbolic link")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	resolvedDestination, err := resolvePathAllowMissing(absDestination)
	if err != nil {
		return fmt.Errorf("inventory_target_unsafe: resolve output target: %w", err)
	}
	resolvedRoot, err := filepath.EvalSymlinks(rootAbsolute)
	if err != nil {
		return err
	}
	if relative, relErr := filepath.Rel(resolvedRoot, resolvedDestination); relErr == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		repositoryPath := filepath.ToSlash(relative)
		if strings.EqualFold(repositoryPath, "docs/repo/plans/README.md") || strings.EqualFold(repositoryPath, plancatalog.LegacyRegisterPath) {
			return fmt.Errorf("inventory_target_is_plan: output cannot create or overwrite planning authority %s", repositoryPath)
		}
		if _, formal := plancatalog.FormalPathKind(resolvedRoot, repositoryPath); formal {
			return fmt.Errorf("inventory_target_is_plan: output cannot create or overwrite formal plan path %s", repositoryPath)
		}
	}
	protected := []string{"docs/repo/plans/README.md", "docs/repo/plans/plan-register.md"}
	candidates, err := plancatalog.Enumerate(root)
	if err != nil {
		return fmt.Errorf("enumerate protected plan targets: %w", err)
	}
	for _, candidate := range candidates {
		protected = append(protected, candidate.Path)
	}
	for _, path := range protected {
		canonical, err := filepath.Abs(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			return err
		}
		resolvedProtected, err := resolvePathAllowMissing(canonical)
		if err != nil {
			return err
		}
		if sameFileOrPath(resolvedDestination, resolvedProtected) {
			return fmt.Errorf("inventory_target_is_plan: output cannot overwrite planning authority %s", path)
		}
	}
	absDestination = resolvedDestination
	var data []byte
	if strings.EqualFold(filepath.Ext(absDestination), ".json") {
		data, err = json.MarshalIndent(inventory, "", "  ")
		data = append(data, '\n')
	} else {
		data, err = state.EncodeYAML(inventory)
	}
	if err != nil {
		return err
	}
	if !samePath(filepath.Join(parentPath, filepath.Base(absDestination)), absDestination) {
		return fmt.Errorf("inventory_target_unsafe: destination parent identity changed during creation")
	}
	if err := validateBoundInventoryParent(parentPath, parentRoot, parentInfo); err != nil {
		return err
	}
	targetName := filepath.Base(absDestination)
	if targetName == "." || targetName == string(filepath.Separator) {
		return fmt.Errorf("inventory_target_unsafe: output target must be a file")
	}
	existingInfo, statErr := parentRoot.Lstat(targetName)
	if statErr == nil && (!existingInfo.Mode().IsRegular() || existingInfo.Mode()&os.ModeSymlink != 0) {
		return fmt.Errorf("inventory_target_unsafe: existing output is not a regular file")
	}
	if statErr != nil && !os.IsNotExist(statErr) {
		return statErr
	}
	var existingFile *os.File
	var existingData []byte
	if statErr == nil {
		existingFile, err = parentRoot.Open(targetName)
		if err != nil {
			return err
		}
		defer func() {
			if existingFile != nil {
				_ = existingFile.Close()
			}
		}()
		openedInfo, err := existingFile.Stat()
		if err != nil || !os.SameFile(existingInfo, openedInfo) {
			return fmt.Errorf("inventory output identity changed while binding existing bytes")
		}
		existingData, err = io.ReadAll(existingFile)
		if err != nil {
			return err
		}
		currentInfo, err := parentRoot.Lstat(targetName)
		if err != nil || !os.SameFile(existingInfo, currentInfo) {
			return fmt.Errorf("inventory output identity changed while reading existing bytes")
		}
	}
	temporaryName, temporary, err := createInventoryTemp(parentRoot)
	if err != nil {
		return err
	}
	temporaryInfo, err := temporary.Stat()
	if err != nil {
		_ = temporary.Close()
		return err
	}
	defer parentRoot.Remove(temporaryName)
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := invokeInventoryHook(hook, "inventory-parent-bound"); err != nil {
		return err
	}
	if err := validateBoundInventoryParent(parentPath, parentRoot, parentInfo); err != nil {
		return err
	}
	var quarantine string
	if statErr == nil {
		quarantine, err = inventorySibling(parentRoot, targetName, "inventory-replaced")
		if err != nil {
			return err
		}
		if err := parentRoot.Rename(targetName, quarantine); err != nil {
			return err
		}
		capturedInfo, captureErr := parentRoot.Lstat(quarantine)
		if captureErr != nil || !os.SameFile(existingInfo, capturedInfo) {
			return fmt.Errorf("inventory output identity changed during replacement; retained captured entry at %s", filepath.Join(parentPath, quarantine))
		}
		if _, err := existingFile.Seek(0, io.SeekStart); err != nil {
			return err
		}
		capturedData, err := io.ReadAll(existingFile)
		if err != nil || !bytes.Equal(existingData, capturedData) {
			return fmt.Errorf("inventory output bytes changed during replacement; retained captured entry at %s", filepath.Join(parentPath, quarantine))
		}
		openedInfo, err := existingFile.Stat()
		if err != nil || !os.SameFile(existingInfo, openedInfo) {
			return fmt.Errorf("inventory output descriptor identity changed during replacement; retained captured entry at %s", filepath.Join(parentPath, quarantine))
		}
		if err := invokeInventoryHook(hook, "inventory-output-quarantined"); err != nil {
			return err
		}
	}
	if err := parentRoot.Link(temporaryName, targetName); err != nil {
		if quarantine != "" {
			if _, targetErr := parentRoot.Lstat(targetName); os.IsNotExist(targetErr) {
				_ = parentRoot.Link(quarantine, targetName)
			}
		}
		return fmt.Errorf("install inventory output without replacement: %w", err)
	}
	installedInfo, err := parentRoot.Lstat(targetName)
	if err != nil || !os.SameFile(temporaryInfo, installedInfo) {
		return fmt.Errorf("inventory output identity changed during installation")
	}
	installedData, err := parentRoot.ReadFile(targetName)
	if err != nil || !bytes.Equal(installedData, data) {
		return fmt.Errorf("inventory output bytes changed during installation")
	}
	if err := parentRoot.Remove(temporaryName); err != nil {
		return err
	}
	if quarantine != "" {
		capturedInfo, err := parentRoot.Lstat(quarantine)
		if err != nil || !os.SameFile(existingInfo, capturedInfo) {
			return fmt.Errorf("replaced inventory evidence changed before cleanup; retained at %s", filepath.Join(parentPath, quarantine))
		}
		if _, err := existingFile.Seek(0, io.SeekStart); err != nil {
			return err
		}
		capturedData, err := io.ReadAll(existingFile)
		if err != nil || !bytes.Equal(existingData, capturedData) {
			return fmt.Errorf("replaced inventory bytes changed before cleanup; retained at %s", filepath.Join(parentPath, quarantine))
		}
		openedInfo, err := existingFile.Stat()
		if err != nil || !os.SameFile(existingInfo, openedInfo) {
			return fmt.Errorf("replaced inventory descriptor identity changed before cleanup; retained at %s", filepath.Join(parentPath, quarantine))
		}
		if err := parentRoot.Remove(quarantine); err != nil {
			return err
		}
		if err := existingFile.Close(); err != nil {
			return err
		}
		existingFile = nil
	}
	if err := syncInventoryDirectory(parentRoot); err != nil {
		return err
	}
	if err := validateBoundInventoryParent(parentPath, parentRoot, parentInfo); err != nil {
		return fmt.Errorf("inventory output was installed through the bound parent but requested parent authority changed: %w", err)
	}
	return nil
}

func createInventoryTemp(root *os.Root) (string, *os.File, error) {
	for attempt := 0; attempt < 8; attempt++ {
		name, err := inventoryRandomName(".plan-inventory-")
		if err != nil {
			return "", nil, err
		}
		file, err := root.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
		if os.IsExist(err) {
			continue
		}
		return name, file, err
	}
	return "", nil, fmt.Errorf("could not reserve an inventory temporary file")
}

func inventorySibling(root *os.Root, target, purpose string) (string, error) {
	for attempt := 0; attempt < 8; attempt++ {
		name, err := inventoryRandomName("." + target + "." + purpose + "-")
		if err != nil {
			return "", err
		}
		if _, err := root.Lstat(name); os.IsNotExist(err) {
			return name, nil
		} else if err != nil {
			return "", err
		}
	}
	return "", fmt.Errorf("could not reserve an inventory replacement quarantine")
}

func inventoryRandomName(prefix string) (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(nonce), nil
}

func validateBoundInventoryParent(path string, root *os.Root, identity os.FileInfo) error {
	pathInfo, err := os.Lstat(path)
	if err != nil || identity == nil || !os.SameFile(identity, pathInfo) {
		return fmt.Errorf("inventory_target_unsafe: output parent path identity changed")
	}
	boundInfo, err := root.Stat(".")
	if err != nil || !os.SameFile(identity, boundInfo) {
		return fmt.Errorf("inventory_target_unsafe: output parent descriptor identity changed")
	}
	return nil
}

func syncInventoryDirectory(root *os.Root) error {
	directory, err := root.Open(".")
	if err != nil {
		return err
	}
	defer directory.Close()
	if err := directory.Sync(); err != nil && runtime.GOOS != "windows" {
		return err
	}
	return nil
}

func invokeInventoryHook(hook func(string) error, phase string) error {
	if hook == nil {
		return nil
	}
	return hook(phase)
}

func sameFileOrPath(left, right string) bool {
	leftInfo, leftErr := os.Stat(left)
	rightInfo, rightErr := os.Stat(right)
	if leftErr == nil && rightErr == nil && os.SameFile(leftInfo, rightInfo) {
		return true
	}
	return samePath(left, right)
}

func rejectRepoLocalSymlinkTraversal(root, target string) error {
	relative, err := filepath.Rel(root, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return nil
	}
	current := root
	for _, part := range strings.Split(relative, string(filepath.Separator)) {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("inventory_target_unsafe: output traverses symbolic link %s", current)
		}
	}
	return nil
}

func resolvePathAllowMissing(path string) (string, error) {
	current := filepath.Clean(path)
	suffix := []string{}
	for {
		if _, err := os.Lstat(current); err == nil {
			break
		} else if !os.IsNotExist(err) {
			return "", err
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("no existing ancestor for %s", path)
		}
		suffix = append(suffix, filepath.Base(current))
		current = parent
	}
	resolved, err := filepath.EvalSymlinks(current)
	if err != nil {
		return "", err
	}
	for index := len(suffix) - 1; index >= 0; index-- {
		resolved = filepath.Join(resolved, suffix[index])
	}
	return filepath.Clean(resolved), nil
}

func samePath(left, right string) bool {
	if filepath.Separator == '\\' {
		return strings.EqualFold(filepath.Clean(left), filepath.Clean(right))
	}
	return filepath.Clean(left) == filepath.Clean(right)
}
