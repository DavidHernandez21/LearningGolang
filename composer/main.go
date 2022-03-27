package main

import (
	"flag"
	"fmt"
	"log"

	splitCsv "github.com/tolik505/split-csv"
)

func ExampleSplitCsv(inputFilePath, outputDirectory string) {
	splitter := splitCsv.New()
	splitter.FileChunkSize = 50000000 //in bytes (50MB)
	result, err := splitter.Split(inputFilePath, outputDirectory)
	if err != nil {
		log.Fatalf("error in the split function: %v\n", err)
	}
	fmt.Println(result)
	// Output: [testdata/test_1.csv testdata/test_2.csv testdata/test_3.csv]
}

func main() {

	inputFilePath := flag.String("i", "DATA_GIA_GOOGLE_20201023-161045_prod.csv", "input file path")
	outputDirectory := flag.String("0", "./", "output directory")

	flag.Parse()

	ExampleSplitCsv(*inputFilePath, *outputDirectory)
}
