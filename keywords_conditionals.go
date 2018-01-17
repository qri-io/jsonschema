package jsonschema

import (
	"encoding/json"
)

// iif MUST be a valid JSON Schema.
// Instances that successfully validate against this keyword's subschema MUST also be valid against the subschema value of the "then" keyword, if present.
// Instances that fail to validate against this keyword's subschema MUST also be valid against the subschema value of the "else" keyword.
// Validation of the instance against this keyword on its own always succeeds, regardless of the validation outcome of against its subschema.
type iif struct {
	Schema Schema
	then   *then
	els    *els
}

func newIif() Validator {
	return &iif{}
}

// Validate implements the Validator interface for iif
func (i *iif) Validate(data interface{}) error {
	if err := i.Schema.Validate(data); err == nil {
		if i.then != nil {
			s := Schema(*i.then)
			sch := &s
			return sch.Validate(data)
		}
	} else {
		if i.els != nil {
			s := Schema(*i.els)
			sch := &s
			return sch.Validate(data)
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for iif
func (i iif) JSONProp(name string) interface{} {
	return Schema(i.Schema).JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for iif
func (i iif) JSONChildren() (res map[string]JSONPather) {
	return i.Schema.JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for iif
func (i *iif) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*i = iif{Schema: sch}
	return nil
}

// MarshalJSON implements json.Marshaler for iif
func (i iif) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Schema)
}

// then MUST be a valid JSON Schema.
// When present alongside of "if", the instance successfully validates against this keyword if it validates against both the "if"'s subschema and this keyword's subschema.
// When "if" is absent, or the instance fails to validate against its subschema, validation against this keyword always succeeds. Implementations SHOULD avoid attempting to validate against the subschema in these cases.
type then Schema

func newThen() Validator {
	return &then{}
}

// Validate implements the Validator interface for then
func (t *then) Validate(data interface{}) error {
	return nil
}

// JSONProp implements JSON property name indexing for then
func (t then) JSONProp(name string) interface{} {
	return Schema(t).JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for iif
func (t then) JSONChildren() (res map[string]JSONPather) {
	return Schema(t).JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for then
func (t *then) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*t = then(sch)
	return nil
}

// MarshalJSON implements json.Marshaler for then
func (t then) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(t))
}

// els MUST be a valid JSON Schema.
// When present alongside of "if", the instance successfully validates against this keyword if it fails to validate against the "if"'s subschema, and successfully validates against this keyword's subschema.
// When "if" is absent, or the instance successfully validates against its subschema, validation against this keyword always succeeds. Implementations SHOULD avoid attempting to validate against the subschema in these cases.
type els Schema

func newEls() Validator {
	return &els{}
}

// Validate implements the Validator interface for els
func (e *els) Validate(data interface{}) error {
	return nil
}

// JSONProp implements JSON property name indexing for els
func (e els) JSONProp(name string) interface{} {
	return Schema(e).JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for els
func (e els) JSONChildren() (res map[string]JSONPather) {
	return Schema(e).JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for els
func (e *els) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*e = els(sch)
	return nil
}

// MarshalJSON implements json.Marshaler for els
func (e els) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(e))
}
