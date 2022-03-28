package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

type Item struct {
	ID            int
	Name          string
	PackingEffort time.Duration
}

func PrepareItems(done <-chan struct{}) <-chan Item {
	items := make(chan Item)
	itemsToShip := []Item{
		{0, "Shirt", 1 * time.Second},
		{1, "Legos", 1 * time.Second},
		{2, "TV", 5 * time.Second},
		{3, "Bananas", 2 * time.Second},
		{4, "Hat", 1 * time.Second},
		{5, "Phone", 2 * time.Second},
		{6, "Plates", 3 * time.Second},
		{7, "Computer", 5 * time.Second},
		{8, "Pint Glass", 3 * time.Second},
		{9, "Watch", 2 * time.Second},
	}

	go func() {
		for _, item := range itemsToShip {
			select {
			case <-done:
				return
			case items <- item:
			}
		}
		close(items)
	}()
	return items
}

func PackItems(done <-chan struct{}, items <-chan Item, workerID int) <-chan Item {
	packages := make(chan Item)
	go func() {
		for item := range items {
			select {
			case <-done:
				return
			case packages <- item:
				time.Sleep(item.PackingEffort)
				fmt.Printf("Worker #%d: Shipping package no. %d, took %ds to pack\n", workerID, item.ID, item.PackingEffort/time.Second)
			}
		}
		close(packages)
	}()
	return packages
}

func merge(done <-chan struct{}, channels ...<-chan Item) <-chan Item {
	var wg sync.WaitGroup
	wg.Add(len(channels))
	outgoingPackages := make(chan Item)
	multiplex := func(c <-chan Item) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case outgoingPackages <- i:
			}
		}
	}
	for _, c := range channels {
		go multiplex(c)
	}
	go func() {
		wg.Wait()
		close(outgoingPackages)
	}()
	return outgoingPackages
}

func main() {
	done := make(chan struct{})
	defer close(done)
	start := time.Now()
	items := PrepareItems(done)
	workers := flag.Int("w", 4, "number of workers")
	flag.Parse()
	// workers := 4
	workersSlice := make([]<-chan Item, *workers)
	for i := 0; i < *workers; i++ {
		workersSlice[i] = PackItems(done, items, i)
	}
	numPackages := 0
	for item := range merge(done, workersSlice...) {
		fmt.Printf("Got Item %#v\n", item)
		numPackages++
	}

	fmt.Printf("Took %fs to ship %d packages with %d workers\n", time.Since(start).Seconds(), numPackages, *workers)

}
