package jsonschema

// Validator is an interface for anything that can validate.
// JSON-Schema keywords are all examples of validators
type Validator interface {
	// Validate checks decoded JSON data and returns a slice
	// of validation errors (if any)
	Validate(data interface{}) []ValError
}

// ValError represents a single error in an instance of a schema
// The only absolutely-required property is Message.
type ValError struct {
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
func (v ValError) Error() string {
	return v.Message
}

// ValMaker is a function that generates instances of a validator.
// Calls to ValMaker will be passed directly to json.Marshal, so it should be a pointer
type ValMaker func() Validator

// RegisterValidator adds a validator to DefaultValidators.
// Custom Validators should satisfy the validator interface,
// and be able to get cleanly endcode/decode to JSON
func RegisterValidator(propName string, maker ValMaker) {
	// TODO - should this call the function and panic if
	// the result can't be fed to json.Umarshal?
	DefaultValidators[propName] = maker
}

// DefaultValidators is a map of JSON keywords to Validators
// to draw from when decoding schemas
var DefaultValidators = map[string]ValMaker{
	// standard keywords
	"type":  NewType,
	"enum":  NewEnum,
	"const": NewConst,

	// numeric keywords
	"multipleOf":       NewMultipleOf,
	"maximum":          NewMaximum,
	"exclusiveMaximum": NewExclusiveMaximum,
	"minimum":          NewMinimum,
	"exclusiveMinimum": NewExclusiveMinimum,

	// string keywords
	"maxLength": NewMaxLength,
	"minLength": NewMinLength,
	"pattern":   NewPattern,

	// boolean keywords
	"allOf": NewAllOf,
	"anyOf": NewAnyOf,
	"oneOf": NewOneOf,
	"not":   NewNot,

	// array keywords
	"items":           NewItems,
	"additionalItems": NewAdditionalItems,
	"maxItems":        NewMaxItems,
	"minItems":        NewMinItems,
	"uniqueItems":     NewUniqueItems,
	"contains":        NewContains,

	// object keywords
	"maxProperties":        NewMaxProperties,
	"minProperties":        NewMinProperties,
	"required":             NewRequired,
	"properties":           NewProperties,
	"patternProperties":    NewPatternProperties,
	"additionalProperties": NewAdditionalProperties,
	"dependencies":         NewDependencies,
	"propertyNames":        NewPropertyNames,

	// conditional keywords
	"if":   NewIf,
	"then": NewThen,
	"else": NewElse,

	//optional formats
	"format": NewFormat,
}
