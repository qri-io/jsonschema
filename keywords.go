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

// DataType gives the primitive json type of a standard json-decoded value, plus the special case
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
type Type struct {
	strVal bool // set to true if Type decoded from a string, false if an array
	vals   []string
}

// NewType creates a new Type Validator
func NewType() Validator {
	return &Type{}
}

// Validate checks to see if input data satisfies the type constraint
func (t Type) Validate(data interface{}) (errs []ValError) {
	jt := DataType(data)
	for _, typestr := range t.vals {
		if jt == typestr || jt == "integer" && typestr == "number" {
			return nil
		}
	}
	if len(t.vals) == 1 {
		errs = append(errs, ValError{
			Message: fmt.Sprintf(`expected "%v" to be of type %s`, data, t.vals[0]),
		})
		return
	}

	str := ""
	for _, ts := range t.vals {
		str += ts + ","
	}
	errs = append(errs, ValError{
		Message: fmt.Sprintf(`expected "%v" to be one of type: %s`, data, str[:len(str)-1]),
	})
	return
}

// JSONProp implements JSON property name indexing for Type
func (t Type) JSONProp(name string) interface{} {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return nil
	}
	if idx > len(t.vals) || idx < 0 {
		return nil
	}
	return t.vals[idx]
}

// UnmarshalJSON implements the json.Unmarshaler interface for Type
func (t *Type) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*t = Type{strVal: true, vals: []string{single}}
	} else {
		var set []string
		if err := json.Unmarshal(data, &set); err == nil {
			*t = Type{vals: set}
		} else {
			return err
		}
	}

	for _, pr := range t.vals {
		if !primitiveTypes[pr] {
			return fmt.Errorf(`"%s" is not a valid type`, pr)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Type
func (t Type) MarshalJSON() ([]byte, error) {
	if t.strVal {
		return json.Marshal(t.vals[0])
	}
	return json.Marshal(t.vals)
}

// Enum validates successfully against this keyword if its value is equal to one of the
// elements in this keyword's array value.
// Elements in the array SHOULD be unique.
// Elements in the array might be of any value, including null.
type Enum []Const

// NewEnum creates a new Enum Validator
func NewEnum() Validator {
	return &Enum{}
}

// String implements the stringer interface for Enum
func (e Enum) String() string {
	str := "["
	for _, c := range e {
		str += c.String() + ", "
	}
	return str[:len(str)-2] + "]"
}

// Validate implements the Validator interface for Enum
func (e Enum) Validate(data interface{}) []ValError {
	for _, v := range e {
		if err := v.Validate(data); err == nil {
			return nil
		}
	}
	return []ValError{
		{Message: fmt.Sprintf("expected %s to be one of %s", data, e.String())},
	}
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
type Const json.RawMessage

// NewConst creates a new Const Validator
func NewConst() Validator {
	return &Const{}
}

// Validate implements the validate interface for Const
func (c Const) Validate(data interface{}) []ValError {
	var con interface{}
	if err := json.Unmarshal(c, &con); err != nil {
		return []ValError{
			{Message: err.Error()},
		}
	}

	if !reflect.DeepEqual(con, data) {
		return []ValError{
			{Message: fmt.Sprintf(`%s must equal %s`, string(c), data)},
		}
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

// MarshalJSON implements json.Marshaler for Const
func (c Const) MarshalJSON() ([]byte, error) {
	return json.Marshal(json.RawMessage(c))
}
