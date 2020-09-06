package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	fmt.Println("handling connection")
	input := bufio.NewScanner(c)

	var wg sync.WaitGroup
	for input.Scan() {
		wg.Add(1)
		go func(s string) {
			echo(c, s, 1*time.Second)
			wg.Done()
		}(input.Text())
	}

	ch := make(chan struct{})
	go func() {
		wg.Wait()
		switch c := c.(type) {
		case *net.TCPConn:
			c.CloseWrite()
		default:
			log.Println("not net.TCPConn")
		}
		ch <- struct{}{}
	}()
	<-ch

	c.Close() // NOTE: ignoring potential errors from input.Err()
}

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}
