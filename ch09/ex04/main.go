package main

import (
	"fmt"
	"time"
)

func main() {

	ch0 := make(chan time.Time)
	go func() {
		for {
			ch0 <- time.Now()
		}
	}()

	prev := ch0
	for i := 0; ; i++ {
		next := make(chan time.Time)
		go func(i int, n chan<- time.Time, p <-chan time.Time) {
			for t := range p {
				d := time.Now().Sub(t)
				n <- t
				fmt.Printf("%d: %v\n", i, d)
			}
		}(i, next, prev)
		prev = next
	}
}
