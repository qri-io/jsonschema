package main

import (
	"context"
	"fmt"
	"strings"
)

var (
	sr *SchemaRegistry
)

type SchemaRegistry struct {
	schemaLookup  map[string]*Schema
	contextLookup map[string]*Schema
}

func GetSchemaRegistry() *SchemaRegistry {
	if sr == nil {
		sr = &SchemaRegistry{
			schemaLookup:  map[string]*Schema{},
			contextLookup: map[string]*Schema{},
		}
	}
	return sr
}

func (sr *SchemaRegistry) Get(uri string, ctx *context.Context) *Schema {
	uri = strings.TrimRight(uri, "#")
	schema := sr.schemaLookup[uri]
	if schema == nil {
		fetchedSchema := &Schema{}
		err := FetchSchema(ctx, uri, fetchedSchema)
		if err != nil {
			SchemaDebug(fmt.Sprintf("[SchemaRegistry] Fetch error: %s", err.Error()))
			return nil
		}
		if fetchedSchema == nil {
			return nil
		}
		fetchedSchema.docPath = uri
		// TODO(arqu): meta validate schema
		schema = fetchedSchema
		sr.schemaLookup[uri] = schema
	}
	return schema
}

func (sr *SchemaRegistry) GetKnown(uri string) *Schema {
	uri = strings.TrimRight(uri, "#")
	return sr.schemaLookup[uri]
}

func (sr *SchemaRegistry) GetLocal(uri string) *Schema {
	uri = strings.TrimRight(uri, "#")
	return sr.contextLookup[uri]
}

func (sr *SchemaRegistry) Register(sch *Schema) {
	if sch.docPath == "" {
		return
	}
	sr.schemaLookup[sch.docPath] = sch
}

func (sr *SchemaRegistry) RegisterLocal(sch *Schema) {
	if sch.id != "" && IsLocalSchemaId(sch.id) {
		sr.contextLookup[sch.id] = sch
	}

	if sch.HasKeyword("$anchor") {
		anchorKeyword := sch.keywords["$anchor"].(*Anchor)
		anchorUri := sch.docPath + "#" + string(*anchorKeyword)
		if sr.contextLookup == nil {
			sr.contextLookup = map[string]*Schema{}
		}
		sr.contextLookup[anchorUri] = sch
	}
}
