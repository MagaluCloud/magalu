package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"slices"
	"strconv"
	"sync"
	"time"
)

const (
	logServiceName = "mgc-cli"

	maxDrainedResponseBytes = 64 << 10 // 64 KiB
)

type PostHogExporter struct {
	endpoint string
	apiKey   string
	client   *http.Client
	debug    DebugLogger
	// warmed fecha quando a conexão aberta pelo Warm fica pronta ou falha
	warmed chan struct{}
}

// O keep-alive fica ligado para o envio reaproveitar a conexão aberta pelo Warm
func NewPostHogExporter(endpoint, apiKey string) *PostHogExporter {
	return &PostHogExporter{
		endpoint: endpoint,
		apiKey:   apiKey,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy:             http.ProxyFromEnvironment,
				ForceAttemptHTTP2: true,
			},
		},
		debug: func(string, ...any) {},
	}
}

// WithDebug liga o log de cada chamada ao PostHog, sem dados do evento nem o token
func (p *PostHogExporter) WithDebug(debug DebugLogger) *PostHogExporter {
	if debug != nil {
		p.debug = debug
	}
	return p
}

// Warm faz um HEAD sem token em segundo plano só para abrir a conexão.
func (p *PostHogExporter) Warm(ctx context.Context) {
	if p.warmed != nil {
		return
	}
	warmed := make(chan struct{})
	p.warmed = warmed

	var once sync.Once
	markWarmed := func() { once.Do(func() { close(warmed) }) }

	go func() {
		defer markWarmed()
		start := time.Now()
		var connectedIn time.Duration
		trace := &httptrace.ClientTrace{GotConn: func(httptrace.GotConnInfo) {
			connectedIn = time.Since(start)
			markWarmed()
		}}
		req, err := http.NewRequestWithContext(httptrace.WithClientTrace(ctx, trace), http.MethodHead, p.endpoint, nil)
		if err != nil {
			return
		}
		resp, err := p.client.Do(req)
		if err == nil {
			resp.Body.Close()
		}
		p.logCall(http.MethodHead, start, resp, err, "connect_ms", connectedIn.Milliseconds())
	}()
}

type ExportStatusError struct {
	StatusCode int
}

func (e ExportStatusError) Error() string {
	return fmt.Sprintf("unexpected status %d", e.StatusCode)
}

func (p *PostHogExporter) Export(ctx context.Context, event Event) error {
	start := time.Now()
	body, err := newOTLPLogs(event)
	if err != nil {
		return err
	}

	if p.warmed != nil {
		select {
		case <-p.warmed:
		case <-ctx.Done():
			p.logCall(http.MethodPost, start, nil, ctx.Err(), "warm_wait_ms", time.Since(start).Milliseconds())
			return ctx.Err()
		}
	}
	warmWait := time.Since(start)

	var reused bool
	trace := &httptrace.ClientTrace{GotConn: func(info httptrace.GotConnInfo) { reused = info.Reused }}
	ctx = httptrace.WithClientTrace(ctx, trace)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		p.logCall(http.MethodPost, start, nil, err, "warm_wait_ms", warmWait.Milliseconds(), "reused_conn", reused)
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxDrainedResponseBytes))
	p.logCall(http.MethodPost, start, resp, nil, "warm_wait_ms", warmWait.Milliseconds(), "reused_conn", reused)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return ExportStatusError{StatusCode: resp.StatusCode}
	}
	return nil
}

func (p *PostHogExporter) logCall(method string, start time.Time, resp *http.Response, err error, extra ...any) {
	kv := []any{"method", method, "duration_ms", time.Since(start).Milliseconds()}
	if resp != nil {
		kv = append(kv, "status", resp.StatusCode)
	}
	if method == http.MethodPost {
		kv = append(kv, "success", err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode <= 299)
	}
	if err != nil {
		kv = append(kv, "error", err.Error())
	}
	p.debug("telemetry: posthog call", append(kv, extra...)...)
}

type otlpLogs struct {
	ResourceLogs []otlpResourceLogs `json:"resourceLogs"`
}

type otlpResourceLogs struct {
	Resource  otlpResource    `json:"resource"`
	ScopeLogs []otlpScopeLogs `json:"scopeLogs"`
}

type otlpResource struct {
	Attributes []otlpKeyValue `json:"attributes"`
}

type otlpScopeLogs struct {
	LogRecords []otlpLogRecord `json:"logRecords"`
}

type otlpLogRecord struct {
	TimeUnixNano string         `json:"timeUnixNano"`
	SeverityText string         `json:"severityText"`
	Body         otlpAnyValue   `json:"body"`
	Attributes   []otlpKeyValue `json:"attributes"`
}

type otlpKeyValue struct {
	Key   string       `json:"key"`
	Value otlpAnyValue `json:"value"`
}

// otlpAnyValue segue o OTLP JSON, em que inteiros de 64 bits vão como string
type otlpAnyValue struct {
	StringValue *string `json:"stringValue,omitempty"`
	IntValue    *string `json:"intValue,omitempty"`
}

func stringValue(s string) otlpAnyValue {
	return otlpAnyValue{StringValue: &s}
}

func newOTLPLogs(event Event) ([]byte, error) {
	attributes, err := logAttributes(event)
	if err != nil {
		return nil, err
	}

	record := otlpLogRecord{
		TimeUnixNano: strconv.FormatInt(event.Timestamp.Truncate(time.Millisecond).UnixNano(), 10),
		SeverityText: string(event.LogLevel()),
		Body:         stringValue(event.Message()),
		Attributes:   attributes,
	}

	return json.Marshal(otlpLogs{ResourceLogs: []otlpResourceLogs{{
		Resource: otlpResource{Attributes: []otlpKeyValue{
			{Key: "service.name", Value: stringValue(logServiceName)},
		}},
		ScopeLogs: []otlpScopeLogs{{LogRecords: []otlpLogRecord{record}}},
	}}})
}

// logAttributes parte do JSON do modelo, assim os nomes e as omissões são os mesmos
// do evento, e só achata os campos aninhados com "."
func logAttributes(event Event) ([]otlpKeyValue, error) {
	raw, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var fields map[string]any
	if err := decoder.Decode(&fields); err != nil {
		return nil, err
	}

	attributes := []otlpKeyValue{{Key: "distinct_id", Value: stringValue(event.InstallationID)}}
	return appendAttributes(attributes, "", fields)
}

func appendAttributes(attributes []otlpKeyValue, prefix string, fields map[string]any) ([]otlpKeyValue, error) {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	var err error
	for _, key := range keys {
		name := prefix + key
		switch value := fields[key].(type) {
		case string:
			attributes = append(attributes, otlpKeyValue{Key: name, Value: stringValue(value)})
		case json.Number:
			n := value.String()
			attributes = append(attributes, otlpKeyValue{Key: name, Value: otlpAnyValue{IntValue: &n}})
		case map[string]any:
			if attributes, err = appendAttributes(attributes, name+".", value); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported telemetry attribute %q of type %T", name, value)
		}
	}
	return attributes, nil
}
