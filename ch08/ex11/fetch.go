package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func main() {
	urls := os.Args[1:]
	done := make(chan struct{})
	res := make(chan struct{}, len(urls))

	for _, url := range urls {
		fmt.Printf("go fetch: %s\n", url)
		go func(url string) {
			local, n, err := fetch(url, done)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fetch %s: %v\n", url, err)
				return
			}
			fmt.Fprintf(os.Stderr, "%s => %s (%d bytes).\n", url, local, n)
			res <- struct{}{}
		}(url)
	}
	<-res
	close(done)
}

func fetch(url string, done <-chan struct{}) (filename string, n int64, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, err
	}

	req.Cancel = done
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" || local == "." {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	n, err = io.Copy(f, resp.Body)
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	return local, n, err
}
