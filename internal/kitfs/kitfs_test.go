package kitfs

import (
	"strings"
	"testing"
)

func TestEmbeddedResourcesExposeStandardProfile(t *testing.T) {
	text, err := ReadText("profiles/standard.yaml")
	if err != nil {
		t.Fatalf("ReadText standard profile: %v", err)
	}
	if !strings.Contains(text, "selected_components:") {
		t.Fatalf("standard profile did not contain selected components:\n%s", text)
	}
}

func TestEmbeddedResourcesExposeManagedComponentFile(t *testing.T) {
	if !Exists("components/agent-interface/managed/README.md") {
		t.Fatalf("expected embedded agent-interface README to exist")
	}
}

func TestCleanRejectsTraversal(t *testing.T) {
	if _, err := Clean("../manifest.yaml"); err == nil {
		t.Fatalf("expected traversal path to be rejected")
	}
}
