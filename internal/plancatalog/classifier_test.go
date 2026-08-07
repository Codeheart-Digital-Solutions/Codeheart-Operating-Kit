package plancatalog

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestMultiRootRepositoryFixtureCoversDiscoveryV2AcceptanceMatrix(t *testing.T) {
	root := copyPlanFixture(t, "multi-root-repository")
	runGitTest(t, root, "init")
	runGitTest(t, root, "checkout", "-b", "main")
	runGitTest(t, root, "config", "user.email", "test@example.invalid")
	runGitTest(t, root, "config", "user.name", "Plan Catalog Test")
	writeClassifierFile(t, root, ".gitignore", "ignored/\n")
	runGitTest(t, root, "add", ".")
	// Some developer machines globally ignore uppercase Markdown extensions;
	// force this intentionally tracked portability fixture into the index.
	runGitTest(t, root, "add", "-f", "docs/case/case_discovery_doc.MD")
	runGitTest(t, root, "commit", "-m", "multi-root acceptance fixture")
	if err := os.MkdirAll(filepath.Join(root, "embedded", ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	settings, configProblems := LoadRepositorySettings(root)
	if HasErrors(configProblems) || settings.Mode != ModeCanonical || settings.RepositoryID != "multi-root" || settings.DiscoveryVersion != DiscoveryV2 || !reflect.DeepEqual(settings.ExcludedRoots, []string{"external/"}) {
		t.Fatalf("fixture settings=%#v problems=%#v", settings, configProblems)
	}
	first, err := ClassifyLocalIndex(root, settings, settings.Mode, settings.RepositoryID)
	if err != nil {
		t.Fatal(err)
	}
	second, err := ClassifyLocalIndex(root, settings, settings.Mode, settings.RepositoryID)
	if err != nil {
		t.Fatal(err)
	}
	firstJSON, _ := json.Marshal(first)
	secondJSON, _ := json.Marshal(second)
	if string(firstJSON) != string(secondJSON) || first.CandidateSetDigest != second.CandidateSetDigest {
		t.Fatal("persistent multi-root fixture classification was not deterministic")
	}

	byPath := map[string]Candidate{}
	for _, candidate := range first.Candidates {
		byPath[candidate.Path] = candidate
	}
	for _, expected := range []string{
		"docs/root/root_discovery_doc.md",
		"docs/business/procurement/procurement_implementation_doc.md",
		"products/payments/packages/settlement/docs/initiatives/settlement_discovery_doc.md",
		"areas/docs/repeated/docs/rollout/rollout_implementation_doc.md",
		"products/payments/docs/plans/README.md",
	} {
		if byPath[expected].Ownership != OwnershipOwned {
			t.Fatalf("expected owned multi-root candidate %s: %#v", expected, byPath[expected])
		}
	}
	if _, exists := byPath["docs/router/README.md"]; exists {
		t.Fatal("ordinary router README became a family candidate")
	}
	if byPath["docs/metadata-only/misnamed.md"].Signal != SignalMetadata || byPath["docs/filename-only/legacy_discovery_doc.md"].Signal != SignalFilename {
		t.Fatalf("candidate signal matrix missing: %#v", byPath)
	}
	if byPath["external/docs/copied/excluded_discovery_doc.md"].Ownership != OwnershipExcluded || byPath["vendor/docs/copied/ambiguous_discovery_doc.md"].Ownership != OwnershipProspectiveBlocked || byPath["notes/misplaced.md"].Ownership != OwnershipHardUnowned || byPath["embedded/docs/nested/nested_discovery_doc.md"].Ownership != OwnershipHardUnowned {
		t.Fatalf("ownership matrix missing: %#v", byPath)
	}
	for _, code := range []string{
		"metadata_missing", "canonical_filename_missing", "metadata_yaml_invalid",
		"record_kind_path_mismatch", "duplicate_plan_id", "plan_candidate_excluded",
		"plan_documentation_root_ambiguous", "plan_metadata_misplaced", "plan_path_case_mismatch",
		"plan_source_unsafe",
	} {
		if !problemCodeExists(first.Discovery.Problems, code) {
			t.Fatalf("acceptance fixture missing problem %s: %#v", code, first.Discovery.Problems)
		}
	}
	if countRecordPath(first.Discovery.Records, "areas/docs/repeated/docs/rollout/rollout_implementation_doc.md") != 1 {
		t.Fatal("repeated docs segments produced duplicate records")
	}

	previewPath := "products/payments/docs/authoring/preview_discovery_doc.md"
	ignoredPath := "ignored/docs/ignored_discovery_doc.md"
	writeClassifierFile(t, root, previewPath, canonicalClassifierDocument("multi-root.discovery.preview", KindDiscovery, "Preview"))
	writeClassifierFile(t, root, ignoredPath, canonicalClassifierDocument("multi-root.discovery.ignored", KindDiscovery, "Ignored"))
	preview, err := ClassifyUntrackedPreview(root, settings, settings.Mode, settings.RepositoryID)
	if err != nil {
		t.Fatal(err)
	}
	if len(preview.Candidates) != 1 || preview.Candidates[0].Path != previewPath || !preview.Candidates[0].Provenance.Source.Preview {
		t.Fatalf("untracked preview matrix=%#v", preview.Candidates)
	}
}

func TestDiscoveryV2ClassifiesRepositoryWideTrackedDocsDeterministically(t *testing.T) {
	root := t.TempDir()
	runGitTest(t, root, "init")
	runGitTest(t, root, "checkout", "-b", "main")
	runGitTest(t, root, "config", "user.email", "test@example.invalid")
	runGitTest(t, root, "config", "user.name", "Plan Catalog Test")
	files := map[string]string{
		"docs/root/root_discovery_doc.md":                            canonicalClassifierDocument("example.discovery.root", KindDiscovery, "Root Discovery"),
		"products/a/packages/b/docs/deep/deep_implementation_doc.md": canonicalClassifierDocument("example.implementation.deep", KindImplementation, "Deep Implementation"),
		"areas/docs/repeated/docs/again/again_discovery_doc.md":      canonicalClassifierDocument("example.discovery.again", KindDiscovery, "Repeated Docs Discovery"),
		"docs/domain/strange.md":                                     canonicalClassifierDocument("example.discovery.metadata-only", KindDiscovery, "Metadata Only"),
		"docs/domain/README.md":                                      canonicalClassifierDocument("example.family.domain", KindFamily, "Domain Family"),
		"docs/domain-wrong-id/README.md":                             canonicalClassifierDocument("example.discovery.not-a-family", KindFamily, "Wrong Identity Kind Family"),
		"docs/domain/malformed.md":                                   strings.Replace(canonicalClassifierDocument("example.discovery.malformed", KindDiscovery, "Malformed"), "  schema_version: 1", "  schema_version: [", 1),
		"docs/domain-bad/README.md":                                  strings.Replace(canonicalClassifierDocument("example.family.malformed", KindFamily, "Malformed Family"), "  schema_version: 1", "  schema_version: [", 1),
		"docs/case/README.MD":                                        canonicalClassifierDocument("example.family.case-extension", KindFamily, "Case Extension"),
		"docs/case/readme.md":                                        canonicalClassifierDocument("example.family.case-name", KindFamily, "Case Name"),
		"docs/router/README.md":                                      "# Ordinary router\n\nNo plan metadata.\n",
		"docs/examples/fenced.md":                                    "```md\n" + MetadataBeginMarker + "\n```\n",
		"docs/examples/quoted.md":                                    "> " + MetadataBeginMarker + "\n\n    " + MetadataBeginMarker + "\n",
		"docs/legacy/legacy_discovery_doc.md":                        "Last updated: 2026-08-06T00:00:00Z (UTC)\nCreated: 2026-08-06\nStatus: draft\n\n# Legacy\n",
		"vendor/docs/copied/copied_discovery_doc.md":                 canonicalClassifierDocument("example.discovery.copied", KindDiscovery, "Copied"),
		"third_party/docs/excluded/excluded_discovery_doc.md":        canonicalClassifierDocument("example.discovery.excluded", KindDiscovery, "Excluded"),
		"notes/misplaced.md":                                         canonicalClassifierDocument("example.discovery.misplaced", KindDiscovery, "Misplaced"),
		"docs/empty/_discovery_doc.md":                               "# Empty reserved basename is not a signal\n",
	}
	for name, content := range files {
		writeClassifierFile(t, root, name, content)
	}
	runGitTest(t, root, "add", ".")
	runGitTest(t, root, "commit", "-m", "multi-root fixture")

	settings := RepositorySettings{Mode: ModeCanonical, RepositoryID: "example", DiscoveryVersion: DiscoveryV2, ExcludedRoots: []string{"third_party/"}}
	first, err := ClassifyLocalIndex(root, settings, ModeCanonical, "example")
	if err != nil {
		t.Fatal(err)
	}
	second, err := ClassifyLocalIndex(root, settings, ModeCanonical, "example")
	if err != nil {
		t.Fatal(err)
	}
	firstJSON, _ := json.Marshal(first)
	secondJSON, _ := json.Marshal(second)
	if string(firstJSON) != string(secondJSON) || first.CandidateSetDigest != second.CandidateSetDigest {
		t.Fatal("classifier output was not byte-deterministic")
	}

	owned := []string{}
	byPath := map[string]Candidate{}
	for _, candidate := range first.Candidates {
		byPath[candidate.Path] = candidate
		if candidate.Ownership == OwnershipOwned {
			owned = append(owned, candidate.Path)
		}
	}
	sort.Strings(owned)
	for _, expected := range []string{
		"docs/root/root_discovery_doc.md",
		"products/a/packages/b/docs/deep/deep_implementation_doc.md",
		"areas/docs/repeated/docs/again/again_discovery_doc.md",
		"docs/domain/strange.md",
		"docs/domain/README.md",
		"docs/legacy/legacy_discovery_doc.md",
	} {
		if !stringSliceContains(owned, expected) {
			t.Fatalf("owned candidate %s missing from %#v", expected, owned)
		}
	}
	for _, absent := range []string{"docs/router/README.md", "docs/examples/fenced.md", "docs/examples/quoted.md", "docs/empty/_discovery_doc.md"} {
		if _, exists := byPath[absent]; exists {
			t.Fatalf("non-signal path became a candidate: %s", absent)
		}
	}
	if byPath["docs/domain/strange.md"].Signal != SignalMetadata || byPath["docs/domain/README.md"].ExpectedKind != KindFamily || !byPath["docs/domain/README.md"].FamilyQualified {
		t.Fatalf("metadata/family signals=%#v %#v", byPath["docs/domain/strange.md"], byPath["docs/domain/README.md"])
	}
	if byPath["docs/domain-wrong-id/README.md"].FamilyQualified || byPath["docs/domain-wrong-id/README.md"].ExpectedKind != "" {
		t.Fatalf("README with a non-family identity gained family authority: %#v", byPath["docs/domain-wrong-id/README.md"])
	}
	for _, malformed := range []string{"docs/domain/malformed.md", "docs/domain-bad/README.md"} {
		if byPath[malformed].FamilyQualified || byPath[malformed].ExpectedKind != "" {
			t.Fatalf("malformed metadata gained filename/family authority: %#v", byPath[malformed])
		}
	}
	if byPath["vendor/docs/copied/copied_discovery_doc.md"].Ownership != OwnershipProspectiveBlocked || byPath["third_party/docs/excluded/excluded_discovery_doc.md"].Ownership != OwnershipExcluded || byPath["notes/misplaced.md"].Ownership != OwnershipHardUnowned {
		t.Fatalf("ownership classification incorrect: vendor=%#v excluded=%#v misplaced=%#v", byPath["vendor/docs/copied/copied_discovery_doc.md"], byPath["third_party/docs/excluded/excluded_discovery_doc.md"], byPath["notes/misplaced.md"])
	}
	for _, code := range []string{"canonical_filename_missing", "metadata_missing", "metadata_yaml_invalid", "plan_documentation_root_ambiguous", "plan_candidate_excluded", "plan_metadata_misplaced", "plan_path_case_mismatch"} {
		if !problemCodeExists(first.Discovery.Problems, code) {
			t.Fatalf("problem %s missing from %#v", code, first.Discovery.Problems)
		}
	}
	if countRecordPath(first.Discovery.Records, "areas/docs/repeated/docs/again/again_discovery_doc.md") != 1 {
		t.Fatal("repeated docs segments caused duplicate discovery")
	}
}

func TestFormalPathKindV2RequiresAuthoritativeIndexCandidate(t *testing.T) {
	root := t.TempDir()
	runGitTest(t, root, "init")
	tracked := "docs/tracked/tracked_discovery_doc.md"
	untracked := "docs/untracked/README.md"
	writeClassifierFile(t, root, tracked, canonicalClassifierDocument("example.discovery.tracked", KindDiscovery, "Tracked"))
	writeClassifierFile(t, root, untracked, canonicalClassifierDocument("example.family.untracked", KindFamily, "Untracked"))
	runGitTest(t, root, "add", tracked)
	settings := RepositorySettings{Mode: ModeCanonical, RepositoryID: "example", DiscoveryVersion: DiscoveryV2, ExcludedRoots: []string{}}
	if kind, ok := FormalPathKindWithSettings(root, tracked, settings); !ok || kind != KindDiscovery {
		t.Fatalf("tracked formal kind=%q ok=%t", kind, ok)
	}
	if kind, ok := FormalPathKindWithSettings(root, untracked, settings); ok || kind != "" {
		t.Fatalf("untracked family gained authority: kind=%q ok=%t", kind, ok)
	}
}

func TestSharedClassifierHardBoundariesNeverReadUnsafeContent(t *testing.T) {
	reads := 0
	reader := func(blob GitBlob) ([]byte, error) {
		reads++
		return nil, nil
	}
	blobs := []GitBlob{
		{Path: "docs/link/link_discovery_doc.md", Mode: GitModeSymlink, ObjectID: strings.Repeat("a", 40)},
		{Path: "docs/submodule/submodule_discovery_doc.md", Mode: GitModeGitlink, ObjectID: strings.Repeat("b", 40)},
		{Path: ".codeheart/kit/docs/managed_discovery_doc.md", Mode: GitModeRegular, ObjectID: strings.Repeat("c", 40)},
	}
	policy := DiscoveryPolicy{Version: DiscoveryV2, ExcludedRoots: []string{}, AmbiguitySegments: append([]string{}, conventionalAmbiguitySegments...), AuthoritativeUniverse: "git-regular-blobs"}
	classified, err := ClassifyGitBlobs(blobs, reader, policy, ModeCanonical, "example", nil)
	if err != nil {
		t.Fatal(err)
	}
	if reads != 0 || len(classified.Discovery.Records) != 0 {
		t.Fatalf("unsafe blobs were read or cataloged: reads=%d records=%#v", reads, classified.Discovery.Records)
	}
	if len(classified.Candidates) != 3 {
		t.Fatalf("hard-unowned filename signals disappeared: %#v", classified.Candidates)
	}
	for _, candidate := range classified.Candidates {
		if candidate.Ownership != OwnershipHardUnowned {
			t.Fatalf("unsafe candidate=%#v", candidate)
		}
	}
}

func TestLocalIndexKeepsFilenameCandidateWhenWorktreeFileIsMissing(t *testing.T) {
	root := t.TempDir()
	runGitTest(t, root, "init")
	name := "docs/missing/missing_discovery_doc.md"
	writeClassifierFile(t, root, name, canonicalClassifierDocument("example.discovery.missing", KindDiscovery, "Missing"))
	runGitTest(t, root, "add", ".")
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(name))); err != nil {
		t.Fatal(err)
	}
	settings := RepositorySettings{Mode: ModeCanonical, RepositoryID: "example", DiscoveryVersion: DiscoveryV2, ExcludedRoots: []string{}}
	classified, err := ClassifyLocalIndex(root, settings, ModeCanonical, "example")
	if err != nil {
		t.Fatal(err)
	}
	if len(classified.Candidates) != 1 || classified.Candidates[0].Signal != SignalBoth || classified.Complete || !problemCodeExists(classified.Discovery.Problems, "plan_source_unsafe") {
		t.Fatalf("missing indexed candidate=%#v", classified)
	}
}

func TestLocalIndexKeepsMetadataCandidateFromIndexWhenWorktreeFileIsMissing(t *testing.T) {
	root := t.TempDir()
	runGitTest(t, root, "init")
	name := "docs/missing/unconventional.md"
	writeClassifierFile(t, root, name, canonicalClassifierDocument("example.discovery.missing-metadata", KindDiscovery, "Missing Metadata Candidate"))
	runGitTest(t, root, "add", ".")
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(name))); err != nil {
		t.Fatal(err)
	}
	settings := RepositorySettings{Mode: ModeCanonical, RepositoryID: "example", DiscoveryVersion: DiscoveryV2, ExcludedRoots: []string{}}
	classified, err := ClassifyLocalIndex(root, settings, ModeCanonical, "example")
	if err != nil {
		t.Fatal(err)
	}
	if len(classified.Candidates) != 1 || classified.Candidates[0].Signal != SignalMetadata || classified.Complete || len(classified.Discovery.Records) != 0 || !problemCodeExists(classified.Discovery.Problems, "plan_source_unsafe") {
		t.Fatalf("missing metadata candidate=%#v", classified)
	}
}

func TestLocalIndexKeepsMetadataCandidateFromIndexWhenWorktreeFileIsSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("ordinary Windows test users may not have symlink privilege")
	}
	root := t.TempDir()
	runGitTest(t, root, "init")
	name := "docs/replaced/unconventional.md"
	writeClassifierFile(t, root, name, canonicalClassifierDocument("example.discovery.replaced-metadata", KindDiscovery, "Replaced Metadata Candidate"))
	runGitTest(t, root, "add", ".")
	external := filepath.Join(t.TempDir(), "external.md")
	if err := os.WriteFile(external, []byte("outside content must not be read\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(name))); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, filepath.Join(root, filepath.FromSlash(name))); err != nil {
		t.Fatal(err)
	}
	settings := RepositorySettings{Mode: ModeCanonical, RepositoryID: "example", DiscoveryVersion: DiscoveryV2, ExcludedRoots: []string{}}
	classified, err := ClassifyLocalIndex(root, settings, ModeCanonical, "example")
	if err != nil {
		t.Fatal(err)
	}
	if len(classified.Candidates) != 1 || classified.Candidates[0].Signal != SignalMetadata || classified.Candidates[0].Ownership != OwnershipHardUnowned || classified.Complete || len(classified.Discovery.Records) != 0 || !problemCodeExists(classified.Discovery.Problems, "plan_source_unsafe") {
		t.Fatalf("symlink-replaced metadata candidate=%#v", classified)
	}
}

func TestUnsafeOrdinaryMarkdownDoesNotBecomeCatalogProblem(t *testing.T) {
	root := t.TempDir()
	runGitTest(t, root, "init")
	missing := "docs/ordinary.md"
	nested := "embedded/docs/ordinary.md"
	writeClassifierFile(t, root, missing, "# Ordinary documentation\n")
	writeClassifierFile(t, root, nested, "# Embedded ordinary documentation\n")
	runGitTest(t, root, "add", ".")
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(missing))); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "embedded", ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	settings := RepositorySettings{Mode: ModeCanonical, RepositoryID: "example", DiscoveryVersion: DiscoveryV2, ExcludedRoots: []string{}}
	classified, err := ClassifyLocalIndex(root, settings, ModeCanonical, "example")
	if err != nil {
		t.Fatal(err)
	}
	if !classified.Complete || len(classified.Candidates) != 0 || len(classified.Discovery.Problems) != 0 {
		t.Fatalf("unrelated unsafe Markdown polluted catalog evidence: %#v", classified)
	}
}

func TestUnsafeFamilyIndexEvidenceCannotSatisfySafeChildReference(t *testing.T) {
	root := t.TempDir()
	runGitTest(t, root, "init")
	familyPath := "docs/family/README.md"
	childPath := "docs/family/child_discovery_doc.md"
	writeClassifierFile(t, root, familyPath, canonicalClassifierDocument("example.family.group", KindFamily, "Family"))
	child := strings.Replace(canonicalClassifierDocument("example.discovery.child", KindDiscovery, "Child"), "  purpose: Exercise repository-wide discovery.\n", "  purpose: Exercise repository-wide discovery.\n  family: example.family.group\n", 1)
	writeClassifierFile(t, root, childPath, child)
	runGitTest(t, root, "add", ".")
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(familyPath))); err != nil {
		t.Fatal(err)
	}
	settings := RepositorySettings{Mode: ModeCanonical, RepositoryID: "example", DiscoveryVersion: DiscoveryV2, ExcludedRoots: []string{}}
	classified, err := ClassifyLocalIndex(root, settings, ModeCanonical, "example")
	if err != nil {
		t.Fatal(err)
	}
	if len(classified.Discovery.Records) != 1 || classified.Discovery.Records[0].Path != childPath || !problemCodeExists(classified.Discovery.Problems, "family_record_missing") {
		t.Fatalf("unsafe family index evidence gained authority: %#v", classified)
	}
}

func TestIndexFallbackIgnoresReplacementObjects(t *testing.T) {
	root := t.TempDir()
	runGitTest(t, root, "init")
	name := "docs/replaced/unconventional.md"
	writeClassifierFile(t, root, name, canonicalClassifierDocument("example.discovery.original-index", KindDiscovery, "Original Index Blob"))
	runGitTest(t, root, "add", ".")
	originalOID := runGitTest(t, root, "rev-parse", ":"+name)
	replacementPath := filepath.Join(root, "replacement-content.md")
	if err := os.WriteFile(replacementPath, []byte("# Replacement without plan metadata\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	replacementOID := runGitTest(t, root, "hash-object", "-w", "--", replacementPath)
	runGitTest(t, root, "replace", originalOID, replacementOID)
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(name))); err != nil {
		t.Fatal(err)
	}
	settings := RepositorySettings{Mode: ModeCanonical, RepositoryID: "example", DiscoveryVersion: DiscoveryV2, ExcludedRoots: []string{}}
	classified, err := ClassifyLocalIndex(root, settings, ModeCanonical, "example")
	if err != nil {
		t.Fatal(err)
	}
	if len(classified.Candidates) != 1 || classified.Candidates[0].Signal != SignalMetadata || classified.Candidates[0].Provenance.Source.ObjectID != originalOID || len(classified.Discovery.Records) != 0 || !problemCodeExists(classified.Discovery.Problems, "plan_source_unsafe") {
		t.Fatalf("replacement ref changed exact index evidence: %#v", classified)
	}
}

func TestDuplicateIDIncludesCandidateWithInvalidHeader(t *testing.T) {
	policy := DiscoveryPolicy{Version: DiscoveryV2, ExcludedRoots: []string{}, AmbiguitySegments: append([]string{}, conventionalAmbiguitySegments...), AuthoritativeUniverse: "git-regular-blobs"}
	valid := canonicalClassifierDocument("example.discovery.duplicate", KindDiscovery, "Valid")
	invalidHeader := strings.Replace(canonicalClassifierDocument("example.discovery.duplicate", KindDiscovery, "Invalid Header"), "Status: draft", "Status: unknown", 1)
	blobs := []GitBlob{
		{Path: "docs/a/a_discovery_doc.md", Mode: GitModeRegular, ObjectID: strings.Repeat("a", 40)},
		{Path: "docs/b/b_discovery_doc.md", Mode: GitModeRegular, ObjectID: strings.Repeat("b", 40)},
	}
	classified, err := ClassifyGitBlobs(blobs, func(blob GitBlob) ([]byte, error) {
		if blob.Path == blobs[0].Path {
			return []byte(valid), nil
		}
		return []byte(invalidHeader), nil
	}, policy, ModeCanonical, "example", nil)
	if err != nil || !problemCodeExists(classified.Discovery.Problems, "header_invalid") || !problemCodeExists(classified.Discovery.Problems, "duplicate_plan_id") {
		t.Fatalf("partial identity evidence was lost: classification=%#v err=%v", classified, err)
	}
}

func TestHeaderInvalidFamilyREADMEKeepsIdentityWithoutFamilyAuthority(t *testing.T) {
	policy := DiscoveryPolicy{Version: DiscoveryV2, ExcludedRoots: []string{}, AmbiguitySegments: append([]string{}, conventionalAmbiguitySegments...), AuthoritativeUniverse: "git-regular-blobs"}
	data := strings.Replace(canonicalClassifierDocument("example.family.invalid-header", KindFamily, "Invalid Family Header"), "Status: draft", "Status: unknown", 1)
	blob := GitBlob{Path: "docs/family/README.md", Mode: GitModeRegular, ObjectID: strings.Repeat("a", 40)}
	classified, err := ClassifyGitBlobs([]GitBlob{blob}, func(GitBlob) ([]byte, error) { return []byte(data), nil }, policy, ModeCanonical, "example", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(classified.Candidates) != 1 || classified.Candidates[0].FamilyQualified || classified.Candidates[0].ExpectedKind != "" || len(classified.Discovery.Records) != 1 || classified.Discovery.Records[0].FamilyQualified || !problemCodeExists(classified.Discovery.Problems, "header_invalid") || !problemCodeExists(classified.Discovery.Problems, "family_placement_invalid") {
		t.Fatalf("invalid family structure gained authority: %#v", classified)
	}
}

func TestClassifierPortableCollisionsAndLocalNestedRepositoryBoundary(t *testing.T) {
	policy := DiscoveryPolicy{Version: DiscoveryV2, ExcludedRoots: []string{}, AmbiguitySegments: append([]string{}, conventionalAmbiguitySegments...), AuthoritativeUniverse: "git-regular-blobs"}
	doc := []byte(canonicalClassifierDocument("example.discovery.case", KindDiscovery, "Case"))
	blobs := []GitBlob{
		{Path: "docs/Case/plan_discovery_doc.md", Mode: GitModeRegular, ObjectID: strings.Repeat("a", 40)},
		{Path: "docs/case/plan_discovery_doc.md", Mode: GitModeRegular, ObjectID: strings.Repeat("b", 40)},
	}
	classified, err := ClassifyGitBlobs(blobs, func(GitBlob) ([]byte, error) { return doc, nil }, policy, ModeCanonical, "example", nil)
	if err != nil || !problemCodeExists(classified.Discovery.Problems, "plan_path_portable_collision") {
		t.Fatalf("portable collision result=%#v err=%v", classified, err)
	}

	root := t.TempDir()
	runGitTest(t, root, "init")
	writeClassifierFile(t, root, "embedded/docs/plan_discovery_doc.md", string(doc))
	runGitTest(t, root, "add", ".")
	if err := os.MkdirAll(filepath.Join(root, "embedded", ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	settings := RepositorySettings{Mode: ModeCanonical, RepositoryID: "example", DiscoveryVersion: DiscoveryV2, ExcludedRoots: []string{}}
	local, err := ClassifyLocalIndex(root, settings, ModeCanonical, "example")
	if err != nil {
		t.Fatal(err)
	}
	if len(local.Candidates) != 1 || local.Candidates[0].Ownership != OwnershipHardUnowned || local.Candidates[0].Provenance.Boundary != "nested-repository" {
		t.Fatalf("nested boundary=%#v", local.Candidates)
	}
}

func canonicalClassifierDocument(id string, kind Kind, title string) string {
	return "Last updated: 2026-08-06T00:00:00Z (UTC)\nCreated: 2026-08-06\nStatus: draft\n\n# " + title + "\n\n" + MetadataBeginMarker + "\n```yaml\nplan:\n  schema_version: 1\n  id: " + id + "\n  kind: " + string(kind) + "\n  purpose: Exercise repository-wide discovery.\n  first_cataloged: 2026-08-06T00:00:00Z\n  catalog_metadata_updated: 2026-08-06T00:00:00Z\n```\n" + MetadataEndMarker + "\n"
}

func writeClassifierFile(t *testing.T, root, name, content string) {
	t.Helper()
	target := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func problemCodeExists(problems []Problem, code string) bool {
	for _, problem := range problems {
		if problem.Code == code {
			return true
		}
	}
	return false
}

func countRecordPath(records []Record, expected string) int {
	count := 0
	for _, record := range records {
		if record.Path == expected {
			count++
		}
	}
	return count
}

func TestClassifierCandidateOrderIsIndependentOfInputOrder(t *testing.T) {
	policy := DiscoveryPolicy{Version: DiscoveryV2, ExcludedRoots: []string{}, AmbiguitySegments: append([]string{}, conventionalAmbiguitySegments...), AuthoritativeUniverse: "git-regular-blobs"}
	doc := []byte(canonicalClassifierDocument("example.discovery.order", KindDiscovery, "Order"))
	left := []GitBlob{{Path: "z/docs/z_discovery_doc.md", Mode: GitModeRegular, ObjectID: strings.Repeat("b", 40)}, {Path: "a/docs/a_discovery_doc.md", Mode: GitModeRegular, ObjectID: strings.Repeat("a", 40)}}
	right := append([]GitBlob{}, left...)
	sort.Slice(right, func(i, j int) bool { return right[i].Path > right[j].Path })
	reader := func(GitBlob) ([]byte, error) { return doc, nil }
	first, _ := ClassifyGitBlobs(left, reader, policy, ModeLegacy, "example", nil)
	second, _ := ClassifyGitBlobs(right, reader, policy, ModeLegacy, "example", nil)
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("input order changed classification\nfirst=%#v\nsecond=%#v", first, second)
	}
}

func TestClassifierRejectsInvalidUTF8BeforeReadingOrHashing(t *testing.T) {
	invalidPath := "docs/" + string([]byte{0xff}) + "/plan_discovery_doc.md"
	reads := 0
	policy := DiscoveryPolicy{Version: DiscoveryV2, ExcludedRoots: []string{}, AmbiguitySegments: append([]string{}, conventionalAmbiguitySegments...), AuthoritativeUniverse: "git-regular-blobs"}
	classified, err := ClassifyGitBlobs([]GitBlob{{Path: invalidPath, Mode: GitModeRegular, ObjectID: strings.Repeat("a", 40)}}, func(GitBlob) ([]byte, error) {
		reads++
		return nil, nil
	}, policy, ModeCanonical, "example", nil)
	if err != nil {
		t.Fatal(err)
	}
	if reads != 0 || len(classified.Candidates) != 0 || classified.Complete || !problemCodeExists(classified.Discovery.Problems, "plan_path_invalid") {
		t.Fatalf("invalid UTF-8 path reached classification: reads=%d result=%#v", reads, classified)
	}
}

func TestV1FilenameAndSourceSizeCompatibilityRemainUnchanged(t *testing.T) {
	root := t.TempDir()
	name := "docs/repo/plans/legacy/_discovery_doc.md"
	content := strings.Repeat("x", MaxPlanSourceBytes+1)
	writeClassifierFile(t, root, name, content)
	candidates, err := Enumerate(root)
	if err != nil || len(candidates) != 1 || candidates[0].ExpectedKind != KindDiscovery {
		t.Fatalf("v1 filename compatibility changed: candidates=%#v err=%v", candidates, err)
	}
	data, err := readRegularSource(root, name)
	if err != nil || len(data) != MaxPlanSourceBytes+1 {
		t.Fatalf("v1 source-size compatibility changed: bytes=%d err=%v", len(data), err)
	}
	if _, err := readBoundedRegularSource(root, name); ErrorCode(err) != "source_too_large" {
		t.Fatalf("v2 bounded read error=%v", err)
	}
}

func TestProspectiveV2SnapshotAndUntrackedPreviewAreReadOnlyAndSeparated(t *testing.T) {
	root := migrationRepository(t)
	nested := "products/widget/docs/discovery/unconventional.md"
	writeClassifierFile(t, root, nested, canonicalClassifierDocument("example.discovery.nested", KindDiscovery, "Nested"))
	writeClassifierFile(t, root, ".gitignore", "ignored/\n")
	runGitTest(t, root, "add", nested, ".gitignore")
	runGitTest(t, root, "commit", "-m", "add nested plan and ignore rule")

	active, err := LoadRepositorySnapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range active.Records {
		if record.Path == nested {
			t.Fatal("active discovery v1 unexpectedly expanded before activation")
		}
	}
	options := SnapshotOptions{TargetDiscoveryVersion: DiscoveryV2}
	prospective, err := LoadRepositorySnapshotWithOptions(root, options)
	if err != nil {
		t.Fatal(err)
	}
	if prospective.Settings.DiscoveryVersion != DiscoveryV2 || prospective.ConfiguredSettings.DiscoveryVersion != DiscoveryV1 || !prospective.Targeted || prospective.CandidateSetDigest == "" || !recordPathExists(prospective.Records, nested) {
		t.Fatalf("prospective snapshot=%#v", prospective)
	}

	previewPath := "products/widget/docs/authoring/new_discovery_doc.md"
	ignoredPath := "ignored/docs/ignored_discovery_doc.md"
	writeClassifierFile(t, root, previewPath, canonicalClassifierDocument("example.discovery.preview", KindDiscovery, "Preview"))
	writeClassifierFile(t, root, ignoredPath, canonicalClassifierDocument("example.discovery.ignored", KindDiscovery, "Ignored"))
	statusBefore := runGitTest(t, root, "status", "--porcelain=v1", "--untracked-files=all")
	withPreview, err := LoadRepositorySnapshotWithOptions(root, SnapshotOptions{TargetDiscoveryVersion: DiscoveryV2, IncludeUntracked: true})
	if err != nil {
		t.Fatal(err)
	}
	statusAfter := runGitTest(t, root, "status", "--porcelain=v1", "--untracked-files=all")
	if statusBefore != statusAfter {
		t.Fatalf("prospective read changed repository state:\nbefore=%s\nafter=%s", statusBefore, statusAfter)
	}
	if prospective.CandidateSetDigest != withPreview.CandidateSetDigest || len(prospective.Records) != len(withPreview.Records) || len(withPreview.PreviewCandidates) != 1 || withPreview.PreviewCandidates[0].Path != previewPath || !withPreview.PreviewCandidates[0].Provenance.Source.Preview || !recordPathExists(withPreview.PreviewRecords, previewPath) {
		t.Fatalf("preview contaminated authority or omitted evidence: before=%#v after=%#v", prospective, withPreview)
	}
	for _, candidate := range withPreview.PreviewCandidates {
		if candidate.Path == ignoredPath {
			t.Fatal("ignored untracked plan entered authoring preview")
		}
	}
	migration, err := BuildMigrationPlan(root, migrationLedger(t, root))
	if err != nil || len(migration.FilePlan.Actions) != 2 {
		t.Fatalf("preview or prospective candidates entered active migration: plan=%#v err=%v", migration, err)
	}
}

func TestV2ActiveBranchTouchesClassifyNestedMetadataPaths(t *testing.T) {
	root := migrationRepository(t)
	path := "products/widget/docs/discovery/unconventional.md"
	runGitTest(t, root, "checkout", "-b", "feature/nested-plan")
	writeClassifierFile(t, root, path, canonicalClassifierDocument("example.discovery.branch-nested", KindDiscovery, "Branch Nested"))
	runGitTest(t, root, "add", path)
	runGitTest(t, root, "commit", "-m", "add nested metadata plan")
	runGitTest(t, root, "checkout", "main")
	settings, problems := LoadRepositorySettings(root)
	if HasErrors(problems) {
		t.Fatalf("settings problems=%#v", problems)
	}
	settings.DiscoveryVersion = DiscoveryV2
	settings.PolicyDigest = settings.DiscoveryPolicy(false).Digest()
	touches, touchProblems := activeBranchTouchesWithSettings(root, settings)
	if len(touchProblems) != 0 || strings.Join(touches[path], ",") != "refs/heads/feature/nested-plan" {
		t.Fatalf("nested branch touches=%#v problems=%#v", touches, touchProblems)
	}
}

func TestPreviewErrorsDoNotChangeAuthoritativeCompleteness(t *testing.T) {
	root := t.TempDir()
	runGitTest(t, root, "init")
	runGitTest(t, root, "config", "user.email", "test@example.invalid")
	runGitTest(t, root, "config", "user.name", "Plan Catalog Test")
	tracked := "docs/tracked/tracked_discovery_doc.md"
	preview := "docs/preview/preview_discovery_doc.md"
	writeClassifierFile(t, root, tracked, canonicalClassifierDocument("example.discovery.tracked-preview-test", KindDiscovery, "Tracked"))
	runGitTest(t, root, "add", tracked)
	runGitTest(t, root, "commit", "-m", "add tracked plan")
	writeClassifierFile(t, root, preview, strings.Replace(canonicalClassifierDocument("example.discovery.preview-invalid", KindDiscovery, "Invalid Preview"), "  schema_version: 1", "  schema_version: [", 1))
	snapshot, err := LoadRepositorySnapshotWithOptions(root, SnapshotOptions{TargetDiscoveryVersion: DiscoveryV2, IncludeUntracked: true})
	if err != nil {
		t.Fatal(err)
	}
	if !snapshot.Complete || HasErrors(snapshot.Problems) || !HasErrors(snapshot.PreviewProblems) || len(snapshot.Records) != 1 || len(snapshot.PreviewCandidates) != 1 {
		t.Fatalf("preview error contaminated authoritative completeness: %#v", snapshot)
	}
}

func TestPreviewValidationUsesTrackedIdentityFamilyAndPortableContext(t *testing.T) {
	root := t.TempDir()
	runGitTest(t, root, "init")
	runGitTest(t, root, "config", "user.email", "test@example.invalid")
	runGitTest(t, root, "config", "user.name", "Plan Catalog Test")
	familyID := "example.family.preview-context"
	trackedID := "example.discovery.tracked-context"
	writeClassifierFile(t, root, "docs/family/README.md", canonicalClassifierDocument(familyID, KindFamily, "Tracked Family"))
	writeClassifierFile(t, root, "docs/tracked/tracked_discovery_doc.md", canonicalClassifierDocument(trackedID, KindDiscovery, "Tracked Plan"))
	runGitTest(t, root, "add", ".")
	runGitTest(t, root, "commit", "-m", "add tracked preview context")

	child := strings.Replace(canonicalClassifierDocument("example.discovery.preview-child", KindDiscovery, "Preview Child"), "  purpose: Exercise repository-wide discovery.\n", "  purpose: Exercise repository-wide discovery.\n  family: "+familyID+"\n", 1)
	writeClassifierFile(t, root, "docs/authoring/child_discovery_doc.md", child)
	writeClassifierFile(t, root, "docs/authoring/duplicate_discovery_doc.md", canonicalClassifierDocument(trackedID, KindDiscovery, "Duplicate Preview"))
	snapshot, err := LoadRepositorySnapshotWithOptions(root, SnapshotOptions{TargetDiscoveryVersion: DiscoveryV2, IncludeUntracked: true})
	if err != nil {
		t.Fatal(err)
	}
	if !problemCodeExists(snapshot.PreviewProblems, "duplicate_plan_id") {
		t.Fatalf("preview duplicate did not see tracked identity: %#v", snapshot.PreviewProblems)
	}
	for _, problem := range snapshot.PreviewProblems {
		if problem.Code == "family_record_missing" && problem.Path == "docs/authoring/child_discovery_doc.md" {
			t.Fatalf("preview child did not see tracked family context: %#v", snapshot.PreviewProblems)
		}
	}
	if HasErrors(snapshot.Problems) || !snapshot.Complete {
		t.Fatalf("preview context contaminated authoritative result: complete=%t problems=%#v", snapshot.Complete, snapshot.Problems)
	}

	preview := Classification{
		Candidates: []Candidate{{Path: "docs/PLAN_discovery_doc.md", Provenance: CandidateProvenance{Source: GitBlob{Preview: true}}}},
		Discovery:  Discovery{Records: []Record{}, Problems: []Problem{}},
	}
	tracked := []Candidate{{Path: "docs/plan_discovery_doc.md"}}
	portable := ValidatePreviewContext(nil, tracked, preview, ModeCanonical, "example")
	if !problemCodeExists(portable, "plan_path_portable_collision") {
		t.Fatalf("preview portable collision did not see tracked candidate: %#v", portable)
	}
}

func TestCanonicalTargetModeImpliesProspectiveDiscoveryV2(t *testing.T) {
	root := migrationRepository(t)
	snapshot, err := LoadRepositorySnapshotWithOptions(root, SnapshotOptions{TargetCatalogMode: ModeCanonical})
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Settings.DiscoveryVersion != DiscoveryV2 || snapshot.Settings.Mode != ModeCanonical || snapshot.CandidateSetDigest == "" || len(snapshot.Candidates) == 0 {
		t.Fatalf("canonical target did not produce complete discovery-v2 evidence: %#v", snapshot)
	}
}

func TestInventoryKindRemainsFilenameOwnedOnMetadataMismatch(t *testing.T) {
	root := t.TempDir()
	runGitTest(t, root, "init")
	runGitTest(t, root, "config", "user.email", "test@example.invalid")
	runGitTest(t, root, "config", "user.name", "Plan Catalog Test")
	planPath := "docs/repo/plans/mismatch/mismatch_discovery_doc.md"
	writeClassifierFile(t, root, planPath, canonicalClassifierDocument("example.implementation.mismatch", KindImplementation, "Mismatched Kind"))
	runGitTest(t, root, "add", ".")
	runGitTest(t, root, "commit", "-m", "add mismatched plan")
	inventory, err := BuildInventory(root, time.Date(2026, 8, 6, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	record := inventoryRecordByPath(t, inventory, planPath)
	if record.Kind != KindDiscovery || !problemCodeExists(inventory.Problems, "record_kind_path_mismatch") {
		t.Fatalf("inventory kind was not filename-owned: record=%#v problems=%#v", record, inventory.Problems)
	}
}

func TestV2ActiveBranchTouchesKeepBothEligibilityChangingRenamePaths(t *testing.T) {
	root := migrationRepository(t)
	oldPath := "holding/renamed_discovery_doc.md"
	newPath := "products/widget/docs/discovery/renamed_discovery_doc.md"
	writeClassifierFile(t, root, oldPath, canonicalClassifierDocument("example.discovery.renamed", KindDiscovery, "Renamed"))
	runGitTest(t, root, "add", oldPath)
	runGitTest(t, root, "commit", "-m", "add misplaced plan")
	runGitTest(t, root, "checkout", "-b", "feature/rename-into-docs")
	if err := os.MkdirAll(filepath.Dir(filepath.Join(root, filepath.FromSlash(newPath))), 0o755); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "mv", oldPath, newPath)
	runGitTest(t, root, "commit", "-m", "move plan into docs")
	runGitTest(t, root, "checkout", "main")
	settings, problems := LoadRepositorySettings(root)
	if HasErrors(problems) {
		t.Fatalf("settings problems=%#v", problems)
	}
	settings.DiscoveryVersion = DiscoveryV2
	settings.PolicyDigest = settings.DiscoveryPolicy(false).Digest()
	touches, touchProblems := activeBranchTouchesWithSettings(root, settings)
	if len(touchProblems) != 0 || strings.Join(touches[oldPath], ",") != "refs/heads/feature/rename-into-docs" || strings.Join(touches[newPath], ",") != "refs/heads/feature/rename-into-docs" {
		t.Fatalf("rename touches=%#v problems=%#v", touches, touchProblems)
	}
}

func recordPathExists(records []Record, expected string) bool {
	for _, record := range records {
		if record.Path == expected {
			return true
		}
	}
	return false
}

func copyPlanFixture(t *testing.T, name string) string {
	t.Helper()
	source := filepath.Join("..", "..", "tests", "fixtures", "plans", name)
	root := t.TempDir()
	if err := filepath.WalkDir(source, func(sourcePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, sourcePath)
		if err != nil || relative == "." {
			return err
		}
		target := filepath.Join(root, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(sourcePath)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	}); err != nil {
		t.Fatal(err)
	}
	return root
}
