package jsonschema

import (
	"encoding/json"
	"fmt"
	"regexp"
	"unicode/utf8"
)

// maxLength MUST be a non-negative integer.
// A string instance is valid against this keyword if its length is less than, or equal to, the value of this keyword.
// The length of a string instance is defined as the number of its characters as defined by RFC 7159 [RFC7159].
type maxLength int

func newMaxLength() Validator {
	return new(maxLength)
}

// Validate implements the Validator interface for maxLength
func (m maxLength) Validate(data interface{}) error {
	if str, ok := data.(string); ok {
		if utf8.RuneCountInString(str) > int(m) {
			return fmt.Errorf("max length of %d characters exceeded: %s", m, str)
		}
	}
	return nil
}

// minLength MUST be a non-negative integer.
// A string instance is valid against this keyword if its length is greater than, or equal to, the value of this keyword.
// The length of a string instance is defined as the number of its characters as defined by RFC 7159 [RFC7159].
// Omitting this keyword has the same behavior as a value of 0.
type minLength int

func newMinLength() Validator {
	return new(minLength)
}

// Validate implements the Validator interface for minLength
func (m minLength) Validate(data interface{}) error {
	if str, ok := data.(string); ok {
		if utf8.RuneCountInString(str) < int(m) {
			return fmt.Errorf("min length of %d characters required: %s", m, str)
		}
	}
	return nil
}

// pattern MUST be a string. This string SHOULD be a valid regular expression,
// according to the ECMA 262 regular expression dialect.
// A string instance is considered valid if the regular expression matches the instance successfully.
// Recall: regular expressions are not implicitly anchored.
type pattern regexp.Regexp

func newPattern() Validator {
	return &pattern{}
}

// Validate implements the Validator interface for pattern
func (p pattern) Validate(data interface{}) error {
	re := regexp.Regexp(p)
	if str, ok := data.(string); ok {
		if !re.Match([]byte(str)) {
			return fmt.Errorf("regext pattrn %s mismatch on string: %s", re.String(), str)
		}
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for pattern
func (p *pattern) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	ptn, err := regexp.Compile(str)
	if err != nil {
		return err
	}

	*p = pattern(*ptn)
	return nil
}

// MarshalJSON implements json.Marshaler for pattern
func (p pattern) MarshalJSON() ([]byte, error) {
	re := regexp.Regexp(p)
	rep := &re
	return json.Marshal(rep.String())
}
