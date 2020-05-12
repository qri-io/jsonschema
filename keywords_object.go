package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	jptr "github.com/qri-io/jsonpointer"
)

// Properties defines the properties JSON Schema keyword
type Properties map[string]*Schema

// NewProperties allocates a new Properties keyword
func NewProperties() Keyword {
	return &Properties{}
}

// Register implements the Keyword interface for Properties
func (p *Properties) Register(uri string, registry *SchemaRegistry) {
	for _, v := range *p {
		v.Register(uri, registry)
	}
}

// Resolve implements the Keyword interface for Properties
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

// ValidateFromContext implements the Keyword interface for Properties
func (p Properties) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Properties] Validating")
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		subCtx := NewSchemaContextFromSource(*schCtx)
		for key := range p {
			if obj[key] != nil {
				(*schCtx.EvaluatedPropertyNames)[key] = true
				(*schCtx.LocalEvaluatedPropertyNames)[key] = true
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
				if errCountBefore == errCountAfter {
					JoinSets(schCtx.EvaluatedPropertyNames, *subCtx.EvaluatedPropertyNames)
					JoinSets(schCtx.LocalEvaluatedPropertyNames, *subCtx.LocalEvaluatedPropertyNames)
				}
			}
		}
	}
}

// JSONProp implements the JSONPather for Properties
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

// Required defines the required JSON Schema keyword
type Required []string

// NewRequired allocates a new Required keyword
func NewRequired() Keyword {
	return &Required{}
}

// Register implements the Keyword interface for Required
func (r *Required) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for Required
func (r *Required) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ValidateFromContext implements the Keyword interface for Required
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

// JSONProp implements the JSONPather for Required
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

// MaxProperties defines the maxProperties JSON Schema keyword
type MaxProperties int

// NewMaxProperties allocates a new MaxProperties keyword
func NewMaxProperties() Keyword {
	return new(MaxProperties)
}

// Register implements the Keyword interface for MaxProperties
func (m *MaxProperties) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for MaxProperties
func (m *MaxProperties) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ValidateFromContext implements the Keyword interface for MaxProperties
func (m MaxProperties) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[MaxProperties] Validating")
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		if len(obj) > int(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("%d object Properties exceed %d maximum", len(obj), m))
		}
	}
}

// MinProperties defines the minProperties JSON Schema keyword
type MinProperties int

// NewMinProperties allocates a new MinProperties keyword
func NewMinProperties() Keyword {
	return new(MinProperties)
}

// Register implements the Keyword interface for MinProperties
func (m *MinProperties) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for MinProperties
func (m *MinProperties) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ValidateFromContext implements the Keyword interface for MinProperties
func (m MinProperties) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[MinProperties] Validating")
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		if len(obj) < int(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("%d object Properties below %d minimum", len(obj), m))
		}
	}
}

// PatternProperties defines the patternProperties JSON Schema keyword
type PatternProperties []patternSchema

// NewPatternProperties allocates a new PatternProperties keyword
func NewPatternProperties() Keyword {
	return &PatternProperties{}
}

type patternSchema struct {
	key    string
	re     *regexp.Regexp
	schema *Schema
}

// Register implements the Keyword interface for PatternProperties
func (p *PatternProperties) Register(uri string, registry *SchemaRegistry) {
	for _, v := range *p {
		v.schema.Register(uri, registry)
	}
}

// Resolve implements the Keyword interface for PatternProperties
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

// ValidateFromContext implements the Keyword interface for PatternProperties
func (p PatternProperties) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[PatternProperties] Validating")
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		for key, val := range obj {
			for _, ptn := range p {
				if ptn.re.Match([]byte(key)) {
					(*schCtx.EvaluatedPropertyNames)[key] = true
					(*schCtx.LocalEvaluatedPropertyNames)[key] = true
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
						JoinSets(schCtx.EvaluatedPropertyNames, *subCtx.EvaluatedPropertyNames)
						JoinSets(schCtx.LocalEvaluatedPropertyNames, *subCtx.LocalEvaluatedPropertyNames)
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

// JSONProp implements the JSONPather for PatternProperties
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

// MarshalJSON implements the json.Marshaler interface for PatternProperties
func (p PatternProperties) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{}
	for _, prop := range p {
		obj[prop.key] = prop.schema
	}
	return json.Marshal(obj)
}

// AdditionalProperties defines the additionalProperties JSON Schema keyword
type AdditionalProperties Schema

// NewAdditionalProperties allocates a new AdditionalProperties keyword
func NewAdditionalProperties() Keyword {
	return &AdditionalProperties{}
}

// Register implements the Keyword interface for AdditionalProperties
func (ap *AdditionalProperties) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(ap).Register(uri, registry)
}

// Resolve implements the Keyword interface for AdditionalProperties
func (ap *AdditionalProperties) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(ap).Resolve(pointer, uri)
}

// ValidateFromContext implements the Keyword interface for AdditionalProperties
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
			if _, ok := (*schCtx.LocalEvaluatedPropertyNames)[key]; ok {
				continue
			}
			if ap.schemaType == schemaTypeFalse {
				AddErrorCtx(errs, schCtx, "additional properties are not allowed")
				return
			}
			(*schCtx.EvaluatedPropertyNames)[key] = true
			(*schCtx.LocalEvaluatedPropertyNames)[key] = true
			subCtx.ClearContext()
			newPtr = schCtx.InstanceLocation.RawDescendant(key)
			subCtx.InstanceLocation = &newPtr

			subCtx.Instance = obj[key]
			(*Schema)(ap).ValidateFromContext(subCtx, errs)
			JoinSets(schCtx.EvaluatedPropertyNames, *subCtx.EvaluatedPropertyNames)
			JoinSets(schCtx.LocalEvaluatedPropertyNames, *subCtx.LocalEvaluatedPropertyNames)
		}
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface for AdditionalProperties
func (ap *AdditionalProperties) UnmarshalJSON(data []byte) error {
	sch := &Schema{}
	if err := json.Unmarshal(data, sch); err != nil {
		return err
	}
	*ap = (AdditionalProperties)(*sch)
	return nil
}

// PropertyNames defines the propertyNames JSON Schema keyword
type PropertyNames Schema

// NewPropertyNames allocates a new PropertyNames keyword
func NewPropertyNames() Keyword {
	return &PropertyNames{}
}

// Register implements the Keyword interface for PropertyNames
func (p *PropertyNames) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(p).Register(uri, registry)
}

// Resolve implements the Keyword interface for PropertyNames
func (p *PropertyNames) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(p).Resolve(pointer, uri)
}

// ValidateFromContext implements the Keyword interface for PropertyNames
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

// JSONProp implements the JSONPather for PropertyNames
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

// MarshalJSON implements the json.Marshaler interface for PropertyNames
func (p PropertyNames) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(p))
}

// DependentSchemas defines the dependentSchemas JSON Schema keyword
type DependentSchemas map[string]SchemaDependency

// NewDependentSchemas allocates a new DependentSchemas keyword
func NewDependentSchemas() Keyword {
	return &DependentSchemas{}
}

// Register implements the Keyword interface for DependentSchemas
func (d *DependentSchemas) Register(uri string, registry *SchemaRegistry) {
	for _, v := range *d {
		v.schema.Register(uri, registry)
	}
}

// Resolve implements the Keyword interface for DependentSchemas
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

// ValidateFromContext implements the Keyword interface for DependentSchemas
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

// UnmarshalJSON implements the json.Unmarshaler interface for DependentSchemas
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

// JSONProp implements the JSONPather for DependentSchemas
func (d DependentSchemas) JSONProp(name string) interface{} {
	return d[name]
}

// JSONChildren implements the JSONContainer interface for DependentSchemas
func (d DependentSchemas) JSONChildren() (r map[string]JSONPather) {
	r = map[string]JSONPather{}
	for key, val := range d {
		r[key] = val
	}
	return
}

// SchemaDependency is the internal representation of a dependent schema
type SchemaDependency struct {
	schema *Schema
	prop   string
}

// Register implements the Keyword interface for SchemaDependency
func (d *SchemaDependency) Register(uri string, registry *SchemaRegistry) {
	d.schema.Register(uri, registry)
}

// Resolve implements the Keyword interface for SchemaDependency
func (d *SchemaDependency) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return d.schema.Resolve(pointer, uri)
}

// ValidateFromContext implements the Keyword interface for SchemaDependency
func (d *SchemaDependency) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[SchemaDependency] Validating")
	data := map[string]interface{}{}
	ok := false
	if data, ok = schCtx.Instance.(map[string]interface{}); !ok {
		return
	}
	if _, okProp := data[d.prop]; !okProp {
		return
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

// MarshalJSON implements the json.Marshaler interface for SchemaDependency
func (d SchemaDependency) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.schema)
}

// JSONProp implements the JSONPather for SchemaDependency
func (d SchemaDependency) JSONProp(name string) interface{} {
	return d.schema.JSONProp(name)
}

// DependentRequired defines the dependentRequired JSON Schema keyword
type DependentRequired map[string]PropertyDependency

// NewDependentRequired allocates a new DependentRequired keyword
func NewDependentRequired() Keyword {
	return &DependentRequired{}
}

// Register implements the Keyword interface for DependentRequired
func (d *DependentRequired) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for DependentRequired
func (d *DependentRequired) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ValidateFromContext implements the Keyword interface for DependentRequired
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

// UnmarshalJSON implements the json.Unmarshaler interface for DependentRequired
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

// MarshalJSON implements the json.Marshaler interface for DependentRequired
func (d DependentRequired) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{}
	for key, prop := range d {
		obj[key] = prop.dependencies
	}
	return json.Marshal(obj)
}

// JSONProp implements the JSONPather for DependentRequired
func (d DependentRequired) JSONProp(name string) interface{} {
	return d[name]
}

// JSONChildren implements the JSONContainer interface for DependentRequired
func (d DependentRequired) JSONChildren() (r map[string]JSONPather) {
	r = map[string]JSONPather{}
	for key, val := range d {
		r[key] = val
	}
	return
}

// PropertyDependency is the internal representation of a dependent property
type PropertyDependency struct {
	dependencies []string
	prop         string
}

// Register implements the Keyword interface for PropertyDependency
func (p *PropertyDependency) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for PropertyDependency
func (p *PropertyDependency) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ValidateFromContext implements the Keyword interface for PropertyDependency
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

// JSONProp implements the JSONPather for PropertyDependency
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

// UnevaluatedProperties defines the unevaluatedProperties JSON Schema keyword
type UnevaluatedProperties Schema

// NewUnevaluatedProperties allocates a new UnevaluatedProperties keyword
func NewUnevaluatedProperties() Keyword {
	return &UnevaluatedProperties{}
}

// Register implements the Keyword interface for UnevaluatedProperties
func (up *UnevaluatedProperties) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(up).Register(uri, registry)
}

// Resolve implements the Keyword interface for UnevaluatedProperties
func (up *UnevaluatedProperties) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(up).Resolve(pointer, uri)
}

// ValidateFromContext implements the Keyword interface for UnevaluatedProperties
func (up *UnevaluatedProperties) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[UnevaluatedProperties] Validating")
	if obj, ok := schCtx.Instance.(map[string]interface{}); ok {
		subCtx := NewSchemaContextFromSourceClean(*schCtx)
		if schCtx.BaseRelativeLocation != nil {
			newPtr := schCtx.BaseRelativeLocation.RawDescendant("unevaluatedProperties")
			subCtx.BaseRelativeLocation = &newPtr
		}
		newPtr := schCtx.RelativeLocation.RawDescendant("unevaluatedProperties")
		subCtx.RelativeLocation = &newPtr
		for key := range obj {
			if _, ok := (*schCtx.EvaluatedPropertyNames)[key]; ok {
				continue
			}
			if up.schemaType == schemaTypeFalse {
				AddErrorCtx(errs, schCtx, "unevaluated properties are not allowed")
				return
			}
			newPtr = schCtx.InstanceLocation.RawDescendant(key)
			subCtx.InstanceLocation = &newPtr

			subCtx.Instance = obj[key]
			(*Schema)(up).ValidateFromContext(subCtx, errs)
		}
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface for UnevaluatedProperties
func (up *UnevaluatedProperties) UnmarshalJSON(data []byte) error {
	sch := &Schema{}
	if err := json.Unmarshal(data, sch); err != nil {
		return err
	}
	*up = (UnevaluatedProperties)(*sch)
	return nil
}
