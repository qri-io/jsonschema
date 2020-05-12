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

	// other
	"contentEncoding":  true,
	"contentMediaType": true,
	"contentSchema":    true,
	"deprecated":       true,

	// backward compatibility
	"definitions":  true,
	"dependencies": true,
}

var keywordRegistry = map[string]KeyMaker{}
var keywordOrder = map[string]int{}
var keywordInsertOrder = map[string]int{}

// IsKeyword validates if a given prop string is a registered keyword
func IsKeyword(prop string) bool {
	_, ok := keywordRegistry[prop]
	return ok
}

// GetKeyword returns a new instance of the keyword
func GetKeyword(prop string) Keyword {
	if !IsKeyword(prop) {
		return NewVoid()
	}
	return keywordRegistry[prop]()
}

// GetKeywordOrder returns the order index of
// the given keyword or defaults to 1
func GetKeywordOrder(prop string) int {
	if order, ok := keywordOrder[prop]; ok {
		return order
	}
	return 1
}

// GetKeywordInsertOrder returns the insert index of
// the given keyword
func GetKeywordInsertOrder(prop string) int {
	if order, ok := keywordInsertOrder[prop]; ok {
		return order
	}
	// TODO(arqu): this is an arbitrary max
	return 1000
}

// SetKeywordOrder assignes a given order to a keyword
func SetKeywordOrder(prop string, order int) {
	keywordOrder[prop] = order
}

// IsNotSupportedKeyword is a utility function to clarify when
// a given keyword, while expected is not supported
func IsNotSupportedKeyword(prop string) bool {
	_, ok := notSupported[prop]
	return ok
}

// IsRegistryLoaded checks if any keywords are present
func IsRegistryLoaded() bool {
	return keywordRegistry != nil && len(keywordRegistry) > 0
}

// RegisterKeyword registers a keyword with the registry
func RegisterKeyword(prop string, maker KeyMaker) {
	keywordRegistry[prop] = maker
	keywordInsertOrder[prop] = len(keywordInsertOrder)
}

// MaxKeywordErrStringLen sets how long a value can be before it's length is truncated
// when printing error strings
// a special value of -1 disables output trimming
var MaxKeywordErrStringLen = 20

// Keyword is an interface for anything that can validate.
// JSON-Schema keywords are all examples of Keyword
type Keyword interface {
	// ValidateFromContext checks decoded JSON data and writes
	// validation errors (if any) to an outparam slice of KeyErrors
	ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError)

	// Register subscribes all subschema to the parent schema
	// which is later used for ref resolution
	Register(uri string, registry *SchemaRegistry)
	// Resolve resolves the pointer to the respective schema
	// allowing validation to flow through referenced schemas
	Resolve(pointer jptr.Pointer, uri string) *Schema
}

// KeyMaker is a function that generates instances of a Keyword.
// Calls to KeyMaker will be passed directly to json.Marshal,
// so the returned value should be a pointer
type KeyMaker func() Keyword

// KeyError represents a single error in an instance of a schema
// The only absolutely-required property is Message.
type KeyError struct {
	// PropertyPath is a string path that leads to the
	// property that produced the error
	PropertyPath string `json:"propertyPath,omitempty"`
	// InvalidValue is the value that returned the error
	InvalidValue interface{} `json:"invalidValue,omitempty"`
	// Message is a human-readable description of the error
	Message string `json:"message"`
}

// Error implements the error interface for KeyError
func (v KeyError) Error() string {
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

// AddError creates and appends a KeyError to errs
// deprecated but kept for backwards compatibility
func AddError(errs *[]KeyError, propPath string, data interface{}, msg string) {
	*errs = append(*errs, KeyError{
		PropertyPath: propPath,
		InvalidValue: data,
		Message:      msg,
	})
}

// AddErrorCtx creates and appends a KeyError to errs from the provided schema context
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
