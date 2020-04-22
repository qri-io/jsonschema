package keywords

import (
	// kw "github.com/qri-io/jsonschema/keywords"
)

func LoadDraft() {
	// standard keywords
	kw.RegisterKeyword("type", NewVoid)
	kw.RegisterKeyword("enum", NewVoid)
	kw.RegisterKeyword("const", NewVoid)

	// numeric keywords
	kw.RegisterKeyword("multipleOf", NewVoid)
	kw.RegisterKeyword("maximum", NewVoid)
	kw.RegisterKeyword("exclusiveMaximum", NewVoid)
	kw.RegisterKeyword("minimum", NewVoid)
	kw.RegisterKeyword("exclusiveMinimum", NewVoid)
}