package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	jptr "github.com/qri-io/jsonpointer"
)

//
// $schema
//

type SchemaURI string

func NewSchemaURI() Keyword {
	return new(SchemaURI)
}

func (s *SchemaURI) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (s *SchemaURI) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[SchemaURI] Validating")
}

func (s *SchemaURI) Register(uri string, registry *SchemaRegistry) {}

func (s *SchemaURI) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// $id
//

type Id string

func NewId() Keyword {
	return new(Id)
}

func (i *Id) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (i *Id) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Id] Validating")
	// TODO(arqu): make sure ID is valid URI for draft2019
}

func (i *Id) Register(uri string, registry *SchemaRegistry) {}

func (i *Id) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// Description
//

type Description string

func NewDescription() Keyword {
	return new(Description)
}

func (d *Description) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (d *Description) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Description] Validating")
}

func (d *Description) Register(uri string, registry *SchemaRegistry) {}

func (d *Description) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// Title
//

type Title string

func NewTitle() Keyword {
	return new(Title)
}

func (t *Title) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (t *Title) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Title] Validating")
}

func (t *Title) Register(uri string, registry *SchemaRegistry) {}

func (t *Title) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// Comment
//

type Comment string

func NewComment() Keyword {
	return new(Comment)
}

func (c *Comment) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (c *Comment) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Comment] Validating")
}

func (c *Comment) Register(uri string, registry *SchemaRegistry) {}

func (c *Comment) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// Default
//

type Default struct {
	data interface{}
}

func NewDefault() Keyword {
	return &Default{}
}

func (d *Default) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (d *Default) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Default] Validating")
}

func (d *Default) Register(uri string, registry *SchemaRegistry) {}

func (d *Default) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (d *Default) UnmarshalJSON(data []byte) error {
	var defaultData interface{}
	if err := json.Unmarshal(data, &defaultData); err != nil {
		return err
	}
	*d = Default{
		data: defaultData,
	}
	return nil
}

//
// Examples
//

type Examples []interface{}

func NewExamples() Keyword {
	return new(Examples)
}

func (e *Examples) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (e *Examples) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Examples] Validating")
}

func (e *Examples) Register(uri string, registry *SchemaRegistry) {}

func (e *Examples) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// ReadOnly
//

type ReadOnly bool

func NewReadOnly() Keyword {
	return new(ReadOnly)
}

func (r *ReadOnly) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (r *ReadOnly) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[ReadOnly] Validating")
}

func (r *ReadOnly) Register(uri string, registry *SchemaRegistry) {}

func (r *ReadOnly) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// WriteOnly
//

type WriteOnly bool

func NewWriteOnly() Keyword {
	return new(WriteOnly)
}

func (w *WriteOnly) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (w *WriteOnly) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[WriteOnly] Validating")
}

func (w *WriteOnly) Register(uri string, registry *SchemaRegistry) {}

func (w *WriteOnly) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// $ref
//

type Ref struct {
	reference         string
	resolved          *Schema
	resolvedRoot      *Schema
	resolvedFragment  *jptr.Pointer
	fragmentLocalized bool
}

func NewRef() Keyword {
	return new(Ref)
}

func (r *Ref) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (r *Ref) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Ref] Validating")
	if r.resolved == nil {
		r._resolveRef(schCtx)
		if r.resolved == nil {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("failed to resolve schema for ref %s", r.reference))
		}
	}

	subCtx := NewSchemaContextFromSourceClean(*schCtx)
	if r.resolvedRoot != nil {
		subCtx.BaseURI = r.resolvedRoot.docPath
		subCtx.Root = r.resolvedRoot
	}
	if r.resolvedFragment != nil && !r.resolvedFragment.IsEmpty() {
		subCtx.BaseRelativeLocation = r.resolvedFragment
	}
	relLocation := schCtx.RelativeLocation.RawDescendant("$ref")
	subCtx.RelativeLocation = &relLocation
	subCtx.Instance = schCtx.Instance

	r.resolved.ValidateFromContext(subCtx, errs)

	schCtx.UpdateEvaluatedPropsAndItems(subCtx)
}

func (r *Ref) _resolveRef(schCtx *SchemaContext) {
	if IsLocalSchemaId(r.reference) {
		r.resolved = schCtx.LocalRegistry.GetLocal(r.reference)
		if r.resolved != nil {
			return
		}
	}

	docPath := schCtx.BaseURI
	refParts := strings.Split(r.reference, "#")
	address := ""
	if refParts != nil && len(strings.TrimSpace(refParts[0])) > 0 {
		address = refParts[0]
	} else if docPath != "" {
		docPathParts := strings.Split(docPath, "#")
		address = docPathParts[0]
	}
	if len(refParts) > 1 {
		frag := refParts[1]
		if len(frag) > 0 && frag[0] != '/' {
			frag = "/" + frag
			r.fragmentLocalized = true
		}
		fragPointer, err := jptr.Parse(frag)
		if err != nil {
			r.resolvedFragment = &jptr.Pointer{}
		} else {
			r.resolvedFragment = &fragPointer
		}
	} else {
		r.resolvedFragment = &jptr.Pointer{}
	}

	if address != "" {
		if u, err := url.Parse(address); err == nil {
			if !u.IsAbs() {
				address = schCtx.Local.id + address
				if docPath != "" {
					uriFolder := ""
					if docPath[len(docPath)-1] == '/' {
						uriFolder = docPath
					} else {
						corePath := strings.Split(docPath, "#")[0]
						pathComponents := strings.Split(corePath, "/")
						pathComponents = pathComponents[:len(pathComponents)-1]
						uriFolder = strings.Join(pathComponents, "/") + "/"
					}
					address, _ = SafeResolveUrl(uriFolder, address)
				}
			}
		}
		r.resolvedRoot = GetSchemaRegistry().Get(address, schCtx.ApplicationContext)
	} else {
		r.resolvedRoot = schCtx.Root
	}

	if r.resolvedRoot == nil {
		return
	}

	knownSchema := GetSchemaRegistry().GetKnown(r.reference)
	if knownSchema != nil {
		r.resolved = knownSchema
		return
	}

	localUri := schCtx.BaseURI
	if r.resolvedRoot != nil && r.resolvedRoot.docPath != "" {
		localUri = r.resolvedRoot.docPath
		if r.fragmentLocalized && !r.resolvedFragment.IsEmpty() {
			current := r.resolvedFragment.Head()
			sch := schCtx.LocalRegistry.GetLocal("#" + *current)
			if sch != nil {
				r.resolved = sch
				return
			}
		}
	}
	r._resolveLocalRef(localUri)
}

func (r *Ref) _resolveLocalRef(uri string) {
	if r.resolvedFragment.IsEmpty() {
		r.resolved = r.resolvedRoot
		return
	}

	if r.resolvedRoot != nil {
		r.resolved = r.resolvedRoot.Resolve(*r.resolvedFragment, uri)
	}
}

func (r *Ref) Register(uri string, registry *SchemaRegistry) {}

func (r *Ref) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (r *Ref) UnmarshalJSON(data []byte) error {
	var ref string
	if err := json.Unmarshal(data, &ref); err != nil {
		return err
	}
	normalizedRef, _ := url.QueryUnescape(ref)
	*r = Ref{
		reference: normalizedRef,
	}
	return nil
}

func (r Ref) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.reference)
}

//
// $recursiveRef
//

type RecursiveRef struct {
	reference        string
	resolved         *Schema
	resolvedRoot     *Schema
	resolvedFragment *jptr.Pointer

	validatingLocations map[string]bool
}

func NewRecursiveRef() Keyword {
	return new(RecursiveRef)
}

func (r *RecursiveRef) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (r *RecursiveRef) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[RecursiveRef] Validating")
	if r.IsLocationVisited(schCtx.InstanceLocation.String()) {
		// recursion detected aborting further descent
		return
	}

	if r.resolved == nil {
		r._resolveRef(schCtx, errs)
		if r.resolved == nil {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("failed to resolve schema for ref %s", r.reference))
		}
	}

	subCtx := NewSchemaContextFromSourceClean(*schCtx)
	if r.resolvedRoot != nil {
		subCtx.BaseURI = r.resolvedRoot.docPath
		subCtx.Root = r.resolvedRoot
	}
	if r.resolvedFragment != nil && !r.resolvedFragment.IsEmpty() {
		subCtx.BaseRelativeLocation = r.resolvedFragment
	}
	relLocation := schCtx.RelativeLocation.RawDescendant("$recursiveRef")
	subCtx.RelativeLocation = &relLocation

	if r.validatingLocations == nil {
		r.validatingLocations = map[string]bool{}
	}

	r.validatingLocations[schCtx.InstanceLocation.String()] = true
	r.resolved.ValidateFromContext(subCtx, errs)
	r.validatingLocations[schCtx.InstanceLocation.String()] = false

	schCtx.UpdateEvaluatedPropsAndItems(subCtx)
}

func (r *RecursiveRef) IsLocationVisited(location string) bool {
	if r.validatingLocations == nil {
		return false
	}
	v, ok := r.validatingLocations[location]
	if !ok {
		return false
	} else {
		return v
	}
}

func (r *RecursiveRef) _resolveRef(schCtx *SchemaContext, errs *[]KeyError) {
	if schCtx.RecursiveAnchor != nil {
		if schCtx.BaseURI == "" {
			AddErrorCtx(errs, schCtx, fmt.Sprintf("base uri not set"))
			return
		}
		baseSchema := GetSchemaRegistry().Get(schCtx.BaseURI, schCtx.ApplicationContext)
		if baseSchema != nil && baseSchema.HasKeyword("$recursiveAnchor") {
			r.resolvedRoot = schCtx.RecursiveAnchor
		}
	}

	if IsLocalSchemaId(r.reference) {
		r.resolved = schCtx.LocalRegistry.GetLocal(r.reference)
		if r.resolved != nil {
			return
		}
	}

	docPath := schCtx.BaseURI
	if r.resolvedRoot != nil && r.resolvedRoot.docPath != "" {
		docPath = r.resolvedRoot.docPath
	}

	refParts := strings.Split(r.reference, "#")
	address := ""
	if refParts != nil && len(strings.TrimSpace(refParts[0])) > 0 {
		address = refParts[0]
	} else {
		address = docPath
	}

	if len(refParts) > 1 {

		fragPointer, err := jptr.Parse(refParts[1])
		if err != nil {
			r.resolvedFragment = &jptr.Pointer{}
		} else {
			r.resolvedFragment = &fragPointer
		}
	} else {
		r.resolvedFragment = &jptr.Pointer{}
	}

	if r.resolvedRoot == nil {
		if address != "" {
			if u, err := url.Parse(address); err == nil {
				if !u.IsAbs() {
					address = schCtx.Local.id + address
					if docPath != "" {
						uriFolder := ""
						if docPath[len(docPath)-1] == '/' {
							uriFolder = docPath
						} else {
							corePath := strings.Split(docPath, "#")[0]
							pathComponents := strings.Split(corePath, "/")
							pathComponents = pathComponents[:len(pathComponents)-1]
							uriFolder = strings.Join(pathComponents, "/")
						}
						address, _ = SafeResolveUrl(uriFolder, address)
					}
				}
			}
			r.resolvedRoot = GetSchemaRegistry().Get(address, schCtx.ApplicationContext)
		} else {
			r.resolvedRoot = schCtx.Root
		}
	}

	if r.resolvedRoot == nil {
		return
	}

	knownSchema := GetSchemaRegistry().GetKnown(r.reference)
	if knownSchema != nil {
		r.resolved = knownSchema
		return
	}

	localUri := schCtx.BaseURI
	if r.resolvedRoot != nil && r.resolvedRoot.docPath != "" {
		localUri = r.resolvedRoot.docPath
	}
	r._resolveLocalRef(localUri)
}

func (r *RecursiveRef) _resolveLocalRef(uri string) {
	if r.resolvedFragment.IsEmpty() {
		r.resolved = r.resolvedRoot
		return
	}

	if r.resolvedRoot != nil {
		r.resolved = r.resolvedRoot.Resolve(*r.resolvedFragment, uri)
	}
}

func (r *RecursiveRef) Register(uri string, registry *SchemaRegistry) {}

func (r *RecursiveRef) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (r *RecursiveRef) UnmarshalJSON(data []byte) error {
	var ref string
	if err := json.Unmarshal(data, &ref); err != nil {
		return err
	}
	*r = RecursiveRef{
		reference: ref,
	}
	return nil
}

func (r RecursiveRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.reference)
}

//
// $anchor
//

type Anchor string

func NewAnchor() Keyword {
	return new(Anchor)
}

func (a *Anchor) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (a *Anchor) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Anchor] Validating")
}

func (a *Anchor) Register(uri string, registry *SchemaRegistry) {}

func (a *Anchor) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

//
// $recursiveAnchor
//

type RecursiveAnchor Schema

func NewRecursiveAnchor() Keyword {
	return &RecursiveAnchor{}
}

func (r RecursiveAnchor) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (r *RecursiveAnchor) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(r).Register(uri, registry)
}

func (r *RecursiveAnchor) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(r).Resolve(pointer, uri)
}

func (r *RecursiveAnchor) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[RecursiveAnchor] Validating")
	if schCtx.RecursiveAnchor == nil {
		schCtx.RecursiveAnchor = schCtx.Local
	}
}

func (r *RecursiveAnchor) UnmarshalJSON(data []byte) error {
	sch := &Schema{}
	if err := json.Unmarshal(data, sch); err != nil {
		return err
	}
	*r = (RecursiveAnchor)(*sch)
	return nil
}

//
// $defs
//

type Defs map[string]*Schema

func NewDefs() Keyword {
	return &Defs{}
}

func (d Defs) Validate(propPath string, data interface{}, errs *[]KeyError) {}

func (d *Defs) Register(uri string, registry *SchemaRegistry) {
	for _, v := range *d {
		v.Register(uri, registry)
	}
}

func (d *Defs) Resolve(pointer jptr.Pointer, uri string) *Schema {
	if pointer == nil {
		return nil
	}
	current := pointer.Head()
	if current == nil {
		return nil
	}

	if schema, ok := (*d)[*current]; ok {
		return schema.Resolve(pointer.Tail(), uri)
	}

	return nil
}

func (p Defs) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Defs] Validating")
}

func (d Defs) JSONProp(name string) interface{} {
	return d[name]
}

func (d Defs) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for key, sch := range d {
		res[key] = sch
	}
	return
}

//
// VOID
//

type Void struct{}

func NewVoid() Keyword {
	return &Void{}
}

func (vo *Void) Validate(propPath string, data interface{}, errs *[]KeyError) {
	SchemaDebug("[Void] Validating")
	SchemaDebug("[Void] WARNING this is a placeholder and should not be used")
	SchemaDebug("[Void] Void is always true")
}

func (vo *Void) Register(uri string, registry *SchemaRegistry) {}

func (vo *Void) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

func (vo *Void) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Void] Validating")
	SchemaDebug("[Void] WARNING this is a placeholder and should not be used")
	SchemaDebug("[Void] Void is always true")
}
