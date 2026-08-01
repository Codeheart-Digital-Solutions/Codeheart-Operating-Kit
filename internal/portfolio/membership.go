package portfolio

import (
	"errors"
	"fmt"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

type MembershipInput struct {
	ConfigData      []byte
	LockData        []byte
	KitMarker       []byte
	KitSourceMarker []byte
	HomeID          string
	Self            bool
}

type MembershipDecision struct {
	Member       bool
	Incomplete   bool
	RepositoryID string
	Reason       string
	Config       Config
}

var errPortfolioBlockMissing = errors.New("portfolio block missing")

func EvaluateMembership(input MembershipInput) MembershipDecision {
	installedMarker := len(input.KitMarker) > 0
	sourceMarker := len(input.KitSourceMarker) > 0 && looksLikeKitSourceMarker(input.KitSourceMarker)
	if !installedMarker && !sourceMarker {
		return MembershipDecision{Reason: "default_branch_kit_marker_missing"}
	}
	if installedMarker {
		if len(input.LockData) == 0 {
			return MembershipDecision{Reason: "default_branch_kit_lock_missing"}
		}
		if err := validateRemoteLock(input.LockData); err != nil {
			return MembershipDecision{Incomplete: true, Reason: "default_branch_kit_lock_invalid: " + err.Error()}
		}
	} else if err := validateKitSourceMarker(input.KitSourceMarker); err != nil {
		return MembershipDecision{Incomplete: true, Reason: "default_branch_kit_source_marker_invalid: " + err.Error()}
	}
	if len(input.ConfigData) == 0 {
		return MembershipDecision{Reason: "default_branch_portfolio_config_missing"}
	}
	config, err := decodeRemoteConfig(input.ConfigData)
	if err != nil {
		if errors.Is(err, errPortfolioBlockMissing) {
			return MembershipDecision{Reason: "default_branch_portfolio_config_missing"}
		}
		return MembershipDecision{Incomplete: true, Reason: "default_branch_portfolio_config_invalid: " + err.Error()}
	}
	if input.Self {
		if config.Role != RoleCoordinationHome {
			return MembershipDecision{RepositoryID: config.RepositoryID, Reason: "coordination_home_self_role_mismatch", Config: config}
		}
	} else if config.Role != RoleMember {
		return MembershipDecision{RepositoryID: config.RepositoryID, Reason: "repository_role_is_not_member", Config: config}
	}
	if config.SchemaVersion != 2 {
		return MembershipDecision{RepositoryID: config.RepositoryID, Reason: "portfolio_v2_required_for_enrollment", Config: config}
	}
	if config.RepositoryID == "" {
		return MembershipDecision{Reason: "member_repository_id_missing", Config: config}
	}
	if config.CoordinationHomeID != input.HomeID {
		return MembershipDecision{RepositoryID: config.RepositoryID, Reason: "coordination_home_id_mismatch", Config: config}
	}
	return MembershipDecision{Member: true, RepositoryID: config.RepositoryID, Config: config}
}

func decodeRemoteConfig(data []byte) (Config, error) {
	if len(data) == 0 {
		return Config{}, fmt.Errorf("config missing")
	}
	value, err := state.DecodeYAMLMap(data)
	if err != nil {
		return Config{}, err
	}
	if value["component_settings"] == nil {
		value["component_settings"] = map[string]any{}
	}
	if err := state.Validate(state.ConfigV1Schema, value); err != nil {
		return Config{}, err
	}
	portfolio := state.Map(value["portfolio"])
	if portfolio == nil {
		return Config{}, errPortfolioBlockMissing
	}
	version := state.AsInt(portfolio["schema_version"])
	if version == 0 {
		version = 1
	}
	return Config{
		SchemaVersion:      version,
		Role:               Role(state.AsString(portfolio["role"])),
		RepositoryID:       state.AsString(portfolio["member_repository_id"]),
		CoordinationHomeID: state.AsString(portfolio["coordination_home_id"]),
		CompatibilityV1:    version == 1,
	}, nil
}

func validateRemoteLock(data []byte) error {
	value, err := state.DecodeYAMLMap(data)
	if err != nil {
		return err
	}
	schema, err := state.SchemaForLockVersion(state.AsInt(value["schema_version"]))
	if err != nil {
		return err
	}
	return state.Validate(schema, value)
}

func looksLikeKitSourceMarker(data []byte) bool {
	value, err := state.DecodeYAMLMap(data)
	if err != nil {
		return false
	}
	for _, field := range []string{"schema_version", "version", "compatibility", "components", "profiles", "consumer_impact"} {
		if value[field] == nil {
			return false
		}
	}
	return true
}

func validateKitSourceMarker(data []byte) error {
	value, err := state.DecodeYAMLMap(data)
	if err != nil {
		return err
	}
	return state.Validate(state.ContentManifestSchema, value)
}
