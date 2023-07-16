package main

import (
	"fmt"

	"github.com/DavidHernandez21/effective_golnag_v3/url"
)

func main() {
	// const rawurl = "http://foo.com/go"
	const rawurl = "http://d"

	u, err := url.Parse(rawurl)
	if err != nil {
		fmt.Printf("Parse(%q) err = %q", rawurl, err)
		return

	}
	fmt.Printf("Parse(%q).URL = %#v", rawurl, u)

}
