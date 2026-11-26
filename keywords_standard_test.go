package jsonschema

import (
	"context"
	"encoding/json"
	"testing"
)

func TestEnumValidate(t *testing.T) {
	ctx := context.Background()
	colors := `{"enum":["red","green","blue"]}`
	escaped := `{"enum":["a\"b","tab\there","back\\slash","sn\u00f8","smörgås","\u0041"]}`
	mixed := `{"enum":[1,2.5,true,null,"1",[1,"a"],{"a":[1]}]}`

	cases := []struct {
		schema string
		input  string
		errors []string
	}{
		{colors, `"red"`, nil},
		{colors, `"green"`, nil},
		{colors, `"blue"`, nil},
		{colors, `"Red"`, []string{`/: "Red" should be one of ["red", "green", "blue"]`}},
		{colors, `"re"`, []string{`/: "re" should be one of ["red", "green", "blue"]`}},
		{colors, `""`, []string{`/: "" should be one of ["red", "green", "blue"]`}},
		{colors, `"\"red\""`, []string{`/: "\"red\"" should be one of ["red", "green", "blue"]`}},
		{colors, `["red"]`, []string{`/: ["red"] should be one of ["red", "green", "blue"]`}},
		{colors, `null`, []string{`/: should be one of ["red", "green", "blue"]`}},
		{`{"enum":[""]}`, `""`, nil},
		{`{"enum":["red"]}`, `"blue"`, []string{`/: "blue" should be one of ["red"]`}},

		{escaped, `"a\"b"`, nil},
		{escaped, `"tab\there"`, nil},
		{escaped, `"back\\slash"`, nil},
		{escaped, `"snø"`, nil},
		{escaped, `"sm\u00f6rg\u00e5s"`, nil},
		{escaped, `"A"`, nil},
		{escaped, `"a\\\"b"`, []string{`/: "a\\\"b" should be one of ["a\"b", "tab\there", "back\\slash", "sn\u00f8", "smörgås", "\u0041"]`}},
		{escaped, `"tab\\there"`, []string{`/: "tab\\there" should be one of ["a\"b", "tab\there", "back\\slash", "sn\u00f8", "smörgås", "\u0041"]`}},
		{escaped, `"sn\\u00f8"`, []string{`/: "sn\\u00f8" should be one of ["a\"b", "tab\there", "back\\slash", "sn\u00f8", "smörgås", "\u0041"]`}},
		{escaped, `"\\u0041"`, []string{`/: "\\u0041" should be one of ["a\"b", "tab\there", "back\\slash", "sn\u00f8", "smörgås", "\u0041"]`}},

		{mixed, `1`, nil},
		{mixed, `1.0`, nil},
		{mixed, `2.5`, nil},
		{mixed, `true`, nil},
		{mixed, `null`, nil},
		{mixed, `"1"`, nil},
		{mixed, `[1,"a"]`, nil},
		{mixed, `{"a":[1]}`, nil},
		{mixed, `2`, []string{`/: 2 should be one of [1, 2.5, true, null, "1", [1,"a"], {"a":[1]}]`}},
		{mixed, `false`, []string{`/: false should be one of [1, 2.5, true, null, "1", [1,"a"], {"a":[1]}]`}},
		{mixed, `"true"`, []string{`/: "true" should be one of [1, 2.5, true, null, "1", [1,"a"], {"a":[1]}]`}},
		{mixed, `"null"`, []string{`/: "null" should be one of [1, 2.5, true, null, "1", [1,"a"], {"a":[1]}]`}},
		{mixed, `["a",1]`, []string{`/: ["a",1] should be one of [1, 2.5, true, null, "1", [1,"a"], {"a":[1]}]`}},
		{mixed, `{"a":[2]}`, []string{`/: {"a":[2]} should be one of [1, 2.5, true, null, "1", [1,"a"], {"a":[1]}]`}},

		{`{"properties":{"color":` + colors + `}}`, `{"color":"green"}`, nil},
		{`{"properties":{"color":` + colors + `}}`,
			`{"color":"pink"}`,
			[]string{`/color: "pink" should be one of ["red", "green", "blue"]`}},
		{`{"items":` + colors + `}`,
			`["red","pink","blue",3]`,
			[]string{
				`/1: "pink" should be one of ["red", "green", "blue"]`,
				`/3: 3 should be one of ["red", "green", "blue"]`,
			}},
		{`{"oneOf":[
			{"properties":{"kind":{"enum":["a"]}},"required":["kind"]},
			{"properties":{"kind":{"enum":["b"]}},"required":["kind"]}
		]}`, `{"kind":"b"}`, nil},
		{`{"oneOf":[
			{"properties":{"kind":{"enum":["a"]}},"required":["kind"]},
			{"properties":{"kind":{"enum":["b"]}},"required":["kind"]}
		]}`,
			`{"kind":"c"}`,
			[]string{`/: {"kind":"c"} did not match any of the specified OneOf schemas`}},
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

// TestEnumMatchesConst checks that an enum accepts exactly the data that one
// of its values accepts as a const, for values that look alike before or
// after decoding.
func TestEnumMatchesConst(t *testing.T) {
	ctx := context.Background()
	values := []string{
		`""`, `"a"`, `"A"`, `"\u0041"`, `"a "`, ` "a" `, `"a\"b"`, `"a\\\"b"`, `"a\\b"`, `"a\/b"`, `"a/b"`,
		`"tab\there"`, `"tab\\there"`, `"sn\u00f8"`, `"snø"`, `"sn\\u00f8"`, "\"del\x7f\"", `"\u007f"`,
		`"1"`, `1`, `1.0`, `1e0`, `-0`, `0`, `"true"`, `true`, `false`, `"null"`, `null`,
		`[]`, `["a"]`, `[ "a" ]`, `["a","b"]`, `{}`, `{"a":"b"}`, `{ "a" : "b" }`, `{"a":["b"]}`,
		`"`, `"a`, `a"`, `"a"b"`, "\"a\tb\"", "\"\xff\"", `not json`, ``,
	}

	data := []interface{}{}
	for _, value := range values {
		var decoded interface{}
		if err := json.Unmarshal([]byte(value), &decoded); err == nil {
			data = append(data, decoded)
		}
	}
	data = append(data, `a\"b`, `tab\there`, `sn\u00f8`, "a\tb", "\xff", "\ufffd", `"a"`, 1, int64(1), float32(1), []string{"a"})

	accepts := func(keyword Keyword, d interface{}) bool {
		state := NewValidationState(&Schema{})
		keyword.ValidateKeyword(ctx, state, d)
		return state.IsValid()
	}

	for _, value := range values {
		con := Const(value)
		for _, d := range data {
			if got, expect := accepts(&Enum{con}, d), accepts(&con, d); got != expect {
				t.Errorf("enum [%s] with data %#v: expected valid: %v, got: %v", value, d, expect, got)
			}
			// a value that does not match must not hide one that does
			if got, expect := accepts(&Enum{Const(`"no match"`), con}, d), accepts(&con, d); got != expect {
				t.Errorf("enum [\"no match\", %s] with data %#v: expected valid: %v, got: %v", value, d, expect, got)
			}
		}
	}
}

func TestPlainJSONString(t *testing.T) {
	cases := []struct {
		raw   string
		plain bool
	}{
		{`""`, true},
		{`"a"`, true},
		{`"red"`, true},
		{`"a b ~ / {}[]:,'"`, true},
		{"\"del\x7f\"", true},
		// escape sequences
		{`"a\"b"`, false},
		{`"a\\b"`, false},
		{`"a\/b"`, false},
		{`"tab\there"`, false},
		{`"sn\u00f8"`, false},
		// not ASCII
		{`"snø"`, false},
		{"\"\xff\"", false},
		// not a JSON string
		{"\"a\tb\"", false},
		{"\"a\x00b\"", false},
		{`"a"b"`, false},
		{`"`, false},
		{`"a`, false},
		{`a"`, false},
		{``, false},
		{`1`, false},
		{`true`, false},
		{`null`, false},
		{`["a"]`, false},
		{`{"a":"b"}`, false},
		// surrounding whitespace is left to the decoder
		{` "a"`, false},
		{`"a" `, false},
	}

	for _, c := range cases {
		got, plain := plainJSONString([]byte(c.raw))
		if plain != c.plain {
			t.Errorf("%s: expected plain: %v, got: %v", c.raw, c.plain, plain)
			continue
		}
		if !plain {
			continue
		}

		var decoded string
		if err := json.Unmarshal([]byte(c.raw), &decoded); err != nil {
			t.Errorf("%s: plain string does not decode: %s", c.raw, err.Error())
		} else if decoded != string(got) {
			t.Errorf("%s: expected the decoded string '%s', got: '%s'", c.raw, decoded, got)
		}
	}
}

// TestEnumValidateKeyword calls the keyword directly: values that are not
// valid JSON never match, and validating leaves no trace in the validation
// state other than the error for data that is not listed.
func TestEnumValidateKeyword(t *testing.T) {
	ctx := context.Background()
	enum := Enum{Const(`"red"`), Const(`not json`), Const(`1`), Const(`"blue"`)}

	for _, data := range []interface{}{"red", "blue", float64(1)} {
		state := NewValidationState(&Schema{})
		enum.ValidateKeyword(ctx, state, data)
		if !state.IsValid() {
			t.Errorf("%v: expected no errors, got: %v", data, *state.Errs)
		}
		if len(*state.EvaluatedPropertyNames) != 0 || len(*state.LocalEvaluatedPropertyNames) != 0 || len(state.Misc) != 0 {
			t.Errorf("%v: expected the state to be left untouched", data)
		}
	}

	state := NewValidationState(&Schema{})
	state.DescendInstance("a", "b")
	enum.ValidateKeyword(ctx, state, "green")
	expect := `/a/b: "green" should be one of ["red", not json, 1, "blue"]`
	if len(*state.Errs) != 1 || (*state.Errs)[0].Error() != expect {
		t.Errorf("expected the single error '%s', got: %v", expect, *state.Errs)
	}
}
