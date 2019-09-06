package jsonschema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type T int

const (
	Unknown T = iota
	Null
	Bool
	Object
	Array
	Number
	Integer
	String
)

// MarshalJSON implements the json.Marshaler interface for T
func (t T) MarshalJSON() ([]byte, error) {
	return json.Marshal(mapTStr[t])
}

var mapTStr = [8]string{
	Unknown: "unknown",
	Null:    "null",
	Bool:    "boolean",
	Object:  "object",
	Array:   "array",
	Number:  "number",
	Integer: "integer",
	String:  "string",
}

var mapStrT = map[string]T{
	"unknown": Unknown,
	"null":    Null,
	"boolean": Bool,
	"object":  Object,
	"array":   Array,
	"number":  Number,
	"integer": Integer,
	"string":  String,
}

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
		return mapTStr[ts[0]]
	}

	// Calculate string length.
	n := len(ts) - 1
	for i := 0; i < len(ts); i++ {
		n += len(mapTStr[ts[i]])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(mapTStr[ts[0]])
	for _, t := range ts[1:] {
		b.WriteString(",")
		b.WriteString(mapTStr[t])
	}
	return b.String()
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
}

func (v ObjectVal) Type() T {
	return Object
}

func (v ObjectVal) Raw() interface{} {
	return v.raw
}

func (v ObjectVal) Len() int {
	switch reflect.TypeOf(v.raw).Kind() {
	case reflect.Map:
		return reflect.ValueOf(v.raw).Len()
	case reflect.Struct:
		return reflect.ValueOf(v.raw).NumField()
	default:
		panic(fmt.Sprintf("invalid value: %v", v.raw))
	}
}

func (v ObjectVal) Field(name string) (interface{}, bool) {
	var val reflect.Value

	switch reflect.TypeOf(v.raw).Kind() {
	case reflect.Map:
		val = reflect.ValueOf(v.raw).MapIndex(reflect.ValueOf(name))
	case reflect.Struct:
		val = reflect.ValueOf(v.raw).FieldByName(name)
	default:
		panic(fmt.Sprintf("invalid value: %v", v.raw))
	}

	if !val.IsValid() {
		return nil, false
	}
	return val.Interface(), true
}

type Pair struct {
	Key   string
	Value interface{}
}

func (v ObjectVal) Iterator(cancelC chan struct{}) <-chan Pair {
	ch := make(chan Pair, 50)

	go func() {
		defer close(ch)

		switch reflect.TypeOf(v.raw).Kind() {
		case reflect.Map:
			for i := reflect.ValueOf(v.raw).MapRange(); i.Next(); {
				select {
				case ch <- Pair{
					Key:   i.Key().String(),
					Value: i.Value().Interface(),
				}:
				case <-cancelC:
					return
				}
			}
		case reflect.Struct:
			val := reflect.ValueOf(v.raw)
			valType := val.Type()
			for i := 0; i < val.NumField(); i++ {
				name := valType.Field(i).Name
				if unicode.IsLower(rune(name[0])) {
					continue
				}
				select {
				case ch <- Pair{
					Key:   valType.Field(i).Name,
					Value: val.Field(i).Interface(),
				}:
				case <-cancelC:
					return
				}
			}
		}
	}()

	return ch
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
var _ Val = ObjectVal{}
var _ Val = ArrayVal{}

func rawToVal(raw interface{}) Val {
	if raw == nil {
		return NullVal{}
	}

	val := reflect.ValueOf(raw)

	switch dataToT(raw) {
	case Bool:
		return BoolVal(val.Bool())
	case Number:
		return NumberVal(val.Float())
	case Integer:
		return IntVal(int64(val.Float()))
	case String:
		return StringVal(val.String())
	case Object:
		return ObjectVal{raw: raw}
	case Array:
		arr := make([]interface{}, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			arr = append(arr, val.Index(i).Interface())
		}
		return ArrayVal(arr)
	default:
		panic(fmt.Sprintf("unable to handle raw data of type '%v' (%v)", dataToT(raw), raw))
	}
}
