package state

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/kitfs"
)

type Ownership string
type PresencePolicy string
type UpdateStrategy string
type RemovalStrategy string

const (
	OwnershipManaged          Ownership = "managed"
	OwnershipScaffold         Ownership = "scaffold"
	OwnershipTemplate         Ownership = "template"
	OwnershipGeneratedSurface Ownership = "generated-surface"
	OwnershipLocalUser        Ownership = "local-user"
	OwnershipLocalMachine     Ownership = "local-machine"

	PresenceRequired         PresencePolicy = "required"
	PresenceCreateWhenAbsent PresencePolicy = "create-when-absent"
	PresenceOptional         PresencePolicy = "optional"
	PresenceConditional      PresencePolicy = "conditional"

	UpdateReplace        UpdateStrategy = "replace"
	UpdateManagedSection UpdateStrategy = "managed-section"
	UpdatePreserve       UpdateStrategy = "preserve"
	UpdateAppendOnly     UpdateStrategy = "append-only"
	UpdateReportOnly     UpdateStrategy = "report-only"

	RemovalReconcile RemovalStrategy = "reconcile"
	RemovalPreserve  RemovalStrategy = "preserve"
)

type Node struct {
	Target          string          `json:"target"`
	Source          string          `json:"source,omitempty"`
	Component       string          `json:"component,omitempty"`
	Ownership       Ownership       `json:"ownership"`
	Presence        PresencePolicy  `json:"presence_policy"`
	Update          UpdateStrategy  `json:"update_strategy"`
	Removal         RemovalStrategy `json:"removal_strategy"`
	InstallWhen     string          `json:"install_when,omitempty"`
	RouteID         string          `json:"route_id,omitempty"`
	ExpectedSHA256  string          `json:"expected_sha256,omitempty"`
	DirectoryTarget bool            `json:"directory_target,omitempty"`
}

type Graph struct {
	ProfileID          string   `json:"profile_id"`
	ProfileVersion     string   `json:"profile_version"`
	SelectedComponents []string `json:"selected_components"`
	Nodes              []Node   `json:"nodes"`
	DigestSHA256       string   `json:"digest_sha256"`
}

type StrategyDefaults struct {
	Presence PresencePolicy  `json:"presence_policy"`
	Update   UpdateStrategy  `json:"update_strategy"`
	Removal  RemovalStrategy `json:"removal_strategy"`
}

type ProfileDeclaration struct {
	ID                 string                      `json:"id"`
	Version            string                      `json:"version"`
	Name               string                      `json:"name"`
	Description        string                      `json:"description"`
	SelectedComponents []string                    `json:"selected_components"`
	GeneratedSurfaces  []SurfaceDeclaration        `json:"generated_surfaces"`
	StateDefaults      map[string]StrategyDefaults `json:"state_defaults"`
	UpdateCheck        map[string]any              `json:"update_check"`
}

type SurfaceDeclaration struct {
	Path        string          `json:"path"`
	Ownership   string          `json:"ownership"`
	Source      string          `json:"source"`
	InstallWhen string          `json:"install_when"`
	Presence    PresencePolicy  `json:"presence_policy"`
	Update      UpdateStrategy  `json:"update_strategy"`
	Removal     RemovalStrategy `json:"removal_strategy"`
}

type ComponentDeclaration struct {
	ID             string            `json:"id"`
	Version        string            `json:"version"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Profiles       []string          `json:"profiles"`
	ConsumerImpact []string          `json:"consumer_impact"`
	OwnershipModes []string          `json:"ownership_modes"`
	Files          []FileDeclaration `json:"files"`
}

type FileDeclaration struct {
	Source      string          `json:"source"`
	Target      string          `json:"target"`
	Ownership   string          `json:"ownership"`
	InstallWhen string          `json:"install_when"`
	InstallMode string          `json:"install_mode"`
	Presence    PresencePolicy  `json:"presence_policy"`
	Update      UpdateStrategy  `json:"update_strategy"`
	Removal     RemovalStrategy `json:"removal_strategy"`
	RouteID     string          `json:"route_id"`
}

func LoadProfile(profileID string) (ProfileDeclaration, error) {
	resourcePath := "profiles/" + profileID + ".yaml"
	root, err := loadValidatedResource(resourcePath, ProfileSchema)
	if err != nil {
		return ProfileDeclaration{}, err
	}
	var wrapper struct {
		Profile ProfileDeclaration `json:"profile"`
	}
	if err := remarshal(root, &wrapper); err != nil {
		return ProfileDeclaration{}, fmt.Errorf("decode profile %s: %w", profileID, err)
	}
	return wrapper.Profile, nil
}

func LoadComponent(resourcePath string) (ComponentDeclaration, error) {
	root, err := loadValidatedResource(resourcePath, ComponentSchema)
	if err != nil {
		return ComponentDeclaration{}, err
	}
	var wrapper struct {
		Component ComponentDeclaration `json:"component"`
	}
	if err := remarshal(root, &wrapper); err != nil {
		return ComponentDeclaration{}, fmt.Errorf("decode component %s: %w", resourcePath, err)
	}
	return wrapper.Component, nil
}

func ComponentManifestPaths() ([]string, error) {
	paths := []string{}
	err := kitfs.WalkDir("components", func(resourcePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(resourcePath, "/component.yaml") {
			paths = append(paths, resourcePath)
		}
		return nil
	})
	sort.Strings(paths)
	return paths, err
}

// CompileGraph compiles the selected profile and components into the only runtime desired-state graph.
func CompileGraph(profileID string) (Graph, error) {
	profile, err := LoadProfile(profileID)
	if err != nil {
		return Graph{}, err
	}
	selected := map[string]bool{}
	for _, componentID := range profile.SelectedComponents {
		selected[componentID] = true
	}
	nodes := map[string]Node{}
	paths, err := ComponentManifestPaths()
	if err != nil {
		return Graph{}, err
	}
	found := map[string]bool{}
	for _, manifestPath := range paths {
		component, err := LoadComponent(manifestPath)
		if err != nil {
			return Graph{}, err
		}
		if !selected[component.ID] {
			continue
		}
		found[component.ID] = true
		for _, file := range component.Files {
			node, err := nodeFromFile(profile, component.ID, file)
			if err != nil {
				return Graph{}, err
			}
			if err := addNode(nodes, node); err != nil {
				return Graph{}, err
			}
		}
	}
	for componentID := range selected {
		if !found[componentID] {
			return Graph{}, fmt.Errorf("profile %s selects missing component %s", profileID, componentID)
		}
	}
	for _, surface := range profile.GeneratedSurfaces {
		if _, exists := nodes[surface.Path]; exists {
			continue
		}
		defaults, ok := profile.StateDefaults[surface.Ownership]
		if !ok {
			return Graph{}, fmt.Errorf("profile %s has no state defaults for ownership %s", profileID, surface.Ownership)
		}
		presence, update, removal := surface.Presence, surface.Update, surface.Removal
		if presence == "" {
			presence = defaults.Presence
		}
		if update == "" {
			update = defaults.Update
		}
		if removal == "" {
			removal = defaults.Removal
		}
		node := Node{
			Target:          surface.Path,
			Source:          surface.Source,
			Ownership:       Ownership(surface.Ownership),
			Presence:        presence,
			Update:          update,
			Removal:         removal,
			InstallWhen:     surface.InstallWhen,
			RouteID:         routeID(surface.Path),
			DirectoryTarget: strings.HasSuffix(surface.Path, "/"),
		}
		if err := validateTarget(node.Target); err != nil {
			return Graph{}, err
		}
		if !node.DirectoryTarget && kitfs.Exists(surface.Source) {
			data, err := kitfs.ReadFile(surface.Source)
			if err != nil {
				return Graph{}, err
			}
			digest := sha256.Sum256(data)
			node.ExpectedSHA256 = hex.EncodeToString(digest[:])
		}
		nodes[node.Target] = node
	}
	ordered := make([]Node, 0, len(nodes))
	for _, node := range nodes {
		ordered = append(ordered, node)
	}
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Target < ordered[j].Target })
	graph := Graph{
		ProfileID:          profile.ID,
		ProfileVersion:     profile.Version,
		SelectedComponents: append([]string(nil), profile.SelectedComponents...),
		Nodes:              ordered,
	}
	digestInput, err := json.Marshal(graph)
	if err != nil {
		return Graph{}, err
	}
	digest := sha256.Sum256(digestInput)
	graph.DigestSHA256 = hex.EncodeToString(digest[:])
	return graph, nil
}

func nodeFromFile(profile ProfileDeclaration, componentID string, file FileDeclaration) (Node, error) {
	if err := validateTarget(file.Target); err != nil {
		return Node{}, fmt.Errorf("component %s: %w", componentID, err)
	}
	if !kitfs.Exists(file.Source) {
		return Node{}, fmt.Errorf("component %s source %s does not exist", componentID, file.Source)
	}
	defaults, ok := profile.StateDefaults[file.Ownership]
	if !ok {
		return Node{}, fmt.Errorf("profile %s has no state defaults for ownership %s", profile.ID, file.Ownership)
	}
	presence, update, removal := file.Presence, file.Update, file.Removal
	if presence == "" {
		presence = defaults.Presence
	}
	if update == "" {
		update = defaults.Update
	}
	if removal == "" {
		removal = defaults.Removal
	}
	if file.InstallMode == "managed-block" {
		update = UpdateManagedSection
	}
	data, err := kitfs.ReadFile(file.Source)
	if err != nil {
		return Node{}, err
	}
	digest := sha256.Sum256(data)
	return Node{
		Target:         file.Target,
		Source:         file.Source,
		Component:      componentID,
		Ownership:      Ownership(file.Ownership),
		Presence:       presence,
		Update:         update,
		Removal:        removal,
		InstallWhen:    file.InstallWhen,
		RouteID:        firstNonEmpty(file.RouteID, routeID(file.Target)),
		ExpectedSHA256: hex.EncodeToString(digest[:]),
	}, nil
}

func addNode(nodes map[string]Node, node Node) error {
	if existing, ok := nodes[node.Target]; ok {
		if existing.Source == node.Source && existing.Ownership == node.Ownership && existing.Update == node.Update {
			return nil
		}
		return fmt.Errorf("desired-state target conflict for %s between %s and %s", node.Target, existing.Source, node.Source)
	}
	nodes[node.Target] = node
	return nil
}

func validateTarget(target string) error {
	if target == "" || strings.Contains(target, "\\") || strings.HasPrefix(target, "/") {
		return fmt.Errorf("invalid desired-state target %q", target)
	}
	trimmed := strings.TrimSuffix(target, "/")
	cleaned := path.Clean(trimmed)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") || cleaned != trimmed {
		return fmt.Errorf("invalid desired-state target %q", target)
	}
	return nil
}

func routeID(target string) string {
	if target == "AGENTS.md" {
		return "root-agents"
	}
	if target == ".codeheart/kit/README.md" {
		return "kit-fallback"
	}
	if strings.HasSuffix(target, ".md") && strings.HasPrefix(target, ".codeheart/kit/") {
		return strings.TrimSuffix(strings.TrimPrefix(target, ".codeheart/kit/"), ".md")
	}
	return ""
}

func loadValidatedResource(resourcePath, schemaPath string) (map[string]any, error) {
	data, err := kitfs.ReadFile(resourcePath)
	if err != nil {
		return nil, err
	}
	value, err := DecodeAndValidateYAML(schemaPath, data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", resourcePath, err)
	}
	return value, nil
}

func remarshal(value any, target any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
