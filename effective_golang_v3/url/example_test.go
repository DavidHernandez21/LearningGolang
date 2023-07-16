package url_test

import (
	"fmt"
	"log"

	"src/github.com/DavidHernandez21/effective_golang_v3/url"
)

func ExampleURL() {
	u, err := url.Parse("http://foo.com/go")

	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	fmt.Println(u)
	// Output:
	// https://foo.com/go
}

func ExampleURL_fields() {
	u, err := url.Parse("https://foo.com/go")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(u.Scheme)
	fmt.Println(u.Host)
	fmt.Println(u.Path)
	fmt.Println(u)
	// Output:
	// https
	// foo.com
	// go
	// https://foo.com/go

}
