package main

import (
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

type Result struct {
	// The name of the test file.
	Name string
	Err  error
}

// RunSingleTest will run a single test file. It returns a Result.
func RunSingleTest(file string) Result {
	// Run the test file.
	// If there is an error, return it.
	// If not, return nil.
	randInt := rand.Intn(1000)
	time.Sleep(time.Duration(randInt) * time.Millisecond)
	return Result{
		Name: file,
		Err:  nil,
	}
}

// RunTests will run multiple test files concurrently. It returns a channel of results.
// Results are received in the same order as the specified test files.
func RunTests(files []string) <-chan Result {
	results := make(chan Result)

	var (
		counter int                           // The counter will keep track of the index we are waiting on.
		cond    = sync.NewCond(&sync.Mutex{}) // The condition variable we will synchronize via waits and broadcasts.
	)

	// Create a loop to spawn a goroutine for each test file.
	for i, file := range files {
		// I understand that in Go 1.22, capturing the loop vars will no longer be necessary, but I do it now because I am scared of change.
		go func(i int, file string) {
			// Immediately run the test. All tests are running concurrently.
			result := RunSingleTest(file)

			// Grab the condition lock!!! We will want to check the condition. If we don't hold the lock, then race-city, here we come.
			cond.L.Lock()

			// It's only polite.
			defer cond.L.Unlock()

			// Here is our condition. We want our counter to equal the index of the goroutine.
			// If it's not, we wait.
			//
			// Make sure to wait in a loop, because when the routine is signaled again, we have no guarantees that the condition will be true.
			// This is one of those rare times when 'for' > 'if' when expressing conditions.
			for counter != i {
				cond.Wait() // When we enter wait, we release the lock. When wait exits, we have the lock. Easy-breezy.
			}

			// The index is correct, our time has come! Push the result!
			results <- result

			// This was the last index. Let's close up shop!
			if i == len(files)-1 {
				close(results)
				return
			}

			// Increment the counter; we want the next index to come through!
			// Notice that we are not using atomics. We still have the lock at this point.
			// Counter is not under threat!
			counter++

			// We let all the other goroutines know that they can wake up and fight for the lock again.
			cond.Broadcast()
			// cond.Signal()

			// Remember that defer of the unlock? It allows those goroutines that have been woken up via the broadcast
			// a chance at getting that lock. If not, they will be deadlocked. And that sounds bad.
		}(i, file)
	}

	return results
}

func main() {
	// start server for pprof in background
	go func() {
		log.Fatal(http.ListenAndServe("localhost:6060", nil))
	}()

	// Run some tests.
	results := RunTests([]string{"test8", "test200", "test3", "test4", "test5", "test6", "test7"})

	// Print the results.
	for result := range results {
		println(result.Name)
	}

	// simulate load with 100 iterations
	for i := 0; i < 100; i++ {
		results := RunTests([]string{"test8", "test200", "test3", "test4", "test5", "test6", "test7"})
		for result := range results {
			println(result.Name)
		}
	}

}
