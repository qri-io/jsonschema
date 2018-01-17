package jsonschema

import (
	"fmt"
)

// MultipleOf MUST be a number, strictly greater than 0.
// MultipleOf validates that a numeric instance is valid only if division
// by this keyword's value results in an integer.
type MultipleOf float64

func newMultipleOf() Validator {
	return new(MultipleOf)
}

// Validate implements the Validator interface for MultipleOf
func (m MultipleOf) Validate(data interface{}) error {
	if num, ok := data.(float64); ok {
		div := num / float64(m)
		if float64(int(div)) != div {
			return fmt.Errorf("%f must be a multiple of %f", num, m)
		}
	}
	return nil
}

// Maximum MUST be a number, representing an inclusive upper limit
// for a numeric instance.
// If the instance is a number, then this keyword validates only if the instance is less than or exactly equal to "maximum".
type Maximum float64

func newMaximum() Validator {
	return new(Maximum)
}

// Validate implements the Validator interface for Maximum
func (m Maximum) Validate(data interface{}) error {
	if num, ok := data.(float64); ok {
		if num > float64(m) {
			return fmt.Errorf("%f must be less than or equal to %f", num, m)
		}
	}
	return nil
}

// ExclusiveMaximum MUST be number, representing an exclusive upper limit for a numeric instance.
// If the instance is a number, then the instance is valid only if it has a value
// strictly less than (not equal to) "exclusiveMaximum".
type ExclusiveMaximum float64

func newExclusiveMaximum() Validator {
	return new(ExclusiveMaximum)
}

// Validate implements the Validator interface for ExclusiveMaximum
func (m ExclusiveMaximum) Validate(data interface{}) error {
	if num, ok := data.(float64); ok {
		if num >= float64(m) {
			return fmt.Errorf("%f must be less than %f", num, m)
		}
	}
	return nil
}

// Minimum MUST be a number, representing an inclusive lower limit for a numeric instance.
// If the instance is a number, then this keyword validates only if the instance is greater than or exactly equal to "minimum".
type Minimum float64

func newMinimum() Validator {
	return new(Minimum)
}

// Validate implements the Validator interface for Minimum
func (m Minimum) Validate(data interface{}) error {
	if num, ok := data.(float64); ok {
		if num < float64(m) {
			return fmt.Errorf("%f must be greater than or equal to %f", num, m)
		}
	}
	return nil
}

// ExclusiveMinimum MUST be number, representing an exclusive lower limit for a numeric instance.
// If the instance is a number, then the instance is valid only if it has a value strictly greater than (not equal to) "exclusiveMinimum".
type ExclusiveMinimum float64

func newExclusiveMinimum() Validator {
	return new(ExclusiveMinimum)
}

// Validate implements the Validator interface for ExclusiveMinimum
func (m ExclusiveMinimum) Validate(data interface{}) error {
	if num, ok := data.(float64); ok {
		if num <= float64(m) {
			return fmt.Errorf("%f must be greater than %f", num, m)
		}
	}
	return nil
}
