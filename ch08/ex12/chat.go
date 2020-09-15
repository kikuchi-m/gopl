package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
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

type client struct {
	who     string
	message chan<- string
}

var (
	entering = make(chan *client)
	leaving  = make(chan *client)
	messages = make(chan string)
	clients  = make(map[*client]bool)
)

func broadcaster() {
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.message <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.message)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	var sb strings.Builder
	if len(clients) > 0 {
		fmt.Fprint(&sb, "current clients in the room:")
		for c := range clients {
			fmt.Fprintf(&sb, "\n\t%s", c.who)
		}
	} else {
		fmt.Fprint(&sb, "no one in the room")
	}
	ch <- sb.String()
	messages <- who + " has arrived"
	cli := client{who, ch}
	entering <- &cli

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- &cli
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
