// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package memo5

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func Test(t *testing.T) {
	m := New(httpGetBody)
	defer m.Close()

	cnt := 0

	for url := range incomingURLs() {
		start := time.Now()
		done := make(chan struct{})
		if cnt < 2 {
			go func() {
				<-time.After(10 * time.Millisecond)
				close(done)
			}()
		}
		value, err := m.Get(url, done)
		cnt++
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("%s, %s, %d bytes\n",
			url, time.Since(start), len(value.([]byte)))
	}
}

func incomingURLs() <-chan string {
	ch := make(chan string)
	go func() {
		for _, url := range []string{
			"https://golang.org",
			"https://godoc.org",
			"https://play.golang.org",
			"http://gopl.io",

			"https://golang.org",
			"https://godoc.org",
			"https://play.golang.org",
			"http://gopl.io",
		} {
			ch <- url
		}
		close(ch)
	}()
	return ch
}
