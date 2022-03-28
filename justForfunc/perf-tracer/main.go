package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

const (
	output     = "out.png"
	width      = 2048
	height     = 2048
	cpuProfile = "cpu.prof"
)

var numWorkers int = runtime.NumCPU()

func main() {
	// uncomment these lines to generate traces into stdout.
	// trace.Start(os.Stdout)
	// defer trace.Stop()
	timer := StartTimer("main")

	defer timer()

	// log.Println("num of cpu processors", numWorkers)

	// f1, err := os.Create(cpuProfile)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pprof.StartCPUProfile(f1)
	// defer closeFile(f1)
	// defer pprof.StopCPUProfile()
	f1, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f1.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f1); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	f, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}

	defer closeFile(f)

	// img := createCol(width, height)
	img := createCol(width, height)

	if err = png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Fatalf("error closing file '%v': %v", file.Name(), err)
	}
}

// createSeq fills one pixel at a time.
func createSeq(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			m.Set(i, j, pixel(i, j, width, height))
		}
	}
	return m
}

// createPixel creates one goroutine per pixel.
func createPixel(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))
	var w sync.WaitGroup
	w.Add(width * height)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			go func(i, j int) {
				m.Set(i, j, pixel(i, j, width, height))
				w.Done()
			}(i, j)
		}
	}
	w.Wait()
	return m
}

// createCol creates one goroutine per column.
func createCol(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))
	var w sync.WaitGroup
	w.Add(width)
	for i := 0; i < width; i++ {
		go func(i int) {
			for j := 0; j < height; j++ {
				m.Set(i, j, pixel(i, j, width, height))
			}
			w.Done()
		}(i)
	}
	w.Wait()
	return m
}

// createWorkers creates numWorkers workers and uses a channel to pass each pixel.
func createWorkers(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))

	type px struct{ x, y int }
	c := make(chan px)

	var w sync.WaitGroup
	w.Add(numWorkers)
	for n := 0; n < numWorkers; n++ {
		go func() {
			for px := range c {
				m.Set(px.x, px.y, pixel(px.x, px.y, width, height))
			}
			w.Done()
		}()
	}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			c <- px{i, j}
		}
	}
	close(c)
	w.Wait()
	return m
}

// createWorkersBuffered creates numWorkers workers and uses a buffered channel to pass each pixel.
func createWorkersBuffered(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))

	type px struct{ x, y int }
	c := make(chan px, width*height)

	var w sync.WaitGroup
	w.Add(numWorkers)
	for n := 0; n < numWorkers; n++ {
		go func() {
			for px := range c {
				m.Set(px.x, px.y, pixel(px.x, px.y, width, height))
			}
			w.Done()
		}()
	}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			c <- px{i, j}
		}
	}
	close(c)
	w.Wait()
	return m
}

// createColWorkers creates numWorkers workers and uses a channel to pass each column.
func createColWorkers(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))

	c := make(chan int)

	var w sync.WaitGroup
	w.Add(numWorkers)
	for n := 0; n < numWorkers; n++ {
		go func() {
			for i := range c {
				for j := 0; j < height; j++ {
					m.Set(i, j, pixel(i, j, width, height))
				}
			}
			w.Done()
		}()
	}

	for i := 0; i < width; i++ {
		c <- i
	}

	close(c)
	w.Wait()
	return m
}

// createColWorkersBuffered creates numWorkers workers and uses a buffered channel to pass each column.
func createColWorkersBuffered(width, height int) image.Image {
	m := image.NewGray(image.Rect(0, 0, width, height))

	c := make(chan int, width)

	var w sync.WaitGroup
	w.Add(numWorkers)
	for n := 0; n < numWorkers; n++ {
		go func() {
			for i := range c {
				for j := 0; j < height; j++ {
					m.Set(i, j, pixel(i, j, width, height))
				}
			}
			w.Done()
		}()
	}

	for i := 0; i < width; i++ {
		c <- i
	}

	close(c)
	w.Wait()
	return m
}

// pixel returns the color of a Mandelbrot fractal at the given point.
func pixel(i, j, width, height int) color.Color {
	// Play with this constant to increase the complexity of the fractal.
	// In the justforfunc.com video this was set to 4.
	const complexity = 1024

	xi := norm(i, width, -1.0, 2)
	yi := norm(j, height, -1, 1)

	const maxI = 1000
	x, y := 0., 0.

	for i := 0; (x*x+y*y < complexity) && i < maxI; i++ {
		x, y = x*x-y*y+xi, 2*x*y+yi
	}

	return color.Gray{uint8(x)}
}

func norm(x, total int, min, max float64) float64 {
	return (max-min)*float64(x)/float64(total) - max
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
