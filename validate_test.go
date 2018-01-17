package jsonschema

import (
	"encoding/json"
	"fmt"
	"testing"
)

type IsFoo bool

func newIsFoo() Validator {
	return new(IsFoo)
}

func (f IsFoo) Validate(data interface{}) error {
	if str, ok := data.(string); ok {
		if str != "foo" {
			return fmt.Errorf("'%s' is not foo. It should be foo. plz make '%s' == foo. plz", str, str)
		}
	}
	return nil
}

func ExampleCustomValidator() {
	// register a custom validator by supplying a function
	// that creates new instances of your Validator.
	RegisterValidator("foo", newIsFoo)

	schBytes := []byte(`{ "foo": true }`)

	rs := new(RootSchema)
	if err := json.Unmarshal(schBytes, rs); err != nil {
		// Real programs handle errors.
		panic(err)
	}

	err := rs.ValidateBytes([]byte(`"bar"`))
	fmt.Println(err.Error())

	// Output: 'bar' is not foo. It should be foo. plz make 'bar' == foo. plz
}

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
