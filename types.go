package jsonschema

import (
	"fmt"
	"reflect"
	"strings"
)

type T string

const (
	Unknown T = "unknown"
	Null      = "null"
	Bool      = "boolean"
	Object    = "object"
	Array     = "array"
	Number    = "number"
	Integer   = "integer"
	String    = "string"
)

// dataToT resolves data to a JSON type.
func dataToT(data interface{}) T {
	if data == nil {
		return Null
	}

	knd := reflect.TypeOf(data).Kind()
	_ = knd

	switch reflect.TypeOf(data).Kind() {
	case reflect.Bool:
		return Bool
	case reflect.Float64:
		number := reflect.ValueOf(data).Float()
		if float64(int(number)) == number {
			return Integer
		}
		return Number
	case reflect.String:
		return String
	case reflect.Array, reflect.Slice:
		return Array
	case reflect.Map, reflect.Struct:
		return Object
	default:
		return Unknown
	}
}

// Ts is a collection of JSON types.
type Ts []T

func (ts Ts) String() string {
	switch len(ts) {
	case 0:
		return ""
	case 1:
		return string(ts[0])
	}

	// Calculate string length.
	n := len(ts) - 1
	for i := 0; i < len(ts); i++ {
		n += len(ts[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(string(ts[0]))
	for _, t := range ts[1:] {
		b.WriteString(",")
		b.WriteString(string(t))
	}
	return b.String()
}

var mapT = map[T]bool{
	Null:    true,
	Bool:    true,
	Object:  true,
	Array:   true,
	Number:  true,
	Integer: true,
	String:  true,
}

type Val interface {
	Type() T
	Raw() interface{}
}

type NullVal struct{}

func (v NullVal) Type() T {
	return Null
}

func (v NullVal) Raw() interface{} {
	return nil
}

type BoolVal bool

func (v BoolVal) Type() T {
	return Bool
}

func (v BoolVal) Raw() interface{} {
	return bool(v)
}

type NumberVal float64

func (v NumberVal) Type() T {
	return Number
}

func (v NumberVal) Raw() interface{} {
	return float64(v)
}

func (v NumberVal) Value() float64 {
	return float64(v)
}

type IntVal int64

func (v IntVal) Type() T {
	return Integer
}

func (v IntVal) Raw() interface{} {
	return int64(v)
}

func (v IntVal) Value() int64 {
	return int64(v)
}

type StringVal string

func (v StringVal) Type() T {
	return String
}

func (v StringVal) Raw() interface{} {
	return string(v)
}

type ObjectVal struct {
	raw interface{}
	m   map[string]interface{}
}

func (v *ObjectVal) Type() T {
	return Object
}

func (v *ObjectVal) Raw() interface{} {
	return v.raw
}

func (v *ObjectVal) Map() map[string]interface{} {
	v.ensureM()
	return v.m
}

func (v *ObjectVal) ensureM() {
	if v.raw == nil || v.m != nil {
		return
	}

	switch reflect.TypeOf(v.raw).Kind() {
	case reflect.Map:
		switch m := v.raw.(type) {
		case map[string]interface{}:
			v.m = m
		case map[interface{}]interface{}:
			tmp := make(map[string]interface{}, len(m))
			for k, v := range m {
				switch str := k.(type) {
				case string:
					tmp[str] = v
				default:
					return
				}
			}
			v.m = tmp
		}
	case reflect.Struct:
		val := reflect.ValueOf(v.raw)
		tmp := make(map[string]interface{}, val.NumField())
		for i := 0; i < val.NumField(); i++ {
			f := val.Field(i)
			t := val.Type()
			tmp[t.Name()] = f.Interface()
		}
		v.m = tmp
	}
}

type ArrayVal []interface{}

func (v ArrayVal) Type() T {
	return Array
}

func (v ArrayVal) Raw() interface{} {
	return []interface{}(v)
}

var _ Val = NullVal{}
var _ Val = BoolVal(false)
var _ Val = NumberVal(0)
var _ Val = IntVal(0)
var _ Val = StringVal("")
var _ Val = &ObjectVal{}
var _ Val = ArrayVal{}

func dataToVal(data interface{}) Val {
	if data == nil {
		return NullVal{}
	}

	val := reflect.ValueOf(data)

	switch dataToT(data) {
	case Bool:
		return BoolVal(val.Bool())
	case Number:
		return NumberVal(val.Float())
	case Integer:
		return IntVal(int64(val.Float()))
	case String:
		return StringVal(val.String())
	case Object:
		return &ObjectVal{raw: data}
	case Array:
		arr := make([]interface{}, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			arr = append(arr, val.Index(i).Interface())
		}
		return ArrayVal(arr)
	default:
		panic(fmt.Sprintf("unable to handle data of type '%v' (%v)", dataToT(data), data))
	}
}
