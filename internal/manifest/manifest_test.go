package manifest

import (
	"testing"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

func TestLoadStandardProfile(t *testing.T) {
	profile, err := LoadProfile("standard")
	if err != nil {
		t.Fatalf("LoadProfile: %v", err)
	}
	if profile.ID != "standard" {
		t.Fatalf("profile ID = %q, want standard", profile.ID)
	}
	if !contains(profile.SelectedComponents, "agent-interface") {
		t.Fatalf("profile selected components missing agent-interface: %#v", profile.SelectedComponents)
	}
	foundLockSurface := false
	for _, surface := range profile.GeneratedSurfaces {
		if surface.Path == ".codeheart/kit.lock.yaml" && surface.Ownership == "generated-surface" {
			foundLockSurface = true
		}
	}
	if !foundLockSurface {
		t.Fatalf("profile generated surfaces missing kit lockfile surface")
	}
}

func TestLoadSelectedComponentsAndFiles(t *testing.T) {
	components, err := LoadComponents("standard")
	if err != nil {
		t.Fatalf("LoadComponents: %v", err)
	}
	if len(components) == 0 {
		t.Fatalf("expected selected components")
	}
	foundPlanning := false
	for _, component := range components {
		if component.ID == "planning-workflows" {
			foundPlanning = true
			if component.ManifestPath != "components/planning-workflows/component.yaml" {
				t.Fatalf("planning manifest path = %q", component.ManifestPath)
			}
		}
	}
	if !foundPlanning {
		t.Fatalf("selected components missing planning-workflows: %#v", components)
	}

	files, err := IterComponentFiles("standard")
	if err != nil {
		t.Fatalf("IterComponentFiles: %v", err)
	}
	if !hasFile(files, ".codeheart/kit/docs/planning-workflows/README.md", "managed") {
		t.Fatalf("component files missing planning managed README")
	}
	if !hasFile(files, "docs/repo/plans/plan-register.md", "scaffold") {
		t.Fatalf("component files missing plan-register scaffold")
	}
}

func TestLoadEmbeddedContentManifest(t *testing.T) {
	releaseManifest, err := LoadReleaseManifest()
	if err != nil {
		t.Fatalf("LoadReleaseManifest: %v", err)
	}
	if releaseManifest.Version != "0.1.21" {
		t.Fatalf("release version = %q, want 0.1.21", releaseManifest.Version)
	}
	if len(releaseManifest.Assets) != 0 {
		t.Fatalf("embedded content identity must not contain release assets: %#v", releaseManifest.Assets)
	}
	compatibility, ok := releaseManifest.Raw["compatibility"].(map[string]any)
	if !ok || compatibility["commands"] == nil {
		t.Fatalf("content manifest missing compatibility: %#v", releaseManifest.Raw)
	}
	graph, err := state.CompileGraph("standard")
	if err != nil {
		t.Fatal(err)
	}
	profiles := releaseManifest.Raw["profiles"].([]any)
	profile := profiles[0].(map[string]any)
	if profile["graph_sha256"] != graph.DigestSHA256 {
		t.Fatalf("content graph digest = %v, want %s", profile["graph_sha256"], graph.DigestSHA256)
	}
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func hasFile(files []ComponentFile, path string, ownership string) bool {
	for _, file := range files {
		if file.Target == path && file.Ownership == ownership {
			return true
		}
	}
	return false
}
