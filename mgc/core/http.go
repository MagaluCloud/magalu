package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var httpClientKey contextKey = "magalu.cloud/core/Transport"

type HttpClient struct {
	http.Client
}

func NewHttpClient(transport *http.Transport) *HttpClient {
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
		return fmt.Errorf("Error decoding JSON response body: %v", err)
	}
	return nil
}

func DecodeOctet(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	// TODO: read the bytes, parse content disposition and save to file
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error decoding octet stream: %v", err)
	}
	return string(b), nil
}

func GetContentType(resp *http.Response) string {
	headerVal := resp.Header.Get("Content-Type")
	return strings.Split(headerVal, ";")[0]
}
