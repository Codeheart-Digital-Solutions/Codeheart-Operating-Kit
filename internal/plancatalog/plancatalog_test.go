package plancatalog

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

func TestParseCanonicalMetadataAndCurrentImplementationTitleShape(t *testing.T) {
	discovery := mustFixture(t, "valid-discovery.md")
	record, err := ParseDocument("semantic_discovery_doc.md", discovery, KindDiscovery)
	if err != nil {
		t.Fatal(err)
	}
	if record.Header.Title != "Semantic Catalog Discovery" || record.Metadata.ID != "operating-kit.discovery.semantic-catalog" {
		t.Fatalf("unexpected discovery record: %#v", record)
	}
	if record.ContentSHA256 == "" || record.Legacy {
		t.Fatalf("expected canonical hashed record: %#v", record)
	}

	implementation := mustFixture(t, "valid-implementation.md")
	record, err = ParseDocument("semantic_implementation_doc.md", implementation, KindImplementation)
	if err != nil {
		t.Fatal(err)
	}
	if record.Header.Title != "Semantic Catalog Implementation" || record.Header.TitleLine != 7 {
		t.Fatalf("current implementation title compatibility failed: %#v", record.Header)
	}
	if !record.Header.CompatibilityTitle {
		t.Fatal("current implementation H2 title layout was not reported as compatibility evidence")
	}
	problems := ValidateRecords([]Record{record}, ModeCanonical, "operating-kit")
	if !problemExists(problems, "legacy_title_layout", SeverityWarning) {
		t.Fatalf("compatibility title layout was not diagnosed: %#v", problems)
	}
	if string(implementation) != string(mustFixture(t, "valid-implementation.md")) {
		t.Fatal("parsing changed canonical document bytes")
	}
}

func TestUnsupportedMetadataVersionHasStableDiagnostic(t *testing.T) {
	_, err := ParseDocument("future_discovery_doc.md", mustFixture(t, "unsupported-version.md"), KindDiscovery)
	if err == nil || ErrorCode(err) != "metadata_schema_unsupported" {
		t.Fatalf("error=%v code=%q", err, ErrorCode(err))
	}
}

func TestGenericDocumentHeaderWithoutSpecificTitleUsesObservableCompatibility(t *testing.T) {
	data := []byte("Last updated: 2026-07-31T10:00:00Z (UTC)\nCreated: 2026-07-31\nStatus: draft\n\n# Document Header\n")
	record, err := ParseDocument("compatibility_implementation_doc.md", data, KindImplementation)
	if err == nil || ErrorCode(err) != "metadata_missing" || record.Header.Title != "Document Header" || !record.Header.CompatibilityTitle {
		t.Fatalf("record=%#v error=%v code=%q", record, err, ErrorCode(err))
	}
}

func TestParseRejectsMissingMisplacedMultipleAndUnknownMetadata(t *testing.T) {
	valid := string(mustFixture(t, "valid-discovery.md"))
	tests := []struct {
		name string
		data string
		code string
	}{
		{name: "missing", data: strings.Split(valid, MetadataBeginMarker)[0], code: "metadata_missing"},
		{name: "misplaced", data: strings.Replace(valid, "# Semantic Catalog Discovery\n\n"+MetadataBeginMarker, "# Semantic Catalog Discovery\n\nBody first.\n\n"+MetadataBeginMarker, 1), code: "metadata_misplaced"},
		{name: "multiple", data: valid + "\n" + MetadataBeginMarker + "\n", code: "metadata_marker_count"},
		{name: "unknown", data: string(mustFixture(t, "invalid-metadata.md")), code: "metadata_yaml_invalid"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseDocument(test.name+"_discovery_doc.md", []byte(test.data), KindDiscovery)
			if err == nil || ErrorCode(err) != test.code {
				t.Fatalf("error=%v code=%q, want %q", err, ErrorCode(err), test.code)
			}
		})
	}
}

func TestMetadataMarkersRemainInertBehindNonClosingFenceText(t *testing.T) {
	for _, fence := range []string{"```", "~~~"} {
		data := strings.Join([]string{
			fence + "md",
			fence + "not-a-close",
			MetadataBeginMarker,
			fence,
			"",
		}, "\n")
		if hasGenuineMetadataMarker([]byte(data)) {
			t.Fatalf("%q trailing text incorrectly closed the fence", fence)
		}
	}
}

func TestBacktickInFenceInfoDoesNotOpenFence(t *testing.T) {
	data := []byte("```lang`invalid\n" + MetadataBeginMarker + "\n")
	if !hasGenuineMetadataMarker(data) {
		t.Fatal("a CommonMark-invalid backtick fence opener hid a genuine metadata marker")
	}
}

func TestInvalidFamilyRecordDoesNotSatisfyFamilyReference(t *testing.T) {
	familyMetadata := metadataForTest("example.family.invalid", KindFamily, "Invalid family placement")
	childMetadata := metadataForTest("example.discovery.child", KindDiscovery, "Child")
	childMetadata.Family = familyMetadata.ID
	records := []Record{
		{Path: "docs/family/not-a-readme.md", Metadata: &familyMetadata},
		{Path: "docs/family/child_discovery_doc.md", Metadata: &childMetadata, ExpectedKind: KindDiscovery},
	}
	problems := ValidateRecordsForDiscovery(records, ModeCanonical, "example", DiscoveryV2)
	if !problemExists(problems, "family_placement_invalid", SeverityError) || !problemExists(problems, "family_record_missing", SeverityWarning) {
		t.Fatalf("invalid family record satisfied a reference: %#v", problems)
	}
}

func TestV1RecognizedFamilyShapeRemainsCompatibleUntilV2Migration(t *testing.T) {
	metadata := metadataForTest("example.family.compatibility", KindFamily, "Legacy family compatibility")
	record := Record{Path: "docs/repo/plans/family/readme.md", Metadata: &metadata, ExpectedKind: KindFamily, FamilyQualified: true}
	if problemExists(ValidateRecords([]Record{record}, ModeCanonical, "example"), "family_placement_invalid", SeverityError) {
		t.Fatal("v1-recognized family was invalidated before discovery-v2 activation")
	}
	if !problemExists(ValidateRecordsForDiscovery([]Record{record}, ModeCanonical, "example", DiscoveryV2), "family_placement_invalid", SeverityError) {
		t.Fatal("prospective v2 validation did not require exact README.md family authority")
	}
}

func TestMetadataValidationFailuresHaveStableCodes(t *testing.T) {
	valid := string(mustFixture(t, "valid-discovery.md"))
	tests := []struct {
		name string
		data string
		code string
	}{
		{
			name: "invalid UTC timestamp",
			data: strings.Replace(valid, "first_cataloged: 2026-07-31T10:00:00Z", "first_cataloged: 2026-07-31T12:00:00+02:00", 1),
			code: "catalog_timestamp_invalid",
		},
		{
			name: "duplicate relation target",
			data: strings.Replace(valid, "      target: operating-kit.implementation.semantic-catalog", "      target: operating-kit.implementation.semantic-catalog\n    - kind: related\n      target: operating-kit.implementation.semantic-catalog", 1),
			code: "metadata_schema_invalid",
		},
		{
			name: "malformed legacy alias",
			data: strings.Replace(valid, "    - OK-PR-027", "    - invalid alias", 1),
			code: "metadata_schema_invalid",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseDocument("invalid_discovery_doc.md", []byte(test.data), KindDiscovery)
			if err == nil || ErrorCode(err) != test.code {
				t.Fatalf("error=%v code=%q, want %q", err, ErrorCode(err), test.code)
			}
		})
	}
}

func TestIdentityNormalizationAndOwnership(t *testing.T) {
	normalized, err := NormalizeSegment("Operating Kit")
	if err != nil || normalized != "operating-kit" {
		t.Fatalf("normalized=%q err=%v", normalized, err)
	}
	if _, err := NormalizeSegment("Überblick"); err == nil || !strings.Contains(err.Error(), "transliteration") {
		t.Fatalf("expected transliteration diagnostic, got %v", err)
	}
	problems := ValidateIdentity("other.discovery.semantic-catalog", "operating-kit", KindDiscovery)
	if len(problems) != 1 || problems[0].Code != "repository_identity_mismatch" {
		t.Fatalf("ownership problems=%#v", problems)
	}
	if runtime.GOOS == "windows" && strings.Contains(normalized, ":") {
		t.Fatalf("identity is not portable: %q", normalized)
	}
}

func TestAbsentCatalogModeResolvesToLegacy(t *testing.T) {
	mode, ok := ParseCatalogMode("")
	if !ok || mode != ModeLegacy {
		t.Fatalf("mode=%q ok=%t", mode, ok)
	}
}

func TestRenamedPlanKeepsSemanticIdentity(t *testing.T) {
	original, err := ParseDocument("before_discovery_doc.md", mustFixture(t, "valid-discovery.md"), KindDiscovery)
	if err != nil {
		t.Fatal(err)
	}
	renamed, err := ParseDocument("after_discovery_doc.md", mustFixture(t, "renamed-discovery.md"), KindDiscovery)
	if err != nil {
		t.Fatal(err)
	}
	if original.Metadata.ID != renamed.Metadata.ID || original.Header.Title == renamed.Header.Title || original.Path == renamed.Path {
		t.Fatalf("identity changed across title/path rename: original=%#v renamed=%#v", original, renamed)
	}
	if original.Metadata.FirstCataloged != renamed.Metadata.FirstCataloged || original.Header.LastUpdated == renamed.Header.LastUpdated {
		t.Fatalf("catalog chronology and content chronology were conflated: original=%#v renamed=%#v", original, renamed)
	}
}

func TestDiscoverPreservesDuplicateIDsAndCatalogModes(t *testing.T) {
	duplicateRoot := filepath.Join("..", "..", "tests", "fixtures", "plans", "same-id-repository")
	discovery, err := Discover(duplicateRoot, ModeCanonical, "operating-kit")
	if err != nil {
		t.Fatal(err)
	}
	if len(discovery.Records) != 2 {
		t.Fatalf("records=%d", len(discovery.Records))
	}
	if !problemExists(discovery.Problems, "duplicate_plan_id", SeverityError) {
		t.Fatalf("duplicate identity was not preserved and diagnosed: %#v", discovery.Problems)
	}

	mixedRoot := filepath.Join("..", "..", "tests", "fixtures", "plans", "mixed-repository")
	mixed, err := Discover(mixedRoot, ModeMixed, "operating-kit")
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := Discover(mixedRoot, ModeCanonical, "operating-kit")
	if err != nil {
		t.Fatal(err)
	}
	if len(mixed.Records) != 2 {
		t.Fatalf("mixed records=%d", len(mixed.Records))
	}
	if !problemExists(mixed.Problems, "metadata_missing", SeverityWarning) || !problemExists(canonical.Problems, "metadata_missing", SeverityError) {
		t.Fatalf("mode problems mixed=%#v canonical=%#v", mixed.Problems, canonical.Problems)
	}
}

func TestFamilyREADMEQualificationAndDerivedMembership(t *testing.T) {
	root := filepath.Join("..", "..", "tests", "fixtures", "plans", "family-repository")
	result, err := Discover(root, ModeCanonical, "operating-kit")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 3 || problemExists(result.Problems, "family_placement_invalid", SeverityError) || problemExists(result.Problems, "family_record_missing", SeverityWarning) {
		t.Fatalf("family discovery records=%d problems=%#v", len(result.Records), result.Problems)
	}
	members := DerivedFamilyMembers(result.Records)["operating-kit.family.catalog-foundation"]
	if strings.Join(members, ",") != "operating-kit.discovery.catalog,operating-kit.implementation.catalog" {
		t.Fatalf("derived family members=%#v", members)
	}
}

func TestMalformedMetadataNeverFallsBackToLegacyAndInvalidFamilyPlacementIsDiagnosed(t *testing.T) {
	root := t.TempDir()
	planDir := filepath.Join(root, "docs", "repo", "plans", "invalid")
	if err := os.MkdirAll(planDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(planDir, "invalid_discovery_doc.md"), mustFixture(t, "invalid-metadata.md"), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := Discover(root, ModeMixed, "operating-kit")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 0 || problemExists(result.Problems, "metadata_missing", SeverityWarning) || !problemExists(result.Problems, "metadata_yaml_invalid", SeverityError) {
		t.Fatalf("malformed metadata fell back to legacy: records=%#v problems=%#v", result.Records, result.Problems)
	}

	family := Record{
		Path:            "docs/repo/plans/not-a-family/README.md",
		ExpectedKind:    KindFamily,
		FamilyQualified: false,
		Metadata: &Metadata{
			SchemaVersion:          1,
			ID:                     "operating-kit.family.invalid-placement",
			Kind:                   KindFamily,
			Purpose:                "Exercise the second-sibling trigger",
			FirstCataloged:         "2026-07-31T10:00:00Z",
			CatalogMetadataUpdated: "2026-07-31T10:00:00Z",
		},
	}
	if !problemExists(ValidateRecords([]Record{family}, ModeCanonical, "operating-kit"), "family_placement_invalid", SeverityError) {
		t.Fatal("non-qualifying family README was not rejected")
	}
}

func TestNewSchemasCompileAndPortfolioV2RequiresHomeRepositoryIdentity(t *testing.T) {
	for _, schema := range []string{state.PlanMetadataSchema, state.PlanCatalogSchema, state.PlanMigrationSchema, state.PortfolioSourcesSchema, state.PortfolioOverlaySchema} {
		if err := state.Validate(schema, minimalSchemaInstance(schema)); err != nil {
			t.Fatalf("%s: %v", schema, err)
		}
	}
	home := mustPortfolioFixture(t, "kit-config-v2-home.yaml")
	if _, err := state.DecodeAndValidateYAML(state.ConfigV1Schema, home); err != nil {
		t.Fatalf("valid home: %v", err)
	}
	invalid := mustPortfolioFixture(t, "kit-config-v2-home-missing-repository-id.yaml")
	if _, err := state.DecodeAndValidateYAML(state.ConfigV1Schema, invalid); err == nil {
		t.Fatal("v2 home without member_repository_id validated")
	}
	commit := strings.Repeat("b", 40)
	self := CoordinationHomeSelfMember("example-home-repository", "refs/heads/main", commit)
	self.DiscoveryVersion = DiscoveryV2
	self.PolicyDigest = strings.Repeat("d", 64)
	self.CandidateSetDigest = strings.Repeat("e", 64)
	complete := true
	self.Complete = &complete
	if !self.SelfMember || self.RepositoryID != "example-home-repository" {
		t.Fatalf("home self-member=%#v", self)
	}
	catalog := minimalSchemaInstance(state.PlanCatalogSchema).(map[string]any)
	catalog["members"] = []any{self}
	observation := map[string]any{
		"repository_id": "example-home-repository", "plan_id": "example-home-repository.discovery.catalog", "title": "Catalog", "kind": "discovery", "purpose": "purpose", "lifecycle": "active", "canonical_path": "docs/repo/plans/catalog_discovery_doc.md", "ref": "refs/heads/main", "commit": commit, "content_sha256": strings.Repeat("c", 64), "visibility": "default", "verification": "verified", "observed_at": "2026-07-31T10:00:00Z", "stale": false, "conflict": true,
	}
	branchObservation := map[string]any{}
	for key, value := range observation {
		branchObservation[key] = value
	}
	branchObservation["ref"] = "refs/heads/feature/catalog"
	branchObservation["commit"] = strings.Repeat("d", 40)
	branchObservation["visibility"] = "unmerged-branch"
	catalog["observations"] = []any{observation, branchObservation}
	if err := state.Validate(state.PlanCatalogSchema, catalog); err != nil {
		t.Fatalf("same-ID observations and home self-member must remain a valid list: %v", err)
	}
	observations := catalog["observations"].([]any)
	if len(observations) != 2 || reflect.DeepEqual(observations[0], observations[1]) {
		t.Fatalf("parallel same-ID observations were collapsed: %#v", observations)
	}

	v1 := mustPortfolioFixture(t, "kit-config-v1-member.yaml")
	if _, err := state.DecodeAndValidateYAML(state.ConfigV1Schema, v1); err != nil {
		t.Fatalf("portfolio-v1 compatibility fixture: %v", err)
	}
}

func TestRepositorySettingsDiscoveryDefaultsAndExclusionSafety(t *testing.T) {
	base := "schema_version: 1\nselected_profile: standard\nproject_display_name: Example\nselected_setup_folder: .\nlocal_consumer_layer:\n  repo_docs_path: docs/repo/\n  agent_memory_path: docs/agent-memory/\n  user_layer_path: .codeheart/user/\ncomponent_settings: {}\n"
	settings, problems := DecodeRepositorySettings([]byte(base))
	if HasErrors(problems) || settings.DiscoveryVersion != DiscoveryV1 || len(settings.ExcludedRoots) != 0 {
		t.Fatalf("legacy defaults settings=%#v problems=%#v", settings, problems)
	}
	v2 := strings.Replace(base, "component_settings: {}", "component_settings:\n  planning-workflows:\n    plan_catalog_discovery_version: 2\n    plan_catalog_ownership:\n      excluded_roots:\n        - vendor/docs/\n        - generated/docs/", 1)
	settings, problems = DecodeRepositorySettings([]byte(v2))
	if HasErrors(problems) || settings.DiscoveryVersion != DiscoveryV2 || !reflect.DeepEqual(settings.ExcludedRoots, []string{"generated/docs/", "vendor/docs/"}) || len(settings.PolicyDigest) != 64 {
		t.Fatalf("v2 settings=%#v problems=%#v", settings, problems)
	}
	if settings.DiscoveryPolicy(false).Digest() != settings.DiscoveryPolicy(true).Digest() {
		t.Fatal("non-authoritative untracked preview changed the policy digest")
	}
	for name, roots := range map[string]string{
		"owned roots": "      owned_roots: [docs/]",
		"absolute":    "      excluded_roots: [/vendor/docs/]",
		"drive":       "      excluded_roots: ['C:/vendor/docs/']",
		"glob":        "      excluded_roots: ['vendor/*/']",
		"escape":      "      excluded_roots: ['vendor/../docs/']",
		"overlap":     "      excluded_roots: ['vendor/', 'vendor/docs/']",
		"collision":   "      excluded_roots: ['Vendor/docs/', 'vendor/docs/']",
	} {
		t.Run(name, func(t *testing.T) {
			config := strings.Replace(base, "component_settings: {}", "component_settings:\n  planning-workflows:\n    plan_catalog_discovery_version: 2\n    plan_catalog_ownership:\n"+roots, 1)
			_, found := DecodeRepositorySettings([]byte(config))
			if !HasErrors(found) {
				t.Fatalf("unsafe config accepted: %s", config)
			}
		})
	}
	extension := strings.Replace(base, "component_settings: {}", "component_settings:\n  planning-workflows:\n    extension_setting: retained", 1)
	if _, found := DecodeRepositorySettings([]byte(extension)); HasErrors(found) {
		t.Fatalf("unrelated planning extension was rejected: %#v", found)
	}
}

func minimalSchemaInstance(schema string) any {
	switch schema {
	case state.PlanMetadataSchema:
		return map[string]any{"plan": map[string]any{"schema_version": 1, "id": "repo.discovery.slug", "kind": "discovery", "purpose": "purpose", "first_cataloged": "2026-07-31T10:00:00Z", "catalog_metadata_updated": "2026-07-31T10:00:00Z"}}
	case state.PlanCatalogSchema:
		return map[string]any{"schema_version": 2, "discovery_version": 2, "policy_digest": strings.Repeat("d", 64), "candidate_set_digest": strings.Repeat("e", 64), "coordination_home_id": "home", "started_at": "2026-07-31T10:00:00Z", "completed_at": "2026-07-31T10:00:00Z", "complete": true, "members": []any{}, "observations": []any{}, "candidates": []any{}, "errors": []any{}, "metrics": map[string]any{"duration_ms": 0, "source_count": 0, "member_count": 0, "candidate_count": 0, "observation_count": 0, "stale_count": 0, "api_call_count": 0, "max_concurrency": 1}}
	case state.PlanMigrationSchema:
		return map[string]any{"schema_version": 2, "repository_id": "repo", "discovery_version": 2, "target_catalog_mode": "canonical", "policy_digest": strings.Repeat("d", 64), "candidate_set_digest": strings.Repeat("e", 64), "inventory_revision": strings.Repeat("a", 40), "reviewed_at": "2026-07-31T10:00:00Z", "records": []any{}}
	case state.PortfolioSourcesSchema:
		return map[string]any{"schema_version": 1, "sources": []any{}}
	case state.PortfolioOverlaySchema:
		return map[string]any{"schema_version": 1, "coordination_home_id": "home", "families": []any{}, "themes": []any{}, "relations": []any{}, "priorities": []any{}, "analyses": []any{}}
	default:
		return map[string]any{}
	}
}

func mustFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "plans", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func mustPortfolioFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "portfolio", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func problemExists(problems []Problem, code string, severity Severity) bool {
	for _, problem := range problems {
		if problem.Code == code && problem.Severity == severity {
			return true
		}
	}
	return false
}
