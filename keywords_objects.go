package jsonschema

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

// MaxProperties MUST be a non-negative integer.
// An object instance is valid against "MaxProperties" if its number of Properties is less than, or equal to, the value of this keyword.
type MaxProperties int

// NewMaxProperties allocates a new MaxProperties validator
func NewMaxProperties() Validator {
	return new(MaxProperties)
}

// Validate implements the validator interface for MaxProperties
func (m MaxProperties) Validate(data interface{}) []ValError {
	if obj, ok := data.(map[string]interface{}); ok {
		if len(obj) > int(m) {
			return []ValError{
				{Message: fmt.Sprintf("%d object Properties exceed %d maximum", len(obj), m)},
			}
		}
	}
	return nil
}

// minProperties MUST be a non-negative integer.
// An object instance is valid against "minProperties" if its number of Properties is greater than, or equal to, the value of this keyword.
// Omitting this keyword has the same behavior as a value of 0.
type minProperties int

// NewMinProperties allocates a new MinProperties validator
func NewMinProperties() Validator {
	return new(minProperties)
}

// Validate implements the validator interface for minProperties
func (m minProperties) Validate(data interface{}) []ValError {
	if obj, ok := data.(map[string]interface{}); ok {
		if len(obj) < int(m) {
			return []ValError{
				{Message: fmt.Sprintf("%d object Properties below %d minimum", len(obj), m)},
			}
		}
	}
	return nil
}

// Required ensures that for a given object instance, every item in the array is the name of a property in the instance.
// The value of this keyword MUST be an array. Elements of this array, if any, MUST be strings, and MUST be unique.
// Omitting this keyword has the same behavior as an empty array.
type Required []string

// NewRequired allocates a new Required validator
func NewRequired() Validator {
	return &Required{}
}

// Validate implements the validator interface for Required
func (r Required) Validate(data interface{}) []ValError {
	if obj, ok := data.(map[string]interface{}); ok {
		for _, key := range r {
			if val, ok := obj[key]; val == nil && !ok {
				return []ValError{
					{Message: fmt.Sprintf(`"%s" value is required`, key)},
				}
			}
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for Required
func (r Required) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(r) || idx < 0 {
		return nil
	}
	return r[idx]
}

// Properties MUST be an object. Each value of this object MUST be a valid JSON Schema.
// This keyword determines how child instances validate for objects, and does not directly validate
// the immediate instance itself.
// Validation succeeds if, for each name that appears in both the instance and as a name within this
// keyword's value, the child instance for that name successfully validates against the corresponding schema.
// Omitting this keyword has the same behavior as an empty object.
type Properties map[string]*Schema

// NewProperties allocates a new Properties validator
func NewProperties() Validator {
	return &Properties{}
}

// Validate implements the validator interface for Properties
func (p Properties) Validate(data interface{}) []ValError {
	if obj, ok := data.(map[string]interface{}); ok {
		for key, val := range obj {
			if p[key] != nil {
				if err := p[key].Validate(val); err != nil {
					return []ValError{
						{Message: fmt.Sprintf(`"%s" property %s`, key, err)},
					}
				}
			}
		}
	}
	return nil
}

// JSONProp implements JSON property name indexing for Properties
func (p Properties) JSONProp(name string) interface{} {
	return p[name]
}

// JSONChildren implements the JSONContainer interface for Properties
func (p Properties) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for key, sch := range p {
		res[key] = sch
	}
	return
}

// PatternProperties determines how child instances validate for objects, and does not directly validate the immediate instance itself.
// Validation of the primitive instance type against this keyword always succeeds.
// Validation succeeds if, for each instance name that matches any regular expressions that appear as a property name in this
// keyword's value, the child instance for that name successfully validates against each schema that corresponds to a matching
// regular expression.
// Each property name of this object SHOULD be a valid regular expression,
// according to the ECMA 262 regular expression dialect.
// Each property value of this object MUST be a valid JSON Schema.
// Omitting this keyword has the same behavior as an empty object.
type PatternProperties []patternSchema

// NewPatternProperties allocates a new PatternProperties validator
func NewPatternProperties() Validator {
	return &PatternProperties{}
}

type patternSchema struct {
	key    string
	re     *regexp.Regexp
	schema *Schema
}

// Validate implements the validator interface for PatternProperties
func (p PatternProperties) Validate(data interface{}) (errs []ValError) {
	if obj, ok := data.(map[string]interface{}); ok {
		for key, val := range obj {
			for _, ptn := range p {
				if ptn.re.Match([]byte(key)) {
					if ves := ptn.schema.Validate(val); len(ves) > 0 {
						// fmt.Sprintf("object key %s pattern prop %s error: %s", key, ptn.key, err.Error())
						errs = append(errs, ves...)
					}
				}
			}
		}
	}
	return
}

// JSONProp implements JSON property name indexing for PatternProperties
func (p PatternProperties) JSONProp(name string) interface{} {
	for _, pp := range p {
		if pp.key == name {
			return pp.schema
		}
	}
	return nil
}

// JSONChildren implements the JSONContainer interface for PatternProperties
func (p PatternProperties) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, pp := range p {
		res[strconv.Itoa(i)] = pp.schema
	}
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface for PatternProperties
func (p *PatternProperties) UnmarshalJSON(data []byte) error {
	var props map[string]*Schema
	if err := json.Unmarshal(data, &props); err != nil {
		return err
	}

	ptn := make(PatternProperties, len(props))
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

// MarshalJSON implements json.Marshaler for PatternProperties
func (p PatternProperties) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{}
	for _, prop := range p {
		obj[prop.key] = prop.schema
	}
	return json.Marshal(obj)
}

// AdditionalProperties determines how child instances validate for objects, and does not directly validate the immediate instance itself.
// Validation with "AdditionalProperties" applies only to the child values of instance names that do not match any names in "Properties",
// and do not match any regular expression in "PatternProperties".
// For all such Properties, validation succeeds if the child instance validates against the "AdditionalProperties" schema.
// Omitting this keyword has the same behavior as an empty schema.
type AdditionalProperties struct {
	Properties *Properties
	patterns   *PatternProperties
	Schema     *Schema
}

// NewAdditionalProperties allocates a new AdditionalProperties validator
func NewAdditionalProperties() Validator {
	return &AdditionalProperties{}
}

// Validate implements the validator interface for AdditionalProperties
func (ap AdditionalProperties) Validate(data interface{}) []ValError {
	if obj, ok := data.(map[string]interface{}); ok {
	KEYS:
		for key, val := range obj {
			if ap.Properties != nil {
				for propKey := range *ap.Properties {
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
			if ves := ap.Schema.Validate(val); len(ves) > 0 {
				// fmt.Sprintf("object key %s AdditionalProperties error: %s", key, err.Error())
				return ves
			}
		}
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for AdditionalProperties
func (ap *AdditionalProperties) UnmarshalJSON(data []byte) error {
	sch := &Schema{}
	if err := json.Unmarshal(data, sch); err != nil {
		return err
	}
	// fmt.Println("unmarshal:", sch.Ref)
	*ap = AdditionalProperties{Schema: sch}
	return nil
}

// JSONProp implements JSON property name indexing for AdditionalProperties
func (ap *AdditionalProperties) JSONProp(name string) interface{} {
	return ap.Schema.JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for AdditionalProperties
func (ap *AdditionalProperties) JSONChildren() (res map[string]JSONPather) {
	if ap.Schema.Ref != "" {
		return map[string]JSONPather{"$ref": ap.Schema}
	}
	return ap.Schema.JSONChildren()
}

// MarshalJSON implements json.Marshaler for AdditionalProperties
func (ap AdditionalProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(ap.Schema)
}

// Dependencies : [CREF1]
// This keyword specifies rules that are evaluated if the instance is an object and contains a
// certain property.
// This keyword's value MUST be an object. Each property specifies a Dependency.
// Each Dependency value MUST be an array or a valid JSON Schema.
// If the Dependency value is a subschema, and the Dependency key is a property in the instance,
// the entire instance must validate against the Dependency value.
// If the Dependency value is an array, each element in the array, if any, MUST be a string,
// and MUST be unique. If the Dependency key is a property in the instance, each of the items
// in the Dependency value must be a property that exists in the instance.
// Omitting this keyword has the same behavior as an empty object.
type Dependencies map[string]Dependency

// NewDependencies allocates a new Dependencies validator
func NewDependencies() Validator {
	return &Dependencies{}
}

// Validate implements the validator interface for Dependencies
func (d Dependencies) Validate(data interface{}) (errs []ValError) {
	if obj, ok := data.(map[string]interface{}); ok {
		for key, val := range d {
			if obj[key] != nil {
				if ves := val.Validate(obj); len(ves) > 0 {
					errs = append(errs, ves...)
				}
			}
		}
	}
	return
}

// JSONProp implements JSON property name indexing for Dependencies
func (d Dependencies) JSONProp(name string) interface{} {
	return d[name]
}

// JSONChildren implements the JSONContainer interface for Dependencies
// func (d Dependencies) JSONChildren() (res map[string]JSONPather) {
// 	res = map[string]JSONPather{}
// 	for key, dep := range d {
// 		if dep.schema != nil {
// 			res[key] = dep.schema
// 		}
// 	}
// 	return
// }

// Dependency is an instance used only in the Dependencies proprty
type Dependency struct {
	schema *Schema
	props  []string
}

// Validate implements the validator interface for Dependency
func (d Dependency) Validate(data interface{}) (errs []ValError) {
	if obj, ok := data.(map[string]interface{}); ok {
		if d.schema != nil {
			return d.schema.Validate(data)
		} else if len(d.props) > 0 {
			for _, k := range d.props {
				if obj[k] == nil {
					errs = append(errs, ValError{Message: fmt.Sprintf("Dependency property %s is Required", k)})
				}
			}
		}
	}
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface for Dependencies
func (d *Dependency) UnmarshalJSON(data []byte) error {
	props := []string{}
	if err := json.Unmarshal(data, &props); err == nil {
		*d = Dependency{props: props}
		return nil
	}
	sch := &Schema{}
	err := json.Unmarshal(data, sch)

	if err == nil {
		*d = Dependency{schema: sch}
	}
	return err
}

// MarshalJSON implements json.Marshaler for Dependency
func (d Dependency) MarshalJSON() ([]byte, error) {
	if d.schema != nil {
		return json.Marshal(d.schema)
	}
	return json.Marshal(d.props)
}

// PropertyNames checks if every property name in the instance validates against the provided schema
// if the instance is an object.
// Note the property name that the schema is testing will always be a string.
// Omitting this keyword has the same behavior as an empty schema.
type PropertyNames Schema

// NewPropertyNames allocates a new PropertyNames validator
func NewPropertyNames() Validator {
	return &PropertyNames{}
}

// Validate implements the validator interface for PropertyNames
func (p PropertyNames) Validate(data interface{}) (errs []ValError) {
	sch := Schema(p)
	if obj, ok := data.(map[string]interface{}); ok {
		for key := range obj {
			if ves := sch.Validate(key); len(ves) > 0 {
				errs = append(errs, ves...)
			}
		}
	}
	return
}

// JSONProp implements JSON property name indexing for Properties
func (p PropertyNames) JSONProp(name string) interface{} {
	return Schema(p).JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for PropertyNames
func (p PropertyNames) JSONChildren() (res map[string]JSONPather) {
	return Schema(p).JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for PropertyNames
func (p *PropertyNames) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*p = PropertyNames(sch)
	return nil
}

// MarshalJSON implements json.Marshaler for PropertyNames
func (p PropertyNames) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(p))
}
