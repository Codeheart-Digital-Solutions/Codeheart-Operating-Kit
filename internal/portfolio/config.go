package portfolio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

var errRootRegularChanged = errors.New("contained regular file changed during read")

type Role string

const (
	RoleMember           Role = "member"
	RoleCoordinationHome Role = "coordination-home"
)

type ProviderSource struct {
	Kind  string `json:"kind" yaml:"kind"`
	Owner string `json:"owner" yaml:"owner"`
}

type LocalSource struct {
	Kind string `json:"kind" yaml:"kind"`
	Root string `json:"root" yaml:"root"`
}

type Config struct {
	SchemaVersion      int              `json:"schema_version"`
	Role               Role             `json:"role"`
	RepositoryID       string           `json:"member_repository_id"`
	CoordinationHomeID string           `json:"coordination_home_id"`
	ProviderSources    []ProviderSource `json:"provider_sources"`
	LocalSources       []LocalSource    `json:"local_sources"`
	CompatibilityV1    bool             `json:"compatibility_v1"`
	Root               string           `json:"-"`
}

func LoadConfig(root string) (Config, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return Config{}, err
	}
	data, err := readRootRegular(absolute, ConfigPath)
	if err != nil {
		return Config{}, fmt.Errorf("portfolio_config_unavailable: %w", err)
	}
	value, err := state.DecodeYAMLMap(data)
	if err != nil {
		return Config{}, fmt.Errorf("portfolio_config_invalid: %w", err)
	}
	if value["component_settings"] == nil {
		value["component_settings"] = map[string]any{}
	}
	if err := state.ValidateConfig(value); err != nil {
		return Config{}, fmt.Errorf("portfolio_config_invalid: %w", err)
	}
	portfolio := state.Map(value["portfolio"])
	if portfolio == nil {
		return Config{}, fmt.Errorf("portfolio_not_configured: %s has no portfolio block", ConfigPath)
	}
	config := Config{
		SchemaVersion:      state.AsInt(portfolio["schema_version"]),
		Role:               Role(strings.TrimSpace(state.AsString(portfolio["role"]))),
		RepositoryID:       strings.TrimSpace(state.AsString(portfolio["member_repository_id"])),
		CoordinationHomeID: strings.TrimSpace(state.AsString(portfolio["coordination_home_id"])),
		ProviderSources:    []ProviderSource{},
		LocalSources:       []LocalSource{},
		Root:               absolute,
	}
	if config.SchemaVersion == 0 {
		config.SchemaVersion = 1
		config.CompatibilityV1 = true
	}
	discovery := state.Map(portfolio["discovery"])
	for _, item := range state.AnySlice(discovery["sources"]) {
		source := state.Map(item)
		if source == nil {
			continue
		}
		config.ProviderSources = append(config.ProviderSources, ProviderSource{Kind: state.AsString(source["kind"]), Owner: state.AsString(source["owner"])})
	}
	if config.Role != RoleMember && config.Role != RoleCoordinationHome {
		return Config{}, fmt.Errorf("portfolio_role_invalid: %q", config.Role)
	}
	// Discovery scopes are coordination-home authority. A member may retain a
	// machine-local sources file after a role change, but it must never consume
	// that file or committed provider scopes while operating as a member.
	if config.Role == RoleCoordinationHome {
		if localData, readErr := readRootRegular(absolute, LocalSourcesPath); readErr == nil {
			localValue, decodeErr := state.DecodeAndValidateYAML(state.PortfolioSourcesSchema, localData)
			if decodeErr != nil {
				return Config{}, fmt.Errorf("portfolio_local_sources_invalid: %w", decodeErr)
			}
			for _, item := range state.AnySlice(localValue["sources"]) {
				source := state.Map(item)
				config.LocalSources = append(config.LocalSources, LocalSource{Kind: state.AsString(source["kind"]), Root: state.AsString(source["root"])})
			}
		} else if !os.IsNotExist(readErr) {
			return Config{}, fmt.Errorf("portfolio_local_sources_unavailable: %w", readErr)
		}
	}
	if config.SchemaVersion == 2 && config.RepositoryID == "" {
		return Config{}, fmt.Errorf("portfolio_repository_identity_missing: member_repository_id is required")
	}
	if config.SchemaVersion == 2 && config.CoordinationHomeID == "" {
		return Config{}, fmt.Errorf("portfolio_home_identity_missing: coordination_home_id is required")
	}
	sort.SliceStable(config.ProviderSources, func(i, j int) bool {
		if config.ProviderSources[i].Kind != config.ProviderSources[j].Kind {
			return config.ProviderSources[i].Kind < config.ProviderSources[j].Kind
		}
		return config.ProviderSources[i].Owner < config.ProviderSources[j].Owner
	})
	sort.SliceStable(config.LocalSources, func(i, j int) bool { return config.LocalSources[i].Root < config.LocalSources[j].Root })
	return config, nil
}

func readRootRegular(root, relative string) ([]byte, error) {
	parentName := filepath.ToSlash(filepath.Dir(filepath.FromSlash(relative)))
	bound, err := bindLocalDirectory(root, parentName, false)
	if err != nil {
		return nil, err
	}
	defer bound.Close()
	name := filepath.Base(filepath.FromSlash(relative))
	info, err := bound.directory.Lstat(name)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("%s is not a contained regular file", relative)
	}
	file, err := bound.directory.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return nil, fmt.Errorf("%w: %s identity changed while opening", errRootRegularChanged, relative)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err := bound.Verify(); err != nil {
		return nil, err
	}
	current, err := bound.directory.Lstat(name)
	if err != nil || !os.SameFile(info, current) {
		return nil, fmt.Errorf("%w: %s identity changed while reading", errRootRegularChanged, relative)
	}
	opened, err = file.Stat()
	if err != nil || !os.SameFile(current, opened) {
		return nil, fmt.Errorf("%w: %s descriptor identity changed while reading", errRootRegularChanged, relative)
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	verified, err := io.ReadAll(file)
	if err != nil || !bytes.Equal(data, verified) {
		return nil, fmt.Errorf("%w: %s bytes changed while reading", errRootRegularChanged, relative)
	}
	return data, nil
}

func BuildSources(config Config, runner CommandRunner) []Source {
	if config.Role != RoleCoordinationHome {
		return []Source{}
	}
	sources := []Source{}
	for _, local := range config.LocalSources {
		if local.Kind == "local-git-root" {
			sources = append(sources, NewLocalGitSource(local.Root, runner))
		}
	}
	owners := []string{}
	for _, provider := range config.ProviderSources {
		if provider.Kind == "github-owner" {
			owners = append(owners, provider.Owner)
		}
	}
	if len(owners) > 0 {
		sources = append(sources, NewGitHubSource(owners, runner))
	}
	return sources
}
