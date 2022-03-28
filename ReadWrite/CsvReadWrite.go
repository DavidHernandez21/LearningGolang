// package main

// import (
// 	"encoding/csv"
// 	"fmt"
// 	"os"
// )

// func main() {

// 	file, err := os.Open("ric_cliente.csv")

// 	if err != nil {
// 		fmt.Println(err)
// 		// panic("stopping because of an error")
// 	}

// 	reader := csv.NewReader(file)
// 	records, err1 := reader.ReadAll()

// 	if err1 != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Printf("records type %T\n", records)
// 	fmt.Println(records)

// }
