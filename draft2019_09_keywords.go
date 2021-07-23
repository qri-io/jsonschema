package jsonschema

// LoadDraft2019_09 loads the keywords for schema validation
// based on draft2019_09
// this is also the default keyword set loaded automatically
// if no other is loaded
func LoadDraft2019_09() {
	r, release := getGlobalKeywordRegistry()
	defer release()

	r.LoadDraft2019_09()
}

// DefaultIfEmpty populates the KeywordRegistry with the 2019_09
// jsonschema draft specification only if the registry is empty.
func (r *KeywordRegistry) DefaultIfEmpty() {
	if !r.IsRegistryLoaded() {
		r.LoadDraft2019_09()
	}
}

// LoadDraft2019_09 loads the keywords for schema validation
// based on draft2019_09
// this is also the default keyword set loaded automatically
// if no other is loaded
func (r *KeywordRegistry) LoadDraft2019_09() {
	// core keywords
	r.RegisterKeyword("$schema", NewSchemaURI)
	r.RegisterKeyword("$id", NewID)
	r.RegisterKeyword("description", NewDescription)
	r.RegisterKeyword("title", NewTitle)
	r.RegisterKeyword("$comment", NewComment)
	r.RegisterKeyword("examples", NewExamples)
	r.RegisterKeyword("readOnly", NewReadOnly)
	r.RegisterKeyword("writeOnly", NewWriteOnly)
	r.RegisterKeyword("$ref", NewRef)
	r.RegisterKeyword("$recursiveRef", NewRecursiveRef)
	r.RegisterKeyword("$anchor", NewAnchor)
	r.RegisterKeyword("$recursiveAnchor", NewRecursiveAnchor)
	r.RegisterKeyword("$defs", NewDefs)
	r.RegisterKeyword("default", NewDefault)

	r.SetKeywordOrder("$ref", 0)
	r.SetKeywordOrder("$recursiveRef", 0)

	// standard keywords
	r.RegisterKeyword("type", NewType)
	r.RegisterKeyword("enum", NewEnum)
	r.RegisterKeyword("const", NewConst)

	// numeric keywords
	r.RegisterKeyword("multipleOf", NewMultipleOf)
	r.RegisterKeyword("maximum", NewMaximum)
	r.RegisterKeyword("exclusiveMaximum", NewExclusiveMaximum)
	r.RegisterKeyword("minimum", NewMinimum)
	r.RegisterKeyword("exclusiveMinimum", NewExclusiveMinimum)

	// string keywords
	r.RegisterKeyword("maxLength", NewMaxLength)
	r.RegisterKeyword("minLength", NewMinLength)
	r.RegisterKeyword("pattern", NewPattern)

	// boolean keywords
	r.RegisterKeyword("allOf", NewAllOf)
	r.RegisterKeyword("anyOf", NewAnyOf)
	r.RegisterKeyword("oneOf", NewOneOf)
	r.RegisterKeyword("not", NewNot)

	// object keywords
	r.RegisterKeyword("properties", NewProperties)
	r.RegisterKeyword("patternProperties", NewPatternProperties)
	r.RegisterKeyword("additionalProperties", NewAdditionalProperties)
	r.RegisterKeyword("required", NewRequired)
	r.RegisterKeyword("propertyNames", NewPropertyNames)
	r.RegisterKeyword("maxProperties", NewMaxProperties)
	r.RegisterKeyword("minProperties", NewMinProperties)
	r.RegisterKeyword("dependentSchemas", NewDependentSchemas)
	r.RegisterKeyword("dependentRequired", NewDependentRequired)
	r.RegisterKeyword("unevaluatedProperties", NewUnevaluatedProperties)

	r.SetKeywordOrder("properties", 2)
	r.SetKeywordOrder("additionalProperties", 3)
	r.SetKeywordOrder("unevaluatedProperties", 4)

	// array keywords
	r.RegisterKeyword("items", NewItems)
	r.RegisterKeyword("additionalItems", NewAdditionalItems)
	r.RegisterKeyword("maxItems", NewMaxItems)
	r.RegisterKeyword("minItems", NewMinItems)
	r.RegisterKeyword("uniqueItems", NewUniqueItems)
	r.RegisterKeyword("contains", NewContains)
	r.RegisterKeyword("maxContains", NewMaxContains)
	r.RegisterKeyword("minContains", NewMinContains)
	r.RegisterKeyword("unevaluatedItems", NewUnevaluatedItems)

	r.SetKeywordOrder("maxContains", 2)
	r.SetKeywordOrder("minContains", 2)
	r.SetKeywordOrder("additionalItems", 3)
	r.SetKeywordOrder("unevaluatedItems", 4)

	// conditional keywords
	r.RegisterKeyword("if", NewIf)
	r.RegisterKeyword("then", NewThen)
	r.RegisterKeyword("else", NewElse)

	r.SetKeywordOrder("then", 2)
	r.SetKeywordOrder("else", 2)

	//optional formats
	r.RegisterKeyword("format", NewFormat)
}
