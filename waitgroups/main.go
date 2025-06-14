package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

func main() {
	total, max := 10, 3
	var wg sync.WaitGroup
	for i := 0; i < total; i += max {
		limit := max
		if i+max > total {
			limit = total - i
		}

		wg.Add(limit)
		for j := 0; j < limit; j++ {
			go func(j int) {
				defer wg.Done()
				conn, err := net.Dial("tcp", ":8080")
				if err != nil {
					log.Fatalf("could not dial: %v", err)
				}
				bs, err := io.ReadAll(conn)
				if err != nil {
					log.Fatalf("could not read from conn: %v", err)
				}
				// if string(bs) != "success" {
				// 	log.Fatalf("request error, request: %d", i+1+j)
				// }

				if string(bs) == "too many connections" {
					log.Printf("server rejected request %v, too many connections", i+1+j)

				}

				if string(bs) == "success" {
					fmt.Printf("request %d: success\n", i+1+j)
				}

			}(j)
		}
		wg.Wait()
	}
}
