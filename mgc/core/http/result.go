package http

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/getkin/kin-openapi/openapi3"
)

type HttpResult interface {
	core.Result
	Request() *http.Request
	RequestBody() any // pre-Marshal, the original structured body data, if any. If request used an io.Reader, this is nil
	Response() *http.Response
	ResponseBody() any // post-Unmarshal, the decoded structured body data, if any; or io.Reader if not a structured body data
}

type httpResult struct {
	SourceData       core.ResultSource
	RequestData      *MarshalableRequest
	RequestBodyData  any // pre-Marshal, the original structured body data, if any. If request used an io.Reader, this is nil
	ResponseData     *MarshalableResponse
	ResponseBodyData any // post-Unmarshal, the decoded structured body data, if any; or io.Reader if not a structured body data
}

func NewZeroHttpResult() *httpResult {
	return &httpResult{}
}

type httpResultWithValue struct {
	httpResult
	ResultSchema *core.Schema
	ResultValue  core.Value
}

func NewZeroHttpResultWithValue() *httpResultWithValue {
	return &httpResultWithValue{}
}

type httpResultWithReader struct {
	httpResult
	BodyReader io.Reader
}

func NewZeroHttpResultWithReader() *httpResultWithReader {
	return &httpResultWithReader{}
}

type httpResultWithMultipart struct {
	httpResult
	BodyMultipart *multipart.Part
}

func NewZeroHttpResultWithMultipart() *httpResultWithMultipart {
	return &httpResultWithMultipart{}
}

// Takes over a response, unwrap it and create a Result based on it.
//
// The requestBody is the structured body data before it was marshalled to bytes. If the
// request was built from an io.Reader, such as file uploads or multipart, then this should be nil.
//
// The result is heavily dependent on the output of UnwrapResponse(), if data is:
//   - io.Reader: then implements ResultWithReader and Reader() returns it;
//   - multipart.Part: then implements ResultWithMultipart and Multipart() returns it;
//   - else: implements ResultWithValue and the decoded structured data is returned by Value().
//     It may be transformed/converted with getValueFromResponseBody() if it's non-nil.
func NewHttpResult(
	source core.ResultSource,
	schema *core.Schema,
	request *http.Request,
	requestBody any, // pre-Marshal, the original JSON, if any
	response *http.Response,
	getValueFromResponseBody func(responseBody any) (core.Value, error),
) (r HttpResult, err error) {
	result := httpResult{
		SourceData:      source,
		RequestData:     (*MarshalableRequest)(request),
		ResponseData:    (*MarshalableResponse)(response),
		RequestBodyData: requestBody,
	}

	result.ResponseBodyData, err = UnwrapResponse[any](response, request)
	if err != nil {
		return
	}

	switch v := result.ResponseBodyData.(type) {
	case *multipart.Part:
		return &httpResultWithMultipart{result, v}, nil
	case io.Reader:
		return &httpResultWithReader{result, v}, nil
	default:
		var value core.Value
		if getValueFromResponseBody == nil {
			value = v
		} else {
			value, err = getValueFromResponseBody(v)
			if err != nil {
				return
			}
		}
		return &httpResultWithValue{result, schema, value}, nil
	}
}

func (r *httpResult) Source() core.ResultSource {
	return r.SourceData
}

func (r *httpResult) Request() *http.Request {
	return (*http.Request)(r.RequestData)
}

func (r *httpResult) RequestBody() any {
	return r.RequestBodyData
}

func (r *httpResult) Response() *http.Response {
	return (*http.Response)(r.ResponseData)
}

func (r *httpResult) ResponseBody() any {
	return r.ResponseBodyData
}

func (r *httpResult) Encode() ([]byte, error) {
	return json.Marshal(*r)
}

func (r *httpResult) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

var _ HttpResult = (*httpResult)(nil)

func (r *httpResultWithValue) Unwrap() core.Result {
	return &r.httpResult
}

func (r *httpResultWithValue) Schema() *core.Schema {
	return r.ResultSchema
}

func (r *httpResultWithValue) ValidateSchema() error {
	return r.ResultSchema.VisitJSON(r.ResultValue, openapi3.MultiErrors())
}

func (r *httpResultWithValue) Value() core.Value {
	return r.ResultValue
}

var _ HttpResult = (*httpResultWithValue)(nil)
var _ core.ResultWrapper = (*httpResultWithValue)(nil)
var _ core.ResultWithValue = (*httpResultWithValue)(nil)

func (r *httpResultWithReader) Unwrap() core.Result {
	return &r.httpResult
}

func (r *httpResultWithReader) Reader() io.Reader {
	return r.BodyReader
}

var _ HttpResult = (*httpResultWithReader)(nil)
var _ core.ResultWrapper = (*httpResultWithReader)(nil)
var _ core.ResultWithReader = (*httpResultWithReader)(nil)

func (r *httpResultWithMultipart) Unwrap() core.Result {
	return &r.httpResult
}

func (r *httpResultWithMultipart) Multipart() *multipart.Part {
	return r.BodyMultipart
}

var _ HttpResult = (*httpResultWithMultipart)(nil)
var _ core.ResultWrapper = (*httpResultWithMultipart)(nil)
var _ core.ResultWithMultipart = (*httpResultWithMultipart)(nil)
