package jsonschema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// primitiveTypes is a map of strings to check types against
var primitiveTypes = map[string]bool{
	"null":    true,
	"boolean": true,
	"object":  true,
	"array":   true,
	"number":  true,
	"string":  true,
	"integer": true,
}

// DataType gives the primitive json type of a value, plus the special case
// "integer" for when numbers are whole
func DataType(data interface{}) string {
	switch v := data.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case float64:
		if float64(int(v)) == v {
			return "integer"
		}
		return "number"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}

// Type specifies one of the six json primitive types.
// The value of this keyword MUST be either a string or an array.
// If it is an array, elements of the array MUST be strings and MUST be unique.
// String values MUST be one of the six primitive types ("null", "boolean", "object", "array", "number", or "string"), or
// "integer" which matches any number with a zero fractional part.
// An instance validates if and only if the instance is in any of the sets listed for this keyword.
type Type []string

// Validate checks to see if input data satisfies the type constraint
func (t Type) Validate(data interface{}) error {
	jt := DataType(data)
	for _, typestr := range t {
		if jt == typestr || jt == "integer" && typestr == "number" {
			return nil
		}
	}
	if len(t) == 1 {
		return fmt.Errorf(`expected "%v" to be of type %s`, data, t[0])
	}

	str := ""
	for _, ts := range t {
		str += ts + ","
	}
	return fmt.Errorf(`expected "%v" to be one of type: %s`, data, str[:len(str)-1])
}

// JSONProp implements JSON property name indexing for Type
func (t Type) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(t) || idx < 0 {
		return nil
	}
	return t[idx]
}

// UnmarshalJSON implements the json.Unmarshaler interface for Type
func (t *Type) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*t = Type{single}
	} else {
		var set []string
		if err := json.Unmarshal(data, &set); err == nil {
			*t = Type(set)
		} else {
			return err
		}
	}

	for _, pr := range *t {
		if !primitiveTypes[pr] {
			return fmt.Errorf(`"%s" is not a valid type`, pr)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Type
func (t Type) MarshalJSON() ([]byte, error) {
	if len(t) == 1 {
		return json.Marshal(t[0])
	} else if len(t) > 1 {
		return json.Marshal([]string(t))
	}
	return []byte(`""`), nil
}

// Enum validates successfully against this keyword if its value is equal to one of the
// elements in this keyword's array value.
// Elements in the array SHOULD be unique.
// Elements in the array might be of any value, including null.
type Enum []Const

// String implements the stringer interface for Enum
func (e Enum) String() string {
	str := "["
	for _, c := range e {
		str += c.String() + ", "
	}
	return str[:len(str)-2] + "]"
}

// Validate implements the Validator interface for Enum
func (e Enum) Validate(data interface{}) error {
	for _, v := range e {
		if err := v.Validate(data); err == nil {
			return nil
		}
	}
	return fmt.Errorf("expected %s to be one of %s", data)
}

// JSONProp implements JSON property name indexing for Enum
func (e Enum) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(e) || idx < 0 {
		return nil
	}
	return e[idx]
}

// JSONChildren implements the JSONContainer interface for Enum
func (e Enum) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, bs := range e {
		res[strconv.Itoa(i)] = bs
	}
	return
}

// Const MAY be of any type, including null.
// An instance validates successfully against this keyword if its
// value is equal to the value of the keyword.
type Const []byte

// Validate implements the validate interface for Const
func (c Const) Validate(data interface{}) error {
	var con interface{}
	if err := json.Unmarshal(c, &con); err != nil {
		return err
	}

	if !reflect.DeepEqual(con, data) {
		return fmt.Errorf(`%s must equal %s`, string(c), data)
	}
	return nil
}

// JSONProp implements JSON property name indexing for Const
func (c Const) JSONProp(name string) interface{} {
	return nil
}

// String implements the Stringer interface for Const
func (c Const) String() string {
	return string(c)
}

// UnmarshalJSON implements the json.Unmarshaler interface for Const
func (c *Const) UnmarshalJSON(data []byte) error {
	*c = data
	return nil
}
