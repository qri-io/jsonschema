package main

func LoadDraft2019_09() {
	// default keywords
	RegisterKeyword("$id", NewId)

	// standard keywords
	RegisterKeyword("type", NewType)
	RegisterKeyword("enum", NewEnum)
	RegisterKeyword("const", NewConst)

	// numeric keywords
	RegisterKeyword("multipleOf", NewMultipleOf)
	RegisterKeyword("maximum", NewMaximum)
	RegisterKeyword("exclusiveMaximum", NewExclusiveMaximum)
	RegisterKeyword("minimum", NewMinimum)
	RegisterKeyword("exclusiveMinimum", NewExclusiveMinimum)

	// string keywords
	RegisterKeyword("maxLength", NewMaxLength)
	RegisterKeyword("minLength", NewMinLength)
	RegisterKeyword("pattern", NewPattern)
}