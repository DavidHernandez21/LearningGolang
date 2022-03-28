package main

import (
	"fmt"
	"io"
	"os"

	"src/github.com/DavidHernandez21/src/github.com/DavidHernandez21/justForfunc/ioPipes/imgcat"

	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "missing paths of images to cat")
		os.Exit(2)
	}

	for _, path := range os.Args[1:] {
		if err := cat(path); err != nil {
			fmt.Fprintf(os.Stderr, "could not cat %s: %v\n", path, err)
		}
	}
}

type badWriter struct{}

func (badWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("bad writer")
}

func cat(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "could not open image")
	}
	defer f.Close()

	wc := imgcat.NewWriter(badWriter{})
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	return wc.Close()
}
