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
	ComponentSchema        = "schemas/component.schema.json"
	ProfileSchema          = "schemas/profile.schema.json"
	LockV1Schema           = "schemas/kit-lock-v1.schema.json"
	LockV2Schema           = "schemas/kit-lock.schema.json"
	ConfigV1Schema         = "schemas/kit-config-v1.schema.json"
	ConfigV2Schema         = "schemas/kit-config.schema.json"
	ConfigSchema           = ConfigV2Schema
	ContentManifestSchema  = "schemas/content-manifest.schema.json"
	ReleaseCatalogSchema   = "schemas/release-catalog.schema.json"
	PackManifestSchema     = "schemas/pack-manifest.schema.json"
	PlanMetadataSchema     = "schemas/plan-metadata.schema.json"
	PlanCatalogV1Schema    = "schemas/plan-catalog-v1.schema.json"
	PlanInventoryV3Schema  = "schemas/plan-inventory.schema.json"
	PlanInventorySchema    = PlanInventoryV3Schema
	PlanCatalogV2Schema    = "schemas/plan-catalog-v2.schema.json"
	PlanCatalogV3Schema    = "schemas/plan-catalog.schema.json"
	PlanCatalogSchema      = PlanCatalogV3Schema
	PlanMigrationV1Schema  = "schemas/plan-migration-ledger-v1.schema.json"
	PlanMigrationV2Schema  = "schemas/plan-migration-ledger-v2.schema.json"
	PlanMigrationV3Schema  = "schemas/plan-migration-ledger.schema.json"
	PlanMigrationSchema    = PlanMigrationV3Schema
	PortfolioSourcesSchema = "schemas/portfolio-local-sources.schema.json"
	PortfolioOverlaySchema = "schemas/portfolio-strategic-overlay.schema.json"
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

// DecodeAndValidateConfigYAML decodes a config, selects its schema from the raw
// schema_version, and only then validates the complete mapping. Callers must not
// decode into a version-specific struct before this check.
func DecodeAndValidateConfigYAML(data []byte) (map[string]any, error) {
	value, err := DecodeYAMLMap(data)
	if err != nil {
		return nil, err
	}
	if err := ValidateConfig(value); err != nil {
		return nil, err
	}
	return value, nil
}

// ValidateConfig validates a raw JSON-compatible config using its declared version.
func ValidateConfig(value map[string]any) error {
	version, err := declaredSchemaVersion(value, "config")
	if err != nil {
		return err
	}
	schemaPath, err := SchemaForConfigVersion(version)
	if err != nil {
		return err
	}
	return Validate(schemaPath, value)
}

func SchemaForConfigVersion(version int) (string, error) {
	switch version {
	case 1:
		return ConfigV1Schema, nil
	case 2:
		return ConfigV2Schema, nil
	default:
		return "", fmt.Errorf("unsupported config schema version %d", version)
	}
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

func SchemaForPlanCatalogVersion(version int) (string, error) {
	switch version {
	case 1:
		return PlanCatalogV1Schema, nil
	case 2:
		return PlanCatalogV2Schema, nil
	case 3:
		return PlanCatalogV3Schema, nil
	default:
		return "", fmt.Errorf("unsupported plan catalog schema version %d", version)
	}
}

func SchemaForPlanMigrationVersion(version int) (string, error) {
	switch version {
	case 1:
		return PlanMigrationV1Schema, nil
	case 2:
		return PlanMigrationV2Schema, nil
	case 3:
		return PlanMigrationV3Schema, nil
	default:
		return "", fmt.Errorf("unsupported plan migration schema version %d", version)
	}
}

func SchemaForPlanInventoryVersion(version int) (string, error) {
	if version == 3 {
		return PlanInventoryV3Schema, nil
	}
	return "", fmt.Errorf("unsupported plan inventory schema version %d", version)
}

func declaredSchemaVersion(value map[string]any, document string) (int, error) {
	raw, ok := value["schema_version"]
	if !ok {
		return 0, fmt.Errorf("%s schema_version is required", document)
	}
	version, ok := raw.(int)
	if !ok {
		return 0, fmt.Errorf("%s schema_version must be an integer", document)
	}
	return version, nil
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
