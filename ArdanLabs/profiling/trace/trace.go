// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program that performs a series of I/O related tasks to
// better understand tracing in Go.
package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	memProfile = "mem.prof"
)

type (
	item struct {
		XMLName     xml.Name `xml:"item"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
	}

	channel struct {
		XMLName xml.Name `xml:"channel"`
		Items   []item   `xml:"item"`
	}

	document struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}
)

func joinBuilder(sb *strings.Builder, words ...string) (string, error) {

	sb.Reset()
	// sb.Grow(grow)

	for _, word := range words {
		_, err := sb.WriteString(word)

		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

func main() {
	// pprof.StartCPUProfile(os.Stdout)
	// defer pprof.StopCPUProfile()

	// f1, err := os.Create("trace.out")
	// if err != nil {
	// 	log.Fatalf("failed to create trace output file: %v", err)
	// }
	// defer func() {
	// 	if err := f1.Close(); err != nil {
	// 		log.Fatalf("failed to close trace file: %v", err)
	// 	}
	// }()

	// if err := trace.Start(f1); err != nil {
	// 	log.Fatalf("failed to start trace: %v", err)
	// }
	// defer trace.Stop()

	myTimer := StartTimer("freqProcessors")
	defer myTimer()

	f1, err := os.Create(memProfile)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f1.Close() // error handling omitted for example
	// runtime.MemProfile()
	// runtime.MemProfileRate = 1

	docs := make([]string, 4000)
	sb := strings.Builder{}
	sb.Grow(15)
	for i := range docs {
		doc, err := joinBuilder(&sb, "newsfeed-", strconv.Itoa(i), ".xml")
		if err != nil {
			log.Fatalf("failed to join: %v", err)
		}

		docs[i] = doc
		// docs[i] = fmt.Sprintf("newsfeed-%.4d.xml", i)
	}

	topic := "president"
	// n := freq(topic, docs)
	// n := freqConcurrent(topic, docs)
	// n := freqConcurrentSem(topic, docs)
	n := freqProcessors(topic, docs)
	// n := freqProcessorsTasks(topic, docs)
	// n := freqActor(topic, docs)

	log.Printf("Searching %d files, found %s %d times.", len(docs), topic, n)

	if err := pprof.Lookup("heap").WriteTo(f1, 0); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

}

func freq(topic string, docs []string) int {
	var found int

	for _, doc := range docs {
		file := fmt.Sprintf("%s.xml", doc[:8])
		f, err := os.OpenFile(file, os.O_RDONLY, 0)
		if err != nil {
			log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
			return 0
		}

		data, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
			return 0
		}

		var d document
		if err := xml.Unmarshal(data, &d); err != nil {
			log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
			return 0
		}

		for _, item := range d.Channel.Items {
			if strings.Contains(item.Title, topic) {
				found++
				continue
			}

			if strings.Contains(item.Description, topic) {
				found++
			}
		}
	}

	return found
}

func freqConcurrent(topic string, docs []string) int {
	var found int32

	g := len(docs)
	var wg sync.WaitGroup
	wg.Add(g)

	for _, doc := range docs {
		go func(doc string) {
			var lFound int32
			defer func() {
				atomic.AddInt32(&found, lFound)
				wg.Done()
			}()

			file := fmt.Sprintf("%s.xml", doc[:8])
			f, err := os.OpenFile(file, os.O_RDONLY, 0)
			if err != nil {
				log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
				return
			}

			data, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
				return
			}

			var d document
			if err := xml.Unmarshal(data, &d); err != nil {
				log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
				return
			}

			for _, item := range d.Channel.Items {
				if strings.Contains(item.Title, topic) {
					lFound++
					continue
				}

				if strings.Contains(item.Description, topic) {
					lFound++
				}
			}
		}(doc)
	}

	wg.Wait()
	return int(found)
}

func freqConcurrentSem(topic string, docs []string) int {
	var found int32

	g := len(docs)
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan bool, runtime.GOMAXPROCS(0))

	for _, doc := range docs {
		go func(doc string) {
			ch <- true
			{
				var lFound int32
				defer func() {
					atomic.AddInt32(&found, lFound)
					wg.Done()
				}()
				file := fmt.Sprintf("%s.xml", doc[:8])
				f, err := os.OpenFile(file, os.O_RDONLY, 0)
				if err != nil {
					log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
					return
				}

				data, err := io.ReadAll(f)
				f.Close()
				if err != nil {
					log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
					return
				}

				var d document
				if err := xml.Unmarshal(data, &d); err != nil {
					log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
					return
				}

				for _, item := range d.Channel.Items {
					if strings.Contains(item.Title, topic) {
						lFound++
						continue
					}

					if strings.Contains(item.Description, topic) {
						lFound++
					}
				}
			}
			<-ch
		}(doc)
	}

	wg.Wait()
	return int(found)
}

func freqProcessors(topic string, docs []string) int {
	var found int32

	g := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan string, g)

	for i := 0; i < g; i++ {
		go func() {
			var lFound int32
			defer func() {
				atomic.AddInt32(&found, lFound)
				wg.Done()
			}()

			// buffer := bytes.Buffer{}
			// buffer.Grow(26000)

			sb := strings.Builder{}
			sb.Grow(10)
			for doc := range ch {
				// file := fmt.Sprintf("%s.xml", doc[:8])
				file, err := joinBuilder(&sb, doc[:8], ".xml")
				// log.Printf("file: %s", file)
				if err != nil {
					log.Fatalf("Joining Builder : ERROR : %v", err)
				}
				f, err := os.OpenFile(file, os.O_RDONLY, 0)
				if err != nil {
					log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
					return
				}

				// _, err = buffer.ReadFrom(f)
				// if err != nil {
				// 	log.Fatalf("Buffer Reading Document [%s] : ERROR : %v", file, err)
				// }
				data, err := io.ReadAll(f)
				if err != nil {
					log.Printf("io Reading Document [%s] : ERROR : %v", doc, err)
					return
				}
				f.Close()
				// buffer.Reset()

				var d document
				if err := xml.Unmarshal(data, &d); err != nil {
					log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
					return
				}

				for _, item := range d.Channel.Items {
					if strings.Contains(item.Title, topic) {
						lFound++
						continue
					}

					if strings.Contains(item.Description, topic) {
						lFound++
					}
				}
			}
		}()
	}

	for _, doc := range docs {
		ch <- doc
	}
	close(ch)

	wg.Wait()
	return int(found)
}

func freqProcessorsTasks(topic string, docs []string) int {
	var found int32

	g := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan string, g)

	for i := 0; i < g; i++ {
		go func() {
			var lFound int32
			defer func() {
				atomic.AddInt32(&found, lFound)
				wg.Done()
			}()

			for doc := range ch {
				func() {
					file := fmt.Sprintf("%s.xml", doc[:8])
					ctx, tt := trace.NewTask(context.Background(), doc)
					defer tt.End()

					reg := trace.StartRegion(ctx, "OpenFile")
					f, err := os.OpenFile(file, os.O_RDONLY, 0)
					if err != nil {
						log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
						return
					}
					reg.End()

					reg = trace.StartRegion(ctx, "ReadAll")
					data, err := io.ReadAll(f)
					f.Close()
					if err != nil {
						log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
						return
					}
					reg.End()

					reg = trace.StartRegion(ctx, "Unmarshal")
					var d document
					if err := xml.Unmarshal(data, &d); err != nil {
						log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
						return
					}
					reg.End()

					reg = trace.StartRegion(ctx, "Search")
					for _, item := range d.Channel.Items {
						if strings.Contains(item.Title, topic) {
							lFound++
							continue
						}

						if strings.Contains(item.Description, topic) {
							lFound++
						}
					}
					reg.End()
				}()
			}
		}()
	}

	for _, doc := range docs {
		ch <- doc
	}
	close(ch)

	wg.Wait()
	return int(found)
}

func freqActor(topic string, docs []string) int {
	files := make(chan *os.File, 100)
	go func() {
		for _, doc := range docs {
			file := fmt.Sprintf("%s.xml", doc[:8])
			f, err := os.OpenFile(file, os.O_RDONLY, 0)
			if err != nil {
				log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
				break
			}
			files <- f
		}
		close(files)
	}()

	data := make(chan []byte, 100)
	go func() {
		for f := range files {
			d, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				log.Printf("Reading Document [%s] : ERROR : %v", f.Name(), err)
				break
			}
			data <- d
		}
		close(data)
	}()

	rss := make(chan document, 100)
	go func() {
		for dt := range data {
			var d document
			if err := xml.Unmarshal(dt, &d); err != nil {
				log.Printf("Decoding Document : ERROR : %v", err)
				break
			}
			rss <- d
		}
		close(rss)
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	var found int
	go func() {
		for d := range rss {
			for _, item := range d.Channel.Items {
				if strings.Contains(item.Title, topic) {
					found++
					continue
				}

				if strings.Contains(item.Description, topic) {
					found++
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()
	return found
}

func StartTimer(name string) func() {
	t := time.Now()
	log.Println(name, "started")

	return func() {
		// d := time.Now().Sub(t)
		d := time.Since(t)
		log.Println(name, "took", d)
	}

}
