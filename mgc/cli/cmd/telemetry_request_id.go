package cmd

import (
	"net/http"
	"sync"
)

const requestIDHeader = "X-Request-Id"

// requestIDRecorder fica no meio das chamadas HTTP e vai anotando o X-Request-Id
// de cada uma. No fim do comando, sobra só o da última.
type requestIDRecorder struct {
	next http.RoundTripper

	mu   sync.Mutex
	last string
}

var _ http.RoundTripper = (*requestIDRecorder)(nil)

func newRequestIDRecorder(next http.RoundTripper) *requestIDRecorder {
	if next == nil {
		next = http.DefaultTransport
	}
	return &requestIDRecorder{next: next}
}

func (r *requestIDRecorder) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := r.next.RoundTrip(req)

	id := ""
	if resp != nil {
		id = resp.Header.Get(requestIDHeader)
	}
	if id == "" {
		// Quando a resposta não traz o id (ou nem chega), ainda dá para achar na
		// requisição, porque o transport escreve o cabeçalho nela.
		id = req.Header.Get(requestIDHeader)
	}
	if id != "" {
		r.mu.Lock()
		r.last = id
		r.mu.Unlock()
	}

	return resp, err
}

func (r *requestIDRecorder) LastRequestID() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.last
}
