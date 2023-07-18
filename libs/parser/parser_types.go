package parser

import "github.com/getkin/kin-openapi/openapi3"

type HttpMethod string

const (
	GET    HttpMethod = "get"
	PUT    HttpMethod = "put"
	POST   HttpMethod = "post"
	DELETE HttpMethod = "delete"
	PATCH  HttpMethod = "patch"
)

var AllHttpMethods = [5]HttpMethod{GET, PUT, POST, DELETE, PATCH}

type OpenAPIParameterLocation string
type OpenAPIParameterStyle string

const (
	QUERY  OpenAPIParameterLocation = "query"
	HEADER OpenAPIParameterLocation = "header"
	PATH   OpenAPIParameterLocation = "path"
	COOKIE OpenAPIParameterLocation = "cookie"
)

const (
	MATRIX          OpenAPIParameterStyle = "matrix"
	LABEL           OpenAPIParameterStyle = "label"
	FORM            OpenAPIParameterStyle = "form"
	SIMPLE          OpenAPIParameterStyle = "simple"
	SPACE_DELIMITED OpenAPIParameterStyle = "space_delimited"
	PIPE_DELIMITED  OpenAPIParameterStyle = "pipe_delimited"
	DEEP_OBJECT     OpenAPIParameterStyle = "deep_object"
)

type OpenAPIFileInfo struct {
	Name        string
	Extension   string
	Path        string
	Description string
	Version     string
}

type OpenAPITag struct {
	Name        string
	Description string
}

type openAPIContext struct {
	ServerURL            string
	Tags                 []*OpenAPITag
	SecurityRequirements openapi3.SecurityRequirements
}

type OpenAPIModule struct {
	Name                 string
	Description          string
	Version              string
	ServerURL            string
	Tags                 []*OpenAPITag
	SecurityRequirements *openapi3.SecurityRequirements
	Actions              []*OpenAPIAction
}

type OpenAPIAction struct {
	Name              string
	Summary           string
	Description       string
	ServerURL         string
	PathName          string
	HttpMethod        HttpMethod
	Tags              []*OpenAPITag
	Deprecated        bool
	PathParams        []*OpenAPIParameter
	HeaderParams      []*OpenAPIParameter
	RequestBodyParams []*OpenAPIParameter
	Request           *openapi3.RequestBodyRef
	Responses         openapi3.Responses
	Security          *openapi3.SecurityRequirements
}

type openAPIActionContext struct {
	ServerURL            string
	Parameters           openapi3.Parameters
	Summary              string
	Description          string
	Tags                 []*OpenAPITag
	SecurityRequirements openapi3.SecurityRequirements
}

type OpenAPIParameter struct {
	Type        string
	Name        string
	Required    bool
	DisplayName string
	Description string
	Explode     bool
	Default     any
	Example     any
}
