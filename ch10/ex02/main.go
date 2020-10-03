package main

import (
	"fmt"
	"os"

	"local.gopl/archive"
	_ "local.gopl/archive/tar"
	_ "local.gopl/archive/zip"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("usage: %s FILE [FILE...]\n", os.Args[0])
		os.Exit(1)
	}
	for _, name := range os.Args[1:] {
		files, err := archive.List(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v (%s)\n", err, name)
			continue
		}
		fmt.Printf("entries of %s\n", name)
		for _, f := range files {
			fmt.Println(f.Path)
		}
	}
}
