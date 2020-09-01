package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type zone struct {
	name string
	host string
}

type dial struct {
	zone
	conn *bufio.Scanner
}

func main() {
	var zones []zone
	for _, a := range os.Args[1:] {
		nh := strings.Split(a, "=")
		zones = append(zones, zone{nh[0], nh[1]})
	}
	var dials []dial
	for _, h := range zones {
		conn, err := net.Dial("tcp", h.host)
		if err != nil {
			log.Fatal(err)
		} else {
			defer conn.Close()
			dials = append(dials, dial{h, bufio.NewScanner(conn)})
		}
	}
	for {
		for _, d := range dials {
			d.conn.Scan()
			fmt.Printf("%s=%s ", d.name, d.conn.Text())
		}
		fmt.Printf("\r")
	}
}
