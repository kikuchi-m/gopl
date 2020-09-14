package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	worklist := make(chan []string)
	unseenLinks := make(chan string)
	seen := make(map[string]bool)
	nSeen := 0

	go func() {
		worklist <- os.Args[1:]
	}()

	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() {
					nSeen++
					worklist <- foundLinks
				}()
			}
		}()
	}

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

	for list := range worklist {
		if cancelled() && nSeen == len(seen) {
			break
		}
		for _, link := range list {
			if cancelled() {
				break
			}
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func Extract(url string) ([]string, error) {
	var req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Cancel = done
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
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
