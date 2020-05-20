package jsonschema

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	jptr "github.com/qri-io/jsonpointer"
)

// SchemaURI defines the $schema JSON Schema keyword
type SchemaURI string

// NewSchemaURI allocates a new SchemaURI keyword
func NewSchemaURI() Keyword {
	return new(SchemaURI)
}

// ValidateFromContext implements the Keyword interface for SchemaURI
func (s *SchemaURI) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[SchemaURI] Validating")
}

// Register implements the Keyword interface for SchemaURI
func (s *SchemaURI) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for SchemaURI
func (s *SchemaURI) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ID defines the $id JSON Schema keyword
type ID string

// NewID allocates a new Id keyword
func NewID() Keyword {
	return new(ID)
}

// ValidateFromContext implements the Keyword interface for ID
func (i *ID) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Id] Validating")
	// TODO(arqu): make sure ID is valid URI for draft2019
}

// Register implements the Keyword interface for ID
func (i *ID) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for ID
func (i *ID) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// Description defines the description JSON Schema keyword
type Description string

// NewDescription allocates a new Description keyword
func NewDescription() Keyword {
	return new(Description)
}

// ValidateFromContext implements the Keyword interface for Description
func (d *Description) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Description] Validating")
}

// Register implements the Keyword interface for Description
func (d *Description) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for Description
func (d *Description) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// Title defines the title JSON Schema keyword
type Title string

// NewTitle allocates a new Title keyword
func NewTitle() Keyword {
	return new(Title)
}

// ValidateFromContext implements the Keyword interface for Title
func (t *Title) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Title] Validating")
}

// Register implements the Keyword interface for Title
func (t *Title) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for Title
func (t *Title) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// Comment defines the comment JSON Schema keyword
type Comment string

// NewComment allocates a new Comment keyword
func NewComment() Keyword {
	return new(Comment)
}

// ValidateFromContext implements the Keyword interface for Comment
func (c *Comment) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Comment] Validating")
}

// Register implements the Keyword interface for Comment
func (c *Comment) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for Comment
func (c *Comment) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// Default defines the default JSON Schema keyword
type Default struct {
	data interface{}
}

// NewDefault allocates a new Default keyword
func NewDefault() Keyword {
	return &Default{}
}

// ValidateFromContext implements the Keyword interface for Default
func (d *Default) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Default] Validating")
}

// Register implements the Keyword interface for Default
func (d *Default) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for Default
func (d *Default) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for Default
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

// Examples defines the examples JSON Schema keyword
type Examples []interface{}

// NewExamples allocates a new Examples keyword
func NewExamples() Keyword {
	return new(Examples)
}

// ValidateFromContext implements the Keyword interface for Examples
func (e *Examples) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Examples] Validating")
}

// Register implements the Keyword interface for Examples
func (e *Examples) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for Examples
func (e *Examples) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ReadOnly defines the readOnly JSON Schema keyword
type ReadOnly bool

// NewReadOnly allocates a new ReadOnly keyword
func NewReadOnly() Keyword {
	return new(ReadOnly)
}

// ValidateFromContext implements the Keyword interface for ReadOnly
func (r *ReadOnly) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[ReadOnly] Validating")
}

// Register implements the Keyword interface for ReadOnly
func (r *ReadOnly) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for ReadOnly
func (r *ReadOnly) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// WriteOnly defines the writeOnly JSON Schema keyword
type WriteOnly bool

// NewWriteOnly allocates a new WriteOnly keyword
func NewWriteOnly() Keyword {
	return new(WriteOnly)
}

// ValidateFromContext implements the Keyword interface for WriteOnly
func (w *WriteOnly) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[WriteOnly] Validating")
}

// Register implements the Keyword interface for WriteOnly
func (w *WriteOnly) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for WriteOnly
func (w *WriteOnly) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// Ref defines the $ref JSON Schema keyword
type Ref struct {
	reference         string
	resolved          *Schema
	resolvedRoot      *Schema
	resolvedFragment  *jptr.Pointer
	fragmentLocalized bool
}

// NewRef allocates a new Ref keyword
func NewRef() Keyword {
	return new(Ref)
}

// ValidateFromContext implements the Keyword interface for Ref
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

// _resolveRef attempts to resolve the reference from the top-level context
func (r *Ref) _resolveRef(schCtx *SchemaContext) {
	if IsLocalSchemaID(r.reference) {
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
					address, _ = SafeResolveURL(uriFolder, address)
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

	localURI := schCtx.BaseURI
	if r.resolvedRoot != nil && r.resolvedRoot.docPath != "" {
		localURI = r.resolvedRoot.docPath
		if r.fragmentLocalized && !r.resolvedFragment.IsEmpty() {
			current := r.resolvedFragment.Head()
			sch := schCtx.LocalRegistry.GetLocal("#" + *current)
			if sch != nil {
				r.resolved = sch
				return
			}
		}
	}
	r._resolveLocalRef(localURI)
}

// _resolveLocalRef attempts to resolve the reference from a local context
func (r *Ref) _resolveLocalRef(uri string) {
	if r.resolvedFragment.IsEmpty() {
		r.resolved = r.resolvedRoot
		return
	}

	if r.resolvedRoot != nil {
		r.resolved = r.resolvedRoot.Resolve(*r.resolvedFragment, uri)
	}
}

// Register implements the Keyword interface for Ref
func (r *Ref) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for Ref
func (r *Ref) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for Ref
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

// MarshalJSON implements the json.Marshaler interface for Ref
func (r Ref) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.reference)
}

// RecursiveRef defines the $recursiveRef JSON Schema keyword
type RecursiveRef struct {
	reference        string
	resolved         *Schema
	resolvedRoot     *Schema
	resolvedFragment *jptr.Pointer

	validatingLocations map[string]bool
}

// NewRecursiveRef allocates a new RecursiveRef keyword
func NewRecursiveRef() Keyword {
	return new(RecursiveRef)
}

// ValidateFromContext implements the Keyword interface for RecursiveRef
func (r *RecursiveRef) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[RecursiveRef] Validating")
	if r.isLocationVisited(schCtx.InstanceLocation.String()) {
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

func (r *RecursiveRef) isLocationVisited(location string) bool {
	if r.validatingLocations == nil {
		return false
	}
	v, ok := r.validatingLocations[location]
	if !ok {
		return false
	}
	return v
}

// _resolveRef attempts to resolve the reference from the top-level context
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

	if IsLocalSchemaID(r.reference) {
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
						address, _ = SafeResolveURL(uriFolder, address)
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

	localURI := schCtx.BaseURI
	if r.resolvedRoot != nil && r.resolvedRoot.docPath != "" {
		localURI = r.resolvedRoot.docPath
	}
	r._resolveLocalRef(localURI)
}

// _resolveLocalRef attempts to resolve the reference from a local context
func (r *RecursiveRef) _resolveLocalRef(uri string) {
	if r.resolvedFragment.IsEmpty() {
		r.resolved = r.resolvedRoot
		return
	}

	if r.resolvedRoot != nil {
		r.resolved = r.resolvedRoot.Resolve(*r.resolvedFragment, uri)
	}
}

// Register implements the Keyword interface for RecursiveRef
func (r *RecursiveRef) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for RecursiveRef
func (r *RecursiveRef) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for RecursiveRef
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

// MarshalJSON implements the json.Marshaler interface for RecursiveRef
func (r RecursiveRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.reference)
}

// Anchor defines the $anchor JSON Schema keyword
type Anchor string

// NewAnchor allocates a new Anchor keyword
func NewAnchor() Keyword {
	return new(Anchor)
}

// ValidateFromContext implements the Keyword interface for Anchor
func (a *Anchor) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Anchor] Validating")
}

// Register implements the Keyword interface for Anchor
func (a *Anchor) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for Anchor
func (a *Anchor) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// RecursiveAnchor defines the $recursiveAnchor JSON Schema keyword
type RecursiveAnchor Schema

// NewRecursiveAnchor allocates a new RecursiveAnchor keyword
func NewRecursiveAnchor() Keyword {
	return &RecursiveAnchor{}
}

// Register implements the Keyword interface for RecursiveAnchor
func (r *RecursiveAnchor) Register(uri string, registry *SchemaRegistry) {
	(*Schema)(r).Register(uri, registry)
}

// Resolve implements the Keyword interface for RecursiveAnchor
func (r *RecursiveAnchor) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return (*Schema)(r).Resolve(pointer, uri)
}

// ValidateFromContext implements the Keyword interface for RecursiveAnchor
func (r *RecursiveAnchor) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[RecursiveAnchor] Validating")
	if schCtx.RecursiveAnchor == nil {
		schCtx.RecursiveAnchor = schCtx.Local
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface for RecursiveAnchor
func (r *RecursiveAnchor) UnmarshalJSON(data []byte) error {
	sch := &Schema{}
	if err := json.Unmarshal(data, sch); err != nil {
		return err
	}
	*r = (RecursiveAnchor)(*sch)
	return nil
}

// Defs defines the $defs JSON Schema keyword
type Defs map[string]*Schema

// NewDefs allocates a new Defs keyword
func NewDefs() Keyword {
	return &Defs{}
}

// Register implements the Keyword interface for Defs
func (d *Defs) Register(uri string, registry *SchemaRegistry) {
	for _, v := range *d {
		v.Register(uri, registry)
	}
}

// Resolve implements the Keyword interface for Defs
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

// ValidateFromContext implements the Keyword interface for Defs
func (d Defs) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Defs] Validating")
}

// JSONProp implements the JSONPather for Defs
func (d Defs) JSONProp(name string) interface{} {
	return d[name]
}

// JSONChildren implements the JSONContainer interface for Defs
func (d Defs) JSONChildren() (res map[string]JSONPather) {
	res = map[string]JSONPather{}
	for key, sch := range d {
		res[key] = sch
	}
	return
}

// Void is a placeholder definition for a keyword
type Void struct{}

// NewVoid allocates a new Void keyword
func NewVoid() Keyword {
	return &Void{}
}

// Register implements the Keyword interface for Void
func (vo *Void) Register(uri string, registry *SchemaRegistry) {}

// Resolve implements the Keyword interface for Void
func (vo *Void) Resolve(pointer jptr.Pointer, uri string) *Schema {
	return nil
}

// ValidateFromContext implements the Keyword interface for Void
func (vo *Void) ValidateFromContext(schCtx *SchemaContext, errs *[]KeyError) {
	SchemaDebug("[Void] Validating")
	SchemaDebug("[Void] WARNING this is a placeholder and should not be used")
	SchemaDebug("[Void] Void is always true")
}
