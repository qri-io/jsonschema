package jsonschema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// items MUST be either a valid JSON Schema or an array of valid JSON Schemas.
// This keyword determines how child instances validate for arrays, and does not directly validate the
// immediate instance itself.
// * If "items" is a schema, validation succeeds if all elements in the array successfully validate
//   against that schema.
// * If "items" is an array of schemas, validation succeeds if each element of the instance validates
//   against the schema at the same position, if any.
// * Omitting this keyword has the same behavior as an empty schema.
type items struct {
	// need to track weather user specficied a singl object or arry
	// b/c it affects additionalItems validation semantics
	single  bool
	Schemas []*Schema
}

func newItems() Validator {
	return &items{}
}

// Validate implements the Validator interface for items
func (it items) Validate(data interface{}) error {
	if arr, ok := data.([]interface{}); ok {
		if it.single {
			for i, elem := range arr {
				if err := it.Schemas[0].Validate(elem); err != nil {
					return fmt.Errorf("element %d %s", i, err.Error())
				}
			}
		} else {
			for i, vs := range it.Schemas {
				if i < len(arr) {
					if err := vs.Validate(arr[i]); err != nil {
						return fmt.Errorf("element %d %s", i, err.Error())
					}
				}
			}
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for items
func (it items) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(it.Schemas) || idx < 0 {
		return nil
	}
	return it.Schemas[idx]
}

// JSONChildren implements the JSONContainer interface for items
func (it items) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, sch := range it.Schemas {
		res[strconv.Itoa(i)] = sch
	}
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface for items
func (it *items) UnmarshalJSON(data []byte) error {
	s := &Schema{}
	if err := json.Unmarshal(data, s); err == nil {
		*it = items{single: true, Schemas: []*Schema{s}}
		return nil
	}
	ss := []*Schema{}
	if err := json.Unmarshal(data, &ss); err != nil {
		return err
	}
	*it = items{Schemas: ss}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for items
func (it items) MarshalJSON() ([]byte, error) {
	if it.single {
		return json.Marshal(it.Schemas[0])
	}
	return json.Marshal([]*Schema(it.Schemas))
}

// additionalItems determines how child instances validate for arrays, and does not directly validate the immediate
// instance itself.
// If "items" is an array of schemas, validation succeeds if every instance element at a position greater than
// the size of "items" validates against "additionalItems".
// Otherwise, "additionalItems" MUST be ignored, as the "items" schema (possibly the default value of an empty schema) is applied to all elements.
// Omitting this keyword has the same behavior as an empty schema.
type additionalItems struct {
	startIndex int
	Schema     *Schema
}

func newAdditionalItems() Validator {
	return &additionalItems{}
}

// Validate implements the Validator interface for additionalItems
func (a *additionalItems) Validate(data interface{}) error {
	if a.startIndex >= 0 {
		if arr, ok := data.([]interface{}); ok {
			for i, elem := range arr {
				if i < a.startIndex {
					continue
				}
				if err := a.Schema.Validate(elem); err != nil {
					return fmt.Errorf("element %d: %s", i, err.Error())
				}
			}
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for additionalItems
func (a *additionalItems) JSONProp(name string) interface{} {
	return a.Schema.JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for additionalItems
func (a *additionalItems) JSONChildren() (res map[string]JSONPather) {
	if a.Schema == nil {
		return map[string]JSONPather{}
	}
	return a.Schema.JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for additionalItems
func (a *additionalItems) UnmarshalJSON(data []byte) error {
	sch := &Schema{}
	if err := json.Unmarshal(data, sch); err != nil {
		return err
	}
	// begin with -1 as default index to prevent additionalItems from evaluating
	// unless startIndex is explicitly set
	*a = additionalItems{startIndex: -1, Schema: sch}
	return nil
}

// maxItems MUST be a non-negative integer.
// An array instance is valid against "maxItems" if its size is less than, or equal to, the value of this keyword.
type maxItems int

func newMaxItems() Validator {
	return new(maxItems)
}

// Validate implements the Validator interface for maxItems
func (m maxItems) Validate(data interface{}) error {
	if arr, ok := data.([]interface{}); ok {
		if len(arr) > int(m) {
			return fmt.Errorf("%d array items exceeds %d max", len(arr), m)
		}
	}
	return nil
}

// minItems MUST be a non-negative integer.
// An array instance is valid against "minItems" if its size is greater than, or equal to, the value of this keyword.
// Omitting this keyword has the same behavior as a value of 0.
type minItems int

func newMinItems() Validator {
	return new(minItems)
}

// Validate implements the Validator interface for minItems
func (m minItems) Validate(data interface{}) error {
	if arr, ok := data.([]interface{}); ok {
		if len(arr) < int(m) {
			return fmt.Errorf("%d array items below %d minimum", len(arr), m)
		}
	}
	return nil
}

// uniqueItems requires array instance elements be unique
// If this keyword has boolean value false, the instance validates successfully. If it has
// boolean value true, the instance validates successfully if all of its elements are unique.
// Omitting this keyword has the same behavior as a value of false.
type uniqueItems bool

func newUniqueItems() Validator {
	return new(uniqueItems)
}

// Validate implements the Validator interface for uniqueItems
func (u *uniqueItems) Validate(data interface{}) error {
	if arr, ok := data.([]interface{}); ok {
		found := []interface{}{}
		for _, elem := range arr {
			for _, f := range found {
				if reflect.DeepEqual(f, elem) {
					return fmt.Errorf("arry must be unique: %v", arr)
				}
			}
			found = append(found, elem)
		}
	}
	return nil
}

// contains validates that an array instance is valid against "contains" if at
// least one of its elements is valid against the given schema.
type contains Schema

func newContains() Validator {
	return &contains{}
}

// Validate implements the Validator interface for contains
func (c *contains) Validate(data interface{}) error {
	v := Schema(*c)
	if arr, ok := data.([]interface{}); ok {
		for _, elem := range arr {
			if err := v.Validate(elem); err == nil {
				return nil
			}
		}
		return fmt.Errorf("expected %v to contain at least one of: %s", data, c)
	}
	return nil
}

// JSONProp implements JSON property name indexing for contains
func (c contains) JSONProp(name string) interface{} {
	return Schema(c).JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for contains
func (c contains) JSONChildren() (res map[string]JSONPather) {
	return Schema(c).JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for contains
func (c *contains) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*c = contains(sch)
	return nil
}
