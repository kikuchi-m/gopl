package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func AppendStrings(args []string) string {
	var s, sep string
	for i := 1; i < len(args); i++ {
		s += sep + args[i]
		sep = " "
	}
	return s
}

func JoinStrings(args []string) string {
	return strings.Join(args, " ")
}

func main() {
	args := os.Args[1:]

	start1 := time.Now()
	fmt.Println(AppendStrings(args))
	d1 := time.Since(start1).Microseconds()

	start2 := time.Now()
	fmt.Println(JoinStrings(args))
	d2 := time.Since(start2).Microseconds()

	fmt.Println("joining is faster", d1-d2, "microsecs than appengind")
	fmt.Println("(appending:", d1, ", joining:", d2, ")")
}
