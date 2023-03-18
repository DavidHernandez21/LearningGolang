//go:build cli

package main

import (
	"bytes"
	"flag"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type testEnv struct {
	args   string
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func (e *testEnv) run() error {
	s := flag.NewFlagSet("hit", flag.ContinueOnError)
	s.SetOutput(&e.stderr)
	return run(s, &e.stdout, strings.Fields(e.args)...)
}

func TestRun(t *testing.T) {
	t.Parallel()

	happy := map[string]struct {
		in  string
		out string
	}{
		"url": {
			in:  "http://foo",
			out: "Making 100 requests to http://foo with concurrency level of " + strconv.Itoa(runtime.NumCPU()) + " (timeout 10m0s).",
		},
		"n_c": {
			in:  "-n 20 -c 10 http://foo",
			out: "20 requests to http://foo with concurrency level of 10 (timeout 10m0s).",
		},
	}

	for name, tc := range happy {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := &testEnv{args: tc.in}
			if err := e.run(); err != nil {
				t.Fatalf("got %q;\nwant nil", err)
			}

			if out := e.stdout.String(); !strings.Contains(out, tc.out) {
				t.Fatalf("got:\n%s\nwant %q", out, tc.out)
			}
		})
	}

	sad := map[string]string{
		"url/missing":            "",
		"url/err":                "://foo",
		"url/host":               "http://",
		"url/scheme":             "ftp://",
		"c/err":                  "-c=x http://foo",
		"n/err":                  "-n=x http://foo",
		"c/neg":                  "-c=-1 http://foo",
		"n/neg":                  "-n=-1 http://foo",
		"c/zero":                 "-c=0 http://foo",
		"n/zero":                 "-n=0 http://foo",
		"c/greater":              "-n=1 -c=2 http://foo",
		"t/err":                  "-t=x http://foo",
		"t/neg":                  "-t=-1s http://foo",
		"t/zero":                 "-t=0s http://foo",
		"t/greater_than_15_min":  "-t=16m http://foo",
		"tr/err":                 "-tr=x http://foo",
		"tr/neg":                 "-tr=-1 http://foo",
		"tr/greater_than_10_min": "-tr=11m http://foo",
	}

	for name, tc := range sad {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := &testEnv{args: tc}
			if err := e.run(); err == nil {
				t.Fatal("got nil; want error")
			}
			if e.stderr.Len() == 0 {
				t.Fatal("stderr = 0 bytes; want > 0 bytes")
			}
		})
	}

}
