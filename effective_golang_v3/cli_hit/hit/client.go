package hit

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// Client sends HTTP requests and returns an aggregated performance
// result. The fields should not be changed after initializing.
type Client struct {
	Concurrency       int           // concurrency level
	RPS               int           // throttles requests per second
	Timeout           time.Duration // timeout
	TimeoutPerRequest time.Duration // timeout per request
	ShowErrors        bool          // show errors
}

// Option allows to customize the behavior of the Client.
type Option func(*Client)

// WithConcurrency sets the concurrency level.
func WithConcurrency(n int) Option {
	return func(c *Client) {
		c.Concurrency = n
	}
}

// WithRPS sets the requests per second.
func WithRPS(n int) Option {
	return func(c *Client) {
		c.RPS = n
	}
}

// WithTimeout sets the timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.Timeout = d
	}
}

// WithTimeoutPerRequest sets the timeout per request.
func WithTimeoutPerRequest(d time.Duration) Option {
	return func(c *Client) {
		c.TimeoutPerRequest = d
	}
}

// WithShowErrors sets the show errors flag.
func WithShowErrors(b bool) Option {
	return func(c *Client) {
		c.ShowErrors = b
	}
}

// Do sends n GET requests and returns an aggregated result.
// it will use as many goroutines as the number of logical cores in the machine.

// Create a new Client to customize the behavior of the requests.
func Do(ctx context.Context, url string, n int, opts ...Option) (*Result, error) {
	r, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("new http request: %w", err)
	}
	var c Client
	for i := range opts {
		opts[i](&c)
	}

	if c.Timeout > 0 {
		ctx, cancel := context.WithTimeout(ctx, c.Timeout)
		defer cancel()
		return c.Do(ctx, r, n), nil
	}

	return c.Do(ctx, r, n), nil
}

// Do sends n HTTP requests and returns an aggregated result.
func (c *Client) Do(ctx context.Context, r *http.Request, n int) *Result {
	t := time.Now()
	fmt.Printf("request: %s\n", r.URL)
	sum := c.do(ctx, r, n)
	return sum.Finalize(time.Since(t))
}

func (c *Client) do(ctx context.Context, r *http.Request, n int) *Result {

	// this should not happend when using the cli
	if c.Concurrency == 0 {
		c.Concurrency = runtime.NumCPU()
	}

	// ctx, cancel := context.WithTimeout(r.Context(), c.T)
	// defer cancel()
	p := produce(ctx, n, func() *http.Request {
		return r.Clone(ctx)
	})

	if c.RPS > 0 {
		p = throttle(ctx, p, time.Second/time.Duration(c.RPS*c.Concurrency))
	}

	var (
		sum    Result
		client = c.client()
	)
	defer client.CloseIdleConnections()
	results := make([]*Result, n)
	var i int
	for result := range split(p, c.Concurrency, c.send(client)) {
		results[i] = result
		i++
	}
	sum.Merge(results...)
	if c.ShowErrors {
		sum.DisplayErrors(results...)
	}
	return &sum
}

func (c *Client) send(client *http.Client) SendFunc {
	return func(r *http.Request) *Result {
		return Send(client, r)
	}
}

func (c *Client) client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: c.Concurrency,
		},
		Timeout: c.TimeoutPerRequest,
	}
}
