package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	jptr "github.com/qri-io/jsonpointer"
)

var notSupported = map[string]bool{
	"$anchor": true,
	"$recursiveAnchor": true,
	"$defs": true,
	"definitions": true,
	"$ref": true,
	"$recursiveRef": true,
	// "default": true,
	"$vocabulary": true,

	// array keywords
	"unevaluatedItems": true,

	// object keywords
	"dependencies": true,
	"dependentRequired": true,
	// "dependentSchemas": true,
	"unevaluatedProperties": true,

	// other
	"contentEncoding": true,
	"contentMediaType": true,
	"contentSchema": true,
	"deprecated": true,
}

var KeywordRegistry = map[string]KeyMaker{}
var KeywordOrder = map[string]int{}

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

func SetKeywordOrder(prop string, order int) {
	KeywordOrder[prop] = order
}

func IsNotSupportedKeyword(prop string) bool {
	_, ok := notSupported[prop]
	return ok
}

func RegisterKeyword(prop string, maker KeyMaker) {
	KeywordRegistry[prop] = maker
}

// MaxValueErrStringLen sets how long a value can be before it's length is truncated
// when printing error strings
// a special value of -1 disables output trimming
var MaxKeywordErrStringLen = 20

// Validator is an interface for anything that can validate.
// JSON-Schema keywords are all examples of validators
type Keyword interface {
	// Validate checks decoded JSON data and writes
	// validation errors (if any) to an outparam slice of ValErrors
	// propPath indicates the position of data in the json tree
	Validate(propPath string, data interface{}, errs *[]KeyError)
	ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError)

	Register(uri string, registry *SchemaRegistry)
	Resolve(pointer jptr.Pointer, uri string) *Schema
}

// BaseValidator is a foundation for building a validator
type BaseKeyword struct {
	path string
}

// SetPath sets base validator's path
func (b *BaseKeyword) SetPath(path string) {
	b.path = path
}

// Path gives this validator's path
func (b BaseKeyword) Path() string {
	return b.path
}

// AddError is a convenience method for appending a new error to an existing error slice
func (b BaseKeyword) AddError(errs *[]KeyError, propPath string, data interface{}, msg string) {
	*errs = append(*errs, KeyError{
		PropertyPath: propPath,
		RulePath:     b.Path(),
		InvalidValue: data,
		Message:      msg,
	})
}

// ValMaker is a function that generates instances of a validator.
// Calls to ValMaker will be passed directly to json.Marshal,
// so the returned value should be a pointer
type KeyMaker func() Keyword

// ValError represents a single error in an instance of a schema
// The only absolutely-required property is Message.
type KeyError struct {
	// PropertyPath is a string path that leads to the
	// property that produced the error
	PropertyPath string `json:"propertyPath,omitempty"`
	// InvalidValue is the value that returned the error
	InvalidValue interface{} `json:"invalidValue,omitempty"`
	// RulePath is the path to the rule that errored
	RulePath string `json:"rulePath,omitempty"`
	// Message is a human-readable description of the error
	Message string `json:"message"`
}

// Error implements the error interface for ValError
func (v KeyError) Error() string {
	// [propPath]: [value] [message]
	if v.PropertyPath != "" && v.InvalidValue != nil {
		return fmt.Sprintf("%s: %s %s", v.PropertyPath, InvalidValueString(v.InvalidValue), v.Message)
	} else if v.PropertyPath != "" {
		return fmt.Sprintf("%s: %s", v.PropertyPath, v.Message)
	}
	return v.Message
}

// InvalidValueString returns the errored value as a string
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
	*errs = append(*errs, KeyError{
		PropertyPath: schCtx.BaseRelativeLocation.String(),
		InvalidValue: schCtx.Instance,
		Message:      msg,
	})
}

