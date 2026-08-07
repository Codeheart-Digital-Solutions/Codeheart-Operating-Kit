package plancatalog

import (
	"path"
)

func qualifiesAsCanonicalFamily(record Record) bool {
	if path.Base(record.Path) != "README.md" || record.Metadata == nil || record.Metadata.Kind != KindFamily {
		return false
	}
	identity, err := ParseIdentity(record.Metadata.ID)
	return err == nil && identity.Kind == KindFamily
}

func validateFamilyRecord(record Record, discoveryVersion DiscoveryVersion) []Problem {
	if record.Metadata == nil || record.Metadata.Kind != KindFamily {
		return nil
	}
	if discoveryVersion != DiscoveryV2 {
		if !record.FamilyQualified {
			return []Problem{{Code: "family_placement_invalid", Message: "family metadata requires a qualifying legacy family README", Path: record.Path, PlanID: record.Metadata.ID, Severity: SeverityError, Remediation: "retain the recognized v1 family shape until reviewed discovery-v2 migration"}}
		}
		return nil
	}
	if !qualifiesAsCanonicalFamily(record) || record.ExpectedKind != KindFamily || !record.FamilyQualified {
		return []Problem{{Code: "family_placement_invalid", Message: "family authority requires exact README.md plus valid family metadata in an owned docs path", Path: record.Path, PlanID: record.Metadata.ID, Severity: SeverityError, Remediation: "use an exact README.md in the intended owned documentation family directory"}}
	}
	return nil
}
