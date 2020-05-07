package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	jptr "github.com/qri-io/jsonpointer"
)

func Must(jsonString string) *Schema {
	s := &Schema{}
	if err := s.UnmarshalJSON([]byte(jsonString)); err != nil {
		panic(err)
	}
	return s
}

type schemaType int

const (
	schemaTypeObject schemaType = iota
	schemaTypeFalse
	schemaTypeTrue
)

type Schema struct {
	schemaType    schemaType
	docPath       string
	hasRegistered bool
	// isValid       bool

	id string

	extraDefinitions map[string]json.RawMessage
	keywords         map[string]Keyword
	orderedkeywords  []string
}

func NewSchema() Keyword {
	return &Schema{}
}

func (s *Schema) HasKeyword(key string) bool {
	_, ok := s.keywords[key]
	return ok
}

func (s *Schema) Register(uri string, registry *SchemaRegistry) {
	if s.hasRegistered {
		return
	}
	s.hasRegistered = true
	registry.RegisterLocal(s)

	if !IsRegistryLoaded() {
		LoadDraft2019_09()
	}

	address := s.id
	if uri != "" && address != "" {
		address, _ = SafeResolveUrl(uri, address)
	}

	if s.docPath == "" && address != "" && address[0] != '#' {
		docUri := ""
		if u, err := url.Parse(address); err != nil {
			docUri, _ = SafeResolveUrl("https://qri.io", address)
		} else {
			docUri = u.String()
		}
		s.docPath = docUri
		GetSchemaRegistry().Register(s)
		uri = docUri
	}

	for _, keyword := range s.keywords {
		keyword.Register(uri, registry)
	}
}

func (s *Schema) Resolve(pointer jptr.Pointer, uri string) *Schema {
	if pointer.IsEmpty() {
		if s.docPath != "" {
			s.docPath, _ = SafeResolveUrl(uri, s.docPath)
		} else {
			s.docPath = uri
		}
		return s
	}

	current := pointer.Head()

	if s.id != "" {
		if u, err := url.Parse(s.id); err == nil {
			if u.IsAbs() {
				uri = s.id
			} else {
				uri, _ = SafeResolveUrl(uri, s.id)
			}
		}
	}

	keyword := s.keywords[*current]
	var keywordSchema *Schema
	if keyword != nil {
		keywordSchema = keyword.Resolve(pointer.Tail(), uri)
	}

	if keywordSchema != nil {
		return keywordSchema
	}

	found, err := pointer.Eval(s.extraDefinitions)
	if err != nil {
		return nil
	}
	if found == nil {
		return nil
	}

	if foundSchema, ok := found.(*Schema); ok {
		return foundSchema
	}

	return nil
}

func (s Schema) JSONProp(name string) interface{} {
	if keyword, ok := s.keywords[name]; ok {
		return keyword
	}
	return s.extraDefinitions[name]
}

func (s Schema) JSONChildren() map[string]JSONPather {
	ch := map[string]JSONPather{}

	if s.keywords != nil {
		for key, val := range s.keywords {
			if jp, ok := val.(JSONPather); ok {
				ch[key] = jp
			}
		}
	}

	return ch
}

type _schema struct {
	ID string `json:"$id,omitempty"`
}

func (s *Schema) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		if b {
			// boolean true Always passes validation, as if the empty schema {}
			*s = Schema{schemaType: schemaTypeTrue}
			return nil
		}
		// boolean false Always fails validation, as if the schema { "not":{} }
		*s = Schema{schemaType: schemaTypeFalse}
		return nil
	}

	_s := _schema{}
	if err := json.Unmarshal(data, &_s); err != nil {
		return err
	}

	sch := &Schema{
		id:       _s.ID,
		keywords: map[string]Keyword{},
	}

	valprops := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &valprops); err != nil {
		return err
	}

	for prop, rawmsg := range valprops {
		var keyword Keyword
		if IsKeyword(prop) {
			keyword = GetKeyword(prop)
		} else if IsNotSupportedKeyword(prop) {
			SchemaDebug(fmt.Sprintf("[Schema] WARN: '%s' is not supported and will be ignored\n", prop))
			continue
		} else {
			if sch.extraDefinitions == nil {
				sch.extraDefinitions = map[string]json.RawMessage{}
			}
			sch.extraDefinitions[prop] = rawmsg
			continue
		}
		if _, ok := keyword.(*Void); !ok {
			if err := json.Unmarshal(rawmsg, keyword); err != nil {
				return fmt.Errorf("error unmarshaling %s from json: %s", prop, err.Error())
			}
		}
		sch.keywords[prop] = keyword
	}

	keyOrders := make([]_keyOrder, len(sch.keywords))
	i := 0
	for k := range sch.keywords {
		keyOrders[i] = _keyOrder{
			Key:   k,
			Order: GetKeywordOrder(k),
		}
		i++
	}
	sort.SliceStable(keyOrders, func(i, j int) bool {
		if keyOrders[i].Order == keyOrders[j].Order {
			return GetKeywordInsertOrder(keyOrders[i].Key) < GetKeywordInsertOrder(keyOrders[j].Key)
		}
		return keyOrders[i].Order < keyOrders[j].Order
	})
	orderedKeys := make([]string, len(sch.keywords))
	i = 0
	for _, keyOrder := range keyOrders {
		orderedKeys[i] = keyOrder.Key
		i++
	}
	sch.orderedkeywords = orderedKeys

	*s = Schema(*sch)
	return nil
}

type _keyOrder struct {
	Key   string
	Order int
}

func (s *Schema) Validate(propPath string, data interface{}, errs *[]KeyError) {
	appCtx := context.Background()
	schCtx := NewSchemaContext(s, data, &jptr.Pointer{}, &jptr.Pointer{}, &jptr.Pointer{}, &appCtx)
	s.ValidateFromContext(schCtx, errs)
}

func (s *Schema) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if s == nil {
		AddErrorCtx(errs, schCtx, fmt.Sprintf("schema is nil"))
		return
	}
	if s.schemaType == schemaTypeTrue {
		return
	}
	if s.schemaType == schemaTypeFalse {
		AddErrorCtx(errs, schCtx, fmt.Sprintf("schema is always false"))
		return
	}

	s.Register("", schCtx.LocalRegistry)
	schCtx.LocalRegistry.RegisterLocal(s)

	schCtx.Local = s

	refKeyword := s.keywords["$ref"]

	if refKeyword != nil {
		if schCtx.BaseURI == "" {
			schCtx.BaseURI = s.docPath
		} else if s.docPath != "" {
			if u, err := url.Parse(s.docPath); err == nil {
				if u.IsAbs() {
					schCtx.BaseURI = s.docPath
				} else {
					schCtx.BaseURI, _ = SafeResolveUrl(schCtx.BaseURI, s.docPath)
				}
			}
		}
	}

	if schCtx.BaseURI != "" && strings.HasSuffix(schCtx.BaseURI, "#") {
		schCtx.BaseURI = strings.TrimRight(schCtx.BaseURI, "#")
	}

	// TODO(arqu): only on versions bellow draft2019_09
	// if refKeyword != nil {
	// 	refKeyword.ValidateFromContext(schCtx, errs)
	// 	return
	// }

	s.validateSchemakeywords(schCtx, errs)
}

func (s *Schema) validateSchemakeywords(schCtx *SchemaContext, errs *[]KeyError) {
	if s.keywords != nil {
		for _, keyword := range s.orderedkeywords {
			s.keywords[keyword].ValidateFromContext(schCtx, errs)
		}
	}
}

func (s *Schema) ValidateBytes(data []byte) ([]KeyError, error) {
	var doc interface{}
	errs := []KeyError{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return errs, fmt.Errorf("error parsing JSON bytes: %s", err.Error())
	}
	s.Validate("/", doc, &errs)
	return errs, nil
}

func (s *Schema) TopLevelType() string {
	if t, ok := s.keywords["type"].(*Type); ok {
		return t.String()
	}
	return "unknown"
}

func (s Schema) MarshalJSON() ([]byte, error) {
	switch s.schemaType {
	case schemaTypeFalse:
		return []byte("false"), nil
	case schemaTypeTrue:
		return []byte("true"), nil
	default:
		obj := map[string]interface{}{}

		for k, v := range s.keywords {
			obj[k] = v
		}
		for k, v := range s.extraDefinitions {
			obj[k] = v
		}
		return json.Marshal(obj)
	}
}
