package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type City struct {
	Name       string
	Population int
}

func main() {

	defer timeTrack(time.Now(), "main")

	f := openCsv("cities.csv")
	defer f.Close()

	rows := genRows(f)

	filterSmallCity := filterByMinPopulation(8000000)
	// upperCityName := modifyCityName("lower")

	// fan out pattern
	// more than one worker competing to consume the same channel

	//       __ worker1
	// rows /__ worker2
	//      \__ worker3
	//       __ workern...

	// ur1 := upperCityName(filterSmallCity(rows, 1), 1)
	// ur2 := upperCityName(filterSmallCity(rows, 2), 2)
	// ur3 := upperCityName(filterSmallCity(rows))
	// ur4 := upperCityName(filterSmallCity(rows))

	// var mySlice []<-chan City

	workers := 2
	mySlice := make([]<-chan City, 0, workers)

	for i := 0; i < workers; i++ {

		mySlice = append(mySlice, upperCityName(filterSmallCity(rows, i+1), i+1))

	}

	// mySlice := genWorkersSlice(2, upperCityName, filterSmallCity, rows)

	// log.Println(mySlice[2])
	// fan in pattern consolidates the outputs from multiple channels into one
	//
	// worker1 ___
	// worker2 ___\ output
	// worker3 ___/

	for c := range fanIn(mySlice...) {
		log.Println(c)
	}

}

func upperCityName(cities <-chan City, worker int) <-chan City {
	out := make(chan City)
	go func(worker int) {
		i := 0
		for c := range cities {
			out <- City{Name: strings.ToUpper(c.Name), Population: c.Population}
			i += 1
		}
		log.Printf("upperCityName worker %v handled %v cities\n", worker, i)
		close(out)
	}(worker)
	return out
}

// func modifyCityName(modification string) func(<-chan City, int) <-chan City {

// 	return func(cities <-chan City, worker int) <-chan City {
// 		out := make(chan City)

// 		go func(worker int) {
// 			// log.Printf("upperCityName worker %v\n", worker)
// 			i := 0
// 			for c := range cities {
// 				name := c.Name
// 				if modification == "upper" {
// 					name = strings.ToUpper(c.Name)
// 				} else if modification == "lower" {
// 					name = strings.ToLower(c.Name)
// 				}
// 				out <- City{Name: name, Population: c.Population}
// 				i += 1
// 			}
// 			log.Printf("modifyCityName worker %v handled %v cities\n", worker, i)
// 			close(out)
// 		}(worker)
// 		return out

// 	}
// }

func filterByMinPopulation(min int) func(<-chan City, int) <-chan City {
	return func(cities <-chan City, worker int) <-chan City {
		out := make(chan City)

		go func(worker int) {
			i := 0
			for c := range cities {
				if c.Population > min {
					out <- City{Name: c.Name, Population: c.Population}
					i += 1
				}
			}
			log.Printf("filterByMinPopulation worker %v handled %v cities\n", worker, i)
			close(out)

		}(worker)
		return out
	}
}

func genRows(r io.Reader) chan City {
	out := make(chan City)
	go func() {
		reader := csv.NewReader(r)
		_, err := reader.Read()
		if err != nil {
			log.Fatal(err)
		}
		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			populationInt, err := strconv.Atoi(row[9])
			if err != nil {
				continue
			}
			out <- City{
				Name:       row[1],
				Population: populationInt,
			}
		}
		close(out)
	}()
	return out
}

func fanIn(chans ...<-chan City) <-chan City {
	out := make(chan City)
	wg := &sync.WaitGroup{}
	wg.Add(len(chans))
	for _, c := range chans {

		go func(city <-chan City) {

			defer wg.Done()

			for r := range city {
				out <- r
			}

		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	// wg.Wait()
	// close(out)

	return out
}

func openCsv(filePath string) *os.File {

	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	return f
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// func genWorkersSlice(n int, f1, f2 func(<-chan City) <-chan City, rows chan City) []<-chan City {

// 	mySlice := make([]<-chan City, n)

// 	for i := 0; i < n; i++ {

// 		mySlice = append(mySlice, f1(f2(rows)))

// 	}

// 	return mySlice
// }
