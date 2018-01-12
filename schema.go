// Package jsonschema implements draft-handrews-json-schema-validation-00
// JSON Schema (application/schema+json) has several purposes, one of which is JSON instance validation.
// This document specifies a vocabulary for JSON Schema to describe the meaning of JSON documents,
// provide hints for user interfaces working with JSON data,
// and to make assertions about what a valid document must look like.
package jsonschema

import (
	"encoding/json"
)

// Validator is an interface for anything that can validate
type Validator interface {
	Validate(data interface{}) error
}

// RootSchema is a top-level Schema.
type RootSchema struct {
	Schema
	// The "$schema" keyword is both used as a JSON Schema version identifier and the location of a
	// resource which is itself a JSON Schema, which describes any schema written for this particular version.
	// The value of this keyword MUST be a URI [RFC3986] (containing a scheme) and this URI MUST be normalized. The current schema MUST be valid against the meta-schema identified by this URI.
	// If this URI identifies a retrievable resource, that resource SHOULD be of media type "application/schema+json".
	// The "$schema" keyword SHOULD be used in a root schema. It MUST NOT appear in subschemas.
	// [CREF2]
	// Values for this property are defined in other documents and by other parties. JSON Schema implementations SHOULD implement support for current and previous published drafts of JSON Schema vocabularies as deemed reasonable.
	SchemaURI string `json:"$schema"`
}

// Schema is the root JSON-schema struct
// A JSON Schema vocabulary is a set of keywords defined for a particular purpose.
// The vocabulary specifies the meaning of its keywords as assertions, annotations,
// and/or any vocabulary-defined keyword category.
//
// The two companion standards to this document each define a vocabulary:
// One for instance validation, and one for hypermedia annotations.
//
// Vocabularies are the primary mechanism for extensibility within the JSON Schema media type.
// Vocabularies may be defined by any entity.
// Vocabulary authors SHOULD take care to avoid keyword name collisions if the vocabulary is
// intended for broad use, and potentially combined with other vocabularies.
// JSON Schema does not provide any formal namespacing system,
// but also does not constrain keyword names, allowing for any number of namespacing approaches.
//
// Vocabularies may build on each other, such as by defining the behavior of their keywords
// with respect to the behavior of keywords from another vocabulary,  or by using a keyword
// from another vocabulary with a restricted or expanded set of acceptable values.
// Not all such vocabulary re-use will result in a new vocabulary that is compatible with the
// vocabulary on which it is built.
//
// Vocabulary authors SHOULD clearly document what level of compatibility, if any, is expected.
// A schema that itself describes a schema is called a meta-schema.
// Meta-schemas are used to validate JSON Schemas and specify which vocabulary it is using. [CREF1]
// A JSON Schema MUST be an object or a boolean.
type Schema struct {
	// The "$id" keyword defines a URI for the schema,
	// and the base URI that other URI references within the schema are resolved against.
	// A subschema's "$id" is resolved against the base URI of its parent schema.
	// If no parent sets an explicit base with "$id", the base URI is that of the entire document,
	// as determined per RFC 3986 section 5 [RFC3986].
	ID string `json:"$id,omitempty"`
	// Title and description can be used to decorate a user interface with information about
	// the data produced by this user interface.
	// A title will preferably be short.
	Title string `json:"title,omitempty"`
	// Description provides an explanation about the purpose
	// of the instance described by this schema.
	Description string `json:"description,omitempty"`
	// There are no restrictions placed on the value of this keyword.
	// When multiple occurrences of this keyword are applicable to a single sub-instance,
	// implementations SHOULD remove duplicates.
	// This keyword can be used to supply a default JSON value associated with a particular schema.
	// It is RECOMMENDED that a default value be valid against the associated schema.
	Default interface{} `json:"default,omitempty"`
	// The value of this keyword MUST be an array. There are no restrictions placed on the values
	// within the array.
	// When multiple occurrences of this keyword are applicable to a single sub-instance,
	// implementations MUST provide a flat array of all values rather than an array of arrays.
	// This keyword can be used to provide sample JSON values associated with a particular schema,
	// for the purpose of illustrating usage.
	// It is RECOMMENDED that these values be valid against the associated schema.
	// Implementations MAY use the value(s) of "default", if present, as an additional example.
	// If "examples" is absent, "default" MAY still be used in this manner.
	Examples []interface{} `json:"examples,omitempty"`
	// If "readOnly" has a value of boolean true, it indicates that the value of the instance is managed
	// exclusively by the owning authority, and attempts by an application to modify the value of this
	// property are expected to be ignored or rejected by that owning authority.
	// An instance document that is marked as "readOnly for the entire document MAY be ignored if sent
	// to the owning authority, or MAY result in an error, at the authority's discretion.
	// For example, "readOnly" would be used to mark a database-generated serial number as read-only, while "writeOnly" would be used to mark a password input field.
	// These keywords can be used to assist in user interface instance generation.
	// In particular, an application MAY choose to use a widget that hides input values as they are typed for write-only fields.
	// Omitting these keywords has the same behavior as values of false.
	ReadOnly bool `json:"readOnly,omitempty"`
	// If "writeOnly" has a value of boolean true, it indicates that the value is never present when the
	// instance is retrieved from the owning authority.
	// It can be present when sent to the owning authority to update or create the document
	// (or the resource it represents), but it will not be included in any updated or newly created
	// version of the instance.
	// An instance document that is marked as "writeOnly" for the entire document MAY be returned as a
	// blank document of some sort, or MAY produce an error upon retrieval, or have the retrieval request
	// ignored, at the authority's discretion.
	WriteOnly bool `json:"writeOnly,omitempty"`
	// This keyword is reserved for comments from schema authors to readers or maintainers of the schema.
	// The value of this keyword MUST be a string.
	// Implementations MUST NOT present this string to end users.
	// Tools for editing schemas SHOULD support displaying and editing this keyword.
	// The value of this keyword MAY be used in debug or error output which is intended for
	// developers making use of schemas.
	// Schema vocabularies SHOULD allow "$comment" within any object containing vocabulary keywords.
	// Implementations MAY assume "$comment" is allowed unless the vocabulary specifically forbids it.
	// Vocabularies MUST NOT specify any effect of "$comment" beyond what is described in this specification.
	Comment string `json:"comment,omitempty"`
	// Ref is used to reference a schema, and provides the ability to validate recursive
	// structures through self-reference.
	// An object schema with a "$ref" property MUST be interpreted as a "$ref" reference.
	// The value of the "$ref" property MUST be a URI Reference. Resolved against the current URI base,
	// it identifies the URI of a schema to use. All other properties in a "$ref" object MUST be ignored.
	// The URI is not a network locator, only an identifier. A schema need not be downloadable from the
	// address if it is a network-addressable URL, and implementations SHOULD NOT assume they should
	// perform a network operation when they encounter a network-addressable URI.
	// A schema MUST NOT be run into an infinite loop against a schema. For example, if two schemas
	// "#alice" and "#bob" both have an "allOf" property that refers to the other, a naive validator
	// might get stuck in an infinite recursive loop trying to validate the instance.
	// Schemas SHOULD NOT make use of infinite recursive nesting like this; the behavior is undefined.
	Ref string `json:"$ref,omitempty"`

	// Definitions provides a standardized location for schema authors to inline re-usable JSON Schemas
	// into a more general schema. The keyword does not directly affect the validation result.
	Definitions map[string]*Schema `json:"definitions,omitempty"`

	Type  Type  `json:"type,omitempty"`
	Enum  Enum  `json:"enum,omitempty"`
	Const Const `json:"const,omitempty"`

	MultipleOf       *MultipleOf       `json:"multipleOf,omitempty"`
	Maximum          *Maximum          `json:"maximum,omitempty"`
	ExclusiveMaximum *ExclusiveMaximum `json:"exclusiveMaximum,omitempty"`
	Minimum          *Minimum          `json:"minimum,omitempty"`
	ExclusiveMinimum *ExclusiveMinimum `json:"exclusiveMinimum,omitempty"`

	MaxLength *MaxLength `json:"maxLength,omitempty"`
	MinLength *MinLength `json:"minLength,omitempty"`
	Pattern   *Pattern   `json:"pattern,omitempty"`

	AllOf AllOf `json:"allOf,omitempty"`
	AnyOf AnyOf `json:"anyOf,omitempty"`
	OneOf OneOf `json:"oneOf,omitempty"`
	Not   *Not  `json:"not,omitempty"`

	Items           *Items           `json:"items,omitempty"`
	AdditionalItems *AdditionalItems `json:"additionalItems,omitempty"`
	MaxItems        *MaxItems        `json:"maxItems,omitempty"`
	MinItems        *MinItems        `json:"minItems,omitempty"`
	UniqueItems     *UniqueItems     `json:"uniqueItems,omitempty"`
	Contains        *Contains        `json:"contains,omitempty"`

	MaxProperties        *MaxProperties        `json:"maxProperties,omitempty"`
	MinProperties        *MinProperties        `json:"minProperties,omitempty"`
	Required             Required              `json:"required,omitempty"`
	Properties           Properties            `json:"properties,omitempty"`
	PatternProperties    PatternProperties     `json:"patternProperties,omitempty"`
	AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitempty"`
	Dependencies         *Dependencies         `json:"dependencies,omitempty"`
	PropertyNames        *PropertyNames        `json:"propertyNames,omitempty"`

	If   *If   `json:"if,omitempty"`
	Then *Then `json:"then,omitempty"`
	Else *Else `json:"else,omitempty"`
}

// _schema is an internal struct for encoding & decoding purposes
type _schema Schema

// UnmarshalJSON implements the json.Unmarshaler interface for Schema
func (s *Schema) UnmarshalJSON(data []byte) error {
	// support simple true false schemas that always pass or fail
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		if b {
			// boolean true Always passes validation, as if the empty schema {}
			*s = Schema{}
			return nil
		}
		// boolean false Always fails validation, as if the schema { "not":{} }
		*s = Schema{Not: &Not{}}
		return nil
	}

	sch := &_schema{}
	if err := json.Unmarshal(data, sch); err != nil {
		return err
	}

	if sch.Items != nil && sch.AdditionalItems != nil && !sch.Items.single {
		sch.AdditionalItems.startIndex = len(sch.Items.Schemas)
	}

	if sch.Properties != nil && sch.AdditionalProperties != nil {
		sch.AdditionalProperties.properties = sch.Properties
	}

	if sch.PatternProperties != nil && sch.AdditionalProperties != nil {
		sch.AdditionalProperties.patterns = sch.PatternProperties
	}

	*s = Schema(*sch)
	return nil
}

// Validators returns a schemas non-nil validators as a slice
func (s *Schema) Validators() (vs []Validator) {
	if s.Type != nil {
		vs = append(vs, s.Type)
	}
	if s.Const != nil {
		vs = append(vs, s.Const)
	}
	if s.Enum != nil {
		vs = append(vs, s.Enum)
	}

	if s.MultipleOf != nil {
		vs = append(vs, s.MultipleOf)
	}
	if s.Maximum != nil {
		vs = append(vs, s.Maximum)
	}
	if s.ExclusiveMaximum != nil {
		vs = append(vs, s.ExclusiveMaximum)
	}
	if s.Minimum != nil {
		vs = append(vs, s.Minimum)
	}
	if s.ExclusiveMinimum != nil {
		vs = append(vs, s.ExclusiveMinimum)
	}

	if s.MaxLength != nil {
		vs = append(vs, s.MaxLength)
	}
	if s.MinLength != nil {
		vs = append(vs, s.MinLength)
	}
	if s.Pattern != nil {
		vs = append(vs, s.Pattern)
	}

	if s.AllOf != nil {
		vs = append(vs, s.AllOf)
	}
	if s.AnyOf != nil {
		vs = append(vs, s.AnyOf)
	}
	if s.OneOf != nil {
		vs = append(vs, s.OneOf)
	}
	if s.Not != nil {
		vs = append(vs, s.Not)
	}

	if s.Items != nil {
		vs = append(vs, s.Items)
	}
	if s.AdditionalItems != nil {
		vs = append(vs, s.AdditionalItems)
	}
	if s.MaxItems != nil {
		vs = append(vs, s.MaxItems)
	}
	if s.MinItems != nil {
		vs = append(vs, s.MinItems)
	}
	if s.UniqueItems != nil {
		vs = append(vs, s.UniqueItems)
	}
	if s.Contains != nil {
		vs = append(vs, s.Contains)
	}

	if s.MaxProperties != nil {
		vs = append(vs, s.MaxProperties)
	}
	if s.MinProperties != nil {
		vs = append(vs, s.MinProperties)
	}
	if s.Required != nil {
		vs = append(vs, s.Required)
	}
	if s.Properties != nil {
		vs = append(vs, s.Properties)
	}
	if s.PatternProperties != nil {
		vs = append(vs, s.PatternProperties)
	}
	if s.AdditionalProperties != nil {
		vs = append(vs, s.AdditionalProperties)
	}
	if s.Dependencies != nil {
		vs = append(vs, s.Dependencies)
	}
	if s.PropertyNames != nil {
		vs = append(vs, s.PropertyNames)
	}

	if s.If != nil {
		vs = append(vs, s.If)
	}
	if s.Then != nil {
		vs = append(vs, s.Then)
	}
	if s.Else != nil {
		vs = append(vs, s.Else)
	}

	return
}

// Validate uses the schema to check an instance, returning error on the first error
func (s *Schema) Validate(data interface{}) error {
	for _, v := range s.Validators() {
		if err := v.Validate(data); err != nil {
			return err
		}
	}
	return nil
}

// DataType gives the primitive json type of a value, plus the special case
// "integer" for when numbers are whole
func DataType(data interface{}) string {
	switch v := data.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case float64:
		if float64(int(v)) == v {
			return "integer"
		}
		return "number"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}
