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
	"reflect"
	"strings"
	"time"

	"magalu.cloud/core"
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

func bodyReaderSafe(resp *http.Response) (io.ReadCloser, error) {
	bodyContents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyContents))
	return io.NopCloser(bytes.NewBuffer(bodyContents)), nil
}

func DecodeJSON[T core.Value](resp *http.Response, data *T) error {
	body, err := bodyReaderSafe(resp)
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
	body, err := bodyReaderSafe(resp)
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

type ObjectStorageXMLError struct {
	Message string `xml:"Message"`
	Code    string `xml:"Code"`
}

func (e *HttpError) Error() string {
	return e.Message
}

func (e *HttpError) String() string {
	return fmt.Sprintf("%T{Status: %q, Slug: %q, Message: %q}", e, e.Status, e.Slug, e.Message)
}

func NewHttpErrorFromResponse(resp *http.Response) *HttpError {
	slug := "unknown"
	message := resp.Status

	defer resp.Body.Close()
	payload, _ := io.ReadAll(resp.Body)

	contentType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		logger().Debugw("ignored invalid response", "Content-Type", resp.Header.Get("Content-Type"), "error", err.Error())
	}
	// if contentType == "application/json" {
	// 	fmt.Println("decoding as json")
	// 	data := BaseApiError{}
	// 	if err := json.Unmarshal(payload, &data); err == nil {
	// 		if data.Message != "" {
	// 			message = data.Message
	// 		}
	// 		if data.Slug != "" {
	// 			slug = data.Slug
	// 		}
	// 	}
	if contentType == "application/json" {
		fmt.Println("decoding as xml")
		data := ObjectStorageXMLError{}
		decoder := xml.NewDecoder(bytes.NewBuffer(payload))
		if err := decoder.Decode(&data); err == nil {
			if data.Message != "" {
				message = data.Message
			}
			if data.Code != "" {
				slug = data.Code
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println("wtf is this content type?", contentType)
	}

	return &HttpError{
		Code:    resp.StatusCode,
		Status:  resp.Status,
		Headers: resp.Header,
		Payload: payload,
		Message: message,
		Slug:    slug,
	}
}

func assignToT[T any, U any](t *T, u U) error {
	if uAsT, ok := any(u).(T); ok {
		*t = uAsT
		return nil
	}

	tVal := reflect.ValueOf(t).Elem()
	// Empty name means `any`
	if tVal.Type().Name() != "" {
		return fmt.Errorf("request response of type %T is not convertible to %T", *t, u)
	}

	tVal.Set(reflect.ValueOf(u))
	return nil
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
func UnwrapResponse[T any](resp *http.Response) (result T, err error) {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = NewHttpErrorFromResponse(resp)
		return
	}

	if resp.StatusCode == 204 {
		return
	}

	contentType, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))

	switch {
	default:
		err = assignToT(&result, resp.Body)
		return
	case strings.HasPrefix(contentType, "multipart/"):
		body, bodyErr := bodyReaderSafe(resp)
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
		err = assignToT(&result, nextPart)
		return
	case contentType == "application/json":
		err = DecodeJSON(resp, &result)
	case contentType == "application/xml":
		err = DecodeXML(resp, &result)
	}

	return
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
