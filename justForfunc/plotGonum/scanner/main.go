package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	filePath := flag.String("f", "data.txt", "path of the file")
	outputFile := flag.String("o", "out.png", "path of output file")
	flag.Parse()

	xys, err := readData(*filePath)
	if err != nil {
		log.Fatalf("could not read %s: %v", *filePath, err)
	}
	_ = xys

	err = plotData(*outputFile, xys)
	if err != nil {
		log.Fatalf("could not plot data: %v", err)
	}
}

// type xy struct{ x, y float64 }

func readData(path string) (xys plotter.XYs, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() (plotter.XYs, error) {
		if errClose := f.Close(); errClose != nil {
			err = errClose
			xys = nil
		}

		return xys, err
	}()

	s := bufio.NewScanner(f)
	for s.Scan() {
		var x, y float64
		_, err := fmt.Sscanf(s.Text(), "%f,%f", &x, &y)
		if err != nil {
			log.Printf("discarding bad data point %q: %v", s.Text(), err)
			continue
		}
		xys = append(xys, plotter.XY{
			X: x,
			Y: y,
		})
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("could not scan: %v", err)
	}
	return xys, err
}

func plotData(path string, xys plotter.XYs) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create %s: %v", path, err)
	}

	p := plot.New()
	if err != nil {
		return fmt.Errorf("could not create plot: %v", err)
	}

	// create scatter with all data points
	s, err := plotter.NewScatter(xys)
	if err != nil {
		return fmt.Errorf("could not create scatter: %v", err)
	}
	s.GlyphStyle.Shape = draw.CrossGlyph{}
	s.Color = color.RGBA{R: 255, A: 255}
	p.Add(s)

	var x, c float64
	x = 1.2
	c = -3

	// create fake linear regression result
	l, err := plotter.NewLine(plotter.XYs{
		plotter.XY{
			X: 3,
			Y: 3*x + c,
		},
		plotter.XY{
			X: 20,
			Y: 20*x + c,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create line: %v", err)
	}
	p.Add(l)

	wt, err := p.WriterTo(256, 256, "png")
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}
	_, err = wt.WriteTo(f)
	if err != nil {
		return fmt.Errorf("could not write to %s: %v", path, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close %s: %v", path, err)
	}
	return nil
}
