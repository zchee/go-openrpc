package jsonschema

// Schema is a JSON-Schema following Specification Draft 4 (http://json-schema.org/).
type Schema struct {
	ID                   string
	Schema               URL
	Ref                  *string
	Description          string
	Type                 string
	Nullable             bool
	Format               string
	Title                string
	Default              *JSON
	Maximum              *float64
	ExclusiveMaximum     bool
	Minimum              *float64
	ExclusiveMinimum     bool
	MaxLength            *int64
	MinLength            *int64
	Pattern              string
	MaxItems             *int64
	MinItems             *int64
	UniqueItems          bool
	MultipleOf           *float64
	Enum                 []JSON
	MaxProperties        *int64
	MinProperties        *int64
	Required             []string
	Items                *PropsOrArray
	AllOf                []Schema
	OneOf                []Schema
	AnyOf                []Schema
	Not                  *Schema
	Properties           map[string]Schema
	AdditionalProperties *PropsOrBool
	PatternProperties    map[string]Schema
	Dependencies         Dependencies
	AdditionalItems      *PropsOrBool
	Definitions          Definitions
	ExternalDocs         *ExternalDocumentation
	Example              *JSON
}

// JSON represents any valid JSON value.
// These types are supported: bool, int64, float64, string, []interface{}, map[string]interface{} and nil.
type JSON interface{}

// URL represents a schema url.
type URL string

// PropsOrArray represents a value that can either be a JSONSchemaProps
// or an array of JSONSchemaProps. Mainly here for serialization purposes.
type PropsOrArray struct {
	Schema      *Schema
	JSONSchemas []Schema
}

// PropsOrBool represents JSONSchemaProps or a boolean value.
// Defaults to true for the boolean property.
type PropsOrBool struct {
	Allows bool
	Schema *Schema
}

// Dependencies represent a dependencies property.
type Dependencies map[string]PropsOrStringArray

// PropsOrStringArray represents a JSONSchemaProps or a string array.
type PropsOrStringArray struct {
	Schema   *Schema
	Property []string
}

// Definitions contains the models explicitly defined in this spec.
type Definitions map[string]Schema

// ExternalDocumentation allows referencing an external resource for extended documentation.
type ExternalDocumentation struct {
	Description string
	URL         string
}
