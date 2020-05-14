package jsonschema

import (
	"context"

	jptr "github.com/qri-io/jsonpointer"
)

// SchemaContext holds the validation context
// The aim is to have one global validation context
// and use local sub contexts when evaluating parallel branches
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

	EvaluatedPropertyNames      *map[string]bool
	LocalEvaluatedPropertyNames *map[string]bool
	Misc                        map[string]interface{}

	ApplicationContext *context.Context
}

// NewSchemaContext creates a new SchemaContext with the provided location pointers and data instance
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
		EvaluatedPropertyNames:      &map[string]bool{},
		LocalEvaluatedPropertyNames: &map[string]bool{},
		Misc:                        map[string]interface{}{},
		ApplicationContext:          appCtx,
	}
}

// NewSchemaContextFromSource creates a new SchemaContext from an existing SchemaContext
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
			EvaluatedPropertyNames:      source.EvaluatedPropertyNames,
		LocalEvaluatedPropertyNames: source.LocalEvaluatedPropertyNames,
		Misc:                    map[string]interface{}{},
		ApplicationContext:      source.ApplicationContext,
	}

	return sch
}

// NewSchemaContextFromSourceClean creates a new SchemaContext from an existing SchemaContext but only
// copies the core structures
func NewSchemaContextFromSourceClean(source SchemaContext) *SchemaContext {
	// return NewSchemaContextFromSource(source)
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
		EvaluatedPropertyNames:      &map[string]bool{},
		LocalEvaluatedPropertyNames: &map[string]bool{},
		Misc:                        map[string]interface{}{},
		ApplicationContext:          source.ApplicationContext,
	}

	return sch
}

// ClearContext resets a schema to it's core elements
func (sc *SchemaContext) ClearContext() {
	if len(*sc.EvaluatedPropertyNames) > 0 {
		sc.EvaluatedPropertyNames = &map[string]bool{}
	}
	if len(*sc.LocalEvaluatedPropertyNames) > 0 {
		sc.LocalEvaluatedPropertyNames = &map[string]bool{}
	}
	if len(sc.Misc) > 0 {
		sc.Misc = map[string]interface{}{}
	}
}

// UpdateEvaluatedPropsAndItems is a utility function to join evaluated properties and set the
// current evaluation position index
func (sc *SchemaContext) UpdateEvaluatedPropsAndItems(subCtx *SchemaContext) {
	JoinSets(sc.EvaluatedPropertyNames, *subCtx.EvaluatedPropertyNames)
	JoinSets(sc.LocalEvaluatedPropertyNames, *subCtx.LocalEvaluatedPropertyNames)
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

// JoinSets is a utility function to join two existance check maps
func JoinSets(consumer *map[string]bool, supplier map[string]bool) {
	for k, v := range supplier {
		(*consumer)[k] = v
	}
}
