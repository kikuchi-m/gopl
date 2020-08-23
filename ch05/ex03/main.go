package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	visit(doc)
}

func visit(n *html.Node) {
	if n == nil {
		return
	}
	if n.Data == "script" || n.Data == "style" {
		visit(n.NextSibling)
	} else {
		if n.Type == html.TextNode {
			fmt.Println(n.Data)
		}
		visit(n.FirstChild)
		visit(n.NextSibling)
	}
}
