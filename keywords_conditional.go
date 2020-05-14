package jsonschema

import (
	"encoding/json"

	jptr "github.com/qri-io/jsonpointer"
)

// If defines the if JSON Schema keyword
type If Schema

// NewIf allocates a new If keyword
func NewIf() Keyword {
	return &If{}
}

// Register implements the Keyword interface for If
func (f *If) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(f).Register(uri, registry)
}

// Resolve implements the Keyword interface for If
func (f *If) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ValidateFromContext implements the Keyword interface for If
func (f *If) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[If] Validating")
	thenKW := schCtx.Local.keywords["then"]
	elseKW := schCtx.Local.keywords["else"]

	if thenKW == nil && elseKW == nil {
		// no then or else for if, aborting validation
		return
	}

	subCtx := NewSchemaContextFromSourceClean(*schCtx)
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

// JSONProp implements the JSONPather for If
func (f If) JSONProp(name string) interface{} {
	return Schema(f).JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for If
func (f If) JSONChildren() (res map[string]JSONPather) {
	return Schema(f).JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for If
func (f *If) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*f = If(sch)
	return nil
}

// MarshalJSON implements the json.Marshaler interface for If
func (f If) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(f))
}

// Then defines the then JSON Schema keyword
type Then Schema

// NewThen allocates a new Then keyword
func NewThen() Keyword {
	return &Then{}
}

// Register implements the Keyword interface for Then
func (t *Then) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(t).Register(uri, registry)
}

// Resolve implements the Keyword interface for Then
func (t *Then) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(t).Resolve(pointer, uri)
}

// ValidateFromContext implements the Keyword interface for Then
func (t *Then) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Then] Validating")
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

// JSONProp implements the JSONPather for Then
func (t Then) JSONProp(name string) interface{} {
	return Schema(t).JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for Then
func (t Then) JSONChildren() (res map[string]JSONPather) {
	return Schema(t).JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for Then
func (t *Then) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*t = Then(sch)
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Then
func (t Then) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(t))
}

// Else defines the else JSON Schema keyword
type Else Schema

// NewElse allocates a new Else keyword
func NewElse() Keyword {
	return &Else{}
}

// Register implements the Keyword interface for Else
func (e *Else) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(e).Register(uri, registry)
}

// Resolve implements the Keyword interface for Else
func (e *Else) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(e).Resolve(pointer, uri)
}

// ValidateFromContext implements the Keyword interface for Else
func (e *Else) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Else] Validating")
	ifResult, okIf := schCtx.Misc["ifResult"]
	if !okIf {
		// if not found
		return
	}
	if ifResult.(bool) {
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

// JSONProp implements the JSONPather for Else
func (e Else) JSONProp(name string) interface{} {
	return Schema(e).JSONProp(name)
}

// JSONChildren implements the JSONContainer interface for Else
func (e Else) JSONChildren() (res map[string]JSONPather) {
	return Schema(e).JSONChildren()
}

// UnmarshalJSON implements the json.Unmarshaler interface for Else
func (e *Else) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch); err != nil {
		return err
	}
	*e = Else(sch)
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Else
func (e Else) MarshalJSON() ([]byte, error) {
	return json.Marshal(Schema(e))
}
