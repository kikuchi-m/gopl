package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "element count: %v\n", err)
		os.Exit(1)
	}
	for elms, count := range visit(make(map[string]int), doc) {
		fmt.Println(elms, count)
	}
}

func visit(elms map[string]int, n *html.Node) map[string]int {
	if n == nil {
		return elms
	}
	if n.Type == html.ElementNode {
		elms[n.Data]++
	}
	return visit(visit(elms, n.FirstChild), n.NextSibling)
}
