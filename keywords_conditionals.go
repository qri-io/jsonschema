package jsonschema

import (
	"encoding/json"
)

// If MUST be a valid JSON Schema.
// Instances that successfully validate against this keyword's subschema MUST also be valid against the subschema value of the "then" keyword, if present.
// Instances that fail to validate against this keyword's subschema MUST also be valid against the subschema value of the "else" keyword.
// Validation of the instance against this keyword on its own always succeeds, regardless of the validation outcome of against its subschema.
type If Schema

// Validate implements the Validator interface for If
func (i *If) Validate(data interface{}) error {
	return nil
}

// JSONProp implements JSON property name indexing for If
func (i If) JSONProp(name string) interface{} {
	return Schema(i).JSONProp(name)
}

// UnmarshalJSON implements the json.Unmarshaler interface for If
func (i *If) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*i = If(sch)
	return nil
}

// Then MUST be a valid JSON Schema.
// When present alongside of "if", the instance successfully validates against this keyword if it validates against both the "if"'s subschema and this keyword's subschema.
// When "if" is absent, or the instance fails to validate against its subschema, validation against this keyword always succeeds. Implementations SHOULD avoid attempting to validate against the subschema in these cases.
type Then Schema

// Validate implements the Validator interface for Then
func (t *Then) Validate(data interface{}) error {
	return nil
}

// JSONProp implements JSON property name indexing for Then
func (t Then) JSONProp(name string) interface{} {
	return Schema(t).JSONProp(name)
}

// UnmarshalJSON implements the json.Unmarshaler interface for Then
func (t *Then) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*t = Then(sch)
	return nil
}

// Else MUST be a valid JSON Schema.
// When present alongside of "if", the instance successfully validates against this keyword if it fails to validate against the "if"'s subschema, and successfully validates against this keyword's subschema.
// When "if" is absent, or the instance successfully validates against its subschema, validation against this keyword always succeeds. Implementations SHOULD avoid attempting to validate against the subschema in these cases.
type Else Schema

// Validate implements the Validator interface for Else
func (e *Else) Validate(data interface{}) error {
	return nil
}

// JSONProp implements JSON property name indexing for Else
func (e Else) JSONProp(name string) interface{} {
	return Schema(e).JSONProp(name)
}

// UnmarshalJSON implements the json.Unmarshaler interface for Else
func (e *Else) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*e = Else(sch)
	return nil
}
