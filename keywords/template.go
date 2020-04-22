package keywords

import (
	"fmt"
)

type Template struct {}

func NewTemplate() Keyword {
	return &Template{}
}

func (t *Template) Validate(propPath string, data interface{}, errs *[]KeyError) {
	fmt.Println("WARN: Using Template Validator - always True")
}

func (t *Template) RegisterSubschemas(uri string, registry *SchemaRegistry) {
	
}