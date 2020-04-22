package main

import (
	"fmt"
	"encoding/json"
    "path/filepath"
    "io/ioutil"
)

func main() {
	LoadDraft2019_09()
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

	if jsonData, err := json.MarshalIndent(rs, "", "  "); err != nil {
		panic("marshal schema: " + err.Error())
	} else {
		fmt.Println(string(jsonData))
	}

    TestDraft2019_09()
}

type TestSet struct {
    Description string      `json:"description"`
    Schema      *Schema `json:"schema"`
    Tests       []TestCase  `json:"tests"`
}

type TestCase struct {
    Description string      `json:"description"`
    Data        interface{} `json:"data"`
    Valid       bool        `json:"valid"`
}

func runJSONTests(testFilepaths []string) {
    tests := 0
    passed := 0
    for _, path := range testFilepaths {
        fmt.Println("Testing: " + path)
        base := filepath.Base(path)
        testSets := []*TestSet{}
        data, err := ioutil.ReadFile(path)
        if err != nil {
            fmt.Errorf("error loading test file: %s", err.Error())
            return
        }

        if err := json.Unmarshal(data, &testSets); err != nil {
            fmt.Errorf("error unmarshaling test set %s from JSON: %s", base, err.Error())
            return
        }
        localTests := 0
        localPassed := 0
        for _, ts := range testSets {
            sc := ts.Schema
            for i, c := range ts.Tests {
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
        fmt.Printf("%d/%d tests passed for %s\n",localPassed, localTests, path)
    }
    fmt.Printf("%d/%d tests passed\n", passed, tests)
}

func ReadErrors(errors []KeyError) string {
    result := ""
    for _, err := range errors {
        result += err.Error() + "\n"
    }
    return result
}

func TestDraft2019_09() {
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

 // DefaultSchemaPool["https://json-schema.org/draft/2019-09/schema#"] = &rsch

 runJSONTests([]string{
     // "testdata/draft2019-09/additionalItems.json",
     // "testdata/draft2019-09/additionalProperties.json",
     // "testdata/draft2019-09/allOf.json",
     // "testdata/draft2019-09/anyOf.json",
     "testdata/draft2019-09/boolean_schema.json",
     "testdata/draft2019-09/const.json",
     // "testdata/draft2019-09/contains.json",
     // "testdata/draft2019-09/default.json",
     "testdata/draft2019-09/enum.json",
     "testdata/draft2019-09/exclusiveMaximum.json",
     "testdata/draft2019-09/exclusiveMinimum.json",
     // "testdata/draft2019-09/format.json",
     // "testdata/draft2019-09/if-then-else.json",
     // "testdata/draft2019-09/items.json",
     "testdata/draft2019-09/maximum.json",
     // "testdata/draft2019-09/maxItems.json",
     "testdata/draft2019-09/maxLength.json",
     // "testdata/draft2019-09/maxProperties.json",
     "testdata/draft2019-09/minimum.json",
     // "testdata/draft2019-09/minItems.json",
     "testdata/draft2019-09/minLength.json",
     // "testdata/draft2019-09/minProperties.json",
     "testdata/draft2019-09/multipleOf.json",
     // "testdata/draft2019-09/not.json",
     // "testdata/draft2019-09/oneOf.json",
     "testdata/draft2019-09/pattern.json",
     // "testdata/draft2019-09/patternProperties.json",
     // "testdata/draft2019-09/properties.json",
     // "testdata/draft2019-09/propertyNames.json",
     // "testdata/draft2019-09/required.json",
     "testdata/draft2019-09/type.json",
     // "testdata/draft2019-09/uniqueItems.json",

     // "testdata/draft2019-09/optional/refOfUnknownKeyword.json",
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

     // TODO
     // "testdata/draft2019-09/anchor.json",
     // "testdata/draft2019-09/defs.json",
     // "testdata/draft2019-09/dependentRequired.json",
     // "testdata/draft2019-09/dependentSchemas.json",
     // "testdata/draft2019-09/ref.json",
     // "testdata/draft2019-09/refRemote.json",

     // TODO: requires keeping state of validated items
     // which is something we might not want to support
     // due to performance reasons (esp for large datasets)
     // "testdata/draft2019-09/unevaluatedItems.json",
     // "testdata/draft2019-09/unevaluatedProperties.json",

     // TODO: implement support
     // "testdata/draft2019-09/optional/bignum.json",
     // "testdata/draft2019-09/optional/content.json",
     // "testdata/draft2019-09/optional/ecmascript-regex.json",

     // TODO: iri fails on IPV6 not having [] around the address
     // which was a legal format in draft7
     // introduced: https://github.com/json-schema-org/JSON-Schema-Test-Suite/commit/2146b02555b163da40ae98e60bf36b2c2f8d4bd0#diff-b2ca98716e146559819bc49635a149a9
     // relevant RFC: https://tools.ietf.org/html/rfc3986#section-3.2.2
     // relevant 'net/url' package discussion: https://github.com/golang/go/issues/31024
     // "testdata/draft2019-09/optional/format/iri.json",

 })
}