package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pretty print: %v\n", err)
		os.Exit(1)
	}
	forEachNode(doc, startElement, endElement)
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

var depth int

func indent() int {
	return depth * 2
}

func startElement(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		fmt.Printf("%*s<%s", indent(), "", n.Data)
		for _, a := range n.Attr {
			fmt.Printf(" %s=%s", a.Key, a.Val)
		}
		if n.FirstChild != nil || n.Data == "meta" {
			fmt.Printf(">\n")
		} else {
			fmt.Printf("/>\n")
		}
		depth++
	case html.TextNode:
		text := strings.Trim(n.Data, " \n")
		if text != "" {
			for _, l := range strings.Split(text, "\n") {
				fmt.Printf("%*s%s\n", indent(), "", l)
			}
		}
	case html.CommentNode:
		fmt.Printf("%*s<!-- %s -->\n", indent(), "", n.Data)
	case html.DoctypeNode:
		fmt.Printf("<!DOCTYPE %s>\n", n.Data)
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		if n.FirstChild != nil {
			fmt.Printf("%*s</%s>\n", indent(), "", n.Data)
		}
	}
}
