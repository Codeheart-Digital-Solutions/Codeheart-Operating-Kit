package yamlmini

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type sourceLine struct {
	indent  int
	content string
}

// Parse reads the small YAML subset currently used by the Operating Kit.
func Parse(text string) (any, error) {
	lines := logicalLines(text)
	if len(lines) == 0 {
		return map[string]any{}, nil
	}
	parsed, next, err := parseBlock(lines, 0, 0)
	if err != nil {
		return nil, err
	}
	if next != len(lines) {
		return nil, fmt.Errorf("unsupported YAML structure near %q", lines[next].content)
	}
	return parsed, nil
}

// MustMap parses YAML and returns the root mapping.
func MustMap(text string) (map[string]any, error) {
	value, err := Parse(text)
	if err != nil {
		return nil, err
	}
	mapping, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected YAML mapping root")
	}
	return mapping, nil
}

// Dump writes the supported YAML subset in deterministic key order.
func Dump(value any) string {
	var lines []string
	emit(&lines, normalize(value), 0, "")
	return strings.Join(lines, "\n") + "\n"
}

// EqualNormalized compares values after parse/dump normalization.
func EqualNormalized(left, right any) bool {
	leftParsed, leftErr := Parse(Dump(left))
	rightParsed, rightErr := Parse(Dump(right))
	if leftErr != nil || rightErr != nil {
		return false
	}
	return reflect.DeepEqual(leftParsed, rightParsed)
}

func logicalLines(text string) []sourceLine {
	result := []sourceLine{}
	for _, raw := range strings.Split(text, "\n") {
		line := stripComments(raw)
		if strings.TrimSpace(line) == "" {
			continue
		}
		trimmed := strings.TrimLeft(line, " ")
		result = append(result, sourceLine{
			indent:  len(line) - len(trimmed),
			content: trimmed,
		})
	}
	return result
}

func stripComments(line string) string {
	inQuote := false
	var quote rune
	var previous rune
	for index, char := range line {
		if (char == '\'' || char == '"') && previous != '\\' {
			if inQuote && char == quote {
				inQuote = false
				quote = 0
			} else if !inQuote {
				inQuote = true
				quote = char
			}
		}
		if char == '#' && !inQuote {
			return strings.TrimRight(line[:index], " \t\r")
		}
		previous = char
	}
	return strings.TrimRight(line, " \t\r")
}

func parseBlock(lines []sourceLine, index int, indent int) (any, int, error) {
	if index >= len(lines) {
		return map[string]any{}, index, nil
	}
	current := lines[index]
	if current.indent < indent {
		return map[string]any{}, index, nil
	}
	if isListItem(current.content) {
		return parseList(lines, index, current.indent)
	}
	return parseMap(lines, index, current.indent)
}

func parseList(lines []sourceLine, index int, indent int) ([]any, int, error) {
	result := []any{}
	for index < len(lines) {
		current := lines[index]
		if current.indent != indent || !isListItem(current.content) {
			break
		}
		item := strings.TrimSpace(current.content[1:])
		index++
		switch {
		case item == "":
			child, next, err := parseBlock(lines, index, indent+2)
			if err != nil {
				return nil, index, err
			}
			result = append(result, child)
			index = next
		case strings.Contains(item, ":"):
			key, value, _ := strings.Cut(item, ":")
			mapping := map[string]any{}
			if strings.TrimSpace(value) != "" {
				mapping[strings.TrimSpace(key)] = parseScalar(value)
			} else {
				child, next, err := parseBlock(lines, index, indent+2)
				if err != nil {
					return nil, index, err
				}
				mapping[strings.TrimSpace(key)] = child
				index = next
			}
			for index < len(lines) {
				next := lines[index]
				if next.indent <= indent {
					break
				}
				if next.indent == indent+2 && !isListItem(next.content) {
					subkey, subvalue, ok := strings.Cut(next.content, ":")
					if !ok {
						return nil, index, fmt.Errorf("unsupported YAML mapping item %q", next.content)
					}
					index++
					if strings.TrimSpace(subvalue) != "" {
						mapping[strings.TrimSpace(subkey)] = parseScalar(subvalue)
					} else {
						child, nextIndex, err := parseBlock(lines, index, next.indent+2)
						if err != nil {
							return nil, index, err
						}
						mapping[strings.TrimSpace(subkey)] = child
						index = nextIndex
					}
					continue
				}
				break
			}
			result = append(result, mapping)
		default:
			result = append(result, parseScalar(item))
		}
	}
	return result, index, nil
}

func parseMap(lines []sourceLine, index int, indent int) (map[string]any, int, error) {
	result := map[string]any{}
	for index < len(lines) {
		current := lines[index]
		if current.indent != indent || isListItem(current.content) {
			break
		}
		key, value, ok := strings.Cut(current.content, ":")
		if !ok {
			return nil, index, fmt.Errorf("unsupported YAML mapping item %q", current.content)
		}
		index++
		if strings.TrimSpace(value) != "" {
			result[strings.TrimSpace(key)] = parseScalar(value)
			continue
		}
		child, next, err := parseBlock(lines, index, indent+2)
		if err != nil {
			return nil, index, err
		}
		result[strings.TrimSpace(key)] = child
		index = next
	}
	return result, index, nil
}

func isListItem(content string) bool {
	return content == "-" || strings.HasPrefix(content, "- ")
}

func parseScalar(value string) any {
	trimmed := strings.TrimSpace(value)
	switch trimmed {
	case "":
		return ""
	case "[]":
		return []any{}
	case "{}":
		return map[string]any{}
	case "true":
		return true
	case "false":
		return false
	}
	if len(trimmed) >= 2 {
		first := trimmed[0]
		last := trimmed[len(trimmed)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			if unquoted, err := strconv.Unquote(trimmed); err == nil {
				return unquoted
			}
			return trimmed[1 : len(trimmed)-1]
		}
	}
	if isDigits(trimmed) {
		if parsed, err := strconv.Atoi(trimmed); err == nil {
			return parsed
		}
	}
	return trimmed
}

func isDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func emit(lines *[]string, value any, level int, keyPrefix string) {
	pad := strings.Repeat(" ", level)
	switch item := value.(type) {
	case map[string]any:
		if keyPrefix != "" {
			*lines = append(*lines, fmt.Sprintf("%s%s:", pad, keyPrefix))
			level += 2
			pad = strings.Repeat(" ", level)
		}
		keys := make([]string, 0, len(item))
		for key := range item {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			nested := normalize(item[key])
			switch nested.(type) {
			case map[string]any, []any:
				*lines = append(*lines, fmt.Sprintf("%s%s:", pad, key))
				emit(lines, nested, level+2, "")
			default:
				*lines = append(*lines, fmt.Sprintf("%s%s: %s", pad, key, formatScalar(nested)))
			}
		}
	case []any:
		if keyPrefix != "" {
			*lines = append(*lines, fmt.Sprintf("%s%s:", pad, keyPrefix))
			level += 2
			pad = strings.Repeat(" ", level)
		}
		for _, nested := range item {
			nested = normalize(nested)
			switch nested.(type) {
			case map[string]any, []any:
				*lines = append(*lines, fmt.Sprintf("%s-", pad))
				emit(lines, nested, level+2, "")
			default:
				*lines = append(*lines, fmt.Sprintf("%s- %s", pad, formatScalar(nested)))
			}
		}
	default:
		if keyPrefix != "" {
			*lines = append(*lines, fmt.Sprintf("%s%s: %s", pad, keyPrefix, formatScalar(item)))
		} else {
			*lines = append(*lines, fmt.Sprintf("%s%s", pad, formatScalar(item)))
		}
	}
}

func normalize(value any) any {
	switch item := value.(type) {
	case map[string]string:
		result := map[string]any{}
		for key, nested := range item {
			result[key] = nested
		}
		return result
	case map[string][]string:
		result := map[string]any{}
		for key, nested := range item {
			result[key] = normalize(nested)
		}
		return result
	case []string:
		result := make([]any, len(item))
		for index, nested := range item {
			result[index] = nested
		}
		return result
	case []map[string]any:
		result := make([]any, len(item))
		for index, nested := range item {
			result[index] = nested
		}
		return result
	default:
		return value
	}
}

func formatScalar(value any) string {
	switch item := value.(type) {
	case bool:
		if item {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(item)
	case int64:
		return strconv.FormatInt(item, 10)
	case nil:
		return "null"
	default:
		return quoteIfNeeded(fmt.Sprint(item))
	}
}

func quoteIfNeeded(value string) string {
	if value == "" || strings.TrimSpace(value) != value || strings.ContainsAny(value, ":#[]{}") || value == "true" || value == "false" || value == "null" || strings.Contains(value, "\n") {
		encoded, err := json.Marshal(value)
		if err != nil {
			return strconv.Quote(value)
		}
		return string(encoded)
	}
	return value
}
