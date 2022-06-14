package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/fatih/semgroup"
)

func main() {
	const maxWorkers = 8
	s := semgroup.NewGroup(context.Background(), maxWorkers)

	visitors := [9]int{1, 1, 1, 1, 2, 2, 1, 1, 2}

	for _, v := range visitors {
		v := v

		s.Go(func() error {
			if v != 1 {
				return errors.New("only one visitor is allowed")
			}
			return nil
		})
	}

	// Wait for all visits to complete. Any errors are accumulated.
	if err := s.Wait(); err != nil {
		fmt.Println(err)
	}
}
