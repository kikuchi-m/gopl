package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

var nameFlag = flag.String("name", "", "client name")

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	var n string
	if *nameFlag != "" {
		n = *nameFlag
	}
	r, _ := regexp.Compile("[\n\r]+")
	n = strings.TrimSpace(r.ReplaceAllString(n, " "))
	_, err = conn.Write([]byte(n + "\n"))
	if err != nil {
		log.Fatalf("failed to send name: %v", err)
	}
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalf("error on reading stdin: %v", err)
	}
}
