package main

import (
	jptr "github.com/qri-io/jsonpointer"
)

type SchemaContext struct {
	Local *Schema
	Root *Schema
	RecursiveAnchor *Schema
	Instance interface{}
	LastEvaluatedIndex int
	LocalLastEvaluatedIndex int
	BaseURI string
	InstanceLocation *jptr.Pointer
	RelativeLocation *jptr.Pointer
	BaseRelativeLocation *jptr.Pointer

	LocalRegistry *SchemaRegistry

	EvaluatedPropertyNames map[string]bool
	LocalEvaluatedPropertyNames map[string]bool
	Misc map[string]interface{}
}

func NewSchemaContext(rs *Schema, inst interface{}, brl *jptr.Pointer, rl *jptr.Pointer, il *jptr.Pointer) *SchemaContext {
	return &SchemaContext{
		Root: rs,
		Instance: inst,
		BaseRelativeLocation: brl,
		RelativeLocation: rl,
		InstanceLocation: il,
		LocalRegistry: &SchemaRegistry{},
		LastEvaluatedIndex: -1,
		LocalLastEvaluatedIndex: -1,
		EvaluatedPropertyNames: map[string]bool{},
		LocalEvaluatedPropertyNames: map[string]bool{},
		Misc: map[string]interface{}{},
	}
}

func NewSchemaContextFromSource(source SchemaContext) *SchemaContext {
	sch := &SchemaContext{
		Local: source.Local,
		Root: source.Root,
		RecursiveAnchor: source.RecursiveAnchor,
		Instance: source.Instance,
		LastEvaluatedIndex: source.LastEvaluatedIndex,
		LocalLastEvaluatedIndex: source.LocalLastEvaluatedIndex,
		BaseURI: source.BaseURI,
		// these should probably be separate copies
		InstanceLocation: source.InstanceLocation,
		RelativeLocation: source.RelativeLocation,
		BaseRelativeLocation: source.RelativeLocation,
		LocalRegistry: source.LocalRegistry,

		Misc: map[string]interface{}{},
	}
	hasAdditionalPropertiesKeyword := false
	hasAdditionalItemsKeyword := false
	if _, ok := sch.Local.Keywords["additionalProperties"]; ok {
		hasAdditionalPropertiesKeyword = true
	}
	if _, ok := sch.Local.Keywords["additionalItems"]; ok {
		hasAdditionalItemsKeyword = true
	}
	if hasAdditionalPropertiesKeyword || hasAdditionalItemsKeyword {
		sch.EvaluatedPropertyNames = copySet(source.EvaluatedPropertyNames)
		sch.LocalEvaluatedPropertyNames = copySet(source.LocalEvaluatedPropertyNames)
	} else {
		sch.EvaluatedPropertyNames = map[string]bool{}
		sch.LocalEvaluatedPropertyNames = map[string]bool{}
	}
		
	return sch
}

func NewSchemaContextFromSourceClean(source SchemaContext) *SchemaContext {
	sch := &SchemaContext{
		Local: source.Local,
		Root: source.Root,
		RecursiveAnchor: source.RecursiveAnchor,
		Instance: source.Instance,
		LastEvaluatedIndex: source.LastEvaluatedIndex,
		LocalLastEvaluatedIndex: source.LocalLastEvaluatedIndex,
		BaseURI: source.BaseURI,
		// these should probably be separate copies
		InstanceLocation: source.InstanceLocation,
		RelativeLocation: source.RelativeLocation,
		BaseRelativeLocation: source.RelativeLocation,
		LocalRegistry: source.LocalRegistry,
		EvaluatedPropertyNames: map[string]bool{},
		LocalEvaluatedPropertyNames: map[string]bool{},

		Misc: map[string]interface{}{},
	}

	return sch
}

func (sc *SchemaContext) ClearContext() {
	if len(sc.EvaluatedPropertyNames) > 0 {
		sc.EvaluatedPropertyNames = map[string]bool{}
	}
	if len(sc.LocalEvaluatedPropertyNames) > 0 {
		sc.LocalEvaluatedPropertyNames = map[string]bool{}
	}
}

func copySet(input map[string]bool) map[string]bool {
	copy := make(map[string]bool, len(input))
	for k,v := range input {
		copy[k] = v
	}
	return copy
}

func JoinSets(consumer *map[string]bool, supplier map[string]bool) {
	for k,v := range supplier {
		(*consumer)[k] = v
	}
}
