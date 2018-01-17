# jsonschema
[![Qri](https://img.shields.io/badge/made%20by-qri-magenta.svg?style=flat-square)](https://qri.io)
[![GoDoc](https://godoc.org/github.com/qri-io/jsonschema?status.svg)](http://godoc.org/github.com/qri-io/jsonschema)
[![License](https://img.shields.io/github/license/qri-io/jsonschema.svg?style=flat-square)](./LICENSE)
[![Codecov](https://img.shields.io/codecov/c/github/qri-io/jsonschema.svg?style=flat-square)](https://codecov.io/gh/qri-io/jsonschema)
[![CI](https://img.shields.io/circleci/project/github/qri-io/jsonschema.svg?style=flat-square)](https://circleci.com/gh/qri-io/jsonschema)
[![Go Report Card](https://goreportcard.com/badge/github.com/qri-io/jsonschema)](https://goreportcard.com/report/github.com/qri-io/jsonschema)

### 🚧🚧 Under Construction 🚧🚧
golang implementation of http://json-schema.org/


### Getting Involved

We would love involvement from more people! If you notice any errors or would
like to submit changes, please see our
[Contributing Guidelines](./.github/CONTRIBUTING.md).

### Developing

We've set up a separate document for [developer guidelines](https://github.com/qri-io/jsonschema/blob/master/DEVELOPERS.md)!

## Usage

Here's a quick example pulled from the [godoc](https://godoc.org/github.com/qri-io/jsonschema):

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/qri-io/jsonschema"
)

func main() {
	var schemaData = []byte(`{
      "title": "Person",
      "type": "object",
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

	rs := &jsonschema.RootSchema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}

	var valid = []byte(`{
    "firstName" : "George",
    "lastName" : "Michael"
    }`)
	if err := rs.ValidateBytes(valid); err != nil {
		panic(err)
	}

	var invalidPerson = []byte(`{
    "firstName" : "Prince"
    }`)
	err := rs.ValidateBytes(invalidPerson)
	fmt.Println(err.Error())

	var invalidFriend = []byte(`{
    "firstName" : "Jay",
    "lastName" : "Z",
    "friends" : [{
      "firstName" : "Nas"
      }]
    }`)
	err = rs.ValidateBytes(invalidFriend)
	fmt.Println(err)
}
```