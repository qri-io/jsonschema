package main

import (
	"strings"
)

var (
	sr *SchemaRegistry
)

type SchemaRegistry struct {
	schemaLookup map[string]*Schema
	contextLookup map[string]*Schema
}

func GetSchemaRegistry() *SchemaRegistry {
	if sr == nil {
		sr = &SchemaRegistry{
			schemaLookup: map[string]*Schema{},
			contextLookup: map[string]*Schema{},
		}
	}
	return sr
}

func (sr *SchemaRegistry) Get(uri string) *Schema {
	uri = strings.TrimRight(uri, "#")
	schema := sr.schemaLookup[uri]
	if schema == nil {
		err := FetchSchema(uri, schema)
		if err != nil {
			// TODO: Validate Schema
			schema.DocPath = uri
			sr.schemaLookup[uri] = schema
		} else {
			return nil
		}
	}
	return schema
}

func (sr *SchemaRegistry) MustGet(uri string) *Schema {
	uri = strings.TrimRight(uri, "#")
	return sr.schemaLookup[uri]
}

func (sr *SchemaRegistry) GetLocal(uri string) *Schema {
	uri = strings.TrimRight(uri, "#")
	return sr.contextLookup[uri]
}

func (sr *SchemaRegistry) Register(sch *Schema) {
	if sch.DocPath == "" {
		return
	}
	sr.schemaLookup[sch.DocPath] = sch
}

func (sr *SchemaRegistry) RegisterLocal(sch *Schema) {
	if sch.DocPath == "" {
		return
	}
	if sch.ID != "" && IsLocalSchemaId(sch.ID) {
		sr.contextLookup[sch.ID] = sch
	}

	if sch.Anchor != "" {
		anchorUri := sch.DocPath + "#" + sch.Anchor
		sr.contextLookup[anchorUri] = sch
	}
}
