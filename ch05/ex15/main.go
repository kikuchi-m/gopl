package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("--- max() ---")
	fmt.Println(max())
	fmt.Println(max(1))
	fmt.Println(max(1, 2))
	fmt.Println(max(1, 4, 2, -5, 3))

	fmt.Println("--- max2() ---")
	fmt.Println(max2(1))
	fmt.Println(max2(1, 2))
	fmt.Println(max2(1, 4, 2, -5, 3))

	fmt.Println("--- min() ---")
	fmt.Println(min())
	fmt.Println(min(1))
	fmt.Println(min(1, 2))
	fmt.Println(min(1, 4, 2, -5, 3))

	fmt.Println("--- min2() ---")
	fmt.Println(min2(1))
	fmt.Println(min2(1, 2))
	fmt.Println(min2(1, 4, 2, -5, 3))
}

func max(values ...int) int {
	if len(values) == 0 {
		return math.MaxInt64
	}
	m := math.MinInt32
	for _, v := range values {
		if v > m {
			m = v
		}
	}
	return m
}

func max2(v int, rest ...int) int {
	switch len(rest) {
	case 0:
		return v
	case 1:
		r := rest[0]
		if v > r {
			return v
		} else {
			return r
		}
	}
	return max2(v, max2(rest[0], rest[1:]...))
}

func min(values ...int) int {
	if len(values) == 0 {
		return math.MinInt64
	}
	m := math.MaxInt32
	for _, v := range values {
		if v < m {
			m = v
		}
	}
	return m
}

func min2(v int, rest ...int) int {
	switch len(rest) {
	case 0:
		return v
	case 1:
		r := rest[0]
		if v < r {
			return v
		} else {
			return r
		}
	}
	return min2(v, min2(rest[0], rest[1:]...))
}
