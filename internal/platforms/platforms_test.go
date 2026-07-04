package platforms

import (
	"testing"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/manifest"
)

func TestCandidateAssetPlatforms(t *testing.T) {
	tests := []struct {
		name   string
		goos   string
		goarch string
		want   []string
	}{
		{name: "darwin", goos: "darwin", goarch: "arm64", want: []string{"macos-universal", "macos", "universal"}},
		{name: "windows-x64", goos: "windows", goarch: "amd64", want: []string{"windows-x64", "windows", "universal"}},
		{name: "linux", goos: "linux", goarch: "amd64", want: []string{"linux", "universal"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CandidateAssetPlatforms(test.goos, test.goarch)
			if !equalStrings(got, test.want) {
				t.Fatalf("CandidateAssetPlatforms = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestSelectAssetsUsesPlatformAndUniversalFallbacks(t *testing.T) {
	assets := []manifest.ReleaseAsset{
		{Name: "bootstrap.md", Platform: "universal"},
		{Name: "macos.zip", Platform: "macos"},
		{Name: "macos-universal.zip", Platform: "macos-universal"},
		{Name: "windows.zip", Platform: "windows"},
		{Name: "windows-x64.zip", Platform: "windows-x64"},
		{Name: "linux.zip", Platform: "linux"},
	}
	macAssets := SelectAssets(assets, "darwin", "arm64")
	if !hasAsset(macAssets, "bootstrap.md") || !hasAsset(macAssets, "macos.zip") || !hasAsset(macAssets, "macos-universal.zip") || hasAsset(macAssets, "windows.zip") {
		t.Fatalf("unexpected macOS assets: %#v", macAssets)
	}
	windowsAssets := SelectAssets(assets, "windows", "amd64")
	if !hasAsset(windowsAssets, "bootstrap.md") || !hasAsset(windowsAssets, "windows.zip") || !hasAsset(windowsAssets, "windows-x64.zip") || hasAsset(windowsAssets, "macos.zip") {
		t.Fatalf("unexpected Windows assets: %#v", windowsAssets)
	}
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func hasAsset(assets []manifest.ReleaseAsset, name string) bool {
	for _, asset := range assets {
		if asset.Name == name {
			return true
		}
	}
	return false
}
