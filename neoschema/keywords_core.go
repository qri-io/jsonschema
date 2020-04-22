package main

import (
	"encoding/json"
	"fmt"
	// "reflect"
	// "strconv"
	// "strings"
	jptr "github.com/qri-io/jsonpointer"
)

//
// $id
//

type Id struct {
	value string
}

func NewId() Keyword {
	return &Id{}
}

func (i *Id) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (i *Id) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	// TODO: make sure ID is valid URI for draft2019
}

func (i *Id) Register(uri string, registry *SchemaRegistry) {}

func (i *Id) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (i *Id) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err != nil {
		return fmt.Errorf("$id must be a string")
	}
	*i = Id{value: single,}
	return nil
}

func (i Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}

//
// VOID
//

type Void struct {}

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