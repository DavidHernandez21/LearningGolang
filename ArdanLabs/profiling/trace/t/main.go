package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/google/martian/v3/log"
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
	docs := make([]string, 4000)
	sb := strings.Builder{}
	sb.Grow(90)
	for i := range docs {
		doc, err := joinBuilder(&sb, "newsfeed-", strconv.Itoa(i), ".xml")
		if err != nil {
			log.Errorf("failed to join: %v", err)
		}
		docs[i] = doc
		// docs[i] = fmt.Sprintf("newsfeed-%.4d.xml", i)
	}

	// fmt.Printf("%#v", docs)

	d, err := os.Open("../newsfeed.xml")
	if err != nil {
		log.Errorf("failed to open: %v", err)
	}
	defer d.Close()
	data, err := io.ReadAll(d)
	if err != nil {
		log.Errorf("failed to read: %v", err)
	}

	fmt.Println(len(data))

}
