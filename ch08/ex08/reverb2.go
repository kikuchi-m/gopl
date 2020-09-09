package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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
	input := bufio.NewScanner(c)
	called := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(10 * time.Second):
				c.Close()
				return
			case <-called:
				// do nothing
			}
		}
	}()

	for input.Scan() {
		called <- struct{}{}
		go echo(c, input.Text(), 1*time.Second)
	}
	c.Close() // NOTE: ignoring potential errors from input.Err()
}

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}
