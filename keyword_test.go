package jsonschema

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	jptr "github.com/qri-io/jsonpointer"
)

func TestErrorMessage(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		schema, doc, message string
	}{
		{`{ "const" : "a value" }`, `"a different value"`, `must equal "a value"`},
	}

	for i, c := range cases {
		rs := &Schema{}
		if err := rs.UnmarshalJSON([]byte(c.schema)); err != nil {
			t.Errorf("case %d schema is invalid: %s", i, err.Error())
			continue
		}

		errs, err := rs.ValidateBytes(ctx, []byte(c.doc))
		if err != nil {
			t.Errorf("case %d error validating: %s", i, err)
			continue
		}

		if len(errs) != 1 {
			t.Errorf("case %d didn't return exactly 1 validation error. got: %d", i, len(errs))
			continue
		}

		if errs[0].Message != c.message {
			t.Errorf("case %d error mismatch. expected '%s', got: '%s'", i, c.message, errs[0].Message)
		}
	}
}

type IsFoo bool

func newIsFoo() Keyword {
	return new(IsFoo)
}

func (f *IsFoo) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (f *IsFoo) Register(uri string, registry *SchemaRegistry) {}

func (f *IsFoo) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (f *IsFoo) ValidateKeyword(ctx context.Context, currentState *ValidationState, data interface{}) {
	if str, ok := data.(string); ok {
		if str != "foo" {
			currentState.AddError(data, fmt.Sprintf("should be foo. plz make '%s' == foo. plz", str))
		}
	}
}

func ExampleCustomValidator() {
	// register a custom validator by supplying a function
	// that creates new instances of your Validator.
	RegisterKeyword("foo", newIsFoo)
	ctx := context.Background()

	schBytes := []byte(`{ "foo": true }`)

	rs := new(Schema)
	if err := json.Unmarshal(schBytes, rs); err != nil {
		// Real programs handle errors.
		panic(err)
	}

	errs, err := rs.ValidateBytes(ctx, []byte(`"bar"`))
	if err != nil {
		panic(err)
	}

	fmt.Println(errs[0].Error())
	// Output: /: "bar" should be foo. plz make 'bar' == foo. plz
}

type FooKeyword uint8

func (f *FooKeyword) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (f *FooKeyword) Register(uri string, registry *SchemaRegistry) {}

func (f *FooKeyword) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (f *FooKeyword) ValidateKeyword(ctx context.Context, currentState *ValidationState, data interface{}) {
}

func TestRegisterFooKeyword(t *testing.T) {
	newFoo := func() Keyword {
		return new(FooKeyword)
	}

	RegisterKeyword("foo", newFoo)

	if !IsRegisteredKeyword("foo") {
		t.Errorf("expected %s to be added as a default validator", "foo")
	}
}
