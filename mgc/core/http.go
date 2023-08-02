package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
)

var httpClientKey contextKey = "magalu.cloud/core/Transport"

type HttpClient struct {
	http.Client
}

func NewHttpClient(transport http.RoundTripper) *HttpClient {
	return &HttpClient{http.Client{Transport: transport}}
}

func NewHttpClientContext(parent context.Context, client *HttpClient) context.Context {
	return context.WithValue(parent, httpClientKey, client)
}

func HttpClientFromContext(context context.Context) *HttpClient {
	client, ok := context.Value(httpClientKey).(*HttpClient)
	if !ok {
		log.Printf("Error casting ctx %s to *HttpClient", httpClientKey)
		return nil
	}
	return client
}

func DecodeJSON(resp *http.Response, data any) error {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return fmt.Errorf("Error decoding JSON response body: %s", err)
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

func (e *HttpError) Error() string {
	return e.Message
}

func (e *HttpError) String() string {
	return fmt.Sprintf("%T{Status: %q, Slug: %q, Message: %q}", e, e.Status, e.Slug, e.Message)
}

func NewHttpErrorFromResponse(resp *http.Response) error {
	slug := "unknown"
	message := resp.Status

	defer resp.Body.Close()
	payload, _ := io.ReadAll(resp.Body)

	contentType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		log.Printf("Ignored invalid response Content-Type %q: %s\n", resp.Header.Get("Content-Type"), err.Error())
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

	return &HttpError{
		Code:    resp.StatusCode,
		Status:  resp.Status,
		Headers: resp.Header,
		Payload: payload,
		Message: message,
		Slug:    slug,
	}
}
