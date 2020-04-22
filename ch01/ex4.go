package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
)

func main() {
	lines := make(map[string]map[string]string)
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

	for l, names := range lines {
		if len(names) > 1 {
			fmt.Printf("%s\t%v\n", l, reflect.ValueOf(names).MapKeys())
		}
	}
}

func CollectDups(f *os.File, name string, lines map[string]map[string]string) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		t := input.Text()
		if lines[t] == nil {
			lines[t] = make(map[string]string)
		}
		// when stdin, breaks if same line appeared
		// if _, exists := lines[t][name]; exists {
		// 	break
		// }
		lines[t][name] = name
	}
	// ignore error of input.Error()
}
