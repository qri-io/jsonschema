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

func (f IsFoo) Validate(propPath string, data Val, errs *[]ValError) {
	if str, ok := data.(StringVal); ok {
		if str != "foo" {
			AddError(errs, propPath, data, fmt.Sprintf("should be foo. plz make '%s' == foo. plz", str))
		}
	}
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

	errs, err := rs.ValidateBytes([]byte(`"bar"`))
	if err != nil {
		panic(err)
	}

	fmt.Println(errs[0].Error())
	// Output: /: "bar" should be foo. plz make 'bar' == foo. plz
}

type FooValidator uint8

func (f *FooValidator) Validate(propPath string, data Val, errs *[]ValError) {}

func TestRegisterValidator(t *testing.T) {
	newFoo := func() Validator {
		return new(FooValidator)
	}

	RegisterValidator("foo", newFoo)

	if _, ok := DefaultValidators["foo"]; !ok {
		t.Errorf("expected %s to be added as a default validator", "foo")
	}
}
