package main

import (
	"encoding/json"
	"strconv"

	jptr "github.com/qri-io/jsonpointer"
)

//
// AllOf
//

type AllOf []*Schema

func NewAllOf() Keyword {
	return &AllOf{}
}

func (a *AllOf) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (a *AllOf) Register(uri string, registry *SchemaRegistry) {
	for _, sch := range *a {
		sch.Register(uri, registry)
	}
}

func (a *AllOf) Resolve(pointer jptr.Pointer, uri string) *Schema {
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

	if pos < 0 || pos >= len(*a) {
		return nil
	}

	return (*a)[pos].Resolve(pointer.Tail(), uri)

	return nil
}

func (a *AllOf) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[AllOf] Validating")
	for i, sch := range *a {
		subCtx := NewSchemaContextFromSource(*schCtx)
		if subCtx.BaseRelativeLocation != nil {
			if newPtr, err := schCtx.BaseRelativeLocation.Descendant("allOf/" + strconv.Itoa(i)); err == nil {
				subCtx.BaseRelativeLocation = &newPtr
			}
		}
		if newPtr, err := schCtx.RelativeLocation.Descendant("allOf/" + strconv.Itoa(i)); err == nil {
			subCtx.RelativeLocation = &newPtr
		}
		sch.ValidateFromContext(subCtx, errs)
		schCtx.UpdateEvaluatedPropsAndItems(subCtx)
	}
}

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

func (a AllOf) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, sch := range a {
		res[strconv.Itoa(i)] = sch
	}
	return
}

//
// AnyOf
//

type AnyOf []*Schema

func NewAnyOf() Keyword {
	return &AnyOf{}
}

func (a *AnyOf) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (a *AnyOf) Register(uri string, registry *SchemaRegistry) {
	for _, sch := range *a {
		sch.Register(uri, registry)
	}
}

func (a *AnyOf) Resolve(pointer jptr.Pointer, uri string) *Schema {
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

	if pos < 0 || pos >= len(*a) {
		return nil
	}

	return (*a)[pos].Resolve(pointer.Tail(), uri)

	return nil
}

func (a *AnyOf) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[AnyOf] Validating")
	for i, sch := range *a {
		subCtx := NewSchemaContextFromSource(*schCtx)
		if subCtx.BaseRelativeLocation != nil {
			if newPtr, err := schCtx.BaseRelativeLocation.Descendant("anyOf/" + strconv.Itoa(i)); err == nil {
				subCtx.BaseRelativeLocation = &newPtr
			}
		}
		if newPtr, err := schCtx.RelativeLocation.Descendant("anyOf/" + strconv.Itoa(i)); err == nil {
			subCtx.RelativeLocation = &newPtr
		}
		test := &[]KeyError{}
		sch.ValidateFromContext(subCtx, test)
		if len(*test) == 0 {
			schCtx.UpdateEvaluatedPropsAndItems(subCtx)
			return
		}
	}

	AddErrorCtx(errs, schCtx, "did Not match any specified AnyOf schemas")
}

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

func (a AnyOf) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, sch := range a {
		res[strconv.Itoa(i)] = sch
	}
	return
}

//
// OneOf
//

type OneOf []*Schema

func NewOneOf() Keyword {
	return &OneOf{}
}

func (o *OneOf) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (o *OneOf) Register(uri string, registry *SchemaRegistry) {
	for _, sch := range *o {
		sch.Register(uri, registry)
	}
}

func (o *OneOf) Resolve(pointer jptr.Pointer, uri string) *Schema {
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

	if pos < 0 || pos >= len(*o) {
		return nil
	}

	return (*o)[pos].Resolve(pointer.Tail(), uri)

	return nil
}

func (o *OneOf) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[OneOf] Validating")
	matched := false
	contextCopy := NewSchemaContextFromSource(*schCtx)
	for i, sch := range *o {
		subCtx := NewSchemaContextFromSource(*schCtx)
		if subCtx.BaseRelativeLocation != nil {
			if newPtr, err := schCtx.BaseRelativeLocation.Descendant("anyOf/" + strconv.Itoa(i)); err == nil {
				subCtx.BaseRelativeLocation = &newPtr
			}
		}
		if newPtr, err := schCtx.RelativeLocation.Descendant("anyOf/" + strconv.Itoa(i)); err == nil {
			subCtx.RelativeLocation = &newPtr
		}
		test := &[]KeyError{}
		sch.ValidateFromContext(subCtx, test)
		contextCopy.UpdateEvaluatedPropsAndItems(subCtx)
		if len(*test) == 0 {
			if matched {
				AddErrorCtx(errs, schCtx, "matched more than one specified OneOf schemas")
				return
			}
			matched = true
		}
	}
	if !matched {
		AddErrorCtx(errs, schCtx, "did not match any of the specified OneOf schemas")
	} else {
		schCtx.UpdateEvaluatedPropsAndItems(contextCopy)
	}
}

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

func (o OneOf) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, sch := range o {
		res[strconv.Itoa(i)] = sch
	}
	return
}

//
// Not
//

type Not Schema

func NewNot() Keyword {
	return &Not{}
}

func (n *Not) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (n *Not) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(n).Register(uri, registry)
}

func (n *Not) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(n).Resolve(pointer, uri)
}

func (n *Not) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Not] Validating")
	subCtx := NewSchemaContextFromSource(*schCtx)
	if subCtx.BaseRelativeLocation != nil {
		if newPtr, err := schCtx.BaseRelativeLocation.Descendant("not"); err == nil {
			subCtx.BaseRelativeLocation = &newPtr
		}
	}
	if newPtr, err := schCtx.RelativeLocation.Descendant("not"); err == nil {
		subCtx.RelativeLocation = &newPtr
	}

	test := &[]KeyError{}
	sch := Schema(*n)
	sch.ValidateFromContext(subCtx, test)
	if len(*test) == 0 {
		AddErrorCtx(errs, schCtx, "result was valid, ('not') expected invalid")
	}
}

func (n Not) JSONProp(name string) interface{} {
	return Schema(n).JSONProp(name)
}

func (n Not) JSONChildren() (res map[string]JSONPather) {
	return Schema(n).JSONChildren()
}

func (n *Not) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*n = Not(sch)
	return nil
}

func (n Not) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(n))
}
