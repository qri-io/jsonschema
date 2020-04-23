package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	jptr "github.com/qri-io/jsonpointer"
)

//
// Properties
//

type Properties map[string]*Schema

func NewProperties() Keyword {
	return &Properties{}
}

func (p Properties) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (p *Properties) Register(uri string, registry *SchemaRegistry) {
	for _, v := range *p {
		v.Register(uri, registry)
	}
}

func (p *Properties) Resolve(pointer jptr.Pointer, uri string) *Schema {
	// TODO: implement this
	return nil
}

func (p Properties) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		subCtx := NewSchemaContextFromSourceClean(*schCtx)
		for key, _ := range p {
			if obj[key] != nil {
				if _, ok := schCtx.Local.Keywords["additionalProperties"]; ok {
					schCtx.EvaluatedPropertyNames[key] = true
					schCtx.LocalEvaluatedPropertyNames[key] = true
				}
				subCtx.ClearContext()
				if schCtx.BaseRelativeLocation != nil {
					if newPtr, err := schCtx.BaseRelativeLocation.RawDescendant("properties", key); err == nil {
						subCtx.BaseRelativeLocation = &newPtr
					}
				}
				if newPtr, err := schCtx.RelativeLocation.RawDescendant("properties", key); err == nil {
					subCtx.RelativeLocation = &newPtr
				}
				if newPtr, err := schCtx.InstanceLocation.RawDescendant(key); err == nil {
					subCtx.InstanceLocation = &newPtr
				}
				subCtx.Instance = obj[key]
				errCountBefore := len(*errs)
				p[key].ValidateFromContext(subCtx, errs)
				errCountAfter := len(*errs)
				if _, ok := schCtx.Local.Keywords["additionalProperties"]; ok && errCountBefore == errCountAfter {
					JoinSets(&schCtx.EvaluatedPropertyNames, subCtx.EvaluatedPropertyNames)
					JoinSets(&schCtx.LocalEvaluatedPropertyNames, subCtx.LocalEvaluatedPropertyNames)
				}
			}
		}
	}
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

//
// Required
//

type Required []string

// NewRequired allocates a new Required validator
func NewRequired() Keyword {
	return &Required{}
}

func (r Required) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (r *Required) Register(uri string, registry *SchemaRegistry) {}

func (r *Required) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (r Required) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		for _, key := range r {
			if val, ok := obj[key]; val == nil && !ok {
				AddErrorCtx(errs, schCtx, fmt.Sprintf(`"%s" value is required`, key))
			}
		}
	}
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

//
// MaxProperties
//

type MaxProperties int

func NewMaxProperties() Keyword {
	return new(MaxProperties)
}

func (m MaxProperties) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *MaxProperties) Register(uri string, registry *SchemaRegistry) {}

func (m *MaxProperties) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m MaxProperties) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		if len(obj) > int(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("%d object Properties exceed %d maximum", len(obj), m))
		}
	}
}

//
// MinProperties
//

type MinProperties int

func NewMinProperties() Keyword {
	return new(MinProperties)
}

func (m MinProperties) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *MinProperties) Register(uri string, registry *SchemaRegistry) {}

func (m *MinProperties) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m MinProperties) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		if len(obj) < int(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("%d object Properties below %d minimum", len(obj), m))
		}
	}
}

//
// PatternProperties
//

type PatternProperties []patternSchema

// NewRequired allocates a new Required validator
func NewPatternProperties() Keyword {
	return &PatternProperties{}
}

type patternSchema struct {
	key    string
	re     *regexp.Regexp
	schema *Schema
}

func (p PatternProperties) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (p *PatternProperties) Register(uri string, registry *SchemaRegistry) {
	for _, v := range *p {
		v.schema.Register(uri, registry)
	}
}

func (p *PatternProperties) Resolve(pointer jptr.Pointer, uri string) *Schema {
	// TODO: implement this
	return nil
}

func (p PatternProperties) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		for key, val := range obj {
			for _, ptn := range p {
				if ptn.re.Match([]byte(key)) {
					schCtx.EvaluatedPropertyNames[key] = true
					schCtx.LocalEvaluatedPropertyNames[key] = true
					subCtx := NewSchemaContextFromSource(*schCtx)
					if schCtx.BaseRelativeLocation != nil {
						if newPtr, err := schCtx.BaseRelativeLocation.RawDescendant("patternProperties", key); err == nil {
							subCtx.BaseRelativeLocation = &newPtr
						}
					}
					if newPtr, err := schCtx.RelativeLocation.RawDescendant("patternProperties", key); err == nil {
						subCtx.RelativeLocation = &newPtr
					}
					if newPtr, err := schCtx.InstanceLocation.RawDescendant(key); err == nil {
						subCtx.InstanceLocation = &newPtr
					}
					subCtx.Instance = val
					errCountBefore := len(*errs)
					ptn.schema.ValidateFromContext(subCtx, errs)
					errCountAfter := len(*errs)

					if errCountBefore == errCountAfter {
						// TODO: check if this should be done only if the result is valid
						JoinSets(&schCtx.EvaluatedPropertyNames, subCtx.EvaluatedPropertyNames)
						JoinSets(&schCtx.LocalEvaluatedPropertyNames, subCtx.LocalEvaluatedPropertyNames)
						if schCtx.LastEvaluatedIndex < subCtx.LastEvaluatedIndex {
							schCtx.LastEvaluatedIndex = subCtx.LastEvaluatedIndex
						}
						if schCtx.LastEvaluatedIndex < subCtx.LocalLastEvaluatedIndex {
							schCtx.LastEvaluatedIndex = subCtx.LocalLastEvaluatedIndex
						}
					}
				}
			}
		}
	}
}

func (p PatternProperties) JSONProp(name string) interface{} {
	for _, pp := range p {
		if pp.key == name {
			return pp.schema
		}
	}
	return nil
}

func (p PatternProperties) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, pp := range p {
		res[strconv.Itoa(i)] = pp.schema
	}
	return
}

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

func (p PatternProperties) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{}
	for _, prop := range p {
		obj[prop.key] = prop.schema
	}
	return json.Marshal(obj)
}

//
// AdditionalProperties
//

type AdditionalProperties Schema

func NewAdditionalProperties() Keyword {
	return &AdditionalProperties{}
}

func (ap AdditionalProperties) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (ap *AdditionalProperties) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(ap).Register(uri, registry)
}

func (ap *AdditionalProperties) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(ap).Resolve(pointer, uri)
}

func (ap *AdditionalProperties) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		subCtx := NewSchemaContextFromSourceClean(*schCtx)
		if schCtx.BaseRelativeLocation != nil {
			if newPtr, err := schCtx.BaseRelativeLocation.RawDescendant("additionalProperties"); err == nil {
				subCtx.BaseRelativeLocation = &newPtr
			}
		}
		if newPtr, err := schCtx.RelativeLocation.RawDescendant("additionalProperties"); err == nil {
			subCtx.RelativeLocation = &newPtr
		}
		for key := range obj {
			if _, ok := schCtx.LocalEvaluatedPropertyNames[key]; ok {
				continue
			}
			if ap.schemaType == schemaTypeFalse {
				AddErrorCtx(errs, schCtx, "additional properties are not allowed")
				return
			}
			schCtx.EvaluatedPropertyNames[key] = true
			schCtx.LocalEvaluatedPropertyNames[key] = true
			subCtx.ClearContext()
			if newPtr, err := schCtx.InstanceLocation.RawDescendant(key); err == nil {
				subCtx.InstanceLocation = &newPtr
			}
			subCtx.Instance = obj[key]
			(*Schema)(ap).ValidateFromContext(subCtx, errs)
			JoinSets(&schCtx.EvaluatedPropertyNames, subCtx.EvaluatedPropertyNames)
			JoinSets(&schCtx.LocalEvaluatedPropertyNames, subCtx.LocalEvaluatedPropertyNames)
		}
	}
}

func (ap *AdditionalProperties) UnmarshalJSON(data []byte) error {
	sch := &Schema{}
	if err := json.Unmarshal(data, sch); err != nil {
		return err
	}
	*ap = (AdditionalProperties)(*sch)
	return nil
}

//
// PropertyNames
//

type PropertyNames Schema

func NewPropertyNames() Keyword {
	return &PropertyNames{}
}

func (p PropertyNames) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (p *PropertyNames) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(p).Register(uri, registry)
}

func (p *PropertyNames) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(p).Resolve(pointer, uri)
}

func (p *PropertyNames) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		for key := range obj {
			subCtx := NewSchemaContextFromSource(*schCtx)
			if schCtx.BaseRelativeLocation != nil {
				if newPtr, err := schCtx.BaseRelativeLocation.Descendant("propertyNames"); err == nil {
					subCtx.BaseRelativeLocation = &newPtr
				}
			}
			if newPtr, err := schCtx.RelativeLocation.Descendant("propertyNames"); err == nil {
				subCtx.RelativeLocation = &newPtr
			}
			if newPtr, err := schCtx.InstanceLocation.Descendant(key); err == nil {
				subCtx.InstanceLocation = &newPtr
			}
			subCtx.Instance = key
			(*Schema)(p).ValidateFromContext(subCtx, errs)
		}
	}
}

func (p PropertyNames) JSONProp(name string) interface{} {
	return Schema(p).JSONProp(name)
}

func (p PropertyNames) JSONChildren() (res map[string]JSONPather) {
	return Schema(p).JSONChildren()
}

func (p *PropertyNames) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*p = PropertyNames(sch)
	return nil
}

func (p PropertyNames) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(p))
}

// //
// // Dependencies
// //

// type Dependencies map[string]Dependency

// func NewDependencies() Keyword {
// 	return &Dependencies{}
// }

// func (d Dependencies) Validate(propPath string, data interface{}, errs *[]KeyError) {}

// func (d *Dependencies) Register(uri string, registry *SchemaRegistry) {
// 	for _, v := range *d {
// 		v.schema.Register(uri, registry)
// 	}
// }

// func (d *Dependencies) Resolve(pointer jptr.Pointer, uri string) *Schema {
// 	// TODO: implement this
// 	return nil
// }

// func (d *Dependencies) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
// 	for _, v := range *d {
// 		subCtx := NewSchemaContextFromSource(*schCtx)
// 		if schCtx.BaseRelativeLocation != nil {
// 			if newPtr, err := schCtx.BaseRelativeLocation.Descendant("dependencies"); err == nil {
// 				subCtx.BaseRelativeLocation = &newPtr
// 			}
// 		}
// 		if newPtr, err := schCtx.RelativeLocation.Descendant("dependencies"); err == nil {
// 			subCtx.RelativeLocation = &newPtr
// 		}
// 		v.ValidateFromContext(subCtx, errs)
// 	}
// 	// jp, err := jsonpointer.Parse(propPath)
// 	// if err != nil {
// 	// 	AddError(errs, propPath, nil, "invalid property path")
// 	// 	return
// 	// }

// 	// if obj, ok := data.(map[string]interface{}); ok {
// 	// 	for key, val := range d {
// 	// 		if obj[key] != nil {
// 	// 			d, _ := jp.Descendant(key)
// 	// 			val.Validate(d.String(), obj, errs)
// 	// 		}
// 	// 	}
// 	// }
// 	// return
// }

// func (d Dependencies) JSONProp(name string) interface{} {
// 	return d[name]
// }

//
// DependentSchemas
//

type DependentSchemas map[string]SchemaDependency

func NewDependentSchemas() Keyword {
	return &DependentSchemas{}
}

func (d DependentSchemas) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (d *DependentSchemas) Register(uri string, registry *SchemaRegistry) {
	for _, v := range *d {
		v.schema.Register(uri, registry)
	}
}

func (d *DependentSchemas) Resolve(pointer jptr.Pointer, uri string) *Schema {
	// TODO: implement this
	return nil
}

func (d *DependentSchemas) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	for _, v := range *d {
		subCtx := NewSchemaContextFromSource(*schCtx)
		if schCtx.BaseRelativeLocation != nil {
			if newPtr, err := schCtx.BaseRelativeLocation.Descendant("dependentSchemas"); err == nil {
				subCtx.BaseRelativeLocation = &newPtr
			}
		}
		if newPtr, err := schCtx.RelativeLocation.Descendant("dependentSchemas"); err == nil {
			subCtx.RelativeLocation = &newPtr
		}
		subCtx.Misc["dependencyParent"] = "dependentSchemas"
		v.ValidateFromContext(subCtx, errs)
	}
}

type _dependentSchemas map[string]Schema

func (d *DependentSchemas) UnmarshalJSON(data []byte) error {
	_d := _dependentSchemas{}
	if err := json.Unmarshal(data, &_d); err != nil {
		return err
	}
	ds := DependentSchemas{}
	for k, v := range _d {
		sch := Schema(v)
		ds[k] = SchemaDependency{
			schema: &sch,
			prop:   k,
		}
	}
	*d = ds
	return nil
}

func (d DependentSchemas) JSONProp(name string) interface{} {
	return d[name]
}

func (d DependentSchemas) JSONChildren() (r map[string]JSONPather) {
	r = map[string]JSONPather{}
	for key, val := range d {
		r[key] = val
	}
	return
}

//
// SchemaDependency
//

type SchemaDependency struct {
	schema *Schema
	prop   string
}

func (d SchemaDependency) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (d *SchemaDependency) Register(uri string, registry *SchemaRegistry) {
	d.schema.Register(uri, registry)
}

func (d *SchemaDependency) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return d.schema.Resolve(pointer, uri)
}

func (d *SchemaDependency) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if data, ok := schCtx.Instance.(map[string]interface{}); !ok {
		return
	} else {
		if _, okProp := data[d.prop]; !okProp {
			return
		}
	}
	subCtx := NewSchemaContextFromSource(*schCtx)
	if schCtx.BaseRelativeLocation != nil {
		if newPtr, err := schCtx.BaseRelativeLocation.Descendant(d.prop); err == nil {
			subCtx.BaseRelativeLocation = &newPtr
		}
	}
	if newPtr, err := schCtx.RelativeLocation.Descendant(d.prop); err == nil {
		subCtx.RelativeLocation = &newPtr
	}
	d.schema.ValidateFromContext(subCtx, errs)
}

func (d SchemaDependency) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.schema)
}

func (d SchemaDependency) JSONProp(name string) interface{} {
	return d.schema.JSONProp(name)
}

// //
// // Dependency
// //

// type Dependency struct {
// 	schema *Schema
// 	props  []string
// }

// func NewDependency() Keyword {
// 	return &Dependency{}
// }

// func (d Dependency) Validate(propPath string, data interface{}, errs *[]KeyError) {}

// func (d *Dependency) Register(uri string, registry *SchemaRegistry) {
// 	d.schema.Register(uri, registry)
// }

// func (d *Dependency) Resolve(pointer jptr.Pointer, uri string) *Schema {
// 	// TODO: implement this
// 	return nil
// }

// func (d *Dependency) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
// 	if d.schema != nil {
// 		d.schema.ValidateFromContext(schCtx, errs)
// 	}
// }

// func (d *Dependency) UnmarshalJSON(data []byte) error {
// 	props := []string{}
// 	if err := json.Unmarshal(data, &props); err == nil {
// 		*d = Dependency{props: props}
// 		return nil
// 	}
// 	sch := &Schema{}
// 	err := json.Unmarshal(data, sch)

// 	if err == nil {
// 		*d = Dependency{schema: sch}
// 	}
// 	return err
// }

// func (d Dependency) MarshalJSON() ([]byte, error) {
// 	if d.schema != nil {
// 		return json.Marshal(d.schema)
// 	}
// 	return json.Marshal(d.props)
// }
