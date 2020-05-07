package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	jptr "github.com/qri-io/jsonpointer"
)

//
// Items
//

type Items struct {
	single  bool
	Schemas []*Schema
}

func NewItems() Keyword {
	return &Items{}
}

func (it Items) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (it *Items) Register(uri string, registry *SchemaRegistry) {
	for _, v := range it.Schemas {
		v.Register(uri, registry)
	}
}

func (it *Items) Resolve(pointer jptr.Pointer, uri string) *Schema {
	if pointer == nil {
		return nil
	}
	current := pointer.Head()
	if current == nil {
		return nil
	}

	pos, err := strconv.Atoi(*current)
	if err != nil {
		return nil
	}

	if pos < 0 || pos >= len(it.Schemas) {
		return nil
	}

	return it.Schemas[pos].Resolve(pointer.Tail(), uri)

	return nil
}

func (it Items) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Items] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if it.single {
			subCtx := NewSchemaContextFromSourceClean(*schCtx)
			if schCtx.BaseRelativeLocation != nil {
				newPtr := schCtx.BaseRelativeLocation.RawDescendant("items")
				subCtx.BaseRelativeLocation = &newPtr
			}
			newPtr := schCtx.RelativeLocation.RawDescendant("items")
			subCtx.RelativeLocation = &newPtr
			for i, elem := range arr {
				if _, ok := schCtx.Local.keywords["additionalItems"]; ok {
					schCtx.EvaluatedPropertyNames["0"] = true
					schCtx.LocalEvaluatedPropertyNames["0"] = true
					if schCtx.LastEvaluatedIndex < i {
						schCtx.LastEvaluatedIndex = i
					}
					if schCtx.LocalLastEvaluatedIndex < i {
						schCtx.LocalLastEvaluatedIndex = i
					}
				}
				subCtx.ClearContext()
				newPtr = schCtx.InstanceLocation.RawDescendant(strconv.Itoa(i))
				subCtx.InstanceLocation = &newPtr
				subCtx.Instance = elem
				it.Schemas[0].ValidateFromContext(subCtx, errs)
				if _, ok := schCtx.Local.keywords["additionalItems"]; ok {
					// TODO(arqu): this might clash with additionalProperties
					// should separate items out
					JoinSets(&schCtx.EvaluatedPropertyNames, subCtx.EvaluatedPropertyNames)
					JoinSets(&schCtx.LocalEvaluatedPropertyNames, subCtx.LocalEvaluatedPropertyNames)
				}
			}
		} else {
			subCtx := NewSchemaContextFromSourceClean(*schCtx)
			for i, vs := range it.Schemas {
				if i < len(arr) {
					if _, ok := schCtx.Local.keywords["additionalItems"]; ok {
						schCtx.EvaluatedPropertyNames[strconv.Itoa(i)] = true
						schCtx.LocalEvaluatedPropertyNames[strconv.Itoa(i)] = true
						if schCtx.LastEvaluatedIndex < i {
							schCtx.LastEvaluatedIndex = i
						}
						if schCtx.LocalLastEvaluatedIndex < i {
							schCtx.LocalLastEvaluatedIndex = i
						}
					}
					subCtx.ClearContext()
					if schCtx.BaseRelativeLocation != nil {
						newPtr := schCtx.BaseRelativeLocation.RawDescendant("items", strconv.Itoa(i))
						subCtx.BaseRelativeLocation = &newPtr
					}
					newPtr := schCtx.RelativeLocation.RawDescendant("items", strconv.Itoa(i))
					subCtx.RelativeLocation = &newPtr
					newPtr = schCtx.InstanceLocation.RawDescendant(strconv.Itoa(i))
					subCtx.InstanceLocation = &newPtr

					subCtx.Instance = arr[i]
					vs.ValidateFromContext(subCtx, errs)
					if _, ok := schCtx.Local.keywords["additionalItems"]; ok {
						JoinSets(&schCtx.EvaluatedPropertyNames, subCtx.EvaluatedPropertyNames)
						JoinSets(&schCtx.LocalEvaluatedPropertyNames, subCtx.LocalEvaluatedPropertyNames)
					}
				}
			}
		}
	}
}

func (it Items) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(it.Schemas) || idx < 0 {
		return nil
	}
	return it.Schemas[idx]
}

func (it Items) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, sch := range it.Schemas {
		res[strconv.Itoa(i)] = sch
	}
	return
}

func (it *Items) UnmarshalJSON(data []byte) error {
	s := &Schema{}
	if err := json.Unmarshal(data, s); err == nil {
		*it = Items{single: true, Schemas: []*Schema{s}}
		return nil
	}
	ss := []*Schema{}
	if err := json.Unmarshal(data, &ss); err != nil {
		return err
	}
	*it = Items{Schemas: ss}
	return nil
}

func (it Items) MarshalJSON() ([]byte, error) {
	if it.single {
		return json.Marshal(it.Schemas[0])
	}
	return json.Marshal([]*Schema(it.Schemas))
}

//
// MaxItems
//

type MaxItems int

func NewMaxItems() Keyword {
	return new(MaxItems)
}

func (m MaxItems) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *MaxItems) Register(uri string, registry *SchemaRegistry) {}

func (m *MaxItems) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m MaxItems) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[MaxItems] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if len(arr) > int(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("array length %d exceeds %d max", len(arr), m))
			return
		}
	}
}

//
// MinItems
//

type MinItems int

func NewMinItems() Keyword {
	return new(MinItems)
}

func (m MinItems) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *MinItems) Register(uri string, registry *SchemaRegistry) {}

func (m *MinItems) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m MinItems) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[MinItems] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if len(arr) < int(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("array length %d below %d minimum items", len(arr), m))
			return
		}
	}
}

//
// UniqueItems
//

type UniqueItems bool

func NewUniqueItems() Keyword {
	return new(UniqueItems)
}

func (u UniqueItems) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (u *UniqueItems) Register(uri string, registry *SchemaRegistry) {}

func (u *UniqueItems) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (u UniqueItems) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[UniqueItems] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		found := []interface{}{}
		for _, elem := range arr {
			for _, f := range found {
				if reflect.DeepEqual(f, elem) {
					AddErrorCtx(errs, schCtx, fmt.Sprintf("array items must be unique. duplicated entry: %v", elem))
					return
				}
			}
			found = append(found, elem)
		}
	}
}

//
// Contains
//

type Contains Schema

func NewContains() Keyword {
	return &Contains{}
}

func (c Contains) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (c *Contains) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(c).Register(uri, registry)
}

func (c *Contains) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(c).Resolve(pointer, uri)
}

func (c *Contains) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Contains] Validating")
	v := Schema(*c)
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		valid := false
		matchCount := 0
		subCtx := NewSchemaContextFromSourceClean(*schCtx)
		if schCtx.BaseRelativeLocation != nil {
			newPtr := schCtx.BaseRelativeLocation.RawDescendant("contains")
			subCtx.BaseRelativeLocation = &newPtr
		}
		newPtr := schCtx.RelativeLocation.RawDescendant("contains")
		subCtx.RelativeLocation = &newPtr
		for i, elem := range arr {
			subCtx.ClearContext()
			newPtr = schCtx.InstanceLocation.RawDescendant(strconv.Itoa(i))
			subCtx.InstanceLocation = &newPtr
			subCtx.Instance = elem
			test := &[]KeyError{}
			v.ValidateFromContext(subCtx, test)
			if len(*test) == 0 {
				valid = true
				matchCount++
			}
		}
		if valid {
			schCtx.Misc["containsCount"] = matchCount
		} else {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("must contain at least one of: %v", c))
		}
	}
}

func (c Contains) JSONProp(name string) interface{} {
	return Schema(c).JSONProp(name)
}

func (c Contains) JSONChildren() (res map[string]JSONPather) {
	return Schema(c).JSONChildren()
}

func (c *Contains) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*c = Contains(sch)
	return nil
}

//
// MaxContains
//

type MaxContains int

func NewMaxContains() Keyword {
	return new(MaxContains)
}

func (m MaxContains) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *MaxContains) Register(uri string, registry *SchemaRegistry) {}

func (m *MaxContains) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m MaxContains) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[MaxContains] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if containsCount, ok := schCtx.Misc["containsCount"]; ok {
			if containsCount.(int) > int(m) {
				AddErrorCtx(errs, schCtx, fmt.Sprintf("contained items %d exceeds %d max", len(arr), m))
			}
		}
	}
}

//
// MinContains
//

type MinContains int

func NewMinContains() Keyword {
	return new(MinContains)
}

func (m MinContains) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *MinContains) Register(uri string, registry *SchemaRegistry) {}

func (m *MinContains) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m MinContains) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[MinContains] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if containsCount, ok := schCtx.Misc["containsCount"]; ok {
			if containsCount.(int) < int(m) {
				AddErrorCtx(errs, schCtx, fmt.Sprintf("contained items %d bellow %d min", len(arr), m))
			}
		}
	}
}

//
// AdditionalItems
//

type AdditionalItems Schema

func NewAdditionalItems() Keyword {
	return &AdditionalItems{}
}

func (ai AdditionalItems) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (ai *AdditionalItems) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(ai).Register(uri, registry)
}

func (ai *AdditionalItems) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(ai).Resolve(pointer, uri)
}

func (ai *AdditionalItems) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[AdditionalItems] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if schCtx.LastEvaluatedIndex > -1 && schCtx.LastEvaluatedIndex < len(arr) {
			for i := schCtx.LastEvaluatedIndex + 1; i < len(arr); i++ {
				if ai.schemaType == schemaTypeFalse {
					AddErrorCtx(errs, schCtx, "additional items are not allowed")
					return
				}
				subCtx := NewSchemaContextFromSourceClean(*schCtx)
				if schCtx.BaseRelativeLocation != nil {
					newPtr := schCtx.BaseRelativeLocation.RawDescendant("additionalItems")
					subCtx.BaseRelativeLocation = &newPtr
				}
				newPtr := schCtx.RelativeLocation.RawDescendant("additionalItems")
				subCtx.RelativeLocation = &newPtr
				newPtr = schCtx.InstanceLocation.RawDescendant(strconv.Itoa(i))
				subCtx.InstanceLocation = &newPtr

				subCtx.Instance = arr[i]
				(*Schema)(ai).ValidateFromContext(subCtx, errs)
				JoinSets(&schCtx.EvaluatedPropertyNames, subCtx.EvaluatedPropertyNames)
				JoinSets(&schCtx.LocalEvaluatedPropertyNames, subCtx.LocalEvaluatedPropertyNames)
			}
		}
	}
}

func (ai *AdditionalItems) UnmarshalJSON(data []byte) error {
	sch := &Schema{}
	if err := json.Unmarshal(data, sch); err != nil {
		return err
	}
	*ai = (AdditionalItems)(*sch)
	return nil
}
