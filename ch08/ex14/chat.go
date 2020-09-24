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
	addr cliName
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
			clients[cli.addr] = cli.ch

		case cli := <-leaving:
			ch := cli.ch
			delete(clients, cli.addr)
			close(ch)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()

	input := bufio.NewScanner(conn)
	name := who
	if input.Scan() {
		name = input.Text()
		if name != "" {
			fmt.Printf("connected: %s\n", name)
		} else {
			fmt.Printf("connected, no name: %s\n", who)
			name = who
		}
	} else {
		fmt.Printf("failed to receive name: \n", who)
	}

	ch <- fmt.Sprintf("You are %s (%s)", name, who)
	messages <- fmt.Sprintf("%s has arrived", name)
	cli := client{cliName(who), cliName(name), ch}
	entering <- cli

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
		messages <- fmt.Sprintf("[%s] %s", name, input.Text())
	}
	// NOTE: ignoring potential errors from input.Err()

	act <- false
	leaving <- cli
	messages <- fmt.Sprintf("%s has left", name)
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
