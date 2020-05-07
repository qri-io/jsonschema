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
	if pointer == nil {
		return nil
	}
	current := pointer.Head()
	if current == nil {
		return nil
	}

	if schema, ok := (*p)[*current]; ok {
		return schema.Resolve(pointer.Tail(), uri)
	}

	return nil
}

func (p Properties) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Properties] Validating")
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		subCtx := NewSchemaContextFromSourceClean(*schCtx)
		for key, _ := range p {
			if obj[key] != nil {
				if _, ok := schCtx.Local.keywords["additionalProperties"]; ok {
					schCtx.EvaluatedPropertyNames[key] = true
					schCtx.LocalEvaluatedPropertyNames[key] = true
				}
				subCtx.ClearContext()
				if schCtx.BaseRelativeLocation != nil {
					newPtr := schCtx.BaseRelativeLocation.RawDescendant("properties", key)
					subCtx.BaseRelativeLocation = &newPtr
				}
				newPtr := schCtx.RelativeLocation.RawDescendant("properties", key)
				subCtx.RelativeLocation = &newPtr
				newPtr = schCtx.InstanceLocation.RawDescendant(key)
				subCtx.InstanceLocation = &newPtr

				subCtx.Instance = obj[key]
				errCountBefore := len(*errs)
				p[key].ValidateFromContext(subCtx, errs)
				errCountAfter := len(*errs)
				if _, ok := schCtx.Local.keywords["additionalProperties"]; ok && errCountBefore == errCountAfter {
					JoinSets(&schCtx.EvaluatedPropertyNames, subCtx.EvaluatedPropertyNames)
					JoinSets(&schCtx.LocalEvaluatedPropertyNames, subCtx.LocalEvaluatedPropertyNames)
				}
			}
		}
	}
}

func (p Properties) JSONProp(name string) interface{} {
	return p[name]
}

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

func NewRequired() Keyword {
	return &Required{}
}

func (r Required) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (r *Required) Register(uri string, registry *SchemaRegistry) {}

func (r *Required) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (r Required) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Required] Validating")
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		for _, key := range r {
			if val, ok := obj[key]; val == nil && !ok {
				AddErrorCtx(errs, schCtx, fmt.Sprintf(`"%s" value is required`, key))
			}
		}
	}
}

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
	SchemaDebug("[MaxProperties] Validating")
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
	SchemaDebug("[MinProperties] Validating")
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
	if pointer == nil {
		return nil
	}
	current := pointer.Head()
	if current == nil {
		return nil
	}

	patProp := &patternSchema{}

	for _, v := range *p {
		if v.key == *current {
			patProp = &v
			break
		}
	}

	return patProp.schema.Resolve(pointer.Tail(), uri)
}

func (p PatternProperties) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[PatternProperties] Validating")
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		for key, val := range obj {
			for _, ptn := range p {
				if ptn.re.Match([]byte(key)) {
					schCtx.EvaluatedPropertyNames[key] = true
					schCtx.LocalEvaluatedPropertyNames[key] = true
					subCtx := NewSchemaContextFromSource(*schCtx)
					if schCtx.BaseRelativeLocation != nil {
						newPtr := schCtx.BaseRelativeLocation.RawDescendant("patternProperties", key)
						subCtx.BaseRelativeLocation = &newPtr
					}
					newPtr := schCtx.RelativeLocation.RawDescendant("patternProperties", key)
					subCtx.RelativeLocation = &newPtr
					newPtr = schCtx.InstanceLocation.RawDescendant(key)
					subCtx.InstanceLocation = &newPtr

					subCtx.Instance = val
					errCountBefore := len(*errs)
					ptn.schema.ValidateFromContext(subCtx, errs)
					errCountAfter := len(*errs)

					if errCountBefore == errCountAfter {
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
	SchemaDebug("[AdditionalProperties] Validating")
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		subCtx := NewSchemaContextFromSourceClean(*schCtx)
		if schCtx.BaseRelativeLocation != nil {
			newPtr := schCtx.BaseRelativeLocation.RawDescendant("additionalProperties")
			subCtx.BaseRelativeLocation = &newPtr
		}
		newPtr := schCtx.RelativeLocation.RawDescendant("additionalProperties")
		subCtx.RelativeLocation = &newPtr
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
			newPtr = schCtx.InstanceLocation.RawDescendant(key)
			subCtx.InstanceLocation = &newPtr

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
	SchemaDebug("[PropertyNames] Validating")
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
	if pointer == nil {
		return nil
	}
	current := pointer.Head()
	if current == nil {
		return nil
	}

	if schema, ok := (*d)[*current]; ok {
		return schema.Resolve(pointer.Tail(), uri)
	}

	return nil
}

func (d *DependentSchemas) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[DependentSchemas] Validating")
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
	SchemaDebug("[SchemaDependency] Validating")
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

//
// DependentRequired
//

type DependentRequired map[string]PropertyDependency

func NewDependentRequired() Keyword {
	return &DependentRequired{}
}

func (d DependentRequired) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (d *DependentRequired) Register(uri string, registry *SchemaRegistry) {}

func (d *DependentRequired) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (d *DependentRequired) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[DependentRequired] Validating")
	for _, prop := range *d {
		subCtx := NewSchemaContextFromSource(*schCtx)
		if schCtx.BaseRelativeLocation != nil {
			if newPtr, err := schCtx.BaseRelativeLocation.Descendant("dependentRequired"); err == nil {
				subCtx.BaseRelativeLocation = &newPtr
			}
		}
		if newPtr, err := schCtx.RelativeLocation.Descendant("dependentRequired"); err == nil {
			subCtx.RelativeLocation = &newPtr
		}
		subCtx.Misc["dependencyParent"] = "dependentRequired"
		prop.ValidateFromContext(subCtx, errs)
	}
}

type _dependentRequired map[string][]string

func (d *DependentRequired) UnmarshalJSON(data []byte) error {
	_d := _dependentRequired{}
	if err := json.Unmarshal(data, &_d); err != nil {
		return err
	}
	dr := DependentRequired{}
	for k, v := range _d {
		dr[k] = PropertyDependency{
			dependencies: v,
			prop:         k,
		}
	}
	*d = dr
	return nil
}

func (d DependentRequired) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{}
	for key, prop := range d {
		obj[key] = prop.dependencies
	}
	return json.Marshal(obj)
}

func (d DependentRequired) JSONProp(name string) interface{} {
	return d[name]
}

func (d DependentRequired) JSONChildren() (r map[string]JSONPather) {
	r = map[string]JSONPather{}
	for key, val := range d {
		r[key] = val
	}
	return
}

//
// PropertyDependency
//

type PropertyDependency struct {
	dependencies []string
	prop         string
}

func (p PropertyDependency) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (p *PropertyDependency) Register(uri string, registry *SchemaRegistry) {}

func (p *PropertyDependency) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (p *PropertyDependency) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[PropertyDependency] Validating")
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		if obj[p.prop] == nil {
			return
		}
		for _, dep := range p.dependencies {
			if obj[dep] == nil {
				AddErrorCtx(errs, schCtx, fmt.Sprintf(`"%s" property is required`, dep))
			}
		}
	}
}

func (p PropertyDependency) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(p.dependencies) || idx < 0 {
		return nil
	}
	return p.dependencies[idx]
}
