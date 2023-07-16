package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/DavidHernandez21/effective_golnag_v3/url"
)

type flags struct {
	url string
	n   int
	c   int
	rps int
	t   time.Duration
	tr  time.Duration
	e   bool
}

func (f *flags) parse(s *flag.FlagSet, args ...string) error {

	s.Usage = func() {
		fmt.Fprintln(s.Output(), usageText[1:])
		s.PrintDefaults()
	}

	s.Var(toNumber(&f.n), "n", "Number of requests to make")
	s.Var(toNumber(&f.c), "c", "Concurrency level")
	s.Var(toNumber(&f.rps), "rps", "Throttle requests per second")
	s.DurationVar(&f.t, "t", 10*time.Minute, "Timeout for the entire process")
	s.DurationVar(&f.tr, "tr", 5*time.Minute, "Timeout for each request")
	s.BoolVar(&f.e, "e", false, "Show errors")

	if err := s.Parse(args); err != nil {
		return err
	}

	if err := f.validate(s); err != nil {

		fmt.Fprintln(s.Output(), err)
		s.Usage()
		return err
	}

	f.url = s.Arg(0)

	return nil

}

// validate post-conditions after parsing the flags.
func (f *flags) validate(s *flag.FlagSet) error {
	if err := validateURL(s.Arg(0)); err != nil {
		return fmt.Errorf("url: %w", err)
	}
	if f.c > f.n {
		return errors.New("-c: concurrency level cannot be greater than number of requests")
	}
	if f.t > 15*time.Minute || f.t <= 0 {
		return errors.New("-t: timeout cannot be zero, negative or greater than 15 minutes")
	}

	if f.tr > 10*time.Minute || f.tr < 0 {
		return errors.New("-t: timeout per request cannot be negative or greater than 10 minutes")
	}

	return nil

}

func validateURL(s string) error {
	u, err := url.Parse(s)
	switch {
	case strings.TrimSpace(s) == "":
		err = errors.New("required")
	case err != nil:
		err = errors.New("parse error")
	case u.Scheme != "http":
		err = errors.New("only supported scheme is http")
	case u.Host == "":
		err = errors.New("missing host")
	}
	return err
}

// number is a natural number.
type number int

// toNumber is a convenience function for converting p to *number.
func toNumber(p *int) *number {
	return (*number)(p)
}
func (n *number) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	switch {
	case err != nil:
		err = errors.New("parse error")
	case v <= 0:
		err = errors.New("should be positive")
	}
	*n = number(v)
	return err
}
func (n *number) String() string {
	return strconv.Itoa(int(*n))
}
