package state

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/kitfs"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

const (
	ComponentSchema       = "schemas/component.schema.json"
	ProfileSchema         = "schemas/profile.schema.json"
	LockV1Schema          = "schemas/kit-lock-v1.schema.json"
	LockV2Schema          = "schemas/kit-lock.schema.json"
	ConfigV1Schema        = "schemas/kit-config.schema.json"
	ContentManifestSchema = "schemas/content-manifest.schema.json"
	ReleaseCatalogSchema  = "schemas/release-catalog.schema.json"
	PackManifestSchema    = "schemas/pack-manifest.schema.json"
)

var schemaCache = struct {
	sync.Mutex
	values map[string]*jsonschema.Schema
}{values: map[string]*jsonschema.Schema{}}

// Validate checks a JSON-compatible value against an embedded schema.
func Validate(schemaPath string, value any) error {
	schema, err := compiledSchema(schemaPath)
	if err != nil {
		return err
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal %s instance: %w", schemaPath, err)
	}
	instance, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode %s instance: %w", schemaPath, err)
	}
	if err := schema.Validate(instance); err != nil {
		return fmt.Errorf("validate against %s: %w", schemaPath, err)
	}
	return nil
}

// DecodeAndValidateYAML decodes YAML and validates it against an embedded schema.
func DecodeAndValidateYAML(schemaPath string, data []byte) (map[string]any, error) {
	value, err := DecodeYAMLMap(data)
	if err != nil {
		return nil, err
	}
	if err := Validate(schemaPath, value); err != nil {
		return nil, err
	}
	return value, nil
}

func SchemaForLockVersion(version int) (string, error) {
	switch version {
	case 1:
		return LockV1Schema, nil
	case 2:
		return LockV2Schema, nil
	default:
		return "", fmt.Errorf("unsupported lock schema version %d", version)
	}
}

func compiledSchema(schemaPath string) (*jsonschema.Schema, error) {
	schemaCache.Lock()
	defer schemaCache.Unlock()
	if schema := schemaCache.values[schemaPath]; schema != nil {
		return schema, nil
	}
	data, err := kitfs.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("read schema %s: %w", schemaPath, err)
	}
	document, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode schema %s: %w", schemaPath, err)
	}
	mapping, ok := document.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("schema %s is not an object", schemaPath)
	}
	id, _ := mapping["$id"].(string)
	if id == "" {
		return nil, fmt.Errorf("schema %s has no $id", schemaPath)
	}
	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft2020)
	compiler.AssertFormat()
	if err := compiler.AddResource(id, document); err != nil {
		return nil, fmt.Errorf("register schema %s: %w", schemaPath, err)
	}
	schema, err := compiler.Compile(id)
	if err != nil {
		return nil, fmt.Errorf("compile schema %s: %w", schemaPath, err)
	}
	schemaCache.values[schemaPath] = schema
	return schema, nil
}
