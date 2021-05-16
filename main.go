package main

import (
	"flag"
	"fmt"
	"github.com/smartclash/Sandeadry/parser"
	"net/url"
	"strings"
)

func main() {
	link := flag.String("l", "", "Link to the degree you want to parse MCQs")
	flag.Parse()

	if *link != "" {
		invokeParser(link)
		return
	}

	flag.PrintDefaults()
}

func invokeParser(link *string) {
	parse, err := url.Parse(*link)
	if err != nil {
		fmt.Println("Please enter a proper link")
		flag.PrintDefaults()
		return
	}

	if !strings.EqualFold(parse.Hostname(), "www.sanfoundry.com") {
		fmt.Println("Enter only sanfoundry links")
		flag.PrintDefaults()
		return
	}

	parser.Parser(*link)
}
