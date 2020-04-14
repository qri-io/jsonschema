package jsonschema

import (
	"io/ioutil"
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

	got, err := rs.ValidateBytes([]byte(`"a"`))
	if err != nil {
		t.Errorf("error validating bytes: %s", err.Error())
		return
	}

	if got == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestReferenceTraversal(t *testing.T) {
	sch, err := ioutil.ReadFile("testdata/draft2019-09_schema.json")
	if err != nil {
		t.Errorf("error reading file: %s", err.Error())
		return
	}

	rs := &RootSchema{}
	if err := rs.UnmarshalJSON(sch); err != nil {
		t.Errorf("unexpected unmarshal error: %s", err.Error())
		return
	}

	elements := 0
	expectElements := 120
	refs := 0
	expectRefs := 29
	walkJSON(rs, func(elem JSONPather) error {
		elements++
		if sch, ok := elem.(*Schema); ok {
			if sch.Ref != "" {
				refs++
			}
		}
		return nil
	})

	if elements != expectElements {
		t.Errorf("expected %d elements, got: %d", expectElements, elements)
	}
	if refs != expectRefs {
		t.Errorf("expected %d references, got: %d", expectRefs, refs)
	}

	cases := []struct {
		input    string
		elements int
		refs     int
	}{
		{`{ "not" : { "$ref":"#" }}`, 3, 1},
	}

	for i, c := range cases {
		rs := &RootSchema{}
		if err := rs.UnmarshalJSON([]byte(c.input)); err != nil {
			t.Errorf("unexpected unmarshal error: %s", err.Error())
			return
		}

		elements := 0
		refs := 0
		walkJSON(rs, func(elem JSONPather) error {
			elements++
			if sch, ok := elem.(*Schema); ok {
				if sch.Ref != "" {
					refs++
				}
			}
			return nil
		})
		if elements != c.elements {
			t.Errorf("case %d: expected %d elements, got: %d", i, c.elements, elements)
		}
		if refs != c.refs {
			t.Errorf("case %d: expected %d references, got: %d", i, c.refs, refs)
		}
	}

}
