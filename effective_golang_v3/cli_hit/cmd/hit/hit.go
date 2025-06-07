package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/DavidHernandez21/effective_golang_v3/cli_hit/hit"
)

const (
	bannerText = `
__ __ __ ______
/\ \_\ \ /\ \ /\__ _\
\ \ __ \ \ \ \ \/_/\ \/
 \ \_\ \_\ \ \_\ \ \_\
 \/_/\/_/ \/_/ \/_/
`
	usageText = `
Usage:
	hit [options] url
Options:`
)

func banner() string {
	return bannerText
}

// func usage() string {
// 	return usageText
// }

func main() {
	if err := run(flag.CommandLine, os.Stdout, os.Args[1:]...); err != nil {
		fmt.Fprintln(os.Stderr, "error occurred:", err)
		os.Exit(1)
	}

}

func run(s *flag.FlagSet, out io.Writer, args ...string) error {

	f := &flags{
		n: 100,
		c: runtime.NumCPU(),
	}
	if err := f.parse(s, args...); err != nil {
		return err
	}

	fmt.Fprintln(out, banner())
	fmt.Fprintf(out, "Making %d requests to %s with concurrency level of %d (timeout %v).\n", f.n, f.url, f.c, f.t)

	if f.rps > 0 {
		fmt.Fprintf(out, "Throttling requests to %d per second.\n", f.rps)
	}

	request, err := http.NewRequest(http.MethodGet, f.url, http.NoBody)
	if err != nil {
		return err
	}

	c := &hit.Client{
		Concurrency:       f.c,
		RPS:               f.rps,
		Timeout:           f.t,
		TimeoutPerRequest: f.tr,
		ShowErrors:        f.e,
	}

	ctx, cancel := context.WithTimeout(context.Background(), f.t)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	defer stop()
	sum := c.Do(ctx, request, f.n)
	sum.Fprint(out)

	switch ctx.Err() {
	case context.Canceled:
		return errors.New("interrupted")
	case context.DeadlineExceeded:
		return fmt.Errorf("timed out in %s", f.t)
	}

	return nil

}
