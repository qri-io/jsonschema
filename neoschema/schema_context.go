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
	}
}

func NewSchemaContextFromSource(source SchemaContext) *SchemaContext {
	return &SchemaContext{
		Local: source.Local,
		Root: source.Root,
		RecursiveAnchor: source.RecursiveAnchor,
		Instance: source.Instance,
		LastEvaluatedIndex: source.LastEvaluatedIndex,
		LocalLastEvaluatedIndex: source.LocalLastEvaluatedIndex,
		BaseURI: source.BaseURI,
		InstanceLocation: source.InstanceLocation,
		RelativeLocation: source.RelativeLocation,
		BaseRelativeLocation: source.RelativeLocation,
		LocalRegistry: source.LocalRegistry,
	}
}