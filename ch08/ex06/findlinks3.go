package main

import (
	"flag"
	"fmt"
	"log"

	"gopl.io/ch5/links"
)

type link struct {
	url   string
	depth uint
}

func main() {
	depth := flag.Uint("depth", 3, "depth")
	flag.Parse()
	fmt.Printf("depth: %d\n", *depth)
	urls := flag.Args()
	fmt.Printf("depth: %d, URLs: %s\n", *depth, urls)

	worklist := make(chan []link)  // lists of URLs, may have duplicates
	unseenLinks := make(chan link) // de-duplicated URLs

	go func() {
		if len(urls) > 0 {
			worklist <- urlsToLinks(urls, 0)
		} else {
			close(worklist)
		}
	}()

	for i := 0; i < 20; i++ {
		go func() {
			for l := range unseenLinks {
				// fmt.Printf("%s (%d)\n", l.url, l.depth)
				fmt.Printf("%s\n", l.url)
				if l.depth < *depth {
					foundLinks := crawl(l)
					d := l.depth + 1
					go func() { worklist <- urlsToLinks(foundLinks, d) }()
				}
			}
		}()
	}

	seen := make(map[string]bool)
	for list := range worklist {
		for _, l := range list {
			if !seen[l.url] {
				seen[l.url] = true
				unseenLinks <- l
			}
		}
	}
}

func crawl(l link) []string {
	list, err := links.Extract(l.url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func urlsToLinks(urls []string, depth uint) []link {
	var ls []link
	for _, l := range urls {
		ls = append(ls, link{l, depth})
	}
	return ls
}
