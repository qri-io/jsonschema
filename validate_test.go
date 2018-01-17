package jsonschema

import (
	"testing"
)

type FooValidator uint8

func (f *FooValidator) Validate(data interface{}) error {
	return nil
}

func TestRegisterValidator(t *testing.T) {
	newFoo := func() Validator {
		return new(FooValidator)
	}

	RegisterValidator("foo", newFoo)

	if _, ok := DefaultValidators["foo"]; !ok {
		t.Errorf("expected %s to be added as a default validator", "foo")
	}
}
