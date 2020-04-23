package main

import (
	"fmt"
	jptr "github.com/qri-io/jsonpointer"
)

//
// MultipleOf
//

type MultipleOf float64

func NewMultipleOf() Keyword {
	return new(MultipleOf)
}

func (m MultipleOf) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *MultipleOf) Register(uri string, registry *SchemaRegistry) {}

func (m *MultipleOf) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m MultipleOf) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if num, ok := schCtx.Instance.(float64); ok {
		div := num / float64(m)
		if float64(int(div)) != div {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("must be a multiple of %f", m))
		}
	}
}

//
// Maximum
//

type Maximum float64

func NewMaximum() Keyword {
	return new(Maximum)
}

func (m Maximum) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *Maximum) Register(uri string, registry *SchemaRegistry) {}

func (m *Maximum) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m Maximum) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if num, ok := schCtx.Instance.(float64); ok {
		if num > float64(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("must be less than or equal to %f", m))
		}
	}
}

//
// ExclusiveMaximum
//

type ExclusiveMaximum float64

func NewExclusiveMaximum() Keyword {
	return new(ExclusiveMaximum)
}

func (m ExclusiveMaximum) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *ExclusiveMaximum) Register(uri string, registry *SchemaRegistry) {}

func (m *ExclusiveMaximum) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m ExclusiveMaximum) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if num, ok := schCtx.Instance.(float64); ok {
		if num >= float64(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("%f must be less than %f", num, m))
		}
	}
}

//
// Minimum
//

type Minimum float64

func NewMinimum() Keyword {
	return new(Minimum)
}

func (m Minimum) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *Minimum) Register(uri string, registry *SchemaRegistry) {}

func (m *Minimum) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m Minimum) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if num, ok := schCtx.Instance.(float64); ok {
		if num < float64(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("must be less than or equal to %f", m))
		}
	}
}

//
// ExclusiveMinimum
//

type ExclusiveMinimum float64

func NewExclusiveMinimum() Keyword {
	return new(ExclusiveMinimum)
}

func (m ExclusiveMinimum) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *ExclusiveMinimum) Register(uri string, registry *SchemaRegistry) {}

func (m *ExclusiveMinimum) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m ExclusiveMinimum) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if num, ok := schCtx.Instance.(float64); ok {
		if num <= float64(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("%f must be less than %f", num, m))
		}
	}
}