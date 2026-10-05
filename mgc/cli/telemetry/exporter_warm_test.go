package telemetry

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type warmRequest struct {
	method        string
	proto         string
	authorization string
}

// newHTTP2Server responde ao HEAD só depois de headDelay, para o envio acontecer com
// o HEAD ainda em andamento, e conta as conexões abertas pelo cliente
func newHTTP2Server(t *testing.T, headDelay time.Duration) (*PostHogExporter, *atomic.Int32, func() []warmRequest) {
	t.Helper()
	var (
		mu       sync.Mutex
		requests []warmRequest
		conns    atomic.Int32
	)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests = append(requests, warmRequest{r.Method, r.Proto, r.Header.Get("Authorization")})
		mu.Unlock()
		if r.Method == http.MethodHead {
			time.Sleep(headDelay)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	server.EnableHTTP2 = true
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			conns.Add(1)
		}
	}
	server.StartTLS()
	t.Cleanup(server.Close)

	exporter := NewPostHogExporter(server.URL+"/i/v1/logs", "phc_test")
	exporter.client.Transport.(*http.Transport).TLSClientConfig = server.Client().Transport.(*http.Transport).TLSClientConfig.Clone()

	return exporter, &conns, func() []warmRequest {
		mu.Lock()
		defer mu.Unlock()
		return append([]warmRequest(nil), requests...)
	}
}

func TestPostHogExporterWarm(t *testing.T) {
	testCases := []struct {
		name      string
		warm      bool
		headDelay time.Duration
		wantConns int32
	}{
		{"export reuses the warm connection while the HEAD is pending", true, 300 * time.Millisecond, 1},
		{"export reuses the warm connection after the HEAD", true, 0, 1},
		{"export without warm opens its own connection", false, 0, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exporter, conns, requests := newHTTP2Server(t, tc.headDelay)

			if tc.warm {
				exporter.Warm(context.Background())
			}
			ctx, cancel := context.WithTimeout(context.Background(), ExportTimeout)
			defer cancel()
			if err := exporter.Export(ctx, exportedEvent()); err != nil {
				t.Fatalf("Export: %v", err)
			}

			if got := conns.Load(); got != tc.wantConns {
				t.Errorf("%d connections, want %d", got, tc.wantConns)
			}
			for _, r := range requests() {
				if r.proto != "HTTP/2.0" {
					t.Errorf("%s used %s, want HTTP/2.0", r.method, r.proto)
				}
				if r.method == http.MethodHead && r.authorization != "" {
					t.Error("the warm-up request must not carry the project token")
				}
			}
		})
	}
}

func TestPostHogExporterWarmStalled(t *testing.T) {
	// aceita a conexão e nunca responde, deixando o aquecimento preso no handshake
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { listener.Close() })
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			t.Cleanup(func() { conn.Close() })
		}
	}()

	exporter := NewPostHogExporter("https://"+listener.Addr().String()+"/i/v1/logs", "phc_test")
	exporter.Warm(context.Background())

	budget := 200 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), budget)
	defer cancel()

	start := time.Now()
	err = exporter.Export(ctx, exportedEvent())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want context.DeadlineExceeded", err)
	}
	if elapsed := time.Since(start); elapsed > budget+100*time.Millisecond {
		t.Errorf("export took %v, want about %v", elapsed, budget)
	}
}

type warmerFunc func(ctx context.Context)

func (f warmerFunc) Export(context.Context, Event) error { return nil }
func (f warmerFunc) Warm(ctx context.Context)            { f(ctx) }

func TestMultiExporterWarm(t *testing.T) {
	var warmed atomic.Int32
	warmer := warmerFunc(func(context.Context) { warmed.Add(1) })

	MultiExporter{warmer, NoopExporter{}, warmer}.Warm(context.Background())

	if got := warmed.Load(); got != 2 {
		t.Errorf("%d warm-ups, want 2", got)
	}
}

func TestServiceWarm(t *testing.T) {
	testCases := []struct {
		name       string
		env        map[string]string
		terminal   bool
		state      State
		wantWarmed int32
	}{
		{"enabled", nil, false, State{}, 1},
		{"disabled by DO_NOT_TRACK", map[string]string{EnvDoNotTrack: "1"}, false, State{}, 0},
		{"interactive with the notice pending", nil, true, State{}, 0},
		{"interactive after the notice", nil, true, State{NoticeShown: true}, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var warmed atomic.Int32
			svc := New(Options{
				Exporter:   warmerFunc(func(context.Context) { warmed.Add(1) }),
				Getenv:     envFrom(tc.env),
				IsTerminal: tc.terminal,
				Store:      &fakeStore{state: tc.state},
			})

			svc.Warm(context.Background())

			if got := warmed.Load(); got != tc.wantWarmed {
				t.Errorf("%d warm-ups, want %d", got, tc.wantWarmed)
			}
		})
	}

	t.Run("exporter without warm-up", func(t *testing.T) {
		New(Options{Exporter: NoopExporter{}}).Warm(context.Background())
	})
}
