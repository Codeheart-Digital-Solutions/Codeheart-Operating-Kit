package commands

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/kitfs"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/portfolio"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

var (
	portfolioIdentifierPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`)
	portfolioOwnerPattern      = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?$`)
)

func RunPortfolio(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		return writeArgError(stderr, "portfolio", fmt.Errorf("the following arguments are required: subcommand"))
	}
	switch args[0] {
	case "configure":
		return runPortfolioConfigure(args[1:], stdout, stderr)
	case "scan":
		return runPortfolioScan(args[1:], stdout, stderr)
	default:
		return writeArgError(stderr, "portfolio", fmt.Errorf("invalid subcommand %q (choose from configure, scan)", args[0]))
	}
}

func runPortfolioConfigure(args []string, stdout io.Writer, stderr io.Writer) int {
	specs := map[string]bool{"--role": true, "--member-repository-id": true, "--coordination-home-id": true, "--github-owner": true, "--local-root": true, "--dry-run": false, "--yes": false, "--json": false}
	values, multiple, bools, positionals, err := parseRepeatableValueArgs(args, specs, map[string]bool{"--github-owner": true, "--local-root": true})
	if err != nil {
		return writeArgError(stderr, "portfolio configure", err)
	}
	if bools["--dry-run"] == bools["--yes"] {
		return writeArgError(stderr, "portfolio configure", fmt.Errorf("choose exactly one of --dry-run or --yes"))
	}
	root, err := singleRoot(positionals)
	if err != nil {
		return writeArgError(stderr, "portfolio configure", err)
	}
	role := portfolio.Role(values["--role"])
	if role != portfolio.RoleMember && role != portfolio.RoleCoordinationHome {
		return writeArgError(stderr, "portfolio configure", fmt.Errorf("--role must be member or coordination-home"))
	}
	repositoryID := values["--member-repository-id"]
	homeID := values["--coordination-home-id"]
	if !portfolioIdentifierPattern.MatchString(repositoryID) {
		return writeArgError(stderr, "portfolio configure", fmt.Errorf("--member-repository-id must be a lowercase portable identifier"))
	}
	if !portfolioIdentifierPattern.MatchString(homeID) {
		return writeArgError(stderr, "portfolio configure", fmt.Errorf("--coordination-home-id must be a lowercase portable identifier"))
	}
	if role == portfolio.RoleMember && (len(multiple["--github-owner"]) > 0 || len(multiple["--local-root"]) > 0) {
		return writeArgError(stderr, "portfolio configure", fmt.Errorf("discovery sources are valid only for role coordination-home"))
	}
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return writeArgError(stderr, "portfolio configure", err)
	}
	configPath := filepath.Join(absoluteRoot, filepath.FromSlash(state.ConfigPath))
	configData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: initialized Kit config required: %v\n", err)
		return 1
	}
	config, err := state.DecodeYAMLMap(configData)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", err)
		return 1
	}
	if config["component_settings"] == nil {
		config["component_settings"] = map[string]any{}
	}
	if err := state.ValidateConfig(config); err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", err)
		return 1
	}
	existingPortfolio := state.Map(config["portfolio"])
	if conflict := portfolioIdentityConflict(existingPortfolio, role, repositoryID, homeID); conflict != "" {
		fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: portfolio_identity_conflict: %s\n", conflict)
		return 1
	}
	owners := existingGitHubOwners(existingPortfolio)
	owners = append(owners, multiple["--github-owner"]...)
	owners = uniqueSortedNonempty(owners)
	portfolioValue := map[string]any{"schema_version": 2, "role": string(role), "member_repository_id": repositoryID, "coordination_home_id": homeID}
	if role == portfolio.RoleCoordinationHome {
		sources := []any{}
		for _, owner := range owners {
			if !portfolioOwnerPattern.MatchString(owner) {
				return writeArgError(stderr, "portfolio configure", fmt.Errorf("invalid GitHub owner %q", owner))
			}
			sources = append(sources, map[string]any{"kind": "github-owner", "owner": owner})
		}
		portfolioValue["discovery"] = map[string]any{"sources": sources}
	}
	config["portfolio"] = portfolioValue
	if err := state.ValidateConfig(config); err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", err)
		return 1
	}
	desiredConfig, err := state.EncodeYAML(config)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", err)
		return 1
	}
	actions := []reconcile.Action{}
	actions, err = appendChangedFileAction(actions, absoluteRoot, state.ConfigPath, desiredConfig, "generated-surface")
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", err)
		return 1
	}
	if role == portfolio.RoleCoordinationHome {
		localRoots, localErr := desiredLocalRoots(absoluteRoot, multiple["--local-root"])
		if localErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", localErr)
			return 1
		}
		if localRoots != nil {
			actions, err = appendChangedFileAction(actions, absoluteRoot, portfolio.LocalSourcesPath, localRoots, "local-machine")
			if err != nil {
				fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", err)
				return 1
			}
		}
		readme, resourceErr := kitfs.ReadFile("components/planning-workflows/scaffolds/portfolio-README.md")
		if resourceErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", resourceErr)
			return 1
		}
		pristineOverlay, resourceErr := kitfs.ReadFile("components/planning-workflows/scaffolds/portfolio-strategic-overlay.yaml")
		if resourceErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", resourceErr)
			return 1
		}
		overlay, encodeErr := state.EncodeYAML(map[string]any{"schema_version": 1, "coordination_home_id": homeID, "families": []any{}, "themes": []any{}, "relations": []any{}, "priorities": []any{}, "analyses": []any{}})
		if encodeErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", encodeErr)
			return 1
		}
		actions, err = appendAbsentFileAction(actions, absoluteRoot, "docs/repo/portfolio/README.md", readme, "repo-owned")
		if err == nil {
			actions, err = appendAbsentOrPristineFileAction(actions, absoluteRoot, portfolio.OverlayPath, overlay, pristineOverlay, "repo-owned")
		}
		if err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", err)
			return 1
		}
	}
	observed, err := state.Inspect(absoluteRoot)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", err)
		return 1
	}
	plan, err := reconcile.BuildFilePlan("portfolio configure", absoluteRoot, string(observed.Classification), []state.Classification{observed.Classification}, actions)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", err)
		return 1
	}
	var result reconcile.Result
	if bools["--dry-run"] {
		result = reconcile.Preview(plan)
	} else {
		result, err = reconcile.Apply(plan, reconcile.ApplyOptions{Now: time.Now()})
		if err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit portfolio configure: error: %v\n", err)
			return 1
		}
	}
	if bools["--json"] {
		_ = reconcile.WriteJSON(stdout, result)
	} else {
		reconcile.WriteText(stdout, result)
	}
	if result.OK() {
		return 0
	}
	return 1
}

func runPortfolioScan(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{"--format": true, "--json": false})
	if err != nil {
		return writeArgError(stderr, "portfolio scan", err)
	}
	format := values["--format"]
	if format == "" {
		format = "text"
	}
	if bools["--json"] {
		if format == "text" && values["--format"] != "" {
			return writeArgError(stderr, "portfolio scan", fmt.Errorf("--json conflicts with --format text"))
		}
		format = "json"
	}
	if format != "text" && format != "json" {
		return writeArgError(stderr, "portfolio scan", fmt.Errorf("invalid --format %q (choose from text, json)", format))
	}
	root, err := singleRoot(positionals)
	if err != nil {
		return writeArgError(stderr, "portfolio scan", err)
	}
	config, err := portfolio.LoadConfig(root)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit portfolio scan: error: %v\n", err)
		return 1
	}
	if config.Role != portfolio.RoleCoordinationHome {
		fmt.Fprintln(stderr, "codeheart-operating-kit portfolio scan: error: portfolio_scan_requires_coordination_home")
		return 1
	}
	runner := portfolio.ExecRunner{}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	result, err := portfolio.Scan(ctx, portfolio.ScanOptions{Root: root, Config: config, Sources: portfolio.BuildSources(config, runner), Runner: runner, WriteCache: true})
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit portfolio scan: error: %v\n", err)
		return 1
	}
	if format == "json" {
		if err := writeJSON(stdout, result); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit portfolio scan: error: %v\n", err)
			return 1
		}
	} else {
		catalog := result.Catalog
		fmt.Fprintf(stdout, "Portfolio scan: complete=%t; mixed coverage complete=%t; canonical ready=%t; %d member(s), %d candidate(s), %d canonical observation(s), %d compatibility observation(s), %d stale, %d error(s); duration %dms.\n", catalog.Complete, catalog.MixedCoverageComplete, catalog.CanonicalReady, catalog.Metrics.MemberCount, catalog.Metrics.CandidateCount, catalog.Metrics.ObservationCount, catalog.Metrics.CompatibilityObservationCount, catalog.Metrics.StaleCount, len(catalog.Errors), catalog.Metrics.DurationMS)
		fmt.Fprintln(stdout, "Only pushed remote refs are included; worktree changes, local heads, and unpushed commits are not globally visible.")
		if result.CacheUpdated {
			fmt.Fprintf(stdout, "Updated last complete cache: %s\n", portfolio.CatalogPath)
		} else if result.PreviousPreserved {
			fmt.Fprintln(stdout, "Preserved the previous complete cache byte-for-byte.")
		}
		for _, scanError := range catalog.Errors {
			fmt.Fprintf(stdout, "- %s: %s\n", scanError.Code, scanError.Message)
		}
	}
	if result.Catalog.Complete {
		return 0
	}
	return 1
}

func portfolioIdentityConflict(existing map[string]any, role portfolio.Role, repositoryID, homeID string) string {
	if existing == nil {
		return ""
	}
	for key, requested := range map[string]string{"role": string(role), "member_repository_id": repositoryID, "coordination_home_id": homeID} {
		current := state.AsString(existing[key])
		if current != "" && current != requested {
			return fmt.Sprintf("existing %s %q conflicts with requested %q", key, current, requested)
		}
	}
	return ""
}

func existingGitHubOwners(existing map[string]any) []string {
	owners := []string{}
	for _, item := range state.AnySlice(state.Map(existing["discovery"])["sources"]) {
		source := state.Map(item)
		if state.AsString(source["kind"]) == "github-owner" {
			owners = append(owners, state.AsString(source["owner"]))
		}
	}
	return owners
}

func uniqueSortedNonempty(values []string) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func desiredLocalRoots(root string, supplied []string) ([]byte, error) {
	path := filepath.Join(root, filepath.FromSlash(portfolio.LocalSourcesPath))
	roots := []string{}
	if data, err := os.ReadFile(path); err == nil {
		value, err := state.DecodeAndValidateYAML(state.PortfolioSourcesSchema, data)
		if err != nil {
			return nil, err
		}
		for _, item := range state.AnySlice(value["sources"]) {
			roots = append(roots, state.AsString(state.Map(item)["root"]))
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	for _, value := range supplied {
		absolute, err := filepath.Abs(expandPath(value))
		if err != nil {
			return nil, err
		}
		roots = append(roots, absolute)
	}
	roots = uniqueSortedNonempty(roots)
	if len(roots) == 0 {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, nil
		}
	}
	sources := []any{}
	for _, value := range roots {
		sources = append(sources, map[string]any{"kind": "local-git-root", "root": value})
	}
	return state.EncodeYAML(map[string]any{"schema_version": 1, "sources": sources})
}

func appendChangedFileAction(actions []reconcile.Action, root, target string, desired []byte, owner string) ([]reconcile.Action, error) {
	path := filepath.Join(root, filepath.FromSlash(target))
	existing, err := os.ReadFile(path)
	if err == nil {
		if bytes.Equal(existing, desired) {
			return actions, nil
		}
		digest := sha256.Sum256(existing)
		return append(actions, reconcile.Action{Kind: "replace", Target: target, Owner: owner, Content: desired, Mode: 0o644, ExpectedSHA256: hex.EncodeToString(digest[:])}), nil
	}
	if os.IsNotExist(err) {
		return append(actions, reconcile.Action{Kind: "create", Target: target, Owner: owner, Content: desired, Mode: 0o644}), nil
	}
	return actions, fmt.Errorf("read portfolio target %s: %w", target, err)
}

func appendAbsentFileAction(actions []reconcile.Action, root, target string, desired []byte, owner string) ([]reconcile.Action, error) {
	if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(target))); err == nil {
		return actions, nil
	} else if !os.IsNotExist(err) {
		return actions, fmt.Errorf("inspect portfolio target %s: %w", target, err)
	}
	return append(actions, reconcile.Action{Kind: "create", Target: target, Owner: owner, Content: desired, Mode: 0o644}), nil
}

func appendAbsentOrPristineFileAction(actions []reconcile.Action, root, target string, desired, pristine []byte, owner string) ([]reconcile.Action, error) {
	repositoryRoot, err := os.OpenRoot(root)
	if err != nil {
		return nil, fmt.Errorf("open portfolio repository root: %w", err)
	}
	defer repositoryRoot.Close()
	relative := filepath.FromSlash(target)
	info, err := repositoryRoot.Lstat(relative)
	if os.IsNotExist(err) {
		return append(actions, reconcile.Action{Kind: "create", Target: target, Owner: owner, Content: desired, Mode: 0o644}), nil
	}
	if err != nil {
		return nil, fmt.Errorf("inspect portfolio target %s: %w", target, err)
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return actions, nil
	}
	file, err := repositoryRoot.Open(relative)
	if err != nil {
		return nil, fmt.Errorf("open portfolio target %s: %w", target, err)
	}
	existing, readErr := io.ReadAll(file)
	opened, statErr := file.Stat()
	closeErr := file.Close()
	current, currentErr := repositoryRoot.Lstat(relative)
	if readErr != nil || statErr != nil || currentErr != nil || closeErr != nil || !os.SameFile(info, opened) || !os.SameFile(info, current) {
		return nil, fmt.Errorf("portfolio target %s changed while reading pristine scaffold", target)
	}
	if bytes.Equal(existing, desired) || !bytes.Equal(existing, pristine) {
		return actions, nil
	}
	digest := sha256.Sum256(existing)
	return append(actions, reconcile.Action{Kind: "replace", Target: target, Owner: owner, Content: desired, Mode: 0o644, ExpectedSHA256: hex.EncodeToString(digest[:])}), nil
}
