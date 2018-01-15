package jsonschema

// import (
// 	"fmt"
// )

type JSONPather interface {
	// JSONProp makes validators traversible by JSON-pointers,
	// which is required to support references in JSON schemas.
	// for a given JSON property name the validator must
	// return any matching property of that name
	// or nil if no such subproperty exists.
	// Note this also applies to array values, which are expected to interpret
	// valid numbers as an array index
	JSONProp(name string) interface{}
}

type JSONContainer interface {
	// JSONChildren should return all immidiate children of this element
	JSONChildren() map[string]JSONPather
}

func walkJSON(elem JSONPather, fn func(elem JSONPather) error) error {
	if err := fn(elem); err != nil {
		return err
	}

	if con, ok := elem.(JSONContainer); ok {
		// fmt.Println(con)
		for _, ch := range con.JSONChildren() {
			// fmt.Println("child:", key, ch)
			if err := walkJSON(ch, fn); err != nil {
				return err
			}
		}
	}

	return nil
}
