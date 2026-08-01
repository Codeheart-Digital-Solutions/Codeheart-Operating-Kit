package plancatalog

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	segmentPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`)
	idPattern      = regexp.MustCompile(`^([a-z0-9](?:[a-z0-9-]*[a-z0-9])?)\.(discovery|implementation|family)\.([a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$`)
)

type Identity struct {
	Repository string `json:"repository"`
	Kind       Kind   `json:"kind"`
	Slug       string `json:"slug"`
}

func ParseIdentity(value string) (Identity, error) {
	match := idPattern.FindStringSubmatch(value)
	if match == nil {
		return Identity{}, fmt.Errorf("identity %q must match <repository>.<discovery|implementation|family>.<stable-slug>", value)
	}
	return Identity{Repository: match[1], Kind: Kind(match[2]), Slug: match[3]}, nil
}

func ValidateIdentity(value, repository string, kind Kind) []Problem {
	identity, err := ParseIdentity(value)
	if err != nil {
		return []Problem{{Code: "invalid_plan_id", Message: err.Error(), PlanID: value, Severity: SeverityError, Remediation: "use a lowercase ASCII repository.kind.stable-slug identity"}}
	}
	problems := []Problem{}
	if repository != "" && identity.Repository != repository {
		problems = append(problems, Problem{Code: "repository_identity_mismatch", Message: fmt.Sprintf("plan repository segment %q does not match configured member_repository_id %q", identity.Repository, repository), PlanID: value, Severity: SeverityError, Remediation: "preserve the configured repository namespace or record an explicit identity migration"})
	}
	if kind != "" && identity.Kind != kind {
		problems = append(problems, Problem{Code: "plan_kind_mismatch", Message: fmt.Sprintf("plan identity kind %q does not match metadata kind %q", identity.Kind, kind), PlanID: value, Severity: SeverityError, Remediation: "make the identity kind segment match the canonical record kind"})
	}
	return problems
}

func NormalizeSegment(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("identifier segment is empty")
	}
	for _, r := range value {
		if r > unicode.MaxASCII {
			return "", fmt.Errorf("identifier segment %q requires an explicit meaningful ASCII transliteration", value)
		}
	}
	var builder strings.Builder
	lastHyphen := false
	for _, r := range strings.ToLower(value) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
			lastHyphen = false
			continue
		}
		if builder.Len() > 0 && !lastHyphen {
			builder.WriteByte('-')
			lastHyphen = true
		}
	}
	normalized := strings.Trim(builder.String(), "-")
	if !segmentPattern.MatchString(normalized) {
		return "", fmt.Errorf("identifier segment %q cannot produce a valid lowercase ASCII segment", value)
	}
	return normalized, nil
}
