package jsonschema

import (
	"context"
	"io/ioutil"
	"testing"
)

func TestSchemaDeref(t *testing.T) {
	ctx := context.Background()
	sch := []byte(`{
    "$defs": {
        "a": {"type": "integer"},
        "b": {"$ref": "#/$defs/a"},
        "c": {"$ref": "#/$defs/b"}
    },
    "$ref": "#/$defs/c"
  }`)

	rs := &Schema{}
	if err := rs.UnmarshalJSON(sch); err != nil {
		t.Errorf("unexpected unmarshal error: %s", err.Error())
		return
	}

	got, err := rs.ValidateBytes(ctx, []byte(`"a"`))
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

	rs := &Schema{}
	if err := rs.UnmarshalJSON(sch); err != nil {
		t.Errorf("unexpected unmarshal error: %s", err.Error())
		return
	}

	elements := 0
	expectElements := 14
	refs := 0
	expectRefs := 6
	walkJSON(rs, func(elem JSONPather) error {
		elements++
		if sch, ok := elem.(*Schema); ok {
			if sch.HasKeyword("$ref") {
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
		{`{ "not" : { "$ref":"#" }}`, 2, 0},
	}

	for i, c := range cases {
		rs := &Schema{}
		if err := rs.UnmarshalJSON([]byte(c.input)); err != nil {
			t.Errorf("unexpected unmarshal error: %s", err.Error())
			return
		}

		elements := 0
		refs := 0
		walkJSON(rs, func(elem JSONPather) error {
			elements++
			if sch, ok := elem.(*Schema); ok {
				if sch.HasKeyword("$ref") {
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
