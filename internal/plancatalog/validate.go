package plancatalog

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var legacyAliasPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]*$`)

var relationKinds = map[string]bool{
	"parent":        true,
	"child":         true,
	"supersedes":    true,
	"superseded-by": true,
	"depends-on":    true,
	"blocks":        true,
	"related":       true,
}

func ValidateRecords(records []Record, mode CatalogMode, repositoryID string) []Problem {
	problems := []Problem{}
	byID := map[string][]Record{}
	families := map[string]bool{}
	for _, record := range records {
		if record.Header.LegacyStatus != "" {
			problems = append(problems, Problem{Code: "legacy_header_status", Message: fmt.Sprintf("legacy status %q is represented as lifecycle %q", record.Header.LegacyStatus, record.Header.Lifecycle), Path: record.Path, Severity: SeverityWarning, Remediation: "preserve the historical header during metadata-only migration and semantically review lifecycle before canonical cutover"})
		}
		if record.Header.CompatibilityTitle {
			problems = append(problems, Problem{Code: "legacy_title_layout", Message: "implementation layout uses the generic H1 'Document Header' with compatibility title handling", Path: record.Path, Severity: SeverityWarning, Remediation: "new canonical plans use their human-readable title as the direct H1; retain this compatibility layout only for existing documents"})
		}
		if record.Metadata == nil {
			continue
		}
		metadata := record.Metadata
		for _, problem := range ValidateIdentity(metadata.ID, repositoryID, metadata.Kind) {
			problem.Path = record.Path
			problems = append(problems, problem)
		}
		if record.ExpectedKind != "" && metadata.Kind != record.ExpectedKind {
			problems = append(problems, Problem{Code: "record_kind_path_mismatch", Message: fmt.Sprintf("metadata kind %q does not match canonical path kind %q", metadata.Kind, record.ExpectedKind), Path: record.Path, PlanID: metadata.ID, Severity: SeverityError, Remediation: "use a matching canonical filename or family README"})
		}
		if metadata.SchemaVersion != 1 {
			problems = append(problems, Problem{Code: "metadata_schema_unsupported", Message: fmt.Sprintf("metadata schema version %d is unsupported", metadata.SchemaVersion), Path: record.Path, PlanID: metadata.ID, Severity: SeverityError})
		}
		for field, value := range map[string]string{"first_cataloged": metadata.FirstCataloged, "catalog_metadata_updated": metadata.CatalogMetadataUpdated} {
			if _, err := time.Parse(time.RFC3339, value); err != nil || !strings.HasSuffix(value, "Z") {
				problems = append(problems, Problem{Code: "catalog_timestamp_invalid", Message: fmt.Sprintf("%s must be a UTC RFC3339 timestamp ending in Z", field), Path: record.Path, PlanID: metadata.ID, Severity: SeverityError})
			}
		}
		problems = append(problems, validateRelations(record)...)
		problems = append(problems, validateAliases(record)...)
		byID[metadata.ID] = append(byID[metadata.ID], record)
		if metadata.Kind == KindFamily {
			families[metadata.ID] = true
			if !record.FamilyQualified {
				problems = append(problems, Problem{Code: "family_placement_invalid", Message: "family metadata requires a qualifying second-sibling family README", Path: record.Path, PlanID: metadata.ID, Severity: SeverityError})
			}
		}
	}
	for id, matches := range byID {
		if len(matches) < 2 {
			continue
		}
		paths := make([]string, 0, len(matches))
		for _, match := range matches {
			paths = append(paths, match.Path)
		}
		sort.Strings(paths)
		problems = append(problems, Problem{Code: "duplicate_plan_id", Message: fmt.Sprintf("plan ID %q is used by %s", id, strings.Join(paths, ", ")), PlanID: id, Severity: SeverityError, Remediation: "retain each source observation and resolve the canonical identity conflict"})
	}
	for _, record := range records {
		if record.Metadata == nil || record.Metadata.Family == "" {
			continue
		}
		identity, err := ParseIdentity(record.Metadata.Family)
		if err != nil || identity.Kind != KindFamily {
			problems = append(problems, Problem{Code: "family_reference_invalid", Message: fmt.Sprintf("family %q must be a semantic family plan ID", record.Metadata.Family), Path: record.Path, PlanID: record.Metadata.ID, Severity: SeverityError})
			continue
		}
		if !families[record.Metadata.Family] {
			problems = append(problems, Problem{Code: "family_record_missing", Message: fmt.Sprintf("family %q has no qualifying canonical family record", record.Metadata.Family), Path: record.Path, PlanID: record.Metadata.ID, Severity: SeverityWarning, Remediation: "create the family README when the second-sibling trigger is reached"})
		}
	}
	if mode != ModeLegacy && mode != ModeMixed && mode != ModeCanonical {
		problems = append(problems, Problem{Code: "catalog_mode_invalid", Message: fmt.Sprintf("catalog mode %q is invalid", mode), Severity: SeverityError})
	}
	SortProblems(problems)
	return problems
}

func validateRelations(record Record) []Problem {
	problems := []Problem{}
	seen := map[string]bool{}
	for _, relation := range record.Metadata.Relations {
		key := relation.Kind + "\x00" + relation.Target
		if seen[key] {
			problems = append(problems, Problem{Code: "relation_duplicate", Message: fmt.Sprintf("duplicate %s relation to %s", relation.Kind, relation.Target), Path: record.Path, PlanID: record.Metadata.ID, Severity: SeverityError})
		}
		seen[key] = true
		if !relationKinds[relation.Kind] {
			problems = append(problems, Problem{Code: "relation_kind_invalid", Message: fmt.Sprintf("relation kind %q is not reserved", relation.Kind), Path: record.Path, PlanID: record.Metadata.ID, Severity: SeverityError})
		}
		if _, err := ParseIdentity(relation.Target); err != nil {
			problems = append(problems, Problem{Code: "relation_target_invalid", Message: err.Error(), Path: record.Path, PlanID: record.Metadata.ID, Severity: SeverityError})
		}
	}
	return problems
}

func validateAliases(record Record) []Problem {
	problems := []Problem{}
	seen := map[string]bool{}
	for _, alias := range record.Metadata.LegacyAliases {
		if !legacyAliasPattern.MatchString(alias) {
			problems = append(problems, Problem{Code: "legacy_alias_invalid", Message: fmt.Sprintf("legacy alias %q has an invalid portable shape", alias), Path: record.Path, PlanID: record.Metadata.ID, Severity: SeverityError})
		}
		if seen[alias] {
			problems = append(problems, Problem{Code: "legacy_alias_duplicate", Message: fmt.Sprintf("legacy alias %q is duplicated", alias), Path: record.Path, PlanID: record.Metadata.ID, Severity: SeverityError})
		}
		seen[alias] = true
	}
	return problems
}
