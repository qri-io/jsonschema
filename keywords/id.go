package keywords

import (
	"fmt"
	js "github.com/qri-io/jsonschema"
)

type Id struct {}

func NewId() Keyword {
	return &Id{}
}

func (i *Id) Validate(propPath string, data interface{}, errs *[]KeyError) {
	fmt.Println("WARN: Using Id Validator - always True")
}

func (i *Id) RegisterSubschemas(uri string, registry *js.SchemaRegistry) {
	
}