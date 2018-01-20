package jsonschema

// Validator is an interface for anything that can validate.
// JSON-Schema keywords are all examples of validators
type Validator interface {
	// Validate checks decoded JSON data against a given constraint
	Validate(data interface{}) error
}

// ValMaker is a function that generates instances of a validator.
// Calls to ValMaker will be passed directly to json.Marshal, so it should be a pointer
type ValMaker func() Validator

// DefaultValidators is a map of JSON keywords to Validators
// to draw from when decoding schemas
var DefaultValidators = map[string]ValMaker{
	// standard keywords
	"type":  newTipe,
	"enum":  newEnum,
	"const": newKonst,

	// numeric keywords
	"multipleOf":       newMultipleOf,
	"maximum":          newMaximum,
	"exclusiveMaximum": newExclusiveMaximum,
	"minimum":          newMinimum,
	"exclusiveMinimum": newExclusiveMinimum,

	// string keywords
	"maxLength": newMaxLength,
	"minLength": newMinLength,
	"pattern":   newPattern,

	// boolean keywords
	"allOf": newAllOf,
	"anyOf": newAnyOf,
	"oneOf": newOneOf,
	"not":   newNot,

	// array keywords
	"items":           newItems,
	"additionalItems": newAdditionalItems,
	"maxItems":        newMaxItems,
	"minItems":        newMinItems,
	"uniqueItems":     newUniqueItems,
	"contains":        newContains,

	// object keywords
	"maxProperties":        newMaxProperties,
	"minProperties":        newMinProperties,
	"required":             newRequired,
	"properties":           newProperties,
	"patternProperties":    newPatternProperties,
	"additionalProperties": newAdditionalProperties,
	"dependencies":         newDependencies,
	"propertyNames":        newPropertyNames,

	// conditional keywords
	"if":   newIif,
	"then": newThen,
	"else": newEls,

	//optional formats
	"format": newFormat,
}

// RegisterValidator adds a validator to DefaultValidators.
// Custom Validators should satisfy the validator interface,
// and be able to get cleanly endcode/decode to JSON
func RegisterValidator(propName string, maker ValMaker) {
	// TODO - should this call the function and panic if
	// the result can't be fed to json.Umarshal?
	DefaultValidators[propName] = maker
}
