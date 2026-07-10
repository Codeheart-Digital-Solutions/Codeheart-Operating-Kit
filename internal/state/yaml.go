package state

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"go.yaml.in/yaml/v3"
)

// DecodeYAML decodes public Operating Kit YAML into JSON-compatible Go values.
func DecodeYAML(data []byte) (any, error) {
	var value any
	if err := yaml.Unmarshal(data, &value); err != nil {
		return nil, fmt.Errorf("decode YAML: %w", err)
	}
	return normalize(value)
}

// DecodeYAMLMap decodes a YAML mapping.
func DecodeYAMLMap(data []byte) (map[string]any, error) {
	value, err := DecodeYAML(data)
	if err != nil {
		return nil, err
	}
	mapping, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected YAML mapping root")
	}
	return mapping, nil
}

// EncodeYAML encodes a JSON-compatible value with deterministic mapping order.
func EncodeYAML(value any) ([]byte, error) {
	normalized, err := normalize(value)
	if err != nil {
		return nil, err
	}
	content, err := yamlNode(normalized)
	if err != nil {
		return nil, err
	}
	doc := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{content}}
	data, err := yaml.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("encode YAML: %w", err)
	}
	return data, nil
}

// DeepCopy returns a JSON-compatible independent copy.
func DeepCopy(value map[string]any) map[string]any {
	result, _ := normalize(value)
	return result.(map[string]any)
}

func normalize(value any) (any, error) {
	switch item := value.(type) {
	case map[string]any:
		result := make(map[string]any, len(item))
		for key, nested := range item {
			value, err := normalize(nested)
			if err != nil {
				return nil, err
			}
			result[key] = value
		}
		return result, nil
	case map[any]any:
		result := make(map[string]any, len(item))
		for key, nested := range item {
			text, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("YAML mapping key %v is not a string", key)
			}
			value, err := normalize(nested)
			if err != nil {
				return nil, err
			}
			result[text] = value
		}
		return result, nil
	case []any:
		result := make([]any, len(item))
		for index, nested := range item {
			value, err := normalize(nested)
			if err != nil {
				return nil, err
			}
			result[index] = value
		}
		return result, nil
	case time.Time:
		return item.UTC().Format(time.RFC3339), nil
	case nil, string, bool, int, float64:
		return item, nil
	case int8:
		return int(item), nil
	case int16:
		return int(item), nil
	case int32:
		return int(item), nil
	case int64:
		if strconv.IntSize == 32 && (item > math.MaxInt32 || item < math.MinInt32) {
			return item, nil
		}
		return int(item), nil
	case uint:
		if uint64(item) > uint64(math.MaxInt) {
			return nil, fmt.Errorf("YAML integer %d exceeds platform int", item)
		}
		return int(item), nil
	case uint8:
		return int(item), nil
	case uint16:
		return int(item), nil
	case uint32:
		if uint64(item) > uint64(math.MaxInt) {
			return nil, fmt.Errorf("YAML integer %d exceeds platform int", item)
		}
		return int(item), nil
	case uint64:
		if item > uint64(math.MaxInt) {
			return nil, fmt.Errorf("YAML integer %d exceeds platform int", item)
		}
		return int(item), nil
	case float32:
		return float64(item), nil
	case json.Number:
		if integer, err := item.Int64(); err == nil {
			return normalize(integer)
		}
		decimal, err := item.Float64()
		if err != nil {
			return nil, err
		}
		return decimal, nil
	default:
		data, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("normalize YAML value %T: %w", item, err)
		}
		var result any
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
		return result, nil
	}
}

func yamlNode(value any) (*yaml.Node, error) {
	switch item := value.(type) {
	case map[string]any:
		node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		keys := make([]string, 0, len(item))
		for key := range item {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			keyNode := &yaml.Node{}
			keyNode.SetString(key)
			valueNode, err := yamlNode(item[key])
			if err != nil {
				return nil, err
			}
			node.Content = append(node.Content, keyNode, valueNode)
		}
		return node, nil
	case []any:
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, nested := range item {
			valueNode, err := yamlNode(nested)
			if err != nil {
				return nil, err
			}
			node.Content = append(node.Content, valueNode)
		}
		return node, nil
	case string:
		node := &yaml.Node{}
		node.SetString(item)
		return node, nil
	case nil:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}, nil
	default:
		node := &yaml.Node{}
		if err := node.Encode(item); err != nil {
			return nil, fmt.Errorf("encode YAML scalar %T: %w", item, err)
		}
		return node, nil
	}
}
