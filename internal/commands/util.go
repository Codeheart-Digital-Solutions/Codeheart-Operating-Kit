package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var purposeLabels = map[string]string{
	"private-automation": "Personal automation",
	"company-automation": "Company operations",
	"software-product":   "Software or product development",
}

func expandPath(value string) string {
	if value == "" {
		return value
	}
	if value == "~" || strings.HasPrefix(value, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			if value == "~" {
				return home
			}
			return filepath.Join(home, value[2:])
		}
	}
	return value
}

func writeJSON(stdout io.Writer, value any) error {
	encoder := json.NewEncoder(stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func parseValueArgs(args []string, specs map[string]bool) (map[string]string, map[string]bool, []string, error) {
	values := map[string]string{}
	bools := map[string]bool{}
	positionals := []string{}
	for index := 0; index < len(args); index++ {
		arg := args[index]
		if !strings.HasPrefix(arg, "--") {
			positionals = append(positionals, arg)
			continue
		}
		flagName := arg
		inlineValue := ""
		hasInlineValue := false
		if name, value, ok := strings.Cut(arg, "="); ok {
			flagName = name
			inlineValue = value
			hasInlineValue = true
		}
		requiresValue, known := specs[flagName]
		if !known {
			return nil, nil, nil, fmt.Errorf("unknown option %s", flagName)
		}
		if !requiresValue {
			if hasInlineValue {
				return nil, nil, nil, fmt.Errorf("option %s does not accept a value", flagName)
			}
			bools[flagName] = true
			continue
		}
		if hasInlineValue {
			values[flagName] = inlineValue
			continue
		}
		if index+1 >= len(args) {
			return nil, nil, nil, fmt.Errorf("option %s requires a value", flagName)
		}
		index++
		if strings.HasPrefix(args[index], "-") {
			return nil, nil, nil, fmt.Errorf("option %s requires a value", flagName)
		}
		values[flagName] = args[index]
	}
	return values, bools, positionals, nil
}

func validatePurpose(purpose string) error {
	if purpose == "" {
		return nil
	}
	if _, ok := purposeLabels[purpose]; ok {
		return nil
	}
	return fmt.Errorf("invalid purpose %q", purpose)
}

func writeArgError(stderr io.Writer, command string, err error) int {
	fmt.Fprintf(stderr, "codeheart-operating-kit %s: error: %s\n", command, err)
	return 2
}

func osMkdirAll(path string) error {
	return os.MkdirAll(path, 0o755)
}

func joinRoot(root string, relative string) string {
	return filepath.Join(root, filepath.FromSlash(relative))
}

func localFileURLPath(parsed *url.URL) string {
	return localFileURLPathForGOOS(parsed, runtime.GOOS)
}

func localFileURLPathForGOOS(parsed *url.URL, goos string) string {
	path, err := url.PathUnescape(parsed.Path)
	if err != nil {
		path = parsed.Path
	}
	if goos != "windows" {
		return path
	}
	if parsed.Host != "" && parsed.Host != "localhost" {
		return `\\` + parsed.Host + strings.ReplaceAll(path, "/", `\`)
	}
	if len(path) >= 3 && path[0] == '/' && path[2] == ':' {
		path = path[1:]
	}
	return strings.ReplaceAll(path, "/", `\`)
}

func mapsToAny(values []map[string]any) []any {
	result := make([]any, len(values))
	for index, value := range values {
		result[index] = value
	}
	return result
}
