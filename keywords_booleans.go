package jsonschema

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// AllOf MUST be a non-empty array. Each item of the array MUST be a valid JSON Schema.
// An instance validates successfully against this keyword if it validates successfully against all schemas defined by this keyword's value.
type AllOf []*Schema

// Validate implements the validator interface for AllOf
func (a AllOf) Validate(data interface{}) error {
	for i, sch := range a {
		if err := sch.Validate(data); err != nil {
			return fmt.Errorf("allOf element %d error: %s", i, err.Error())
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for AllOf
func (a AllOf) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(a) || idx < 0 {
		return nil
	}
	return a[idx]
}

// JSONChildren implements the JSONContainer interface for AllOf
func (a AllOf) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, sch := range a {
		res[strconv.Itoa(i)] = sch
	}
	return
}

// AnyOf MUST be a non-empty array. Each item of the array MUST be a valid JSON Schema.
// An instance validates successfully against this keyword if it validates successfully against at
// least one schema defined by this keyword's value.
type AnyOf []*Schema

// Validate implements the validator interface for AnyOf
func (a AnyOf) Validate(data interface{}) error {
	for _, sch := range a {
		if err := sch.Validate(data); err == nil {
			return nil
		}
	}
	return fmt.Errorf("value did not match any specified anyOf schemas: %v", data)
}

// JSONProp implements JSON property name indexing for AnyOf
func (a AnyOf) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(a) || idx < 0 {
		return nil
	}
	return a[idx]
}

// JSONChildren implements the JSONContainer interface for AnyOf
func (a AnyOf) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, sch := range a {
		res[strconv.Itoa(i)] = sch
	}
	return
}

// OneOf MUST be a non-empty array. Each item of the array MUST be a valid JSON Schema.
// An instance validates successfully against this keyword if it validates successfully against exactly one schema defined by this keyword's value.
type OneOf []*Schema

// Validate implements the validator interface for OneOf
func (o OneOf) Validate(data interface{}) error {
	matched := false
	for _, sch := range o {
		if err := sch.Validate(data); err == nil {
			if matched {
				return fmt.Errorf("value matched more than one specified oneOf schemas")
			}
			matched = true
		}
	}
	if !matched {
		return fmt.Errorf("value did not match any of the specified oneOf schemas")
	}
	return nil
}

// JSONProp implements JSON property name indexing for OneOf
func (o OneOf) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(o) || idx < 0 {
		return nil
	}
	return o[idx]
}

// JSONChildren implements the JSONContainer interface for OneOf
func (o OneOf) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, sch := range o {
		res[strconv.Itoa(i)] = sch
	}
	return
}

// Not MUST be a valid JSON Schema.
// An instance is valid against this keyword if it fails to validate successfully against the schema defined
// by this keyword.
type Not Schema

// Validate implements the validator interface for Not
func (n *Not) Validate(data interface{}) error {
	sch := Schema(*n)
	if sch.Validate(data) == nil {
		// TODO - make this error actually make sense
		return fmt.Errorf("not clause")
	}
	return nil
}

// JSONProp implements JSON property name indexing for Not
func (n Not) JSONProp(name string) interface{} {
	return Schema(n).JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for Not
func (n Not) JSONChildren() (res map[string]JSONPather) {
	if n.Ref != "" {
		s := Schema(n)
		return map[string]JSONPather{"$ref": &s}
	}
	return Schema(n).JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for Not
func (n *Not) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*n = Not(sch)
	return nil
}
