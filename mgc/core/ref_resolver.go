package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	mgcSchema "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/jsonpointer"
)

// Similar to JSON Pointers - [RFC6901]
//
// The difference is that it allow an URL to be specified before the JSON Pointer,
// so both are valid:
//   - /path/to/element
//   - http://some.url.com#/path/to/element
//
// [RFC6901]: https://datatracker.ietf.org/doc/html/rfc6901
type RefPath string

type RefPathResolver interface {
	// validate the given string, convert it to RefPath and resolve it
	Resolve(path string) (result any, err error)
	ResolvePath(path RefPath) (result any, err error)
}

// refResolverContextKey is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type refResolverContextKey string

// theRefResolverContextKey is the key for sdk.Grouper values in Contexts. It is
// unexported; clients use NewGrouperContext() and GrouperFromContext()
// instead of using this key directly.
var theRefResolverContextKey refResolverContextKey = "github.com/MagaluCloud/magalu/mgc/core/RefPathResolver"

func NewRefPathResolverContext(parent context.Context, refResolver RefPathResolver) context.Context {
	return context.WithValue(parent, theRefResolverContextKey, refResolver)
}

func RefPathResolverFromContext(ctx context.Context) RefPathResolver {
	if value, ok := ctx.Value(theRefResolverContextKey).(RefPathResolver); !ok {
		return nil
	} else {
		return value
	}
}

func (path RefPath) Split() (parentPath RefPath, field string) {
	if path == "" || path == "/" {
		return
	}

	i := strings.LastIndex(string(path), "/")
	if i < 0 {
		return
	}

	parentPath = path[:i]
	field = jsonpointer.Unescape(string(path[i+1:]))
	return
}

var errorInvalidStart = errors.New(`JSON pointer must be empty or start with a "/"`)
var errorDocumentUrlUnsupported = errors.New("document URL is unsupported")

func (path RefPath) Validate() (err error) {
	_, p, _ := strings.Cut(string(path), "#")
	if p == "" {
		p = string(path)
	}
	if !strings.HasPrefix(p, "/") {
		return &RefPathResolveError{path, errorInvalidStart}
	}
	return nil
}

func (path RefPath) Add(parts ...string) RefPath {
	for _, part := range parts {
		path += RefPath("/" + jsonpointer.Escape(part))
	}
	return path
}

func (path RefPath) SplitUrl() (url string, p RefPath) {
	before, after, found := strings.Cut(string(path), "#")
	if !found {
		return "", path
	}
	return before, RefPath(after)
}

func (p *RefPath) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*p = RefPath(s)
	return p.Validate()
}

var _ json.Unmarshaler = (*RefPath)(nil)

type RefPathResolveError struct {
	Path RefPath
	Err  error
}

func (e *RefPathResolveError) Error() string {
	leaf := e.Err
	path := e.Path
	for {
		if e, ok := leaf.(*RefPathResolveError); ok {
			path = e.Path
			leaf = e.Err
		} else {
			break
		}
	}
	return fmt.Sprintf("could not resolve %q: %s", path, leaf.Error())
}

func (e *RefPathResolveError) Unwrap() error {
	return e.Err
}

var _ error = (*RefPathResolveError)(nil)

type MissingFieldError string

func (e MissingFieldError) Error() string {
	return fmt.Sprintf("missing field: %q", string(e))
}

var _ error = (*MissingFieldError)(nil)

// Wraps MultiRefPathResolver and automatically pass CurrentUrl
type BoundRefPathResolver struct {
	Resolver *MultiRefPathResolver

	CurrentUrl string // if RefPath's URL is this placeholder given to MultiRefPathResolver, use this URL instead.
}

var _ RefPathResolver = (*BoundRefPathResolver)(nil)

func NewBoundRefResolver(currentUrl string, resolver *MultiRefPathResolver) *BoundRefPathResolver {
	return &BoundRefPathResolver{resolver, currentUrl}
}

func (r *BoundRefPathResolver) Resolve(ref string) (result any, err error) {
	return r.Resolver.Resolve(ref, r.CurrentUrl)
}

func (r *BoundRefPathResolver) ResolvePath(path RefPath) (result any, err error) {
	return r.Resolver.ResolvePath(path, r.CurrentUrl)
}

// Discovers the document from the pointerPath and call the specific document resolver.
//
// if RefPath's URL is CurrentUrlPlaceholder, then use the given currentUrl. If the placeholder
// if empty, currentUrl is also used.
//
// if CurrentUrlPlaceholder is non-empty and RefPath's URL is empty, then use EmptyDocumentUrl.
type MultiRefPathResolver struct {
	Resolvers map[string]RefPathResolver // maps url => resolver

	CurrentUrlPlaceholder string // if RefPath's URL is this placeholder, use the given currentUrl
	EmptyDocumentUrl      string // if RefPath's URL is empty, then use this URL.
}

func NewMultiRefPathResolver() *MultiRefPathResolver {
	return &MultiRefPathResolver{
		Resolvers: map[string]RefPathResolver{},
	}
}

func (r *MultiRefPathResolver) Resolve(ref string, currentUrl string) (result any, err error) {
	path := RefPath(ref)
	err = path.Validate()
	if err != nil {
		return
	}
	return r.ResolvePath(path, currentUrl)
}

func (r *MultiRefPathResolver) ResolvePath(path RefPath, currentUrl string) (result any, err error) {
	url, p := path.SplitUrl()
	if url == r.CurrentUrlPlaceholder {
		url = currentUrl
	} else if url == "" {
		url = r.EmptyDocumentUrl
	}

	docResolver := r.Resolvers[url]
	if docResolver == nil {
		return nil, &RefPathResolveError{path, fmt.Errorf("unknown document %q", url)}
	}
	return docResolver.ResolvePath(p)
}

// Associates the given URL to a document resolver
func (r *MultiRefPathResolver) Add(url string, docResolver RefPathResolver) error {
	if _, ok := r.Resolvers[url]; ok {
		return fmt.Errorf("document resolver already added: %q", url)
	}
	r.Resolvers[url] = docResolver
	return nil
}

// Resolve within a single document, walks from root group to leaf Executors and their sub-fields
//
// Most users will want BoundRefPathResolver instead.
type DocumentRefPathResolver struct {
	cache   sync.Map // JSON Pointer string => resolved value
	getRoot func() (any, error)
}

var _ RefPathResolver = (*DocumentRefPathResolver)(nil)

// NOTE: getRoot() will be called multiple times, make sure it caches its own results, ex utils.NewLazyLoaderWithError()
func NewDocumentRefPathResolver(getRoot utils.LoadWithError[any]) *DocumentRefPathResolver {
	return &DocumentRefPathResolver{
		cache:   sync.Map{},
		getRoot: getRoot,
	}
}

func (r *DocumentRefPathResolver) Resolve(ref string) (result any, err error) {

	path := RefPath(ref)
	err = path.Validate()
	if err != nil {
		return
	}
	return r.ResolvePath(path)
}

func (r *DocumentRefPathResolver) ResolvePath(path RefPath) (result any, err error) {
	if path == "" || path == "/" {
		var root any
		root, err = r.getRoot()
		if err != nil {
			return
		}

		result = root
		return
	}

	if cachedResult, ok := r.cache.Load(path); ok {
		return cachedResult, nil
	}

	url, _ := path.SplitUrl()
	if url != "" {
		err = &RefPathResolveError{path, errorDocumentUrlUnsupported}
		return
	}

	parentPath, field := path.Split()
	parent, err := r.ResolvePath(parentPath)
	if err != nil {
		err = &RefPathResolveError{path, err}
		return
	}
	if field == "" { // not expected to have "parent/", but okay...
		result = parent
	} else {
		result, err = resolveField(field, parent)
	}

	if err == nil {
		r.cache.Store(path, result)
	} else {
		err = &RefPathResolveError{path, err}
	}

	return
}

func resolveField(field string, doc any) (result any, err error) {
	switch v := doc.(type) {
	case Grouper:
		return resolveGrouper(field, v)

	case Executor:
		return resolveExecutor(field, v)

	case Linker:
		return resolveLinker(field, v)

	case Links:
		return resolveLinkerMap(field, v)

	case map[string]Executor:
		return resolveExecutorMap(field, v)

	default:
		return resolveGeneric(field, v)
	}
}

func resolveGeneric(field string, doc any) (result any, err error) {
	result, _, err = jsonpointer.GetForToken(doc, field)
	if err != nil {
		return
	}

	return
}

func resolveDescriptor(field string, desc Descriptor) (result any, err error) {
	switch field {
	case "name":
		return desc.Name(), nil
	case "description":
		return desc.Description(), nil
	}

	return nil, MissingFieldError(field)
}

func resolveGrouper(field string, group Grouper) (result any, err error) {
	result, err = group.GetChildByName(field)
	if err != nil {
		// if no child with names clobbering Descriptor fields, resolve those:
		result, err = resolveDescriptor(field, group)
	}
	return
}

func resolveExecutor(name string, exec Executor) (result any, err error) {
	switch name {
	case "parametersSchema":
		return exec.ParametersSchema(), nil
	case "configsSchema":
		return exec.ConfigsSchema(), nil
	case "resultSchema":
		return exec.ResultSchema(), nil
	case "links":
		return exec.Links(), nil
	case "related":
		return exec.Related(), nil
	default:
		return resolveDescriptor(name, exec)
	}
}

func resolveLinker(field string, linker Linker) (result any, err error) {
	switch field {
	case "name":
		return linker.Name(), nil
	case "description":
		return linker.Description(), nil
	case "parametersSchema", "additionalParametersSchema":
		return linker.AdditionalParametersSchema(), nil
	case "configsSchema", "additionalConfigsSchema":
		return linker.AdditionalConfigsSchema(), nil
	case "resultSchema":
		return linker.ResultSchema(), nil
	default:
		return nil, MissingFieldError(field)
	}
}

func resolveLinkerMap(field string, m Links) (result Linker, err error) {
	if result, ok := m[field]; ok {
		return result, nil
	}
	return nil, MissingFieldError(field)
}

func resolveExecutorMap(field string, m map[string]Executor) (result Executor, err error) {
	if result, ok := m[field]; ok {
		return result, nil
	}
	return nil, MissingFieldError(field)
}

const extensionResolvedKey = "github.com/MagaluCloud/magalu/mgc/core/resolved"

func isSchemaResolved(schema *mgcSchema.Schema) bool {
	return schema.Extensions[extensionResolvedKey] == "true"
}

func markSchemaResolved(schema *mgcSchema.Schema) {
	if schema.Extensions == nil {
		schema.Extensions = map[string]any{}
	}
	schema.Extensions[extensionResolvedKey] = "true"
}

// This changes schema *IN PLACE* so future resolutions will use the same pointers
//
// This guarantees the Schema and internal elements are fully resolved.
func ResolveSchemaChildren(r RefPathResolver, schema *mgcSchema.Schema) (result *mgcSchema.Schema, err error) {
	if isSchemaResolved(schema) {
		return schema, nil
	}

	var pendingResolution mgcSchema.SchemaRefs

	addPending := func(ref *mgcSchema.SchemaRef) {
		if ref != nil && ref.Value == nil {
			pendingResolution = append(pendingResolution, ref)
		}
	}
	addPendingSlice := func(refs []*mgcSchema.SchemaRef) {
		for _, ref := range refs {
			addPending(ref)
		}
	}

	for _, ref := range schema.Properties {
		addPending(ref)
	}

	addPending(schema.Not)
	addPending(schema.Items)
	addPending(schema.AdditionalProperties.Schema)
	addPendingSlice(schema.OneOf)
	addPendingSlice(schema.AnyOf)
	addPendingSlice(schema.AllOf)

	for _, ref := range pendingResolution {
		value, err := ResolveSchemaRef(r, ref)
		if err != nil {
			return nil, err
		}
		ref.Value = (*openapi3.Schema)(value)
	}

	markSchemaResolved(schema)
	return schema, nil
}

func ResolveSchema(refResolver RefPathResolver, ref string) (result *mgcSchema.Schema, err error) {
	path := RefPath(ref)
	err = path.Validate()
	if err != nil {
		return
	}
	return ResolveSchemaPath(refResolver, path)
}

func ResolveSchemaPath(refResolver RefPathResolver, path RefPath) (result *mgcSchema.Schema, err error) {
	v, err := refResolver.ResolvePath(path)
	if err != nil {
		return
	}

	switch t := v.(type) {
	case *mgcSchema.Schema:
		return ResolveSchemaChildren(refResolver, t)
	case *openapi3.Schema:
		return ResolveSchemaChildren(refResolver, (*mgcSchema.Schema)(t))
	case *mgcSchema.SchemaRef: // same as openapi3.SchemaRef
		return ResolveSchemaRef(refResolver, t)
	default:
		return nil, fmt.Errorf("could not resolve %q: expected type %T, got %T", path, result, t)
	}
}

// This changes schemaRef *IN PLACE* so future resolutions will use the same pointers
//
// This guarantees the SchemaRef, Schema and internal elements are fully resolved.
func ResolveSchemaRef(refResolver RefPathResolver, schemaRef *mgcSchema.SchemaRef) (result *mgcSchema.Schema, err error) {
	if schemaRef.Value != nil {
		return ResolveSchemaChildren(refResolver, (*mgcSchema.Schema)(schemaRef.Value))
	}

	ref := schemaRef.Ref
	if ref == "" {
		return nil, fmt.Errorf("empty schema reference")
	}

	result, err = ResolveSchema(refResolver, ref)
	if err != nil {
		return
	}
	schemaRef.Ref = "" // force kin-openapi/openapi3 to MarshalJSON()/MarshalYAML() the full object
	schemaRef.Value = (*openapi3.Schema)(result)
	return
}

func ResolveExecutor(refResolver RefPathResolver, ref string) (exec Executor, err error) {
	path := RefPath(ref)
	err = path.Validate()
	if err != nil {
		return
	}
	return ResolveExecutorPath(refResolver, path)
}

func ResolveExecutorPath(refResolver RefPathResolver, path RefPath) (exec Executor, err error) {
	v, err := refResolver.ResolvePath(path)
	if err != nil {
		return
	}
	exec, ok := v.(Executor)
	if !ok {
		err = fmt.Errorf("path must be an executor, %q is %#v", path, v)
		return
	}

	return exec, nil
}
