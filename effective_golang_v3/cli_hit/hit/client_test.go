package hit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestClientDo(t *testing.T) {
	t.Parallel()

	const wantHits, wantErrors = 10, 0

	var (
		gotHits atomic.Int32
		server  = newtestServer(t, func(_ http.ResponseWriter, _ *http.Request) {
			gotHits.Add(1)
		})
		request = newRequest(t, http.MethodGet, server.URL)
	)

	client := &Client{
		Concurrency: 1,
	}

	sum := client.Do(context.Background(), request, wantHits)

	if got := gotHits.Load(); got != wantHits {
		t.Errorf("hits=%d; want %d", got, wantHits)
	}

	if got := sum.Requests; got != wantHits {
		t.Errorf("Requests=%d; want %d", got, wantHits)
	}

	if got := sum.Errors; got != wantErrors {
		t.Errorf("Errors=%d; want %d", got, wantErrors)
	}
}

func newtestServer(tb testing.TB, handler http.HandlerFunc) *httptest.Server {
	tb.Helper()

	server := httptest.NewServer(handler)
	tb.Cleanup(server.Close)

	return server
}

func newRequest(tb testing.TB, method, url string) *http.Request {
	tb.Helper()

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		tb.Fatalf("NewRequest(%q, %q) err=%q; want nil", method, url, err)
	}

	return request
}
