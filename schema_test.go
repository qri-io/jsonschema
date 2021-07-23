package jsonschema

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func ExampleBasic() {
	ctx := context.Background()
	var schemaData = []byte(`{
	"title": "Person",
	"type": "object",
	"$id": "https://qri.io/schema/",
	"$comment" : "sample comment",
	"properties": {
	    "firstName": {
	        "type": "string"
	    },
	    "lastName": {
	        "type": "string"
	    },
	    "age": {
	        "description": "Age in years",
	        "type": "integer",
	        "minimum": 0
	    },
	    "friends": {
	    	"type" : "array",
	    	"items" : { "title" : "REFERENCE", "$ref" : "#" }
	    }
	},
	"required": ["firstName", "lastName"]
	}`)

	rs := &Schema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}

	var valid = []byte(`{
		"firstName" : "George",
		"lastName" : "Michael"
		}`)
	errs, err := rs.ValidateBytes(ctx, valid)
	if err != nil {
		panic(err)
	}

	if len(errs) > 0 {
		fmt.Println(errs[0].Error())
	}

	var invalidPerson = []byte(`{
		"firstName" : "Prince"
		}`)

	errs, err = rs.ValidateBytes(ctx, invalidPerson)
	if err != nil {
		panic(err)
	}
	if len(errs) > 0 {
		fmt.Println(errs[0].Error())
	}

	var invalidFriend = []byte(`{
		"firstName" : "Jay",
		"lastName" : "Z",
		"friends" : [{
			"firstName" : "Nas"
			}]
		}`)
	errs, err = rs.ValidateBytes(ctx, invalidFriend)
	if err != nil {
		panic(err)
	}
	if len(errs) > 0 {
		fmt.Println(errs[0].Error())
	}

	// Output: /: {"firstName":"Prince... "lastName" value is required
	// /friends/0: {"firstName":"Nas"} "lastName" value is required
}

func TestTopLevelType(t *testing.T) {
	schemaObject := []byte(`{
    "title": "Car",
    "type": "object",
    "properties": {
        "color": {
            "type": "string"
        }
    },
    "required": ["color"]
}`)
	rs := &Schema{}
	if err := json.Unmarshal(schemaObject, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}
	if rs.TopLevelType() != "object" {
		t.Errorf("error: schemaObject should be an object")
	}

	schemaArray := []byte(`{
    "title": "Cities",
    "type": "array",
    "items" : { "title" : "REFERENCE", "$ref" : "#" }
}`)
	rs = &Schema{}
	if err := json.Unmarshal(schemaArray, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}
	if rs.TopLevelType() != "array" {
		t.Errorf("error: schemaArray should be an array")
	}

	schemaUnknown := []byte(`{
    "title": "Typeless",
    "items" : { "title" : "REFERENCE", "$ref" : "#" }
}`)
	rs = &Schema{}
	if err := json.Unmarshal(schemaUnknown, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}
	if rs.TopLevelType() != "unknown" {
		t.Errorf("error: schemaUnknown should have unknown type")
	}
}

func TestParseUrl(t *testing.T) {
	// Easy case, id is a standard URL
	schemaObject := []byte(`{
    "title": "Car",
    "type": "object",
    "$id": "http://example.com/root.json"
}`)
	rs := &Schema{}
	if err := json.Unmarshal(schemaObject, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}

	// Tricky case, id is only a URL fragment
	schemaObject = []byte(`{
    "title": "Car",
    "type": "object",
    "$id": "#/properites/firstName"
}`)
	rs = &Schema{}
	if err := json.Unmarshal(schemaObject, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}

	// Another tricky case, id is only an empty fragment
	schemaObject = []byte(`{
    "title": "Car",
    "type": "object",
    "$id": "#"
}`)
	rs = &Schema{}
	if err := json.Unmarshal(schemaObject, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}
}

func TestMust(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				if err.Error() != "unexpected end of JSON input" {
					t.Errorf("expected panic error to equal: %s", "unexpected end of JSON input")
				}
			} else {
				t.Errorf("must paniced with a non-error")
			}
		} else {
			t.Errorf("expected invalid call to Must to panic")
		}
	}()

	// Valid call to Must shouldn't panic
	rs := Must(`{}`)
	if rs == nil {
		t.Errorf("expected parse of empty schema to return *RootSchema, got nil")
		return
	}

	// This should panic, checked in defer above
	Must(``)
}

func TestDraft3(t *testing.T) {
	runJSONTests(t, []string{
		"testdata/draft3/additionalProperties.json",
		"testdata/draft3/default.json",
		"testdata/draft3/format.json",
		"testdata/draft3/items.json",
		"testdata/draft3/maxItems.json",
		"testdata/draft3/maxLength.json",
		"testdata/draft3/minItems.json",
		"testdata/draft3/minLength.json",
		"testdata/draft3/pattern.json",
		"testdata/draft3/patternProperties.json",
		"testdata/draft3/properties.json",
		"testdata/draft3/uniqueItems.json",

		// disabled due to changes in spec
		// "testdata/draft3/disallow.json",
		// "testdata/draft3/divisibleBy.json",
		// "testdata/draft3/enum.json",
		// "testdata/draft3/extends.json",
		// "testdata/draft3/maximum.json",
		// "testdata/draft3/minimum.json",
		// "testdata/draft3/ref.json",
		// "testdata/draft3/refRemote.json",
		// "testdata/draft3/required.json",
		// "testdata/draft3/type.json",
		// "testdata/draft3/optional/format.json",
		// "testdata/draft3/optional/zeroTerminatedFloats.json",

		// TODO(arqu): implement this
		// "testdata/draft3/additionalItems.json",
		// "testdata/draft3/dependencies.json",

		// wont fix
		// "testdata/draft3/optional/bignum.json",
		// "testdata/draft3/optional/ecmascript-regex.json",
	})
}

func TestDraft4(t *testing.T) {
	runJSONTests(t, []string{
		"testdata/draft4/additionalItems.json",
		"testdata/draft4/allOf.json",
		"testdata/draft4/anyOf.json",
		"testdata/draft4/default.json",
		"testdata/draft4/enum.json",
		"testdata/draft4/format.json",
		"testdata/draft4/maxItems.json",
		"testdata/draft4/maxLength.json",
		"testdata/draft4/maxProperties.json",
		"testdata/draft4/minItems.json",
		"testdata/draft4/minLength.json",
		"testdata/draft4/minProperties.json",
		"testdata/draft4/multipleOf.json",
		"testdata/draft4/not.json",
		"testdata/draft4/oneOf.json",
		"testdata/draft4/optional/format.json",
		"testdata/draft4/pattern.json",
		"testdata/draft4/patternProperties.json",
		"testdata/draft4/properties.json",
		"testdata/draft4/required.json",
		"testdata/draft4/type.json",
		"testdata/draft4/uniqueItems.json",

		// disabled due to changes in spec
		// "testdata/draft4/maximum.json",
		// "testdata/draft4/minimum.json",
		// "testdata/draft4/ref.json",
		// "testdata/draft4/refRemote.json",
		// "testdata/draft4/optional/zeroTerminatedFloats.json",

		// TODO(arqu): implement this
		// "testdata/draft4/definitions.json",
		// "testdata/draft4/dependencies.json",
		// "testdata/draft4/items.json",

		// wont fix
		// "testdata/draft4/additionalProperties.json",
		// "testdata/draft4/optional/bignum.json",
		// "testdata/draft4/optional/ecmascript-regex.json",
	})
}

func TestDraft6(t *testing.T) {
	runJSONTests(t, []string{
		"testdata/draft6/additionalItems.json",
		"testdata/draft6/allOf.json",
		"testdata/draft6/anyOf.json",
		"testdata/draft6/boolean_schema.json",
		"testdata/draft6/const.json",
		"testdata/draft6/contains.json",
		"testdata/draft6/default.json",
		"testdata/draft6/enum.json",
		"testdata/draft6/exclusiveMaximum.json",
		"testdata/draft6/exclusiveMinimum.json",
		"testdata/draft6/format.json",
		"testdata/draft6/maximum.json",
		"testdata/draft6/maxItems.json",
		"testdata/draft6/maxLength.json",
		"testdata/draft6/maxProperties.json",
		"testdata/draft6/minimum.json",
		"testdata/draft6/minItems.json",
		"testdata/draft6/minLength.json",
		"testdata/draft6/minProperties.json",
		"testdata/draft6/multipleOf.json",
		"testdata/draft6/not.json",
		"testdata/draft6/oneOf.json",
		"testdata/draft6/pattern.json",
		"testdata/draft6/patternProperties.json",
		"testdata/draft6/properties.json",
		"testdata/draft6/propertyNames.json",
		"testdata/draft6/required.json",
		"testdata/draft6/type.json",
		"testdata/draft6/uniqueItems.json",

		"testdata/draft6/optional/format.json",
		"testdata/draft6/optional/zeroTerminatedFloats.json",

		// TODO(arqu): implement this
		// "testdata/draft6/definitions.json",
		// "testdata/draft6/dependencies.json",
		// "testdata/draft6/items.json",
		// "testdata/draft6/ref.json",

		// wont fix
		// "testdata/draft6/additionalProperties.json",
		// "testdata/draft6/refRemote.json",
		// "testdata/draft6/optional/bignum.json",
		// "testdata/draft6/optional/ecmascript-regex.json",
	})
}

func TestDraft7(t *testing.T) {
	path := "testdata/draft-07_schema.json"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("error reading %s: %s", path, err.Error())
		return
	}

	rsch := &Schema{}
	if err := json.Unmarshal(data, rsch); err != nil {
		t.Errorf("error unmarshaling schema: %s", err.Error())
		return
	}

	runJSONTests(t, []string{
		"testdata/draft7/additionalItems.json",
		"testdata/draft7/allOf.json",
		"testdata/draft7/anyOf.json",
		"testdata/draft7/boolean_schema.json",
		"testdata/draft7/const.json",
		"testdata/draft7/contains.json",
		"testdata/draft7/default.json",
		"testdata/draft7/enum.json",
		"testdata/draft7/exclusiveMaximum.json",
		"testdata/draft7/exclusiveMinimum.json",
		"testdata/draft7/format.json",
		"testdata/draft7/if-then-else.json",
		"testdata/draft7/maximum.json",
		"testdata/draft7/maxItems.json",
		"testdata/draft7/maxLength.json",
		"testdata/draft7/maxProperties.json",
		"testdata/draft7/minimum.json",
		"testdata/draft7/minItems.json",
		"testdata/draft7/minLength.json",
		"testdata/draft7/minProperties.json",
		"testdata/draft7/multipleOf.json",
		"testdata/draft7/not.json",
		"testdata/draft7/oneOf.json",
		"testdata/draft7/pattern.json",
		"testdata/draft7/patternProperties.json",
		"testdata/draft7/properties.json",
		"testdata/draft7/propertyNames.json",
		"testdata/draft7/required.json",
		"testdata/draft7/type.json",
		"testdata/draft7/uniqueItems.json",

		"testdata/draft7/optional/zeroTerminatedFloats.json",
		"testdata/draft7/optional/format/date-time.json",
		"testdata/draft7/optional/format/date.json",
		"testdata/draft7/optional/format/email.json",
		"testdata/draft7/optional/format/hostname.json",
		"testdata/draft7/optional/format/idn-email.json",
		"testdata/draft7/optional/format/idn-hostname.json",
		"testdata/draft7/optional/format/ipv4.json",
		"testdata/draft7/optional/format/ipv6.json",
		"testdata/draft7/optional/format/iri-reference.json",
		"testdata/draft7/optional/format/json-pointer.json",
		"testdata/draft7/optional/format/regex.json",
		"testdata/draft7/optional/format/relative-json-pointer.json",
		"testdata/draft7/optional/format/time.json",
		"testdata/draft7/optional/format/uri-reference.json",
		"testdata/draft7/optional/format/uri-template.json",
		"testdata/draft7/optional/format/uri.json",

		// TODO(arqu): implement this
		// "testdata/draft7/definitions.json",
		// "testdata/draft7/dependencies.json",
		// "testdata/draft7/items.json",
		// "testdata/draft7/ref.json",

		// wont fix
		// "testdata/draft7/additionalProperties.json",
		// "testdata/draft7/refRemote.json",
		// "testdata/draft7/optional/bignum.json",
		// "testdata/draft7/optional/content.json",
		// "testdata/draft7/optional/ecmascript-regex.json",
		// "testdata/draft7/optional/format/iri.json",
	})
}

func TestDraft2019_09(t *testing.T) {
	path := "testdata/draft2019-09_schema.json"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("error reading %s: %s", path, err.Error())
		return
	}

	rsch := &Schema{}
	if err := json.Unmarshal(data, rsch); err != nil {
		t.Errorf("error unmarshaling schema: %s", err.Error())
		return
	}

	runJSONTests(t, []string{
		"testdata/draft2019-09/additionalItems.json",
		// "testdata/draft2019-09/additionalProperties.json",
		"testdata/draft2019-09/allOf.json",
		"testdata/draft2019-09/anchor.json",
		"testdata/draft2019-09/anyOf.json",
		"testdata/draft2019-09/boolean_schema.json",
		"testdata/draft2019-09/const.json",
		"testdata/draft2019-09/contains.json",
		"testdata/draft2019-09/default.json",
		"testdata/draft2019-09/defs.json",
		"testdata/draft2019-09/dependentRequired.json",
		"testdata/draft2019-09/dependentSchemas.json",
		"testdata/draft2019-09/enum.json",
		"testdata/draft2019-09/exclusiveMaximum.json",
		"testdata/draft2019-09/exclusiveMinimum.json",
		"testdata/draft2019-09/format.json",
		"testdata/draft2019-09/if-then-else.json",
		"testdata/draft2019-09/items.json",
		"testdata/draft2019-09/maximum.json",
		"testdata/draft2019-09/maxItems.json",
		"testdata/draft2019-09/maxLength.json",
		"testdata/draft2019-09/maxProperties.json",
		"testdata/draft2019-09/minimum.json",
		"testdata/draft2019-09/minItems.json",
		"testdata/draft2019-09/minLength.json",
		"testdata/draft2019-09/minProperties.json",
		"testdata/draft2019-09/multipleOf.json",
		"testdata/draft2019-09/not.json",
		"testdata/draft2019-09/oneOf.json",
		"testdata/draft2019-09/pattern.json",
		"testdata/draft2019-09/patternProperties.json",
		"testdata/draft2019-09/properties.json",
		"testdata/draft2019-09/propertyNames.json",
		"testdata/draft2019-09/ref.json",
		"testdata/draft2019-09/required.json",
		"testdata/draft2019-09/type.json",
		// "testdata/draft2019-09/unevaluatedProperties.json",
		// "testdata/draft2019-09/unevaluatedItems.json",
		"testdata/draft2019-09/uniqueItems.json",

		"testdata/draft2019-09/optional/zeroTerminatedFloats.json",
		"testdata/draft2019-09/optional/format/date-time.json",
		"testdata/draft2019-09/optional/format/date.json",
		"testdata/draft2019-09/optional/format/email.json",
		"testdata/draft2019-09/optional/format/hostname.json",
		"testdata/draft2019-09/optional/format/idn-email.json",
		"testdata/draft2019-09/optional/format/idn-hostname.json",
		"testdata/draft2019-09/optional/format/ipv4.json",
		"testdata/draft2019-09/optional/format/ipv6.json",
		"testdata/draft2019-09/optional/format/iri-reference.json",
		"testdata/draft2019-09/optional/format/json-pointer.json",
		"testdata/draft2019-09/optional/format/regex.json",
		"testdata/draft2019-09/optional/format/relative-json-pointer.json",
		"testdata/draft2019-09/optional/format/time.json",
		"testdata/draft2019-09/optional/format/uri-reference.json",
		"testdata/draft2019-09/optional/format/uri-template.json",
		"testdata/draft2019-09/optional/format/uri.json",

		// TODO(arqu): investigate further, test is modified because
		// if does not formally validate and simply returns
		// when no then or else is present
		"testdata/draft2019-09/unevaluatedProperties_modified.json",
		"testdata/draft2019-09/unevaluatedItems_modified.json",

		// TODO(arqu): investigate further, test is modified because of inconsistent
		// expectations from spec on how evaluated properties are tracked between
		// additionalProperties and unevaluatedProperties
		"testdata/draft2019-09/additionalProperties_modified.json",

		// wont fix
		// "testdata/draft2019-09/refRemote.json",
		// "testdata/draft2019-09/optional/bignum.json",
		// "testdata/draft2019-09/optional/content.json",
		// "testdata/draft2019-09/optional/ecmascript-regex.json",
		// "testdata/draft2019-09/optional/refOfUnknownKeyword.json",

		// TODO(arqu): iri fails on IPV6 not having [] around the address
		// which was a legal format in draft7
		// introduced: https://github.com/json-schema-org/JSON-Schema-Test-Suite/commit/2146b02555b163da40ae98e60bf36b2c2f8d4bd0#diff-b2ca98716e146559819bc49635a149a9
		// relevant RFC: https://tools.ietf.org/html/rfc3986#section-3.2.2
		// relevant 'net/url' package discussion: https://github.com/golang/go/issues/31024
		// "testdata/draft2019-09/optional/format/iri.json",
	})
}

// TestSet is a json-based set of tests
// JSON-Schema comes with a lovely JSON-based test suite:
// https://github.com/json-schema-org/JSON-Schema-Test-Suite
type TestSet struct {
	Description string     `json:"description"`
	Schema      *Schema    `json:"schema"`
	Tests       []TestCase `json:"tests"`
}

type TestCase struct {
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
	Valid       bool        `json:"valid"`
}

func runJSONTests(t *testing.T, testFilepaths []string) {
	tests := 0
	passed := 0
	ctx := context.Background()
	for _, path := range testFilepaths {
		t.Run(path, func(t *testing.T) {
			base := filepath.Base(path)
			testSets := []*TestSet{}
			data, err := ioutil.ReadFile(path)
			if err != nil {
				t.Errorf("error loading test file: %s", err.Error())
				return
			}

			if err := json.Unmarshal(data, &testSets); err != nil {
				t.Errorf("error unmarshaling test set %s from JSON: %s", base, err.Error())
				return
			}

			for _, ts := range testSets {
				sc := ts.Schema
				for i, c := range ts.Tests {
					tests++

					// Ensure we can register keywords in go routines
					RegisterKeyword(fmt.Sprintf("content-encoding-%d", tests), newContentEncoding)

					validationState := sc.Validate(ctx, c.Data)
					if validationState.IsValid() != c.Valid {
						t.Errorf("%s: %s test case %d: %s. error: %s", base, ts.Description, i, c.Description, *validationState.Errs)
					} else {
						passed++
					}
				}
			}
		})
	}
	t.Logf("%d/%d tests passed", passed, tests)
}

func TestDataType(t *testing.T) {
	type customObject struct{}
	type customNumber float64

	cases := []struct {
		data   interface{}
		expect string
	}{
		{nil, "null"},
		{float64(4), "integer"},
		{float64(4.0), "integer"},
		{float64(4.5), "number"},
		{customNumber(4.5), "number"},
		{true, "boolean"},
		{"foo", "string"},
		{[]interface{}{}, "array"},
		{[0]interface{}{}, "array"},
		{map[string]interface{}{}, "object"},
		{struct{}{}, "object"},
		{customObject{}, "object"},
		{uint8(42), "integer"},
		{uint16(42), "integer"},
		{uint32(42), "integer"},
		{uint64(42), "integer"},
		{int8(42), "integer"},
		{int16(42), "integer"},
		{int32(42), "integer"},
		{int64(42), "integer"},
		{float32(42), "integer"},
		{float32(42.0), "integer"},
		{float32(42.5), "number"},
		// special cases which should pass with type hints
		{"true", "boolean"},
		{4.0, "number"},
	}

	for i, c := range cases {
		got := DataTypeWithHint(c.data, c.expect)
		if got != c.expect {
			t.Errorf("case %d result mismatch. expected: '%s', got: '%s'", i, c.expect, got)
		}
	}
}

func TestJSONCoding(t *testing.T) {
	cases := []string{
		"testdata/coding/false.json",
		"testdata/coding/true.json",
		"testdata/coding/std.json",
		"testdata/coding/booleans.json",
		"testdata/coding/conditionals.json",
		"testdata/coding/numeric.json",
		"testdata/coding/objects.json",
		"testdata/coding/strings.json",
	}

	for i, c := range cases {
		data, err := ioutil.ReadFile(c)
		if err != nil {
			t.Errorf("case %d error reading file: %s", i, err.Error())
			continue
		}

		rs := &Schema{}
		if err := json.Unmarshal(data, rs); err != nil {
			t.Errorf("case %d error unmarshaling from json: %s", i, err.Error())
			continue
		}

		output, err := json.MarshalIndent(rs, "", "  ")
		if err != nil {
			t.Errorf("case %d error marshaling to JSON: %s", i, err.Error())
			continue
		}

		if !bytes.Equal(data, output) {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(string(data), string(output), true)
			if len(diffs) == 0 {
				t.Logf("case %d bytes were unequal but computed no difference between results", i)
				continue
			}

			t.Errorf("case %d %s mismatch:\n", i, c)
			t.Errorf("diff:\n%s", dmp.DiffPrettyText(diffs))
			t.Errorf("expected:\n%s", string(data))
			t.Errorf("got:\n%s", string(output))
			continue
		}
	}
}

func TestValidateBytes(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		schema string
		input  string
		errors []string
	}{
		{`true`, `"just a string yo"`, nil},
		{`{"type":"array", "items": {"type":"string"}}`,
			`[1,false,null]`,
			[]string{
				`/0: 1 type should be string, got integer`,
				`/1: false type should be string, got boolean`,
				`/2: type should be string, got null`,
			}},
		{`{
		"type": "object",
		"properties" : {
		},
		"additionalProperties" : false
	}`,
			`{
	"port": 80
}`,
			[]string{
				`/port: {"port":80} additional properties are not allowed`,
			}},
	}

	for i, c := range cases {
		rs := &Schema{}
		if err := rs.UnmarshalJSON([]byte(c.schema)); err != nil {
			t.Errorf("case %d error parsing %s", i, err.Error())
			continue
		}

		errors, err := rs.ValidateBytes(ctx, []byte(c.input))
		if err != nil {
			t.Errorf("case %d error validating: %s", i, err.Error())
			continue
		}

		if len(errors) != len(c.errors) {
			t.Errorf("case %d: error length mismatch. expected: '%d', got: '%d'", i, len(c.errors), len(errors))
			t.Errorf("%v", errors)
			continue
		}

		for j, e := range errors {
			if e.Error() != c.errors[j] {
				t.Errorf("case %d: validation error %d mismatch. expected: '%s', got: '%s'", i, j, c.errors[j], e.Error())
				continue
			}
		}
	}
}

func BenchmarkAdditionalItems(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make([]interface{}, sampleSize)
			for i := 0; i < sampleSize; i++ {
				data[i] = float64(i)
			}
			return `{
				"items": {},
				"additionalItems": false
			}`, data
		},
	)
}

func BenchmarkAdditionalProperties(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make(map[string]interface{}, sampleSize)
			for i := 0; i < sampleSize; i++ {
				p := fmt.Sprintf("p%v", i)
				data[p] = struct{}{}
			}
			d, err := json.Marshal(data)
			if err != nil {
				b.Errorf("unable to marshal data: %v", err)
				return "", nil
			}
			return `{
				"properties": ` + string(d) + `,
				"additionalProperties": false
			}`, data
		},
	)
}

func BenchmarkConst(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make(map[string]interface{}, sampleSize)
			for i := 0; i < sampleSize; i++ {
				data[fmt.Sprintf("p%v", i)] = fmt.Sprintf("p%v", 2*i)
			}
			d, err := json.Marshal(data)
			if err != nil {
				b.Errorf("unable to marshal data: %v", err)
				return "", nil
			}
			return `{
				"const": ` + string(d) + `
			}`, data
		},
	)
}

func BenchmarkContains(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make([]interface{}, sampleSize)
			for i := 0; i < sampleSize; i++ {
				data[i] = float64(i)
			}
			return `{
				"contains": { "const": ` + strconv.Itoa(sampleSize-1) + ` }
			}`, data
		},
	)
}

func BenchmarkDependencies(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make(map[string]interface{}, sampleSize)
			deps := []string{}
			for i := 0; i < sampleSize; i++ {
				p := fmt.Sprintf("p%v", i)
				data[p] = fmt.Sprintf("p%v", 2*i)
				if i != 0 {
					deps = append(deps, p)
				}
			}
			d, err := json.Marshal(deps)
			if err != nil {
				b.Errorf("unable to marshal data: %v", err)
				return "", nil
			}
			return `{
				"dependencies": {"p0": ` + string(d) + `}
			}`, data
		},
	)
}

func BenchmarkEnum(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make([]interface{}, sampleSize)
			for i := 0; i < sampleSize; i++ {
				data[i] = float64(i)
			}
			d, err := json.Marshal(data)
			if err != nil {
				b.Errorf("unable to marshal data: %v", err)
				return "", nil
			}
			return `{
				"enum": ` + string(d) + `
			}`, float64(sampleSize / 2)
		},
	)
}

func BenchmarkMaximum(b *testing.B) {
	runBenchmark(b, func(sampleSize int) (string, interface{}) {
		return `{
			"maximum": 3
		}`, float64(2)
	})
}

func BenchmarkMinimum(b *testing.B) {
	runBenchmark(b, func(sampleSize int) (string, interface{}) {
		return `{
			"minimum": 3
		}`, float64(4)
	})
}

func BenchmarkExclusiveMaximum(b *testing.B) {
	runBenchmark(b, func(sampleSize int) (string, interface{}) {
		return `{
			"exclusiveMaximum": 3
		}`, float64(2)
	})
}

func BenchmarkExclusiveMinimum(b *testing.B) {
	runBenchmark(b, func(sampleSize int) (string, interface{}) {
		return `{
			"exclusiveMinimum": 3
		}`, float64(4)
	})
}

func BenchmarkMaxItems(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make([]interface{}, sampleSize)
			for i := 0; i < sampleSize; i++ {
				data[i] = float64(i)
			}
			return `{
				"maxItems": 10000
			}`, data
		},
	)
}

func BenchmarkMinItems(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make([]interface{}, sampleSize)
			for i := 0; i < sampleSize; i++ {
				data[i] = float64(i)
			}
			return `{
				"minItems": 1
			}`, data
		},
	)
}

func BenchmarkMaxLength(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make([]rune, sampleSize)
			for i := 0; i < sampleSize; i++ {
				data[i] = 'a'
			}
			return `{
				"maxLength": ` + strconv.Itoa(sampleSize) + `
			}`, string(data)
		},
	)
}

func BenchmarkMinLength(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make([]rune, sampleSize)
			for i := 0; i < sampleSize; i++ {
				data[i] = 'a'
			}
			return `{
				"minLength": 1
			}`, string(data)
		},
	)
}

func BenchmarkMaxProperties(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make(map[string]interface{}, sampleSize)
			for i := 0; i < sampleSize; i++ {
				data[fmt.Sprintf("p%v", i)] = fmt.Sprintf("p%v", 2*i)
			}
			return `{
				"maxProperties": ` + strconv.Itoa(sampleSize) + `
			}`, data
		},
	)
}

func BenchmarkMinProperties(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make(map[string]interface{}, sampleSize)
			for i := 0; i < sampleSize; i++ {
				data[fmt.Sprintf("p%v", i)] = fmt.Sprintf("p%v", 2*i)
			}
			return `{
				"minProperties": 1
			}`, data
		},
	)
}

func BenchmarkMultipleOf(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			return `{
				"multipleOf": 2
			}`, float64(42)
		},
	)
}

func BenchmarkPattern(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make([]rune, sampleSize)
			for i := 0; i < sampleSize; i++ {
				data[i] = 'a'
			}
			return `{
				"pattern": "^a*$"
			}`, string(data)
		},
	)
}

func BenchmarkType(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make(map[string]interface{}, sampleSize)
			var schema strings.Builder

			for i := 0; i < sampleSize; i++ {
				propNull := fmt.Sprintf("n%v", nil)
				propBool := fmt.Sprintf("b%v", i)
				propInt := fmt.Sprintf("i%v", i)
				propFloat := fmt.Sprintf("f%v", i)
				propStr := fmt.Sprintf("s%v", i)
				propArr := fmt.Sprintf("a%v", i)
				propObj := fmt.Sprintf("o%v", i)

				data[propBool] = true
				data[propInt] = float64(42)
				data[propFloat] = float64(42.5)
				data[propStr] = "foobar"
				data[propArr] = []interface{}{interface{}(1), interface{}(2), interface{}(3)}
				data[propObj] = struct{}{}

				schema.WriteString(fmt.Sprintf(`"%v": { "type": "null" },`, propNull))
				schema.WriteString(fmt.Sprintf(`"%v": { "type": "boolean" },`, propBool))
				schema.WriteString(fmt.Sprintf(`"%v": { "type": "integer" },`, propInt))
				schema.WriteString(fmt.Sprintf(`"%v": { "type": "number" },`, propFloat))
				schema.WriteString(fmt.Sprintf(`"%v": { "type": "string" },`, propStr))
				schema.WriteString(fmt.Sprintf(`"%v": { "type": "array" },`, propArr))
				schema.WriteString(fmt.Sprintf(`"%v": { "type": "object" }`, propObj))

				if i != sampleSize-1 {
					schema.WriteString(",")
				}
			}

			return `{
				"type": "object",
				"properties": { ` + schema.String() + ` }
			}`, data
		},
	)
}

func runBenchmark(b *testing.B, dataFn func(sampleSize int) (string, interface{})) {
	ctx := context.Background()
	for _, sampleSize := range []int{1, 10, 100, 1000} {
		b.Run(fmt.Sprintf("sample size %v", sampleSize), func(b *testing.B) {
			schema, data := dataFn(sampleSize)
			if data == nil {
				b.Skip("data == nil, skipping")
				return
			}

			var validator Schema
			if err := json.Unmarshal([]byte(schema), &validator); err != nil {
				b.Errorf("error parsing schema: %s", err.Error())
				return
			}

			currentState := NewValidationState(&validator)
			currentState.ClearState()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				validator.ValidateKeyword(ctx, currentState, data)
			}
			b.StopTimer()

			if !currentState.IsValid() {
				b.Errorf("error running benchmark: %s", *currentState.Errs)
			}
		})
	}
}
