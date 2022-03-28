package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Task struct {
	ID        int    `json:"id,omitempty"`
	UserID    int    `json:"user_id,omitempty"`
	Title     string `json:"title,omitempty"`
	Completed bool   `json:"completed,omitempty"`
}

var timeout int

func init() {
	flag.IntVar(&timeout, "timeout", 5, "Timeout in seconds (default is 5).")

}

func (t *Task) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&t)
}

func main() {

	timer := StartTimer("main")
	flag.Parse()
	defer timer()
	var t Task
	// wg := &sync.WaitGroup{}
	// wg.Add(100)

	var semMaxWeight int64 = 10
	var semAcquisitionWeight int64 = 5

	sem := semaphore.NewWeighted(semMaxWeight)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	errGroup, ctx := errgroup.WithContext(ctx)
	for i := 0; i < 100; i++ {
		fmt.Println(runtime.NumGoroutine())
		if err := sem.Acquire(ctx, semAcquisitionWeight); err != nil {
			log.Fatal(err)
		}
		i := i
		errGroup.Go(func() error {
			// defer wg.Done()
			defer sem.Release(semAcquisitionWeight)
			res, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", i))
			if err != nil {
				// log.Printf("error request: %v\n", err)
				return fmt.Errorf("request: %w", err)
			}
			defer res.Body.Close()
			// if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			// 	return err
			// }
			if err := t.Decode(res.Body); err != nil {
				return fmt.Errorf("Decode json: %w", err)
			}
			fmt.Println(t.Title)
			return nil
		})
	}
	if err := errGroup.Wait(); err != nil {
		log.Fatal(err)
	}
}

func StartTimer(name string) func() {
	t := time.Now()
	fmt.Println(name, "started")

	return func() {
		// d := time.Now().Sub(t)
		d := time.Since(t)
		fmt.Println(name, "took", d)
	}

}
