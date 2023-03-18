package hit

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Result is a request's result.
type Result struct {
	RPS      float64       // RPS is the requests per second
	Requests int           // Requests is the number of requests made
	Errors   int           // Errors is the number of errors occurred
	Bytes    int64         // Bytes is the number of bytes downloaded
	Duration time.Duration // Duration is a single or all requests' duration
	Fastest  time.Duration // Fastest request result duration among others
	Slowest  time.Duration // Slowest request result duration among others
	Status   uint16        // Status is a request's HTTP status code
	Error    error         // Error is not nil if the request is failed
}

// Merge this Result with another.
func (r *Result) Merge(results ...*Result) {
	r.Requests += len(results)
	for i := range results {
		r.Bytes += results[i].Bytes
		if r.Fastest == 0 || results[i].Duration < r.Fastest {
			r.Fastest = results[i].Duration
		}
		if results[i].Duration > r.Slowest {
			r.Slowest = results[i].Duration
		}
		switch {
		case results[i].Error != nil:
			fallthrough
		case results[i].Status >= http.StatusBadRequest:
			r.Errors++
		}
	}
}

// Finalize the total duration and calculate RPS.
func (r *Result) Finalize(total time.Duration) *Result {
	r.Duration = total
	r.RPS = float64(r.Requests) / total.Seconds()
	return r
}

// Fprint the result to an io.Writer.
func (r *Result) Fprint(out io.Writer) {
	p := func(format string, args ...any) {
		fmt.Fprintf(out, format, args...)
	}
	p("\nSummary:\n")
	p("\tSuccess : %.0f%%\n", r.success())
	p("\tRPS : %.1f\n", r.RPS)
	p("\tRequests : %d\n", r.Requests)
	p("\tErrors : %d\n", r.Errors)
	p("\tBytes : %d\n", r.Bytes)
	p("\tDuration : %s\n", round(r.Duration))
	if r.Requests > 1 {
		p("\tFastest : %s\n", round(r.Fastest))
		p("\tSlowest : %s\n", round(r.Slowest))
	}
}
func (r *Result) success() float64 {
	rr, e := float64(r.Requests), float64(r.Errors)
	return (rr - e) / rr * 100
}
func round(t time.Duration) time.Duration {
	return t.Round(time.Microsecond)
}

func (r *Result) DisplayErrors(results ...*Result) {

	for i := range results {
		if results[i].Error != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", results[i].Error)
		}
	}
}
