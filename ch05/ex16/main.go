package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(join(", "))
	fmt.Println(join(", ", "a"))
	fmt.Println(join(", ", "a", "b", "c"))
}

func join(sep string, s ...string) string {
	return strings.Join(s, sep)
}
