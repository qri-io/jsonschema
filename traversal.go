package jsonschema

import (
	"strconv"

	jptr "github.com/qri-io/jsonpointer"
)

// JSONPather makes validators traversible by JSON-pointers,
// which is required to support references in JSON schemas.
type JSONPather interface {
	// JSONProp take a string references for a given JSON property
	// implementations must return any matching property of that name
	// or nil if no such subproperty exists.
	// Note this also applies to array values, which are expected to interpret
	// valid numbers as an array index
	JSONProp(name string) interface{}
}

// JSONContainer is an interface that enables tree traversal by listing
// the immideate children of an object
type JSONContainer interface {
	// JSONChildren should return all immidiate children of this element
	JSONChildren() map[string]JSONPather
}

func walkJSON(elem JSONPather, fn func(elem JSONPather) error) error {
	if err := fn(elem); err != nil {
		return err
	}

	if con, ok := elem.(JSONContainer); ok {
		for _, ch := range con.JSONChildren() {
			if err := walkJSON(ch, fn); err != nil {
				return err
			}
		}
	}

	return nil
}

// resolveSchemaIndex resolves a pointer into a list of schemas: the head of
// the pointer picks a schema by its index, and the rest of the pointer is
// resolved by that schema. It returns nil when the pointer is empty or its
// head is not an index into schemas.
func resolveSchemaIndex(schemas []*Schema, pointer jptr.Pointer, uri string) *Schema {
	if pointer == nil {
		return nil
	}
	current := pointer.Head()
	if current == nil {
		return nil
	}

	pos, err := strconv.Atoi(*current)
	if err != nil {
		return nil
	}

	if pos < 0 || pos >= len(schemas) {
		return nil
	}

	return schemas[pos].Resolve(pointer.Tail(), uri)
}

// resolveSchemaKey resolves a pointer into a map of schemas: the head of the
// pointer picks a schema by its key, and the rest of the pointer is resolved
// by that schema. It returns nil when the pointer is empty or its head is not
// a key of schemas.
func resolveSchemaKey(schemas map[string]*Schema, pointer jptr.Pointer, uri string) *Schema {
	if pointer == nil {
		return nil
	}
	current := pointer.Head()
	if current == nil {
		return nil
	}

	if schema, ok := schemas[*current]; ok {
		return schema.Resolve(pointer.Tail(), uri)
	}

	return nil
}
