package core

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HttpClientLogger struct {
	Transport         http.RoundTripper
	LogRequest        func(logger *HttpClientLogger, req *http.Request)
	LogResponse       func(logger *HttpClientLogger, req *http.Request, resp *http.Response, err error)
	LogRequestPrefix  string
	LogResponsePrefix string
	RedactHeader      func(logger *HttpClientLogger, req *http.Request, key string, value string) string
}

func NewDefaultHttpClientLogger(transport http.RoundTripper) *HttpClientLogger {
	return &HttpClientLogger{
		Transport:         transport,
		LogRequest:        DefaultHttpClientLogRequest,
		LogResponse:       DefaultHttpClientLogResponse,
		LogRequestPrefix:  "HTTP >",
		LogResponsePrefix: "HTTP <",
		RedactHeader:      DefaultHttpClientRedactHeader,
	}
}

func (t *HttpClientLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.LogRequest != nil {
		t.LogRequest(t, req)
	}

	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	resp, err := transport.RoundTrip(req)

	if t.LogResponse != nil {
		t.LogResponse(t, req, resp, err)
	}

	return resp, err
}

func DefaultHttpClientLogRequest(logger *HttpClientLogger, req *http.Request) {
	log.Printf("%s %s %s %s\n", logger.LogRequestPrefix, req.Method, req.URL, req.Proto)
	for k, lst := range req.Header {
		for _, v := range lst {
			if logger.RedactHeader != nil {
				v = logger.RedactHeader(logger, req, k, v)
			}
			log.Printf("%s %s: %s\n", logger.LogRequestPrefix, k, v)
		}
	}
	log.Printf("%s\n", logger.LogRequestPrefix)
}

func DefaultHttpClientLogResponse(logger *HttpClientLogger, req *http.Request, resp *http.Response, err error) {
	if resp == nil {
		if err == nil {
			err = errors.New("Unknown Error")
		}
		log.Printf("%s %s %s error=%s\n", logger.LogResponsePrefix, req.Method, req.URL, err.Error())
		return
	}
	errMsg := ""
	if err != nil {
		errMsg = fmt.Sprintf("; error: %s", err.Error())
	}
	log.Printf("%s %s%s\n", logger.LogResponsePrefix, resp.Status, errMsg)
	for k, lst := range resp.Header {
		for _, v := range lst {
			if logger.RedactHeader != nil {
				v = logger.RedactHeader(logger, req, k, v)
			}
			log.Printf("%s %s: %s\n", logger.LogResponsePrefix, k, v)
		}
	}
}

func DefaultHttpClientRedactHeader(logger *HttpClientLogger, req *http.Request, key string, value string) string {
	if strings.ToLower(key) == "authorization" {
		parts := strings.SplitN(value, " ", 2)
		if len(parts) > 1 {
			if strings.ToLower(parts[0]) == "bearer" {
				parts[1] = fmt.Sprintf("<redacted %d chars>", len(parts[1]))
			}
		}
		return strings.Join(parts, " ")
	}
	return value
}
