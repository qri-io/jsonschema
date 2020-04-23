package main

import (
	"encoding/json"
	jptr "github.com/qri-io/jsonpointer"
)

//
// If
//

type If Schema

func NewIf() Keyword {
	return &If{}
}

func (f *If) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (f *If) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(f).Register(uri, registry)
}

func (f *If) Resolve(pointer jptr.Pointer, uri string) *Schema {
	// TODO: check if this should be nil
	return (*Schema)(f).Resolve(pointer, uri)
}

func (f *If) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {

	thenKW := schCtx.Local.Keywords["then"]
	elseKW := schCtx.Local.Keywords["else"]

	if (thenKW == nil && elseKW == nil) {
		// no then or else for if, aborting validation
		return
	}

	subCtx := NewSchemaContextFromSource(*schCtx)
	if subCtx.BaseRelativeLocation != nil {
		if newPtr, err := schCtx.BaseRelativeLocation.Descendant("if"); err == nil {
			subCtx.BaseRelativeLocation = &newPtr
		}
	}
	if newPtr, err := schCtx.RelativeLocation.Descendant("if"); err == nil {
		subCtx.RelativeLocation = &newPtr
	}
	test := &[]KeyError{}
	sch := Schema(*f)
	sch.ValidateFromContext(subCtx, test)

	schCtx.Misc["ifResult"] = (len(*test) == 0)
}

func (f If) JSONProp(name string) interface{} {
	return Schema(f).JSONProp(name)
}

func (f If) JSONChildren() (res map[string]JSONPather) {
	return Schema(f).JSONChildren()
}

func (f *If) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*f = If(sch)
	return nil
}

func (f If) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(f))
}

//
// Then
//

type Then Schema

func NewThen() Keyword {
	return &Then{}
}

func (t *Then) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (t *Then) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(t).Register(uri, registry)
}

func (t *Then) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(t).Resolve(pointer, uri)
}

func (t *Then) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	ifResult, okIf := schCtx.Misc["ifResult"]
	if !okIf {
		// if not found
		return
	}
	if !(ifResult.(bool)) {
		// if was false
		return
	}

	subCtx := NewSchemaContextFromSource(*schCtx)
	if subCtx.BaseRelativeLocation != nil {
		if newPtr, err := schCtx.BaseRelativeLocation.Descendant("then"); err == nil {
			subCtx.BaseRelativeLocation = &newPtr
		}
	}
	if newPtr, err := schCtx.RelativeLocation.Descendant("then"); err == nil {
		subCtx.RelativeLocation = &newPtr
	}
	sch := Schema(*t)
	sch.ValidateFromContext(subCtx, errs)
}


func (t Then) JSONProp(name string) interface{} {
	return Schema(t).JSONProp(name)
}

func (t Then) JSONChildren() (res map[string]JSONPather) {
	return Schema(t).JSONChildren()
}

func (t *Then) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*t = Then(sch)
	return nil
}

func (t Then) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(t))
}

//
// Else
//

type Else Schema

func NewElse() Keyword {
	return &Else{}
}

func (e *Else) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (e *Else) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(e).Register(uri, registry)
}

func (e *Else) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(e).Resolve(pointer, uri)
}

func (e *Else) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	ifResult, okIf := schCtx.Misc["ifResult"]
	if !okIf {
		// if not found
		return
	}
	if (ifResult.(bool)) {
		// if was true
		return
	}

	subCtx := NewSchemaContextFromSource(*schCtx)
	if subCtx.BaseRelativeLocation != nil {
		if newPtr, err := schCtx.BaseRelativeLocation.Descendant("else"); err == nil {
			subCtx.BaseRelativeLocation = &newPtr
		}
	}
	if newPtr, err := schCtx.RelativeLocation.Descendant("else"); err == nil {
		subCtx.RelativeLocation = &newPtr
	}
	sch := Schema(*e)
	sch.ValidateFromContext(subCtx, errs)
}

func (e Else) JSONProp(name string) interface{} {
	return Schema(e).JSONProp(name)
}

func (e Else) JSONChildren() (res map[string]JSONPather) {
	return Schema(e).JSONChildren()
}

func (e *Else) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*e = Else(sch)
	return nil
}

func (e Else) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(e))
}
