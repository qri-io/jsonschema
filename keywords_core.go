package main

import (
	"encoding/json"
	"fmt"

	jptr "github.com/qri-io/jsonpointer"
)

//
// $schema
//

type SchemaURI string

func NewSchemaURI() Keyword {
	return new(SchemaURI)
}

func (s *SchemaURI) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (s *SchemaURI) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {}

func (s *SchemaURI) Register(uri string, registry *SchemaRegistry) {}

func (s *SchemaURI) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// $id
//

type Id string

func NewId() Keyword {
	return new(Id)
}

func (i *Id) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (i *Id) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	// TODO: make sure ID is valid URI for draft2019
}

func (i *Id) Register(uri string, registry *SchemaRegistry) {}

func (i *Id) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// Description
//

type Description string

func NewDescription() Keyword {
	return new(Description)
}

func (d *Description) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (d *Description) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {}

func (d *Description) Register(uri string, registry *SchemaRegistry) {}

func (d *Description) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// Title
//

type Title string

func NewTitle() Keyword {
	return new(Title)
}

func (t *Title) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (t *Title) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {}

func (t *Title) Register(uri string, registry *SchemaRegistry) {}

func (t *Title) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// Comment
//

type Comment string

func NewComment() Keyword {
	return new(Comment)
}

func (c *Comment) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (c *Comment) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {}

func (c *Comment) Register(uri string, registry *SchemaRegistry) {}

func (c *Comment) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// Default
//

type Default Schema

func NewDefault() Keyword {
	return &Default{}
}

func (d *Default) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (d *Default) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {}

func (d *Default) Register(uri string, registry *SchemaRegistry) {}

func (d *Default) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (d Default) JSONProp(name string) interface{} {
	return Schema(d).JSONProp(name)
}

func (d Default) JSONChildren() (res map[string]JSONPather) {
	return Schema(d).JSONChildren()
}

func (d *Default) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*d = Default(sch)
	return nil
}

func (d Default) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(d))
}

//
// Examples
//

type Examples []interface{}

func NewExamples() Keyword {
	return new(Examples)
}

func (e *Examples) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (e *Examples) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {}

func (e *Examples) Register(uri string, registry *SchemaRegistry) {}

func (e *Examples) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// ReadOnly
//

type ReadOnly bool

func NewReadOnly() Keyword {
	return new(ReadOnly)
}

func (r *ReadOnly) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (r *ReadOnly) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {}

func (r *ReadOnly) Register(uri string, registry *SchemaRegistry) {}

func (r *ReadOnly) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// WriteOnly
//

type WriteOnly bool

func NewWriteOnly() Keyword {
	return new(WriteOnly)
}

func (w *WriteOnly) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (w *WriteOnly) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {}

func (w *WriteOnly) Register(uri string, registry *SchemaRegistry) {}

func (w *WriteOnly) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// VOID
//

type Void struct{}

func NewVoid() Keyword {
	return &Void{}
}

func (vo *Void) Validate(propPath string, data interface{}, errs *[]KeyError) {
	fmt.Println("WARN: Using Void Validator - always True")
}

func (vo *Void) Register(uri string, registry *SchemaRegistry) {}

func (vo *Void) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (vo *Void) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	fmt.Println("WARN: Using Void Validator - always True")
}
