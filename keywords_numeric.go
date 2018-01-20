package jsonschema

import (
	"fmt"
)

// multipleOf MUST be a number, strictly greater than 0.
// multipleOf validates that a numeric instance is valid only if division
// by this keyword's value results in an integer.
type multipleOf float64

func newMultipleOf() Validator {
	return new(multipleOf)
}

// Validate implements the Validator interface for multipleOf
func (m multipleOf) Validate(data interface{}) []ValError {
	if num, ok := data.(float64); ok {
		div := num / float64(m)
		if float64(int(div)) != div {
			return []ValError{
				{Message: fmt.Sprintf("%f must be a multiple of %f", num, m)},
			}
		}
	}
	return nil
}

// maximum MUST be a number, representing an inclusive upper limit
// for a numeric instance.
// If the instance is a number, then this keyword validates only if the instance is less than or exactly equal to "maximum".
type maximum float64

func newMaximum() Validator {
	return new(maximum)
}

// Validate implements the Validator interface for maximum
func (m maximum) Validate(data interface{}) []ValError {
	if num, ok := data.(float64); ok {
		if num > float64(m) {
			return []ValError{
				{Message: fmt.Sprintf("%f must be less than or equal to %f", num, m)},
			}
		}
	}
	return nil
}

// exclusiveMaximum MUST be number, representing an exclusive upper limit for a numeric instance.
// If the instance is a number, then the instance is valid only if it has a value
// strictly less than (not equal to) "exclusivemaximum".
type exclusiveMaximum float64

func newExclusiveMaximum() Validator {
	return new(exclusiveMaximum)
}

// Validate implements the Validator interface for exclusiveMaximum
func (m exclusiveMaximum) Validate(data interface{}) []ValError {
	if num, ok := data.(float64); ok {
		if num >= float64(m) {
			return []ValError{
				{Message: fmt.Sprintf("%f must be less than %f", num, m)},
			}
		}
	}
	return nil
}

// minimum MUST be a number, representing an inclusive lower limit for a numeric instance.
// If the instance is a number, then this keyword validates only if the instance is greater than or exactly equal to "minimum".
type minimum float64

func newMinimum() Validator {
	return new(minimum)
}

// Validate implements the Validator interface for minimum
func (m minimum) Validate(data interface{}) []ValError {
	if num, ok := data.(float64); ok {
		if num < float64(m) {
			return []ValError{
				{Message: fmt.Sprintf("%f must be greater than or equal to %f", num, m)},
			}
		}
	}
	return nil
}

// exclusiveMinimum MUST be number, representing an exclusive lower limit for a numeric instance.
// If the instance is a number, then the instance is valid only if it has a value strictly greater than (not equal to) "exclusiveminimum".
type exclusiveMinimum float64

func newExclusiveMinimum() Validator {
	return new(exclusiveMinimum)
}

// Validate implements the Validator interface for exclusiveMinimum
func (m exclusiveMinimum) Validate(data interface{}) []ValError {
	if num, ok := data.(float64); ok {
		if num <= float64(m) {
			return []ValError{
				{Message: fmt.Sprintf("%f must be greater than %f", num, m)},
			}
		}
	}
	return nil
}
