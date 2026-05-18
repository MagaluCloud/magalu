package http

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/core/xml"
	"github.com/andybalholm/brotli"
)

// contextKey is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type contextKey string

var httpClientKey contextKey = "github.com/MagaluCloud/magalu/mgc/core/Transport"

type Client struct {
	http.Client
}

func NewClient(transport http.RoundTripper) *Client {
	return &Client{http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects, return an error to use the last response.
			return http.ErrUseLastResponse
		}}}
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

// DecompressResponse wraps resp.Body with the appropriate decoder based on the
// Content-Encoding header, then clears the encoding metadata so downstream
// consumers see a transparent stream. Net/http only auto-decompresses gzip
// when the caller did not set Accept-Encoding itself; we advertise it
// explicitly, so the decode is our responsibility.
func DecompressResponse(resp *http.Response) error {
	if resp == nil || resp.Body == nil {
		return nil
	}
	encoding := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Encoding")))
	if encoding == "" || encoding == "identity" {
		return nil
	}

	original := resp.Body
	var wrapped io.ReadCloser
	switch encoding {
	case "gzip", "x-gzip":
		gzr, err := gzip.NewReader(original)
		if err != nil {
			return fmt.Errorf("error creating gzip reader: %w", err)
		}
		wrapped = &decompressedBody{Reader: gzr, closers: []io.Closer{gzr, original}}
	case "deflate":
		// The HTTP spec defines "deflate" as zlib (RFC 1950), but Microsoft's IIS
		// historically sends raw DEFLATE (RFC 1951), a behavior subsequently
		// adopted by other servers. We inspect the header without consuming bytes
		// to select the appropriate reader.
		buffered := bufio.NewReader(original)
		var dr io.ReadCloser
		if hasZlibHeader(buffered) {
			zr, err := zlib.NewReader(buffered)
			if err != nil {
				return fmt.Errorf("error creating deflate reader: %w", err)
			}
			dr = zr
		} else {
			dr = flate.NewReader(buffered)
		}
		wrapped = &decompressedBody{Reader: dr, closers: []io.Closer{dr, original}}
	case "br":
		br := brotli.NewReader(original)
		wrapped = &decompressedBody{Reader: br, closers: []io.Closer{original}}
	default:
		return nil
	}

	resp.Body = wrapped
	resp.Header.Del("Content-Encoding")
	resp.Header.Del("Content-Length")
	resp.ContentLength = -1
	resp.Uncompressed = true
	return nil
}

// hasZlibHeader inspects the first 2 bytes to determine whether the
// stream is zlib (RFC 1950) or raw DEFLATE (RFC 1951). A valid zlib
// header uses the deflate compression method (low nibble of CMF == 8)
// and (CMF * 256 + FLG) must be a multiple of 31.
func hasZlibHeader(br *bufio.Reader) bool {
	header, err := br.Peek(2)
	if err != nil || len(header) < 2 {
		return false
	}
	cmf, flg := header[0], header[1]
	if cmf&0x0F != 8 {
		return false
	}
	return (uint16(cmf)<<8|uint16(flg))%31 == 0
}

type decompressedBody struct {
	io.Reader
	closers []io.Closer
}

func (d *decompressedBody) Close() error {
	var firstErr error
	for _, c := range d.closers {
		if err := c.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
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
	decoder.UseNumber()
	err = decoder.Decode(data)
	if err != nil {
		return fmt.Errorf("error decoding JSON response body: %w", err)
	}

	err = convertJSONNumbers(reflect.ValueOf(data).Elem())
	if err != nil {
		return fmt.Errorf("error converting JSON numbers: %w", err)
	}
	return nil
}

func convertJSONNumbers(v reflect.Value) error {
	if !v.IsValid() {
		return nil
	}
	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return nil
		}
		return convertJSONNumbers(v.Elem())
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return convertJSONNumbers(v.Elem())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if err := convertJSONNumbers(v.Field(i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		if v.IsNil() {
			return nil
		}
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			if val.Kind() == reflect.Interface {
				if val.IsNil() {
					continue
				}
				valElem := val.Elem()
				if valElem.Type() == reflect.TypeOf(json.Number("")) {
					num := valElem.Interface().(json.Number)
					if i, err := strconv.ParseInt(string(num), 10, 64); err == nil {
						v.SetMapIndex(key, reflect.ValueOf(i))
					} else if f, err := strconv.ParseFloat(string(num), 64); err == nil {
						v.SetMapIndex(key, reflect.ValueOf(f))
					}
				} else {
					newVal := reflect.New(valElem.Type()).Elem()
					newVal.Set(valElem)
					if err := convertJSONNumbers(newVal); err != nil {
						return err
					}
					v.SetMapIndex(key, newVal)
				}
			} else {
				newVal := reflect.New(val.Type()).Elem()
				newVal.Set(val)
				if err := convertJSONNumbers(newVal); err != nil {
					return err
				}
				v.SetMapIndex(key, newVal)
			}
		}
	case reflect.Slice:
		if v.IsNil() {
			return nil
		}
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			if item.Kind() == reflect.Interface {
				if item.IsNil() {
					continue
				}
				elemValue := item.Elem()
				if elemValue.Type() == reflect.TypeOf(json.Number("")) {
					num := elemValue.Interface().(json.Number)
					if item.CanSet() {
						if i, err := strconv.ParseInt(string(num), 10, 64); err == nil {
							item.Set(reflect.ValueOf(i))
						} else if f, err := strconv.ParseFloat(string(num), 64); err == nil {
							item.Set(reflect.ValueOf(f))
						}
					}
				} else {
					if err := convertJSONNumbers(item); err != nil {
						return err
					}
				}
			} else {
				if err := convertJSONNumbers(item); err != nil {
					return err
				}
			}
		}
	default:
		if v.Type() == reflect.TypeOf(json.Number("")) && v.CanSet() {
			num := v.Interface().(json.Number)
			if v.CanSet() {
				if i, err := strconv.ParseInt(string(num), 10, 64); err == nil {
					v.Set(reflect.ValueOf(i))
				} else if f, err := strconv.ParseFloat(string(num), 64); err == nil {
					v.Set(reflect.ValueOf(f))
				}
			}
		}
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
	TraceID   string `json:"traceID"`
}

func (e *IdentifiableHttpError) Unwrap() error {
	return e.HttpError
}

func (e *IdentifiableHttpError) Error() string {
	msg := "\n Status: " + e.HttpError.Error()
	if e.RequestID != "" {
		msg += "\n Request ID: " + e.RequestID
	}
	if e.TraceID != "" {
		msg += "\n MGC Trace ID: " + e.TraceID
	}
	if e.Payload != nil {
		msg += "\n\n" + string(e.Payload)
	}
	return msg
}

func (e *HttpError) Error() string {
	msg := e.Message
	if e.Status != msg {
		msg = e.Status + " - " + msg
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
	a := IdentifiableHttpError{
		HttpError: httpError,
	}
	if response != nil {
		if id := response.Header.Get("X-Request-Id"); id != "" {
			a.RequestID = id
		}
		if id := response.Header.Get("X-Mgc-Trace-Id"); id != "" {
			a.TraceID = id
		}
	}
	return &a
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
		defaultTransport.MaxIdleConns = 1000   //500
		defaultTransport.MaxConnsPerHost = 500 //200
		defaultTransport.IdleConnTimeout = 30 * time.Second
	}
	return defaultTransport
}
