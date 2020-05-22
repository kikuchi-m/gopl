package main

import (
	"fmt"
	"os"
)

func main() {
	ss := os.Args[1:]
	ss = RemoveAdjacent(ss)
	fmt.Println("result", ss)
}

func RemoveAdjacent(ss []string) []string {
	if len(ss) < 2 {
		return ss
	}
	cur := 0
	for i := 1; i < len(ss); i++ {
		if ss[cur] != ss[i] {
			cur++
			ss[cur] = ss[i]
		}
	}
	return ss[:cur+1]
}
