package release

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type PreparedUpgrade struct {
	Catalog LoadedCatalog
	Asset   CatalogAsset
	Pack    VerifiedPack
	WorkDir string
}

func PlatformForRuntime(goos, goarch string) (string, error) {
	switch {
	case goos == "darwin" && (goarch == "arm64" || goarch == "amd64"):
		return "macos-universal", nil
	case goos == "windows" && goarch == "amd64":
		return "windows-x64", nil
	default:
		return "", fmt.Errorf("unsupported upgrade platform %s/%s", goos, goarch)
	}
}

func PrepareUpgrade(catalogLocation, version, workDir string, smokeTest bool) (PreparedUpgrade, error) {
	platform, err := PlatformForRuntime(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return PreparedUpgrade{}, err
	}
	catalog, err := LoadCatalog(catalogLocation)
	if err != nil {
		return PreparedUpgrade{}, err
	}
	asset, err := catalog.Select(version, platform)
	if err != nil {
		return PreparedUpgrade{}, err
	}
	archive := filepath.Join(workDir, asset.Name)
	if err := FetchAsset(catalog, asset, archive); err != nil {
		return PreparedUpgrade{}, err
	}
	extract := filepath.Join(workDir, "extract")
	pack, err := VerifyPack(archive, extract, asset, VerifyOptions{Version: version, Platform: platform, Command: "codeheart-operating-kit", SmokeTest: smokeTest})
	if err != nil {
		return PreparedUpgrade{}, err
	}
	return PreparedUpgrade{Catalog: catalog, Asset: asset, Pack: pack, WorkDir: workDir}, nil
}

func RequireForwardUpgrade(current, target string) error {
	comparison, err := CompareVersions(target, current)
	if err != nil {
		return err
	}
	if comparison <= 0 {
		return fmt.Errorf("upgrade target %s must be newer than installed version %s", target, current)
	}
	return nil
}

func CompareVersions(left, right string) (int, error) {
	parse := func(value string) ([]int, error) {
		value = strings.TrimPrefix(value, "v")
		parts := strings.Split(value, ".")
		if len(parts) != 3 {
			return nil, fmt.Errorf("version %q is not semantic major.minor.patch", value)
		}
		result := make([]int, 3)
		for index, part := range parts {
			parsed, err := strconv.Atoi(part)
			if err != nil || parsed < 0 {
				return nil, fmt.Errorf("version %q is invalid", value)
			}
			result[index] = parsed
		}
		return result, nil
	}
	leftParts, err := parse(left)
	if err != nil {
		return 0, err
	}
	rightParts, err := parse(right)
	if err != nil {
		return 0, err
	}
	for index := range leftParts {
		if leftParts[index] > rightParts[index] {
			return 1, nil
		}
		if leftParts[index] < rightParts[index] {
			return -1, nil
		}
	}
	return 0, nil
}

func DefaultInstalledBinary() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(path)
}
