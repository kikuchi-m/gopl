package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	names := os.Args[1:]
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tags: %v\n", err)
		os.Exit(1)
	}
	elms := ElementsByTagName(doc, names...)
	fmt.Println(elms)
}

func ElementsByTagName(doc *html.Node, names ...string) []*html.Node {
	var elms []*html.Node
	var traverse func(n *html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && contains(n.Data, names) {
			elms = append(elms, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return elms
}

func contains(name string, names []string) bool {
	for _, n := range names {
		if name == n {
			return true
		}
	}
	return false
}
