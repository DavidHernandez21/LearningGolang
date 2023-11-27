package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

func readOrders(name string) [][]string {

	f, err := os.Open(name)

	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", name, err.Error())
	}

	defer f.Close()

	r := csv.NewReader(f)

	r.Comma = ','

	rows, err := r.ReadAll()

	if err != nil {
		log.Fatalln("Cannot read CSV data: ", err.Error())
	}

	return rows

}

func calculate(rows [][]string) [][]string {

	sum := 0.0
	nb := 0

	for i := range rows {

		if i == 0 {
			rows[0] = append(rows[0], "Total")
			continue
		}

		item := rows[i][2]

		// price, err := strconv.Atoi(strings.ReplaceAll(rows[i][3], ".", ""))
		price, err := strconv.ParseFloat(rows[i][3], 32)

		if err != nil {
			log.Fatalf("Cannot retrieve price of %s: %s\n", item, err)
		}

		quantity, err := strconv.Atoi(rows[i][4])

		if err != nil {
			log.Fatalf("Cannot retrive quantity of %s: %s\n", item, err)
		}

		total := price * float64(quantity)

		rows[i] = append(rows[i], strconv.FormatFloat(total, 'f', 2, 64))

		sum += total

		if item == "Ball Pen" {
			nb += quantity
		}
	}

	rows = append(rows, []string{"", "", "", "Sum", "", strconv.FormatFloat(sum, 'f', 2, 64)})
	rows = append(rows, []string{"", "", "", "Ball Pens", fmt.Sprint(nb), ""})

	return rows
}

func writeOrders(name string, rows [][]string, wg *sync.WaitGroup) {

	defer wg.Done()

	f, err := os.Create(name)

	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", name, err.Error())
	}

	defer func() {
		e := f.Close()
		if e != nil {
			log.Fatalf("Cannot close '%s': %s\n", name, e.Error())
		}
	}()

	w := csv.NewWriter(f)
	err = w.WriteAll(rows)
	log.Printf("CSV '%s' successfully create", name)

}

func main() {

	var wg = sync.WaitGroup{}

	// rowsChan := make(chan [][]string)

	rows := readOrders("data.csv")

	rows = calculate(rows)

	wg.Add(1)
	go writeOrders("ordersReport.csv", rows, &wg)

	wg.Wait()

}
