package jsonschema

import "testing"

func TestErrorMessage(t *testing.T) {
	cases := []struct {
		schema, doc, message string
	}{
		{`{ "const" : "a value" }`, `"a different value"`, `must equal "a value"`},
	}

	for i, c := range cases {
		rs := &RootSchema{}
		if err := rs.UnmarshalJSON([]byte(c.schema)); err != nil {
			t.Errorf("case %d schema is invalid: %s", i, err.Error())
			continue
		}

		errs, err := rs.ValidateBytes([]byte(c.doc))
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
