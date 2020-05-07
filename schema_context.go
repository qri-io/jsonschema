package main

import (
	"context"

	jptr "github.com/qri-io/jsonpointer"
)

type SchemaContext struct {
	Local                   *Schema
	Root                    *Schema
	RecursiveAnchor         *Schema
	Instance                interface{}
	LastEvaluatedIndex      int
	LocalLastEvaluatedIndex int
	BaseURI                 string
	InstanceLocation        *jptr.Pointer
	RelativeLocation        *jptr.Pointer
	BaseRelativeLocation    *jptr.Pointer

	LocalRegistry *SchemaRegistry

	EvaluatedPropertyNames      map[string]bool
	LocalEvaluatedPropertyNames map[string]bool
	Misc                        map[string]interface{}

	ApplicationContext *context.Context
}

func NewSchemaContext(rs *Schema, inst interface{}, brl *jptr.Pointer, rl *jptr.Pointer, il *jptr.Pointer, appCtx *context.Context) *SchemaContext {
	return &SchemaContext{
		Root:                        rs,
		Instance:                    inst,
		BaseRelativeLocation:        brl,
		RelativeLocation:            rl,
		InstanceLocation:            il,
		LocalRegistry:               &SchemaRegistry{},
		LastEvaluatedIndex:          -1,
		LocalLastEvaluatedIndex:     -1,
		EvaluatedPropertyNames:      map[string]bool{},
		LocalEvaluatedPropertyNames: map[string]bool{},
		Misc:                        map[string]interface{}{},
		ApplicationContext:          appCtx,
	}
}

func NewSchemaContextFromSource(source SchemaContext) *SchemaContext {
	sch := &SchemaContext{
		Local:                   source.Local,
		Root:                    source.Root,
		RecursiveAnchor:         source.RecursiveAnchor,
		Instance:                source.Instance,
		LastEvaluatedIndex:      source.LastEvaluatedIndex,
		LocalLastEvaluatedIndex: source.LocalLastEvaluatedIndex,
		BaseURI:                 source.BaseURI,
		InstanceLocation:        source.InstanceLocation,
		RelativeLocation:        source.RelativeLocation,
		BaseRelativeLocation:    source.RelativeLocation,
		LocalRegistry:           source.LocalRegistry,
		Misc:                    map[string]interface{}{},
		ApplicationContext:      source.ApplicationContext,
	}
	hasAdditionalPropertiesKeyword := false
	hasAdditionalItemsKeyword := false
	if _, ok := sch.Local.keywords["additionalProperties"]; ok {
		hasAdditionalPropertiesKeyword = true
	}
	if _, ok := sch.Local.keywords["additionalItems"]; ok {
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
		Local:                       source.Local,
		Root:                        source.Root,
		RecursiveAnchor:             source.RecursiveAnchor,
		Instance:                    source.Instance,
		LastEvaluatedIndex:          source.LastEvaluatedIndex,
		LocalLastEvaluatedIndex:     source.LocalLastEvaluatedIndex,
		BaseURI:                     source.BaseURI,
		InstanceLocation:            source.InstanceLocation,
		RelativeLocation:            source.RelativeLocation,
		BaseRelativeLocation:        source.RelativeLocation,
		LocalRegistry:               source.LocalRegistry,
		EvaluatedPropertyNames:      map[string]bool{},
		LocalEvaluatedPropertyNames: map[string]bool{},
		Misc:                        map[string]interface{}{},
		ApplicationContext:          source.ApplicationContext,
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

func (sc *SchemaContext) UpdateEvaluatedPropsAndItems(subCtx *SchemaContext) {
	JoinSets(&sc.EvaluatedPropertyNames, subCtx.EvaluatedPropertyNames)
	JoinSets(&sc.LocalEvaluatedPropertyNames, subCtx.LocalEvaluatedPropertyNames)
	if subCtx.LastEvaluatedIndex > sc.LastEvaluatedIndex {
		sc.LastEvaluatedIndex = subCtx.LastEvaluatedIndex
	}
	if subCtx.LocalLastEvaluatedIndex > sc.LastEvaluatedIndex {
		sc.LastEvaluatedIndex = subCtx.LocalLastEvaluatedIndex
	}
}

func copySet(input map[string]bool) map[string]bool {
	copy := make(map[string]bool, len(input))
	for k, v := range input {
		copy[k] = v
	}
	return copy
}

func JoinSets(consumer *map[string]bool, supplier map[string]bool) {
	for k, v := range supplier {
		(*consumer)[k] = v
	}
}
