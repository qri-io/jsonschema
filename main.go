package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

func main() {
	LoadDraft2019_09()
	// runDraft2019_09()
	// testVB()
	testEM()
}

type mainTestSet struct {
	Description string         `json:"description"`
	Schema      *Schema        `json:"schema"`
	Tests       []mainTestCase `json:"tests"`
}

type mainTestCase struct {
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
	Valid       bool        `json:"valid"`
}

func testEM() {
	cases := []struct {
		schema, doc, message string
	}{
		{`{ "const" : "a value" }`, `"a different value"`, `must equal "a value"`},
	}

	for i, c := range cases {
		rs := &Schema{}
		if err := rs.UnmarshalJSON([]byte(c.schema)); err != nil {
			fmt.Printf("case %d schema is invalid: %s\n", i, err.Error())
			continue
		}

		errs, err := rs.ValidateBytes([]byte(c.doc))
		if err != nil {
			fmt.Printf("case %d error validating: %s\n", i, err)
			continue
		}

		if len(errs) != 1 {
			fmt.Printf("case %d didn't return exactly 1 validation error. got: %d\n", i, len(errs))
			continue
		}

		if errs[0].Message != c.message {
			fmt.Printf("case %d error mismatch. expected '%s', got: '%s'\n", i, c.message, errs[0].Message)
		}
		fmt.Printf("TestEM done\n")
	}
}

func mainRunJSONTests(testFilepaths []string) {
	tests := 0
	passed := 0
	debug := false
	for _, path := range testFilepaths {
		fmt.Println("Testing: " + path)
		base := filepath.Base(path)
		testSets := []*mainTestSet{}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Errorf("error loading test file: %s", err.Error())
			return
		}

		if err := json.Unmarshal(data, &testSets); err != nil {
			fmt.Printf("error unmarshaling test set %s from JSON: %s\n", base, err.Error())
			return
		}
		localTests := 0
		localPassed := 0
		for _, ts := range testSets {
			if debug {
				fmt.Println("\tTest set: " + ts.Description)
			}
			sc := ts.Schema
			for i, c := range ts.Tests {
				if debug {
					buff, _ := json.MarshalIndent(sc, "", " ")
					fmt.Println(string(buff))
					fmt.Println("\tCase: " + strconv.Itoa(i))
				}
				tests++
				localTests++
				got := []KeyError{}
				sc.Validate("/", c.Data, &got)
				valid := len(got) == 0
				if valid != c.Valid {
					fmt.Printf("%s: %s test case %d: %s. error: %s \n", base, ts.Description, i, c.Description, got)
				} else {
					passed++
					localPassed++
				}
			}
		}
		fmt.Printf("%d/%d tests passed for %s\n", localPassed, localTests, path)
	}
	fmt.Printf("%d/%d tests passed\n", passed, tests)
}

func testVB() {
	cases := []struct {
		schema string
		input  string
		errors []string
	}{
		// {`true`, `"just a string yo"`, nil},
		{`{"type":"array", "items": {"type":"string"}}`,
			`[1,false,null]`,
			[]string{
				`/0: 1 type should be string`,
				`/1: false type should be string`,
				`/2: type should be string`,
			}},
	}

	for i, c := range cases {
		rs := &Schema{}
		if err := rs.UnmarshalJSON([]byte(c.schema)); err != nil {
			fmt.Printf("case %d error parsing %s\n", i, err.Error())
			continue
		}

		errors, err := rs.ValidateBytes([]byte(c.input))
		if err != nil {
			fmt.Printf("case %d error validating: %s\n", i, err.Error())
			continue
		}

		if len(errors) != len(c.errors) {
			fmt.Printf("case %d: error length mismatch. expected: '%d', got: '%d'\n", i, len(c.errors), len(errors))
			fmt.Printf("%v", errors)
			continue
		}

		for j, e := range errors {
			if e.Error() != c.errors[j] {
				fmt.Printf("case %d: validation error %d mismatch. expected: '%s', got: '%s'\n", i, j, c.errors[j], e.Error())
				continue
			}
		}
	}
}

func runDraft2019_09() {
	path := "testdata/draft2019-09_schema.json"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Errorf("error reading %s: %s", path, err.Error())
		return
	}

	rsch := &Schema{}
	if err := json.Unmarshal(data, rsch); err != nil {
		fmt.Errorf("error unmarshaling schema: %s", err.Error())
		return
	}

	mainRunJSONTests([]string{
		// "testdata/draft2019-09/additionalItems.json",
		// "testdata/draft2019-09/additionalProperties.json",
		// "testdata/draft2019-09/allOf.json",
		// "testdata/draft2019-09/anchor.json",
		// "testdata/draft2019-09/anyOf.json",
		// "testdata/draft2019-09/boolean_schema.json",
		// "testdata/draft2019-09/const.json",
		// "testdata/draft2019-09/contains.json",
		// "testdata/draft2019-09/default.json",
		// "testdata/draft2019-09/defs.json",
		// "testdata/draft2019-09/dependentRequired.json",
		// "testdata/draft2019-09/dependentSchemas.json",
		// "testdata/draft2019-09/enum.json",
		// "testdata/draft2019-09/exclusiveMaximum.json",
		// "testdata/draft2019-09/exclusiveMinimum.json",
		// "testdata/draft2019-09/format.json",
		// "testdata/draft2019-09/if-then-else.json",
		// "testdata/draft2019-09/items.json",
		// "testdata/draft2019-09/maximum.json",
		// "testdata/draft2019-09/maxItems.json",
		// "testdata/draft2019-09/maxLength.json",
		// "testdata/draft2019-09/maxProperties.json",
		// "testdata/draft2019-09/minimum.json",
		// "testdata/draft2019-09/minItems.json",
		// "testdata/draft2019-09/minLength.json",
		// "testdata/draft2019-09/minProperties.json",
		// "testdata/draft2019-09/multipleOf.json",
		// "testdata/draft2019-09/not.json",
		// "testdata/draft2019-09/oneOf.json",
		// "testdata/draft2019-09/pattern.json",
		// "testdata/draft2019-09/patternProperties.json",
		// "testdata/draft2019-09/properties.json",
		// "testdata/draft2019-09/propertyNames.json",
		// "testdata/draft2019-09/required.json",
		// "testdata/draft2019-09/type.json",
		// "testdata/draft2019-09/uniqueItems.json",

		// "testdata/draft2019-09/optional/zeroTerminatedFloats.json",
		// "testdata/draft2019-09/optional/format/date-time.json",
		// "testdata/draft2019-09/optional/format/date.json",
		// "testdata/draft2019-09/optional/format/email.json",
		// "testdata/draft2019-09/optional/format/hostname.json",
		// "testdata/draft2019-09/optional/format/idn-email.json",
		// "testdata/draft2019-09/optional/format/idn-hostname.json",
		// "testdata/draft2019-09/optional/format/ipv4.json",
		// "testdata/draft2019-09/optional/format/ipv6.json",
		// "testdata/draft2019-09/optional/format/iri-reference.json",
		// "testdata/draft2019-09/optional/format/json-pointer.json",
		// "testdata/draft2019-09/optional/format/regex.json",
		// "testdata/draft2019-09/optional/format/relative-json-pointer.json",
		// "testdata/draft2019-09/optional/format/time.json",
		// "testdata/draft2019-09/optional/format/uri-reference.json",
		// "testdata/draft2019-09/optional/format/uri-template.json",
		// "testdata/draft2019-09/optional/format/uri.json",

		// TODO(arqu): finalize implementations
		// "testdata/draft2019-09/ref.json",

		// TODO(arqu): requires keeping state of validated items
		// which is something we might not want to support
		// due to performance reasons (esp for large/deeply nested schemas)
		// "testdata/draft2019-09/unevaluatedItems.json",
		// "testdata/draft2019-09/unevaluatedProperties.json",

		// TODO(arqu): wont fix
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
