package jsonschema

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

var sr *SchemaRegistry

func init() {
	// sets the global sr variable the first time the module is used;
	// this fixes a concurrency issue.
	sr = &SchemaRegistry{
		mtx:           sync.RWMutex{},
		schemaLookup:  map[string]*Schema{},
		contextLookup: map[string]*Schema{},
	}
}

// SchemaRegistry maintains a lookup table between schema string references
// and actual schemas
type SchemaRegistry struct {
	mtx           sync.RWMutex
	schemaLookup  map[string]*Schema
	contextLookup map[string]*Schema
}

// GetSchemaRegistry provides an accessor to a globally available schema registry
func GetSchemaRegistry() *SchemaRegistry {
	return sr
}

// Get fetches a schema from the top level context registry or fetches it from a remote
func (sr *SchemaRegistry) Get(ctx context.Context, uri string) *Schema {
	uri = strings.TrimRight(uri, "#")

	sr.mtx.RLock()
	schema := sr.schemaLookup[uri]
	sr.mtx.RUnlock()

	if schema == nil {
		fetchedSchema := &Schema{}
		err := FetchSchema(ctx, uri, fetchedSchema)
		if err != nil {
			schemaDebug(fmt.Sprintf("[SchemaRegistry] Fetch error: %s", err.Error()))
			return nil
		}
		if fetchedSchema == nil {
			return nil
		}
		fetchedSchema.docPath = uri
		// TODO(arqu): meta validate schema
		schema = fetchedSchema

		sr.mtx.Lock()
		sr.schemaLookup[uri] = schema
		sr.mtx.Unlock()
	}
	return schema
}

// GetKnown fetches a schema from the top level context registry
func (sr *SchemaRegistry) GetKnown(uri string) *Schema {
	uri = strings.TrimRight(uri, "#")

	sr.mtx.RLock()
	defer sr.mtx.RUnlock()

	return sr.schemaLookup[uri]
}

// GetLocal fetches a schema from the local context registry
func (sr *SchemaRegistry) GetLocal(uri string) *Schema {
	uri = strings.TrimRight(uri, "#")

	sr.mtx.RLock()
	defer sr.mtx.RUnlock()

	return sr.contextLookup[uri]
}

// Register registers a schema to the top level context
func (sr *SchemaRegistry) Register(sch *Schema) {
	if sch.docPath == "" {
		return
	}

	sr.mtx.Lock()
	sr.schemaLookup[sch.docPath] = sch
	sr.mtx.Unlock()
}

// RegisterLocal registers a schema to a local context
func (sr *SchemaRegistry) RegisterLocal(sch *Schema) {
	if sch.id != "" && IsLocalSchemaID(sch.id) {
		sr.mtx.Lock()
		sr.contextLookup[sch.id] = sch
		sr.mtx.Unlock()
	}

	if sch.HasKeyword("$anchor") {
		sr.mtx.Lock()
		defer sr.mtx.Unlock()

		anchorKeyword := sch.keywords["$anchor"].(*Anchor)
		anchorURI := sch.docPath + "#" + string(*anchorKeyword)
		if sr.contextLookup == nil {
			sr.contextLookup = map[string]*Schema{}
		}
		sr.contextLookup[anchorURI] = sch
	}
}
