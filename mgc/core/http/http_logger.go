package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	mgcLoggerPkg "github.com/MagaluCloud/magalu/mgc/core/logger"
	"go.uber.org/zap"
)

type ClientLogger struct {
	Transport http.RoundTripper
}

func NewDefaultClientLogger(transport http.RoundTripper) *ClientLogger {
	return &ClientLogger{
		Transport: transport,
	}
}

// 1 or "progressive" = log progressive;
// 2 or "final" = accumulates it all and logs at the end.
const logPayloadEnvVar = "MGC_SDK_LOG_HTTP_PAYLOAD"

type payloadLoggerFn func(
	parent io.Reader,
	message string,
	logger *zap.SugaredLogger,
) io.ReadCloser

var payloadLogger *payloadLoggerFn

func getPayloadLogger() payloadLoggerFn {
	if payloadLogger == nil {
		payloadLogger = new(payloadLoggerFn)
		switch value := strings.ToLower(os.Getenv(logPayloadEnvVar)); value {
		case "", "0":
			break

		case "1", "progressive":
			*payloadLogger = func(parent io.Reader, message string, logger *zap.SugaredLogger) io.ReadCloser {
				return mgcLoggerPkg.NewProgressiveLoggerReader(parent, func(readData mgcLoggerPkg.LogReadData) {
					logger.Infow(message, "body", readData)
				})
			}

		case "2", "final":
			*payloadLogger = func(parent io.Reader, message string, logger *zap.SugaredLogger) io.ReadCloser {
				return mgcLoggerPkg.NewFinalLoggerReader(parent, func(readData mgcLoggerPkg.LogReadData) {
					logger.Infow(message, "body", readData)
				})
			}

		default:
			logger().Warnw(logPayloadEnvVar+": unknown value", "value", value)
		}
	}
	return *payloadLogger
}

func (t *ClientLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	log := logger().With("method", req.Method, "url", req.URL, "protocol", req.Proto)
	t.logRequest(log, req)
	if req.Body != nil {
		if payloadLogger := getPayloadLogger(); payloadLogger != nil {
			newReq := *req
			newReq.Body = payloadLogger(req.Body, "read request body", log)
			if req.GetBody != nil {
				getBody := req.GetBody
				newReq.GetBody = func() (io.ReadCloser, error) {
					r, err := getBody()
					if err != nil {
						return r, err
					}
					return payloadLogger(r, "read (new) request body", log), nil
				}
			}
			req = &newReq
		}
	}

	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	resp, err := transport.RoundTrip(req)
	if resp != nil && resp.Body != nil {
		if payloadLogger := getPayloadLogger(); payloadLogger != nil {
			newResp := *resp
			newResp.Body = payloadLogger(resp.Body, "read response body", log)
			resp = &newResp
		}
	}

	t.logResponse(log, req, resp, err)

	return resp, err
}

func (t *ClientLogger) logRequest(log *zap.SugaredLogger, req *http.Request) {
	log.Infow("request", "headers", LogHttpHeaders(req.Header))
}

func (t *ClientLogger) logResponse(log *zap.SugaredLogger, req *http.Request, resp *http.Response, err error) {
	if resp == nil {
		if err == nil {
			err = errors.New("Unknown Error")
		}
		log.Infow("request error", "error", err)
		return
	}
	log = log.With("headers", LogHttpHeaders(resp.Header))
	respBody, _ := io.ReadAll(resp.Body)
	if err != nil {
		log.Infow("response with error", "status", resp.Status, "error", err)
	} else {
		log.Infow("response", "status", resp.Status, "body", string(respBody))
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
}

type LogRequest http.Request

func (r LogRequest) MarshalJSON() ([]byte, error) {
	var url any
	if r.URL != nil {
		url = r.URL.String()
	}
	return json.Marshal(map[string]any{
		"method":   r.Method,
		"url":      url,
		"protocol": r.Proto,
		"headers":  LogHttpHeaders(r.Header),
	})
}

type LogResponse http.Response

func (r LogResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"status":   r.Status,
		"protocol": r.Proto,
		"headers":  LogHttpHeaders(r.Header),
	})
}
