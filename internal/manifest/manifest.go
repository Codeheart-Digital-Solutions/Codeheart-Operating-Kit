package manifest

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/kitfs"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/yamlmini"
)

type GeneratedSurface struct {
	Path        string
	Ownership   string
	Source      string
	InstallWhen string
}

type Profile struct {
	ID                 string
	Version            string
	Name               string
	Description        string
	SelectedComponents []string
	GeneratedSurfaces  []GeneratedSurface
	Raw                map[string]any
}

type ComponentFile struct {
	Source      string
	Target      string
	Ownership   string
	InstallWhen string
	Component   string
}

type Component struct {
	ID             string
	Version        string
	Name           string
	Description    string
	Profiles       []string
	ConsumerImpact []string
	OwnershipModes []string
	Files          []ComponentFile
	ManifestPath   string
	Raw            map[string]any
}

type ReleaseAsset struct {
	Name     string
	URL      string
	SHA256   string
	Platform string
}

type ReleaseManifest struct {
	SchemaVersion  int
	Version        string
	ReleasedAt     string
	Assets         []ReleaseAsset
	ConsumerImpact []string
	Raw            map[string]any
}

func LoadYAMLResource(resourcePath string) (map[string]any, error) {
	text, err := kitfs.ReadText(resourcePath)
	if err != nil {
		return nil, err
	}
	return yamlmini.MustMap(text)
}

func LoadProfile(profileID string) (Profile, error) {
	root, err := LoadYAMLResource("profiles/" + profileID + ".yaml")
	if err != nil {
		return Profile{}, err
	}
	profileMap, err := requiredMap(root, "profile")
	if err != nil {
		return Profile{}, err
	}
	profile := Profile{
		ID:                 stringValue(profileMap["id"]),
		Version:            stringValue(profileMap["version"]),
		Name:               stringValue(profileMap["name"]),
		Description:        stringValue(profileMap["description"]),
		SelectedComponents: stringSlice(profileMap["selected_components"]),
		Raw:                profileMap,
	}
	for _, value := range anySlice(profileMap["generated_surfaces"]) {
		mapping, ok := value.(map[string]any)
		if !ok {
			continue
		}
		profile.GeneratedSurfaces = append(profile.GeneratedSurfaces, GeneratedSurface{
			Path:        stringValue(mapping["path"]),
			Ownership:   stringValue(mapping["ownership"]),
			Source:      stringValue(mapping["source"]),
			InstallWhen: stringValue(mapping["install_when"]),
		})
	}
	return profile, nil
}

func ComponentManifestPaths() ([]string, error) {
	paths := []string{}
	err := kitfs.WalkDir("components", func(resourcePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if strings.HasSuffix(resourcePath, "/component.yaml") {
			paths = append(paths, resourcePath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	return paths, nil
}

func LoadComponent(resourcePath string) (Component, error) {
	root, err := LoadYAMLResource(resourcePath)
	if err != nil {
		return Component{}, err
	}
	componentMap, err := requiredMap(root, "component")
	if err != nil {
		return Component{}, err
	}
	component := Component{
		ID:             stringValue(componentMap["id"]),
		Version:        stringValue(componentMap["version"]),
		Name:           stringValue(componentMap["name"]),
		Description:    stringValue(componentMap["description"]),
		Profiles:       stringSlice(componentMap["profiles"]),
		ConsumerImpact: stringSlice(componentMap["consumer_impact"]),
		OwnershipModes: stringSlice(componentMap["ownership_modes"]),
		ManifestPath:   resourcePath,
		Raw:            componentMap,
	}
	for _, value := range anySlice(componentMap["files"]) {
		mapping, ok := value.(map[string]any)
		if !ok {
			continue
		}
		component.Files = append(component.Files, ComponentFile{
			Source:      stringValue(mapping["source"]),
			Target:      stringValue(mapping["target"]),
			Ownership:   stringValue(mapping["ownership"]),
			InstallWhen: stringValue(mapping["install_when"]),
			Component:   component.ID,
		})
	}
	return component, nil
}

func LoadComponents(profileID string) ([]Component, error) {
	profile, err := LoadProfile(profileID)
	if err != nil {
		return nil, err
	}
	selected := map[string]bool{}
	for _, componentID := range profile.SelectedComponents {
		selected[componentID] = true
	}
	paths, err := ComponentManifestPaths()
	if err != nil {
		return nil, err
	}
	components := []Component{}
	for _, path := range paths {
		component, err := LoadComponent(path)
		if err != nil {
			return nil, err
		}
		if selected[component.ID] {
			components = append(components, component)
		}
	}
	return components, nil
}

func IterComponentFiles(profileID string) ([]ComponentFile, error) {
	components, err := LoadComponents(profileID)
	if err != nil {
		return nil, err
	}
	files := []ComponentFile{}
	for _, component := range components {
		files = append(files, component.Files...)
	}
	return files, nil
}

func LoadReleaseManifest() (ReleaseManifest, error) {
	root, err := LoadYAMLResource("manifest.yaml")
	if err != nil {
		return ReleaseManifest{}, err
	}
	result := ReleaseManifest{
		SchemaVersion:  intValue(root["schema_version"]),
		Version:        stringValue(root["version"]),
		ReleasedAt:     stringValue(root["released_at"]),
		ConsumerImpact: stringSlice(root["consumer_impact"]),
		Raw:            root,
	}
	for _, value := range anySlice(root["assets"]) {
		mapping, ok := value.(map[string]any)
		if !ok {
			continue
		}
		result.Assets = append(result.Assets, ReleaseAsset{
			Name:     stringValue(mapping["name"]),
			URL:      stringValue(mapping["url"]),
			SHA256:   stringValue(mapping["sha256"]),
			Platform: stringValue(mapping["platform"]),
		})
	}
	return result, nil
}

func requiredMap(root map[string]any, key string) (map[string]any, error) {
	value, ok := root[key]
	if !ok {
		return nil, fmt.Errorf("missing YAML key %q", key)
	}
	mapping, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("YAML key %q is not a mapping", key)
	}
	return mapping, nil
}

func anySlice(value any) []any {
	switch item := value.(type) {
	case []any:
		return item
	case []string:
		result := make([]any, len(item))
		for index, nested := range item {
			result[index] = nested
		}
		return result
	default:
		return nil
	}
}

func stringSlice(value any) []string {
	result := []string{}
	for _, nested := range anySlice(value) {
		result = append(result, stringValue(nested))
	}
	return result
}

func stringValue(value any) string {
	switch item := value.(type) {
	case string:
		return item
	case fmt.Stringer:
		return item.String()
	case nil:
		return ""
	default:
		return fmt.Sprint(item)
	}
}

func intValue(value any) int {
	switch item := value.(type) {
	case int:
		return item
	default:
		return 0
	}
}
