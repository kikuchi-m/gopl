package main

import (
	"fmt"
	"time"
)

var cnt = 0

func main() {
	aChan := make(chan int)
	anChan := make(chan int)

	go p2(aChan, anChan)
	anChan <- 0
	go p2(anChan, aChan)

	last := 0
	for {
		<-time.After(1 * time.Second)
		c := cnt
		fmt.Printf("%d (add %d)\r", c, c-last)
		last = c
	}
}

func p2(send chan<- int, recv <-chan int) {
	for {
		c := <-recv
		c++
		cnt = c
		send <- c
	}
}
