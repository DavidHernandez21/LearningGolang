// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"strings"

// 	"github.com/jung-kurt/gofpdf"
// )

// func main() {

// 	file := "cover_letter.txt"

// 	content, err := ioutil.ReadFile(file)

// 	if err != nil {
// 		log.Fatalf("%s file not found", file)
// 	}

// 	pdf := gofpdf.New("P", "mm", "A4", "")

// 	// pdf.AddUTF8Font("dejavu", "", "DejaVuSansCondensed.ttf", )

// 	pdf.AddPage()

// 	pdf.SetFont("helvetica", "", 14)

// 	pdf.MultiCell(190, 5, strings.TrimRight(string(content), "\r\n"), "0", "J", false)

// 	// strings.TrimRight(string(content), "\r\n")

// 	err1 := pdf.OutputFileAndClose("cover_letter.pdf")

// 	if err1 != nil {
// 		log.Fatalf("error while creating the file: %s", err1)
// 	}

// 	fmt.Println("pdf file created")

// }
