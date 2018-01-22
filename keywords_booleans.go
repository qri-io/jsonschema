package jsonschema

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// allOf MUST be a non-empty array. Each item of the array MUST be a valid JSON Schema.
// An instance validates successfully against this keyword if it validates successfully against all schemas defined by this keyword's value.
type allOf []*Schema

func newAllOf() Validator {
	return &allOf{}
}

// Validate implements the validator interface for allOf
func (a allOf) Validate(data interface{}) (errs []ValError) {
	for _, sch := range a {
		if ves := sch.Validate(data); len(ves) > 0 {
			errs = append(errs, ves...)
		}
	}
	return
}

// JSONProp implements JSON property name indexing for allOf
func (a allOf) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(a) || idx < 0 {
		return nil
	}
	return a[idx]
}

// JSONChildren implements the JSONContainer interface for allOf
func (a allOf) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, sch := range a {
		res[strconv.Itoa(i)] = sch
	}
	return
}

// anyOf MUST be a non-empty array. Each item of the array MUST be a valid JSON Schema.
// An instance validates successfully against this keyword if it validates successfully against at
// least one schema defined by this keyword's value.
type anyOf []*Schema

func newAnyOf() Validator {
	return &anyOf{}
}

// Validate implements the validator interface for anyOf
func (a anyOf) Validate(data interface{}) []ValError {
	for _, sch := range a {
		if err := sch.Validate(data); err == nil {
			return nil
		}
	}
	return []ValError{
		{Message: fmt.Sprintf("value did not match any specified anyOf schemas: %v", data)},
	}
}

// JSONProp implements JSON property name indexing for anyOf
func (a anyOf) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(a) || idx < 0 {
		return nil
	}
	return a[idx]
}

// JSONChildren implements the JSONContainer interface for anyOf
func (a anyOf) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, sch := range a {
		res[strconv.Itoa(i)] = sch
	}
	return
}

// oneOf MUST be a non-empty array. Each item of the array MUST be a valid JSON Schema.
// An instance validates successfully against this keyword if it validates successfully against exactly one schema defined by this keyword's value.
type oneOf []*Schema

func newOneOf() Validator {
	return &oneOf{}
}

// Validate implements the validator interface for oneOf
func (o oneOf) Validate(data interface{}) []ValError {
	matched := false
	for _, sch := range o {
		if err := sch.Validate(data); err == nil {
			if matched {
				return []ValError{
					{Message: fmt.Sprintf("value matched more than one specified oneOf schemas")},
				}
			}
			matched = true
		}
	}
	if !matched {
		return []ValError{
			{Message: fmt.Sprintf("value did not match any of the specified oneOf schemas")},
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for oneOf
func (o oneOf) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(o) || idx < 0 {
		return nil
	}
	return o[idx]
}

// JSONChildren implements the JSONContainer interface for oneOf
func (o oneOf) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, sch := range o {
		res[strconv.Itoa(i)] = sch
	}
	return
}

// not MUST be a valid JSON Schema.
// An instance is valid against this keyword if it fails to validate successfully against the schema defined
// by this keyword.
type not Schema

func newNot() Validator {
	return &not{}
}

// Validate implements the validator interface for not
func (n *not) Validate(data interface{}) []ValError {
	sch := Schema(*n)
	if sch.Validate(data) == nil {
		// TODO - make this error actually make sense
		return []ValError{
			{Message: fmt.Sprintf("not clause")},
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for not
func (n not) JSONProp(name string) interface{} {
	return Schema(n).JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for not
func (n not) JSONChildren() (res map[string]JSONPather) {
	if n.Ref != "" {
		s := Schema(n)
		return map[string]JSONPather{"$ref": &s}
	}
	return Schema(n).JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for not
func (n *not) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*n = not(sch)
	return nil
}

// MarshalJSON implements json.Marshaller for not
func (n not) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(n))
}
