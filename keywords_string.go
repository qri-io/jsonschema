package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"unicode/utf8"

	jptr "github.com/qri-io/jsonpointer"
)

//
// MaxLength
//

type MaxLength int

func NewMaxLength() Keyword {
	return new(MaxLength)
}

func (m MaxLength) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *MaxLength) Register(uri string, registry *SchemaRegistry) {}

func (m *MaxLength) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m MaxLength) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[MaxLength] Validating")
	if str, ok := schCtx.Instance.(string); ok {
		if utf8.RuneCountInString(str) > int(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("max length of %d characters exceeded: %s", m, str))
		}
	}
}

//
// MinLength
//

type MinLength int

func NewMinLength() Keyword {
	return new(MinLength)
}

func (m MinLength) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (m *MinLength) Register(uri string, registry *SchemaRegistry) {}

func (m *MinLength) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (m MinLength) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[MinLength] Validating")
	if str, ok := schCtx.Instance.(string); ok {
		if utf8.RuneCountInString(str) < int(m) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("max length of %d characters exceeded: %s", m, str))
		}
	}
}

//
// Pattern
//

type Pattern regexp.Regexp

func NewPattern() Keyword {
	return &Pattern{}
}

func (p Pattern) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (p *Pattern) Register(uri string, registry *SchemaRegistry) {}

func (p *Pattern) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (p Pattern) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Pattern] Validating")
	re := regexp.Regexp(p)
	if str, ok := schCtx.Instance.(string); ok {
		if !re.Match([]byte(str)) {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("regexp pattern %s mismatch on string: %s", re.String(), str))
		}
	}
}

func (p *Pattern) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	ptn, err := regexp.Compile(str)
	if err != nil {
		return err
	}

	*p = Pattern(*ptn)
	return nil
}

func (p Pattern) MarshalJSON() ([]byte, error) {
	re := regexp.Regexp(p)
	rep := &re
	return json.Marshal(rep.String())
}
