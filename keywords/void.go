package keywords

import (
	"fmt"
)

type Void struct {}

func NewVoid() Keyword {
	return &Void{}
}

func (vo *Void) Validate(propPath string, data interface{}, errs *[]KeyError) {
	fmt.Println("WARN: Using Void Validator - always True")
}

func (vo *Void) RegisterSubschemas(uri string, registry *SchemaRegistry) {
	
}