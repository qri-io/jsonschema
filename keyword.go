package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	jptr "github.com/qri-io/jsonpointer"
)

var notSupported = map[string]bool{
	// core
	"$vocabulary": true,

	// array keywords
	"unevaluatedItems": true,

	// object keywords
	"unevaluatedProperties": true,

	// other
	"contentEncoding":  true,
	"contentMediaType": true,
	"contentSchema":    true,
	"deprecated":       true,

	// backward compatibilit
	"definitions":  true,
	"dependencies": true,
}

var KeywordRegistry = map[string]KeyMaker{}
var KeywordOrder = map[string]int{}
var KeywordInsertOrder = map[string]int{}

func IsKeyword(prop string) bool {
	_, ok := KeywordRegistry[prop]
	return ok
}

func GetKeyword(prop string) Keyword {
	if !IsKeyword(prop) {
		return NewVoid()
	}
	return KeywordRegistry[prop]()
}

func GetKeywordOrder(prop string) int {
	if order, ok := KeywordOrder[prop]; ok {
		return order
	}
	return 1
}

func GetKeywordInsertOrder(prop string) int {
	if order, ok := KeywordInsertOrder[prop]; ok {
		return order
	}
	// TODO(arqu): this is an arbitrary max
	return 1000
}

func SetKeywordOrder(prop string, order int) {
	KeywordOrder[prop] = order
}

func IsNotSupportedKeyword(prop string) bool {
	_, ok := notSupported[prop]
	return ok
}

func IsRegistryLoaded() bool {
	return KeywordRegistry != nil && len(KeywordRegistry) > 0
}

func RegisterKeyword(prop string, maker KeyMaker) {
	KeywordRegistry[prop] = maker
	KeywordInsertOrder[prop] = len(KeywordInsertOrder)
}

var MaxKeywordErrStringLen = 20

type Keyword interface {
	Validate(propPath string, data interface{}, errs *[]KeyError)
	ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError)

	Register(uri string, registry *SchemaRegistry)
	Resolve(pointer jptr.Pointer, uri string) *Schema
}

type BaseKeyword struct {
	path string
}

func (b *BaseKeyword) SetPath(path string) {
	b.path = path
}

func (b BaseKeyword) Path() string {
	return b.path
}

func (b BaseKeyword) AddError(errs *[]KeyError, propPath string, data interface{}, msg string) {
	*errs = append(*errs, KeyError{
		PropertyPath: propPath,
		RulePath:     b.Path(),
		InvalidValue: data,
		Message:      msg,
	})
}

type KeyMaker func() Keyword

type KeyError struct {
	PropertyPath string      `json:"propertyPath,omitempty"`
	InvalidValue interface{} `json:"invalidValue,omitempty"`
	RulePath     string      `json:"rulePath,omitempty"`
	Message      string      `json:"message"`
}

func (v KeyError) Error() string {
	if v.PropertyPath != "" && v.InvalidValue != nil {
		return fmt.Sprintf("%s: %s %s", v.PropertyPath, InvalidValueString(v.InvalidValue), v.Message)
	} else if v.PropertyPath != "" {
		return fmt.Sprintf("%s: %s", v.PropertyPath, v.Message)
	}
	return v.Message
}

func InvalidValueString(data interface{}) string {
	bt, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	bt = bytes.Replace(bt, []byte{'\n', '\r'}, []byte{' '}, -1)
	if MaxKeywordErrStringLen != -1 && len(bt) > MaxKeywordErrStringLen {
		bt = append(bt[:MaxKeywordErrStringLen], []byte("...")...)
	}
	return string(bt)
}

func (k KeyError) String() string {
	return fmt.Sprintf("for: '%s' msg:'%s'", k.InvalidValue, k.Message)
}

func AddError(errs *[]KeyError, propPath string, data interface{}, msg string) {
	*errs = append(*errs, KeyError{
		PropertyPath: propPath,
		InvalidValue: data,
		Message:      msg,
	})
}

func AddErrorCtx(errs *[]KeyError, schCtx *SchemaContext, msg string) {
	instancePath := schCtx.InstanceLocation.String()
	if len(instancePath) == 0 {
		instancePath = "/"
	}
	*errs = append(*errs, KeyError{
		PropertyPath: instancePath,
		InvalidValue: schCtx.Instance,
		Message:      msg,
	})
}
