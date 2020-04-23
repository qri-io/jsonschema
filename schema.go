package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"sort"

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
	DocPath       string
	HasRegistered bool
	isValid       bool

	ID     string `json:"$id,omitempty"`
	Anchor string

	extraDefinitions map[string]json.RawMessage
	Keywords         map[string]Keyword
	orderedKeywords  []string
}

func NewSchema() Keyword {
	return &Schema{}
}

func (s *Schema) Path() string {
	if s.DocPath != "" {
		return s.DocPath
	}
	if s.ID != "" {
		s.DocPath = s.ID
	}
	return s.DocPath
}

func (s *Schema) Register(uri string, registry *SchemaRegistry) {
	if s.HasRegistered {
		return
	}
	s.HasRegistered = true
	registry.RegisterLocal(s)

	address := s.ID
	if uri != "" && address != "" {
		address, _ = SafeResolveUrl(uri, address)
	}

	if s.DocPath == "" && address != "" && address[0] != '#' {
		docUri := ""
		if u, err := url.Parse(address); err != nil {
			docUri, _ = SafeResolveUrl("https://qri.io", address)
		} else {
			docUri = u.String()
		}
		s.DocPath = docUri
		GetSchemaRegistry().Register(s)
		uri = docUri
	}

	for _, keyword := range s.Keywords {
		keyword.Register(uri, registry)
	}
}

func (s *Schema) Resolve(pointer jptr.Pointer, uri string) *Schema {
	if pointer.IsEmpty() {
		if s.DocPath != "" {
			s.DocPath, _ = SafeResolveUrl(uri, s.DocPath)
		} else {
			s.DocPath = uri
		}
		return s
	}

	if _, err := url.Parse(s.ID); err == nil {
		if filepath.IsAbs(s.ID) {
			uri = s.ID
		} else {
			uri, _ = SafeResolveUrl(uri, s.ID)
		}
	}

	// TODO: grok and finish this

	return nil
}

func (s Schema) JSONProp(name string) interface{} {
	if keyword, ok := s.Keywords[name]; ok {
		return keyword
	}
	return s.extraDefinitions[name]
}

func (s Schema) JSONChildren() map[string]JSONPather {
	ch := map[string]JSONPather{}

	if s.Keywords != nil {
		for key, val := range s.Keywords {
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
		ID:       _s.ID,
		Keywords: map[string]Keyword{},
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
			fmt.Printf("WARN: '%s' is not supported and will be ignored\n", prop)
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
		sch.Keywords[prop] = keyword
	}

	keyOrders := make([]_keyOrder, len(sch.Keywords))
	i := 0
	for k := range sch.Keywords {
		keyOrders[i] = _keyOrder{
			Key:   k,
			Order: GetKeywordOrder(k),
		}
		i++
	}
	sort.SliceStable(keyOrders, func(i, j int) bool {
		return keyOrders[i].Order < keyOrders[j].Order
	})
	orderedKeys := make([]string, len(sch.Keywords))
	i = 0
	for _, keyOrder := range keyOrders {
		orderedKeys[i] = keyOrder.Key
		i++
	}
	sch.orderedKeywords = orderedKeys

	*s = Schema(*sch)
	return nil
}

type _keyOrder struct {
	Key   string
	Order int
}

func (s *Schema) Validate(propPath string, data interface{}, errs *[]KeyError) {
	schCtx := NewSchemaContext(s, data, &jptr.Pointer{}, &jptr.Pointer{}, &jptr.Pointer{})
	s.ValidateFromContext(schCtx, errs)
}

func (s *Schema) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	if s.schemaType == schemaTypeTrue {
		return
	}
	if s.schemaType == schemaTypeFalse {
		AddErrorCtx(errs, schCtx, fmt.Sprintf("schema is always false"))
		return
	}
	// IsValid := false
	schCtx.Local = s

	// TODO: handle non draft2019-09 ref resolution
	// ref := s.Ref
	// if ref != "" {
	// 	if schCtx.BaseURI == "" {
	// 		schCtx.BaseURI = s.DocPath
	// 	} else if s.DocPath != "" {
	// 		if filepath.IsAbs(s.DocPath) {
	// 			schCtx.BaseURI = s.DocPath
	// 		} else {
	// 			schCtx.BaseURI, _ = SafeResolveUrl(schCtx.BaseURI, s.DocPath)
	// 		}
	// 	}
	// }

	// if schCtx.BaseURI != "" && strings.HasSuffix(schCtx.BaseURI, "#") {
	// 	schCtx.BaseURI = strings.TrimRight(schCtx.BaseURI, "#")
	// }

	// TODO: handle non draft2019-09 ref resolution
	// if ref != "" {}

	// if s.Keywords != nil {
	// for _, keyword := range s.orderedKeywords() {
	// 	keyword.ValidateFromContext(schCtx, errs)
	// }
	// }
	s.validateSchemaKeywords(schCtx, errs)
}

func (s *Schema) validateSchemaKeywords(schCtx *SchemaContext, errs *[]KeyError) {
	if s.Keywords != nil {
		for _, keyword := range s.orderedKeywords {
			s.Keywords[keyword].ValidateFromContext(schCtx, errs)
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
	if t, ok := s.Keywords["type"].(*Type); ok {
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

		for k, v := range s.Keywords {
			obj[k] = v
		}
		for k, v := range s.extraDefinitions {
			obj[k] = v
		}
		return json.Marshal(obj)
	}
}
