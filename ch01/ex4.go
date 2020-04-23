package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	lines := make(map[string]map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		CollectDups(os.Stdin, "stdin", lines)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ex4: %v\n", err)
			}
			CollectDups(f, arg, lines)
			f.Close()
		}
	}

	for l, nameCount := range lines {
		if Total(nameCount) > 1 {
			fmt.Printf("%s\t%v\n", l, UniqNames(nameCount))
		}
	}
}

func CollectDups(f *os.File, name string, lines map[string]map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		t := input.Text()
		if lines[t] == nil {
			lines[t] = make(map[string]int)
		}
		// when stdin, breaks if same line appeared
		// if _, exists := lines[t][name]; exists {
		// 	break
		// }
		lines[t][name]++
	}
	// ignore error of input.Error()
}

func Total(nameCount map[string]int) int {
	total := 0
	for _, c := range nameCount {
		total += c
	}
	return total
}

func UniqNames(names map[string]int) []string {
	var res []string
	for n, _ := range names {
		res = append(res, n)
	}
	return res
}
