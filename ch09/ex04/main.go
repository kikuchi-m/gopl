package main

import (
	"flag"
	"fmt"
	"time"
)

// Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz
// memory 64G
// swap 2G

var sizes = []int{
	1e5,
	1e6,
	2e6,
	5e6,
	1e7,
	15e6,
	17e6,

	// available when executing solitarily
	// 175e5,   // 17500000: 4.026785105s (creating pipe: 20.84946375s)
	// 178e5,   // 17800000: 4.149323902s (creating pipe: 20.735688705s)
	// 1785e4,  // 17850000: 4.313809994s (creating pipe: 20.882241036s)
}

var sFlag = flag.Uint("s", 0, "size")

func main() {
	flag.Parse()

	ss := sizes
	if *sFlag != 0 {
		ss = []int{int(*sFlag)}
	}

	for _, s := range ss {
		ch0 := make(chan time.Time)
		ps := time.Now()
		p := pipe(s, ch0)
		pd := time.Now().Sub(ps)
		ch0 <- time.Now()
		t := <-p
		fmt.Printf("%9d: %s (creating pipe: %s)\n", s, time.Now().Sub(t), pd)
	}
}

func pipe(size int, ch0 <-chan time.Time) (exit <-chan time.Time) {
	exit = ch0
	for i := 0; i < size; i++ {
		next := make(chan time.Time)
		go func(n chan<- time.Time, p <-chan time.Time) {
			n <- <-p
		}(next, exit)
		exit = next
	}
	return
}
