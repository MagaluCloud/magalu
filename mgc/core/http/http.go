package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/core/xml"
)

// contextKey is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type contextKey string

var httpClientKey contextKey = "magalu.cloud/core/Transport"

type Client struct {
	http.Client
}

func NewClient(transport http.RoundTripper) *Client {
	return &Client{http.Client{Transport: transport}}
}

func NewClientContext(parent context.Context, client *Client) context.Context {
	return context.WithValue(parent, httpClientKey, client)
}

func ClientFromContext(context context.Context) *Client {
	client, ok := context.Value(httpClientKey).(*Client)
	if !ok {
		logger().Debugf("Error casting ctx %s to *HttpClient", httpClientKey)
		return nil
	}
	return client
}

func BodyReaderSafe(resp *http.Response) (io.ReadCloser, error) {
	bodyContents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyContents))
	return io.NopCloser(bytes.NewBuffer(bodyContents)), nil
}

func DecodeJSON[T core.Value](resp *http.Response, data *T) error {
	body, err := BodyReaderSafe(resp)
	if err != nil {
		return fmt.Errorf("error when reading response body: %w", err)
	}
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(data)
	if err != nil {
		return fmt.Errorf("error decoding JSON response body: %w", err)
	}
	return nil
}

func DecodeXML[T core.Value](resp *http.Response, data *T) error {
	body, err := BodyReaderSafe(resp)
	if err != nil {
		return fmt.Errorf("error when reading response body: %w", err)
	}
	decoder := xml.NewDecoder(body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(data)
	if err != nil {
		return fmt.Errorf("error decoding XML response body: %w", err)
	}
	return nil
}

type HttpError struct {
	Code    int
	Status  string
	Headers http.Header
	Payload []byte
	Message string // MGC reports this in the json body
	Slug    string // MGC reports this in the json body
}

type BaseApiError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type IdentifiableHttpError struct {
	*HttpError
	RequestID string `json:"requestID"`
}

func (e *IdentifiableHttpError) Unwrap() error {
	return e.HttpError
}

func (e *IdentifiableHttpError) Error() string {
	msg := e.HttpError.Error()
	if e.RequestID != "" {
		msg += " (request-id: " + e.RequestID + ")"
	}
	return msg
}

func (e *HttpError) Error() string {
	msg := e.Message
	if e.Status != msg {
		msg = e.Status + " - " + msg
	}
	if e.Slug != "" {
		msg = "(" + e.Slug + ")" + " " + msg
	}

	return msg
}

func (e *HttpError) String() string {
	return fmt.Sprintf("%T{Status: %q, Slug: %q, Message: %q}", e, e.Status, e.Slug, e.Message)
}

func NewHttpErrorFromResponse(resp *http.Response, req *http.Request) *IdentifiableHttpError {
	slug := "unknown"
	message := resp.Status

	defer resp.Body.Close()
	payload, _ := io.ReadAll(resp.Body)

	contentType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		logger().Debugw("ignored invalid response", "Content-Type", resp.Header.Get("Content-Type"), "error", err.Error())
	}
	if contentType == "application/json" {
		data := BaseApiError{}
		if err := json.Unmarshal(payload, &data); err == nil {
			if data.Message != "" {
				message = data.Message
			}
			if data.Slug != "" {
				slug = data.Slug
			}
		}
	}

	httpError := &HttpError{
		Code:    resp.StatusCode,
		Status:  resp.Status,
		Headers: resp.Header,
		Payload: payload,
		Message: message,
		Slug:    slug,
	}

	return NewIdentifiableHttpError(httpError, req, resp)

}

func NewIdentifiableHttpError(httpError *HttpError, request *http.Request, response *http.Response) *IdentifiableHttpError {
	requestId := ""
	if response != nil {
		if id := response.Header.Get("X-Request-Id"); id != "" {
			requestId = id
		}
	}
	if id := request.Header.Get("X-Request-Id"); id != "" {
		requestId = id
	}

	return &IdentifiableHttpError{
		HttpError: httpError,
		RequestID: requestId,
	}
}

// Handles the response, and tries to convert the data to T
//
// If the Content-Type header starts with "multipart/", then a pointer to multipart.Part
// is returned as data.
//
// If the Content-Type header is one of:
//   - application/json
//   - application/xml
//
// Then it will be decoded.
//
// If the Content-Type is none of the above, then io.ReadCloser (resp.Body) is returned.
//
// To avoid errors when the result type isn't known, UnwrapResponse[any] can be used.
func UnwrapResponse[T any](resp *http.Response, req *http.Request) (result T, err error) {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = NewHttpErrorFromResponse(resp, req)
		return
	}

	if resp.StatusCode == 204 {
		return
	}

	contentType, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))

	switch {
	default:
		err = utils.AssignToT(&result, resp.Body)
		return
	case strings.HasPrefix(contentType, "multipart/"):
		body, bodyErr := BodyReaderSafe(resp)
		if bodyErr != nil {
			err = fmt.Errorf("error when reading response body: %w", bodyErr)
			return
		}
		// TODO: do we have multi-part downloads? or just single?
		// If multi, then we need to return a multipart reader...
		// return multipart.NewReader(resp.Body, params["boundary"]), nil
		r := multipart.NewReader(body, params["boundary"])
		nextPart, npErr := r.NextPart()
		err = npErr
		if err != nil {
			return
		}
		err = utils.AssignToT(&result, nextPart)
		return
	case contentType == "application/json":
		err = DecodeJSON(resp, &result)
	case contentType == "application/xml":
		err = DecodeXML(resp, &result)
	}

	return
}

// Checks if the response's StatusCode is less than 200 or greater equal than 300. If so, returns an error of type *HttpError
func ExtractErr(resp *http.Response, req *http.Request) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return NewHttpErrorFromResponse(resp, req)
	}
	return nil
}

var defaultTransport *http.Transport

func DefaultTransport() http.RoundTripper {
	if defaultTransport == nil {
		defaultTransport = (http.DefaultTransport).(*http.Transport)
		defaultTransport.MaxIdleConns = 10
		defaultTransport.IdleConnTimeout = 30 * time.Second
	}
	return defaultTransport
}
