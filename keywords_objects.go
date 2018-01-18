package jsonschema

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

// maxProperties MUST be a non-negative integer.
// An object instance is valid against "maxProperties" if its number of properties is less than, or equal to, the value of this keyword.
type maxProperties int

func newMaxProperties() Validator {
	return new(maxProperties)
}

// Validate implements the validator interface for maxProperties
func (m maxProperties) Validate(data interface{}) error {
	if obj, ok := data.(map[string]interface{}); ok {
		if len(obj) > int(m) {
			return fmt.Errorf("%d object properties exceed %d maximum", len(obj), m)
		}
	}
	return nil
}

// minProperties MUST be a non-negative integer.
// An object instance is valid against "minProperties" if its number of properties is greater than, or equal to, the value of this keyword.
// Omitting this keyword has the same behavior as a value of 0.
type minProperties int

func newMinProperties() Validator {
	return new(minProperties)
}

// Validate implements the validator interface for minProperties
func (m minProperties) Validate(data interface{}) error {
	if obj, ok := data.(map[string]interface{}); ok {
		if len(obj) < int(m) {
			return fmt.Errorf("%d object properties below %d minimum", len(obj), m)
		}
	}
	return nil
}

// required ensures that for a given object instance, every item in the array is the name of a property in the instance.
// The value of this keyword MUST be an array. Elements of this array, if any, MUST be strings, and MUST be unique.
// Omitting this keyword has the same behavior as an empty array.
type required []string

func newRequired() Validator {
	return &required{}
}

// Validate implements the validator interface for required
func (r required) Validate(data interface{}) error {
	if obj, ok := data.(map[string]interface{}); ok {
		for _, key := range r {
			if val, ok := obj[key]; val == nil && !ok {
				return fmt.Errorf(`"%s" value is required`, key)
			}
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for required
func (r required) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(r) || idx < 0 {
		return nil
	}
	return r[idx]
}

// properties MUST be an object. Each value of this object MUST be a valid JSON Schema.
// This keyword determines how child instances validate for objects, and does not directly validate
// the immediate instance itself.
// Validation succeeds if, for each name that appears in both the instance and as a name within this
// keyword's value, the child instance for that name successfully validates against the corresponding schema.
// Omitting this keyword has the same behavior as an empty object.
type properties map[string]*Schema

func newProperties() Validator {
	return &properties{}
}

// Validate implements the validator interface for properties
func (p properties) Validate(data interface{}) error {
	if obj, ok := data.(map[string]interface{}); ok {
		for key, val := range obj {
			if p[key] != nil {
				if err := p[key].Validate(val); err != nil {
					return fmt.Errorf(`"%s" property %s`, key, err)
				}
			}
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for properties
func (p properties) JSONProp(name string) interface{} {
	return p[name]
}

// JSONChildren implements the JSONContainer interface for properties
func (p properties) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for key, sch := range p {
		res[key] = sch
	}
	return
}

// patternProperties determines how child instances validate for objects, and does not directly validate the immediate instance itself.
// Validation of the primitive instance type against this keyword always succeeds.
// Validation succeeds if, for each instance name that matches any regular expressions that appear as a property name in this
// keyword's value, the child instance for that name successfully validates against each schema that corresponds to a matching
// regular expression.
// Each property name of this object SHOULD be a valid regular expression,
// according to the ECMA 262 regular expression dialect.
// Each property value of this object MUST be a valid JSON Schema.
// Omitting this keyword has the same behavior as an empty object.
type patternProperties []patternSchema

func newPatternProperties() Validator {
	return &patternProperties{}
}

type patternSchema struct {
	key    string
	re     *regexp.Regexp
	schema *Schema
}

// Validate implements the validator interface for patternProperties
func (p patternProperties) Validate(data interface{}) error {
	if obj, ok := data.(map[string]interface{}); ok {
		for key, val := range obj {
			for _, ptn := range p {
				if ptn.re.Match([]byte(key)) {
					if err := ptn.schema.Validate(val); err != nil {
						return fmt.Errorf("object key %s pattern prop %s error: %s", key, ptn.key, err.Error())
					}
				}
			}
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for patternProperties
func (p patternProperties) JSONProp(name string) interface{} {
	for _, pp := range p {
		if pp.key == name {
			return pp.schema
		}
	}
	return nil
}

// JSONChildren implements the JSONContainer interface for patternProperties
func (p patternProperties) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, pp := range p {
		res[strconv.Itoa(i)] = pp.schema
	}
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface for patternProperties
func (p *patternProperties) UnmarshalJSON(data []byte) error {
	var props map[string]*Schema
	if err := json.Unmarshal(data, &props); err != nil {
		return err
	}

	ptn := make(patternProperties, len(props))
	i := 0
	for key, sch := range props {
		re, err := regexp.Compile(key)
		if err != nil {
			return fmt.Errorf("invalid pattern: %s: %s", key, err.Error())
		}
		ptn[i] = patternSchema{
			key:    key,
			re:     re,
			schema: sch,
		}
		i++
	}

	*p = ptn
	return nil
}

// MarshalJSON implements json.Marshaler for patternProperties
func (p patternProperties) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{}
	for _, prop := range p {
		obj[prop.key] = prop.schema
	}
	return json.Marshal(obj)
}

// additionalProperties determines how child instances validate for objects, and does not directly validate the immediate instance itself.
// Validation with "additionalProperties" applies only to the child values of instance names that do not match any names in "properties",
// and do not match any regular expression in "patternproperties".
// For all such properties, validation succeeds if the child instance validates against the "additionalProperties" schema.
// Omitting this keyword has the same behavior as an empty schema.
type additionalProperties struct {
	properties *properties
	patterns   *patternProperties
	Schema     *Schema
}

func newAdditionalProperties() Validator {
	return &additionalProperties{}
}

// Validate implements the validator interface for additionalProperties
func (ap additionalProperties) Validate(data interface{}) error {
	if obj, ok := data.(map[string]interface{}); ok {
	KEYS:
		for key, val := range obj {
			if ap.properties != nil {
				for propKey := range *ap.properties {
					if propKey == key {
						continue KEYS
					}
				}
			}
			if ap.patterns != nil {
				for _, ptn := range *ap.patterns {
					if ptn.re.Match([]byte(key)) {
						continue KEYS
					}
				}
			}
			if err := ap.Schema.Validate(val); err != nil {
				return fmt.Errorf("object key %s additionalProperties error: %s", key, err.Error())
			}
		}
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for additionalProperties
func (ap *additionalProperties) UnmarshalJSON(data []byte) error {
	sch := &Schema{}
	if err := json.Unmarshal(data, sch); err != nil {
		return err
	}
	// fmt.Println("unmarshal:", sch.Ref)
	*ap = additionalProperties{Schema: sch}
	return nil
}

// JSONProp implements JSON property name indexing for additionalProperties
func (ap *additionalProperties) JSONProp(name string) interface{} {
	return ap.Schema.JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for additionalProperties
func (ap *additionalProperties) JSONChildren() (res map[string]JSONPather) {
	if ap.Schema.Ref != "" {
		return map[string]JSONPather{"$ref": ap.Schema}
	}
	return ap.Schema.JSONChildren()
}

// MarshalJSON implements json.Marshaler for additionalProperties
func (ap additionalProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(ap.Schema)
}

// dependencies : [CREF1]
// This keyword specifies rules that are evaluated if the instance is an object and contains a
// certain property.
// This keyword's value MUST be an object. Each property specifies a dependency.
// Each dependency value MUST be an array or a valid JSON Schema.
// If the dependency value is a subschema, and the dependency key is a property in the instance,
// the entire instance must validate against the dependency value.
// If the dependency value is an array, each element in the array, if any, MUST be a string,
// and MUST be unique. If the dependency key is a property in the instance, each of the items
// in the dependency value must be a property that exists in the instance.
// Omitting this keyword has the same behavior as an empty object.
type dependencies map[string]dependency

func newDependencies() Validator {
	return &dependencies{}
}

// Validate implements the validator interface for dependencies
func (d dependencies) Validate(data interface{}) error {
	if obj, ok := data.(map[string]interface{}); ok {
		for key, val := range d {
			if obj[key] != nil {
				if err := val.Validate(obj); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for dependencies
func (d dependencies) JSONProp(name string) interface{} {
	return d[name]
}

// JSONChildren implements the JSONContainer interface for dependencies
// func (d dependencies) JSONChildren() (res map[string]JSONPather) {
// 	res = map[string]JSONPather{}
// 	for key, dep := range d {
// 		if dep.schema != nil {
// 			res[key] = dep.schema
// 		}
// 	}
// 	return
// }

// dependency is an instance used only in the dependencies proprty
type dependency struct {
	schema *Schema
	props  []string
}

// Validate implements the validator interface for dependency
func (d dependency) Validate(data interface{}) error {
	if obj, ok := data.(map[string]interface{}); ok {
		if d.schema != nil {
			return d.schema.Validate(data)
		} else if len(d.props) > 0 {
			for _, k := range d.props {
				if obj[k] == nil {
					return fmt.Errorf("dependency property %s is required", k)
				}
			}
		}
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for dependencies
func (d *dependency) UnmarshalJSON(data []byte) error {
	props := []string{}
	if err := json.Unmarshal(data, &props); err == nil {
		*d = dependency{props: props}
		return nil
	}
	sch := &Schema{}
	err := json.Unmarshal(data, sch)

	if err == nil {
		*d = dependency{schema: sch}
	}
	return err
}

// MarshalJSON implements json.Marshaler for dependency
func (d dependency) MarshalJSON() ([]byte, error) {
	if d.schema != nil {
		return json.Marshal(d.schema)
	}
	return json.Marshal(d.props)
}

// propertyNames checks if every property name in the instance validates against the provided schema
// if the instance is an object.
// Note the property name that the schema is testing will always be a string.
// Omitting this keyword has the same behavior as an empty schema.
type propertyNames Schema

func newPropertyNames() Validator {
	return &propertyNames{}
}

// Validate implements the validator interface for propertyNames
func (p propertyNames) Validate(data interface{}) error {
	sch := Schema(p)
	if obj, ok := data.(map[string]interface{}); ok {
		for key := range obj {
			if err := sch.Validate(key); err != nil {
				return fmt.Errorf("invalid propertyName: %s", err.Error())
			}
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for properties
func (p propertyNames) JSONProp(name string) interface{} {
	return Schema(p).JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for propertyNames
func (p propertyNames) JSONChildren() (res map[string]JSONPather) {
	return Schema(p).JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for propertyNames
func (p *propertyNames) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*p = propertyNames(sch)
	return nil
}

// MarshalJSON implements json.Marshaler for propertyNames
func (p propertyNames) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(p))
}
