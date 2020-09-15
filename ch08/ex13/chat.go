package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var timeout = flag.Int("s", 300, "timeout (sec)")

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

type cliChan chan<- string
type cliName string
type client struct {
	name cliName
	ch   cliChan
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[cliName]cliChan)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				clients[cli] <- msg
			}

		case cli := <-entering:
			var sb strings.Builder
			if len(clients) > 0 {
				fmt.Fprint(&sb, "current clients in the room:")
				for c := range clients {
					fmt.Fprintf(&sb, "\n\t%s", c)
				}
			} else {
				fmt.Fprint(&sb, "no one in the room")
			}
			cli.ch <- sb.String()
			clients[cli.name] = cli.ch

		case cli := <-leaving:
			ch := cli.ch
			delete(clients, cli.name)
			close(ch)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	cli := client{cliName(who), ch}
	entering <- cli

	input := bufio.NewScanner(conn)
	act := make(chan bool)
	go func() {
		timeout := time.Duration(*timeout) * time.Second
		noMessage := time.NewTimer(timeout)
		for {
			select {
			case a := <-act:
				noMessage.Reset(timeout)
				if !a {
					return
				}
			case <-noMessage.C:
				conn.Close()
			}
		}
	}()

	for input.Scan() {
		act <- true
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	act <- false
	leaving <- cli
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
