package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

func GetContentType(resp *http.Response) string {
	headerVal := resp.Header.Get("Content-Type")
	return strings.Split(headerVal, ";")[0]
}
