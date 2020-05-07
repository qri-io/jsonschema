package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	jptr "github.com/qri-io/jsonpointer"
)

//
// Const
//

type Const json.RawMessage

func NewConst() Keyword {
	return &Const{}
}

func (c Const) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (c Const) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Const] Validating")
	var con interface{}
	if err := json.Unmarshal(c, &con); err != nil {
		AddErrorCtx(errs, schCtx, err.Error())
		return
	}

	if !reflect.DeepEqual(con, schCtx.Instance) {
		AddErrorCtx(errs, schCtx, fmt.Sprintf(`must equal %s`, InvalidValueString(con)))
	}
}

func (c *Const) Register(uri string, registry *SchemaRegistry) {}

func (c *Const) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (c Const) JSONProp(name string) interface{} {
	return nil
}

func (c Const) String() string {
	return string(c)
}

func (c *Const) UnmarshalJSON(data []byte) error {
	*c = data
	return nil
}

func (c Const) MarshalJSON() ([]byte, error) {
	return json.Marshal(json.RawMessage(c))
}

//
// Enum
//

type Enum []Const

func NewEnum() Keyword {
	return &Enum{}
}

func (e Enum) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (e Enum) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Enum] Validating")
	for _, v := range e {
		test := &[]KeyError{}
		v.ValidateFromContext(schCtx, test)
		if len(*test) == 0 {
			return
		}
	}

	AddErrorCtx(errs, schCtx, fmt.Sprintf("should be one of %s", e.String()))
}

func (e *Enum) Register(uri string, registry *SchemaRegistry) {}

func (e *Enum) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

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

func (e Enum) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for i, bs := range e {
		res[strconv.Itoa(i)] = bs
	}
	return
}

func (e Enum) String() string {
	str := "["
	for _, c := range e {
		str += c.String() + ", "
	}
	return str[:len(str)-2] + "]"
}

//
// Type
//

var primitiveTypes = map[string]bool{
	"null":    true,
	"boolean": true,
	"object":  true,
	"array":   true,
	"number":  true,
	"string":  true,
	"integer": true,
}

func DataType(data interface{}) string {
	if data == nil {
		return "null"
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Float64:
		number := reflect.ValueOf(data).Float()
		if float64(int(number)) == number {
			return "integer"
		}
		return "number"
	case reflect.String:
		return "string"
	case reflect.Array, reflect.Slice:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	default:
		return "unknown"
	}
}

func DataTypeWithHint(data interface{}, hint string) string {
	dt := DataType(data)
	if dt == "string" {
		if hint == "boolean" {
			_, err := strconv.ParseBool(data.(string))
			if err == nil {
				return "boolean"
			}
		}
	}
	// deals with traling 0 floats
	if dt == "integer" && hint == "number" {
		return "number"
	}
	return dt
}

type Type struct {
	BaseKeyword
	strVal bool
	vals   []string
}

func NewType() Keyword {
	return &Type{}
}

func (t *Type) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (t *Type) Register(uri string, registry *SchemaRegistry) {}

func (t *Type) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (t Type) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Type] Validating")
	jt := DataType(schCtx.Instance)
	for _, typestr := range t.vals {
		if jt == typestr || jt == "integer" && typestr == "number" {
			return
		}
		if jt == "string" && (typestr == "boolean" || typestr == "number" || typestr == "integer") {
			if DataTypeWithHint(schCtx.Instance, typestr) == typestr {
				return
			}
		}
	}
	if len(t.vals) == 1 {
		AddErrorCtx(errs, schCtx, fmt.Sprintf(`type should be %s`, t.vals[0]))
		return
	}

	str := ""
	for _, ts := range t.vals {
		str += ts + ","
	}

	AddErrorCtx(errs, schCtx, fmt.Sprintf(`type should be one of: %s`, str[:len(str)-1]))
}

func (t Type) String() string {
	if len(t.vals) == 0 {
		return "unknown"
	}
	return strings.Join(t.vals, ",")
}

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

func (t Type) MarshalJSON() ([]byte, error) {
	if t.strVal {
		return json.Marshal(t.vals[0])
	}
	return json.Marshal(t.vals)
}
