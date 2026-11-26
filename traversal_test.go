package jsonschema

import (
	"context"
	"io/ioutil"
	"testing"

	jptr "github.com/qri-io/jsonpointer"
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

func TestKeywordResolve(t *testing.T) {
	const uri = "https://example.com/root.json"
	list := `[{"title": "0"}, {"properties": {"p": {"title": "1/properties/p"}}}]`
	dict := `{"a": {"title": "a"}, "b": {"properties": {"p": {"title": "b/properties/p"}}}}`
	rs := Must(`{
		"allOf": ` + list + `,
		"anyOf": ` + list + `,
		"oneOf": ` + list + `,
		"items": ` + list + `,
		"$defs": ` + dict + `,
		"properties": ` + dict + `
	}`)

	type resolveCase struct {
		keyword string
		pointer jptr.Pointer
		title   string // title of the schema expected back, "" for no schema
	}
	cases := []resolveCase{}
	for _, kw := range []string{"allOf", "anyOf", "oneOf", "items"} {
		cases = append(cases,
			resolveCase{kw, jptr.Pointer{"0"}, "0"},
			resolveCase{kw, jptr.Pointer{"1", "properties", "p"}, "1/properties/p"},
			resolveCase{kw, jptr.Pointer{"2"}, ""},
			resolveCase{kw, jptr.Pointer{"-1"}, ""},
			resolveCase{kw, jptr.Pointer{"a"}, ""},
			resolveCase{kw, jptr.Pointer{}, ""},
			resolveCase{kw, nil, ""},
		)
	}
	for _, kw := range []string{"$defs", "properties"} {
		cases = append(cases,
			resolveCase{kw, jptr.Pointer{"a"}, "a"},
			resolveCase{kw, jptr.Pointer{"b", "properties", "p"}, "b/properties/p"},
			resolveCase{kw, jptr.Pointer{"c"}, ""},
			resolveCase{kw, jptr.Pointer{}, ""},
			resolveCase{kw, nil, ""},
		)
	}

	for _, c := range cases {
		got := rs.keywords[c.keyword].Resolve(c.pointer, uri)
		if c.title == "" {
			if got != nil {
				t.Errorf("%s %q: expected no schema, got: %v", c.keyword, c.pointer.String(), got)
			}
			continue
		}
		if got == nil {
			t.Errorf("%s %q: expected schema titled %q, got nil", c.keyword, c.pointer.String(), c.title)
			continue
		}
		if title, _ := got.keywords["title"].(*Title); title == nil || string(*title) != c.title {
			t.Errorf("%s %q: expected schema titled %q, got: %v", c.keyword, c.pointer.String(), c.title, got.keywords["title"])
		}
		if got.docPath != uri {
			t.Errorf("%s %q: expected docPath %q, got: %q", c.keyword, c.pointer.String(), uri, got.docPath)
		}
	}
}
