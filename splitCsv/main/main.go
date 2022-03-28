package main

import (
	"flag"
	"fmt"
	"log"
	"src/github.com/DavidHernandez21/src/github.com/DavidHernandez21/splitCsv"
)

func ExampleSplitCsv(inputFilePath, outputDirectory, jsonKeyPath string) {
	splitter := splitCsv.New()
	splitter.FileChunkSize = 500000 //in bytes (50MB)
	result, err := splitter.Split(inputFilePath, outputDirectory, jsonKeyPath)
	if err != nil {
		log.Fatalf("error in the split function: %v\n", err)
	}
	fmt.Println(result)
	// Output: [testdata/test_1.csv testdata/test_2.csv testdata/test_3.csv]
}

func main() {

	inputFilePath := flag.String("i", "DATA_GIA_GOOGLE_20201023-161045_prod.csv", "input file path")
	outputDirectory := flag.String("o", "./", "output directory")
	// C:/Users/asbel.hernandez/OneDrive - Accenture/Personal/vs_studio/golang/src/github.com/DavidHernandez21/splitCsv/key.json
	inputJsonKeyPath := flag.String("jsonKey", "key.json", "json key file path")

	flag.Parse()

	ExampleSplitCsv(*inputFilePath, *outputDirectory, *inputJsonKeyPath)
}
