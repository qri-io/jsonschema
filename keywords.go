package jsonschema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Type specifies one of the six json primitive types.
// The value of this keyword MUST be either a string or an array.
// If it is an array, elements of the array MUST be strings and MUST be unique.
// String values MUST be one of the six primitive types ("null", "boolean", "object", "array", "number", or "string"), or
// "integer" which matches any number with a zero fractional part.
// An instance validates if and only if the instance is in any of the sets listed for this keyword.
type Type struct {
	BaseValidator
	strVal bool // set to true if Type decoded from a string, false if an array
	vals   Ts
}

// NewType creates a new Type Validator
func NewType() Validator {
	return &Type{}
}

// String returns the type(s) as a string, or unknown if there is no known type.
func (t Type) String() string {
	if len(t.vals) == 0 {
		return "unknown"
	}
	return t.vals.String()
}

// Validate checks to see if input data satisfies the type constraint
func (t Type) Validate(propPath string, data Val, errs *[]ValError) {
	for _, typestr := range t.vals {
		if data.Type() == typestr || data.Type() == Integer && typestr == Number {
			return
		}
	}
	if len(t.vals) == 1 {
		t.AddError(errs, propPath, data, fmt.Sprintf(`type should be %s`, t.vals[0]))
		return
	}

	str := ""
	for _, ts := range t.vals {
		str += string(ts) + ","
	}

	t.AddError(errs, propPath, data, fmt.Sprintf(`type should be one of: %s`, str[:len(str)-1]))
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
		*t = Type{strVal: true, vals: Ts{T(single)}}
	} else {
		var set Ts
		if err := json.Unmarshal(data, &set); err == nil {
			*t = Type{vals: set}
		} else {
			return err
		}
	}

	for _, pr := range t.vals {
		if !mapT[pr] {
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

// Path gives a jsonpointer path to the validator
func (e Enum) Path() string {
	return ""
}

// Validate implements the Validator interface for Enum
func (e Enum) Validate(propPath string, data Val, errs *[]ValError) {
	for _, v := range e {
		test := &[]ValError{}
		v.Validate(propPath, data, test)
		if len(*test) == 0 {
			return
		}
	}

	AddError(errs, propPath, data, fmt.Sprintf("should be one of %s", e.String()))
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

// Path gives a jsonpointer path to the validator
func (c Const) Path() string {
	return ""
}

// Validate implements the validate interface for Const
func (c Const) Validate(propPath string, data Val, errs *[]ValError) {
	var con interface{}
	if err := json.Unmarshal(c, &con); err != nil {
		AddError(errs, propPath, data, err.Error())
		return
	}

	// Normalize the input to match the constant type.
	var val interface{}
	b, err := json.Marshal(data.Raw())
	if err != nil {
		AddError(errs, propPath, data, err.Error())
		return
	}
	if err := json.Unmarshal(b, &val); err != nil {
		AddError(errs, propPath, data, err.Error())
		return
	}

	if !reflect.DeepEqual(con, val) {
		AddError(errs, propPath, data, fmt.Sprintf(`must equal %s`, InvalidValueString(con)))
	}
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
