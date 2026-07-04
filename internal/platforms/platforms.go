package platforms

import (
	"runtime"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/manifest"
)

func CurrentPlatform() string {
	return PlatformFor(runtime.GOOS)
}

func PlatformFor(goos string) string {
	switch goos {
	case "darwin":
		return "macos"
	case "windows":
		return "windows"
	default:
		return goos
	}
}

func IsSupportedG1Platform(goos string) bool {
	switch PlatformFor(goos) {
	case "macos", "windows":
		return true
	default:
		return false
	}
}

func CandidateAssetPlatforms(goos string, goarch string) []string {
	switch PlatformFor(goos) {
	case "macos":
		return []string{"macos-universal", "macos", "universal"}
	case "windows":
		if goarch == "amd64" || goarch == "x64" {
			return []string{"windows-x64", "windows", "universal"}
		}
		return []string{"windows", "universal"}
	default:
		return []string{PlatformFor(goos), "universal"}
	}
}

func SelectAssets(assets []manifest.ReleaseAsset, goos string, goarch string) []manifest.ReleaseAsset {
	allowed := map[string]bool{}
	for _, platform := range CandidateAssetPlatforms(goos, goarch) {
		allowed[platform] = true
	}
	selected := []manifest.ReleaseAsset{}
	for _, asset := range assets {
		if allowed[asset.Platform] {
			selected = append(selected, asset)
		}
	}
	return selected
}
