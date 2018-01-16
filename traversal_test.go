package jsonschema

import (
	"testing"
)

func TestSchemaDeref(t *testing.T) {
	sch := []byte(`{
    "definitions": {
        "a": {"type": "integer"},
        "b": {"$ref": "#/definitions/a"},
        "c": {"$ref": "#/definitions/b"}
    },
    "$ref": "#/definitions/c"
  }`)

	rs := &RootSchema{}
	if err := rs.UnmarshalJSON(sch); err != nil {
		t.Errorf("unexpected unmarshal error: %s", err.Error())
		return
	}

	got := rs.ValidateBytes([]byte(`"a"`))
	if got == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
