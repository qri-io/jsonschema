package jsonschema

import (
	"fmt"
)

// MultipleOf MUST be a number, strictly greater than 0.
// MultipleOf validates that a numeric instance is valid only if division
// by this keyword's value results in an integer.
type MultipleOf float64

// NewMultipleOf allocates a new MultipleOf validator
func NewMultipleOf() Validator {
	return new(MultipleOf)
}

// Validate implements the Validator interface for MultipleOf
func (m MultipleOf) Validate(data interface{}) []ValError {
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

// Maximum MUST be a number, representing an inclusive upper limit
// for a numeric instance.
// If the instance is a number, then this keyword validates only if the instance is less than or exactly equal to "Maximum".
type Maximum float64

// NewMaximum allocates a new Maximum validator
func NewMaximum() Validator {
	return new(Maximum)
}

// Validate implements the Validator interface for Maximum
func (m Maximum) Validate(data interface{}) []ValError {
	if num, ok := data.(float64); ok {
		if num > float64(m) {
			return []ValError{
				{Message: fmt.Sprintf("%f must be less than or equal to %f", num, m)},
			}
		}
	}
	return nil
}

// ExclusiveMaximum MUST be number, representing an exclusive upper limit for a numeric instance.
// If the instance is a number, then the instance is valid only if it has a value
// strictly less than (not equal to) "Exclusivemaximum".
type ExclusiveMaximum float64

// NewExclusiveMaximum allocates a new ExclusiveMaximum validator
func NewExclusiveMaximum() Validator {
	return new(ExclusiveMaximum)
}

// Validate implements the Validator interface for ExclusiveMaximum
func (m ExclusiveMaximum) Validate(data interface{}) []ValError {
	if num, ok := data.(float64); ok {
		if num >= float64(m) {
			return []ValError{
				{Message: fmt.Sprintf("%f must be less than %f", num, m)},
			}
		}
	}
	return nil
}

// Minimum MUST be a number, representing an inclusive lower limit for a numeric instance.
// If the instance is a number, then this keyword validates only if the instance is greater than or exactly equal to "Minimum".
type Minimum float64

// NewMinimum allocates a new Minimum validator
func NewMinimum() Validator {
	return new(Minimum)
}

// Validate implements the Validator interface for Minimum
func (m Minimum) Validate(data interface{}) []ValError {
	if num, ok := data.(float64); ok {
		if num < float64(m) {
			return []ValError{
				{Message: fmt.Sprintf("%f must be greater than or equal to %f", num, m)},
			}
		}
	}
	return nil
}

// ExclusiveMinimum MUST be number, representing an exclusive lower limit for a numeric instance.
// If the instance is a number, then the instance is valid only if it has a value strictly greater than (not equal to) "ExclusiveMinimum".
type ExclusiveMinimum float64

// NewExclusiveMinimum allocates a new ExclusiveMinimum validator
func NewExclusiveMinimum() Validator {
	return new(ExclusiveMinimum)
}

// Validate implements the Validator interface for ExclusiveMinimum
func (m ExclusiveMinimum) Validate(data interface{}) []ValError {
	if num, ok := data.(float64); ok {
		if num <= float64(m) {
			return []ValError{
				{Message: fmt.Sprintf("%f must be greater than %f", num, m)},
			}
		}
	}
	return nil
}
