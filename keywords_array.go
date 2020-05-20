package jsonschema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	jptr "github.com/qri-io/jsonpointer"
)

// Items defines the items JSON Schema keyword
type Items struct {
	single  bool
	Schemas []*Schema
}

// NewItems allocates a new Items keyword
func NewItems() Keyword {
	return &Items{}
}

// Register implements the Keyword interface for Items
func (it *Items) Register(uri string, registry *SchemaRegistry) {
	for _, v := range it.Schemas {
		v.Register(uri, registry)
	}
}

// Resolve implements the Keyword interface for Items
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

// ValidateFromContext implements the Keyword interface for Items
func (it Items) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	schemaDebug("[Items] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if it.single {
			subCtx := NewSchemaContextFromSource(*schCtx)
			if schCtx.BaseRelativeLocation != nil {
				newPtr := schCtx.BaseRelativeLocation.RawDescendant("items")
				subCtx.BaseRelativeLocation = &newPtr
			}
			newPtr := schCtx.RelativeLocation.RawDescendant("items")
			subCtx.RelativeLocation = &newPtr
			for i, elem := range arr {
				subCtx.ClearContext()
				newPtr = schCtx.InstanceLocation.RawDescendant(strconv.Itoa(i))
				subCtx.InstanceLocation = &newPtr
				subCtx.Instance = elem
				it.Schemas[0].ValidateFromContext(subCtx, errs)
				subCtx.LastEvaluatedIndex = i
				subCtx.LocalLastEvaluatedIndex = i
				// TODO(arqu): this might clash with additional/unevaluated
				// Properties/Items, should separate out
				schCtx.UpdateEvaluatedPropsAndItems(subCtx)
			}
		} else {
			subCtx := NewSchemaContextFromSource(*schCtx)
			for i, vs := range it.Schemas {
				if i < len(arr) {
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
					subCtx.LastEvaluatedIndex = i
					subCtx.LocalLastEvaluatedIndex = i
					schCtx.UpdateEvaluatedPropsAndItems(subCtx)
				}
			}
		}
	}
}

// JSONProp implements the JSONPather for Items
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

// JSONChildren implements the JSONContainer interface for Items
func (it Items) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, sch := range it.Schemas {
		res[strconv.Itoa(i)] = sch
	}
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface for Items
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

// MarshalJSON implements the json.Marshaler interface for Items
func (it Items) MarshalJSON() ([]byte, error) {
	if it.single {
		return json.Marshal(it.Schemas[0])
	}
	return json.Marshal([]*Schema(it.Schemas))
}

// MaxItems defines the maxItems JSON Schema keyword
type MaxItems int

// NewMaxItems allocates a new MaxItems keyword
func NewMaxItems() Keyword {
	return new(MaxItems)
}

// Register implements the Keyword interface for MaxItems
func (m *MaxItems) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for MaxItems
func (m *MaxItems) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ValidateFromContext implements the Keyword interface for MaxItems
func (m MaxItems) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	schemaDebug("[MaxItems] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if len(arr) > int(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("array length %d exceeds %d max", len(arr), m))
			return
		}
	}
}

// MinItems defines the minItems JSON Schema keyword
type MinItems int

// NewMinItems allocates a new MinItems keyword
func NewMinItems() Keyword {
	return new(MinItems)
}

// Register implements the Keyword interface for MinItems
func (m *MinItems) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for MinItems
func (m *MinItems) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ValidateFromContext implements the Keyword interface for MinItems
func (m MinItems) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	schemaDebug("[MinItems] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if len(arr) < int(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("array length %d below %d minimum items", len(arr), m))
			return
		}
	}
}

// UniqueItems defines the uniqueItems JSON Schema keyword
type UniqueItems bool

// NewUniqueItems allocates a new UniqueItems keyword
func NewUniqueItems() Keyword {
	return new(UniqueItems)
}

// Register implements the Keyword interface for UniqueItems
func (u *UniqueItems) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for UniqueItems
func (u *UniqueItems) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ValidateFromContext implements the Keyword interface for UniqueItems
func (u UniqueItems) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	schemaDebug("[UniqueItems] Validating")
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

// Contains defines the contains JSON Schema keyword
type Contains Schema

// NewContains allocates a new Contains keyword
func NewContains() Keyword {
	return &Contains{}
}

// Register implements the Keyword interface for Contains
func (c *Contains) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(c).Register(uri, registry)
}

// Resolve implements the Keyword interface for Contains
func (c *Contains) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(c).Resolve(pointer, uri)
}

// ValidateFromContext implements the Keyword interface for Contains
func (c *Contains) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	schemaDebug("[Contains] Validating")
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

// JSONProp implements the JSONPather for Contains
func (c Contains) JSONProp(name string) interface{} {
	return Schema(c).JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for Contains
func (c Contains) JSONChildren() (res map[string]JSONPather) {
	return Schema(c).JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for Contains
func (c *Contains) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*c = Contains(sch)
	return nil
}

// MaxContains defines the maxContains JSON Schema keyword
type MaxContains int

// NewMaxContains allocates a new MaxContains keyword
func NewMaxContains() Keyword {
	return new(MaxContains)
}

// Register implements the Keyword interface for MaxContains
func (m *MaxContains) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for MaxContains
func (m *MaxContains) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ValidateFromContext implements the Keyword interface for MaxContains
func (m MaxContains) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	schemaDebug("[MaxContains] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if containsCount, ok := schCtx.Misc["containsCount"]; ok {
			if containsCount.(int) > int(m) {
				AddErrorCtx(errs, schCtx, fmt.Sprintf("contained items %d exceeds %d max", len(arr), m))
			}
		}
	}
}

// MinContains defines the minContains JSON Schema keyword
type MinContains int

// NewMinContains allocates a new MinContains keyword
func NewMinContains() Keyword {
	return new(MinContains)
}

// Register implements the Keyword interface for MinContains
func (m *MinContains) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for MinContains
func (m *MinContains) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ValidateFromContext implements the Keyword interface for MinContains
func (m MinContains) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	schemaDebug("[MinContains] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if containsCount, ok := schCtx.Misc["containsCount"]; ok {
			if containsCount.(int) < int(m) {
				AddErrorCtx(errs, schCtx, fmt.Sprintf("contained items %d bellow %d min", len(arr), m))
			}
		}
	}
}

// AdditionalItems defines the additionalItems JSON Schema keyword
type AdditionalItems Schema

// NewAdditionalItems allocates a new AdditionalItems keyword
func NewAdditionalItems() Keyword {
	return &AdditionalItems{}
}

// Register implements the Keyword interface for AdditionalItems
func (ai *AdditionalItems) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(ai).Register(uri, registry)
}

// Resolve implements the Keyword interface for AdditionalItems
func (ai *AdditionalItems) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(ai).Resolve(pointer, uri)
}

// ValidateFromContext implements the Keyword interface for AdditionalItems
func (ai *AdditionalItems) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	schemaDebug("[AdditionalItems] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if schCtx.LastEvaluatedIndex > -1 && schCtx.LastEvaluatedIndex < len(arr) {
			for i := schCtx.LastEvaluatedIndex + 1; i < len(arr); i++ {
				if ai.schemaType == schemaTypeFalse {
					AddErrorCtx(errs, schCtx, "additional items are not allowed")
					return
				}
				subCtx := NewSchemaContextFromSourceClean(*schCtx)
				subCtx.LastEvaluatedIndex = i
				subCtx.LocalLastEvaluatedIndex = i
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
				schCtx.UpdateEvaluatedPropsAndItems(subCtx)
			}
		}
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface for AdditionalItems
func (ai *AdditionalItems) UnmarshalJSON(data []byte) error {
	sch := &Schema{}
	if err := json.Unmarshal(data, sch); err != nil {
		return err
	}
	*ai = (AdditionalItems)(*sch)
	return nil
}

// UnevaluatedItems defines the unevaluatedItems JSON Schema keyword
type UnevaluatedItems Schema

// NewUnevaluatedItems allocates a new UnevaluatedItems keyword
func NewUnevaluatedItems() Keyword {
	return &UnevaluatedItems{}
}

// Register implements the Keyword interface for UnevaluatedItems
func (ui *UnevaluatedItems) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(ui).Register(uri, registry)
}

// Resolve implements the Keyword interface for UnevaluatedItems
func (ui *UnevaluatedItems) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(ui).Resolve(pointer, uri)
}

// ValidateFromContext implements the Keyword interface for UnevaluatedItems
func (ui *UnevaluatedItems) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	schemaDebug("[UnevaluatedItems] Validating")
	if arr, ok := schCtx.Instance.([]interface{}); ok {
		if schCtx.LastEvaluatedIndex < len(arr) {
			for i := schCtx.LastEvaluatedIndex + 1; i < len(arr); i++ {
				if ui.schemaType == schemaTypeFalse {
					AddErrorCtx(errs, schCtx, "unevaluated items are not allowed")
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
				(*Schema)(ui).ValidateFromContext(subCtx, errs)
			}
		}
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface for UnevaluatedItems
func (ui *UnevaluatedItems) UnmarshalJSON(data []byte) error {
	sch := &Schema{}
	if err := json.Unmarshal(data, sch); err != nil {
		return err
	}
	*ui = (UnevaluatedItems)(*sch)
	return nil
}
