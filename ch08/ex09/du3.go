package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type rootDir struct {
	path          string
	size          chan int64
	nFiles, nByes int64
}

var vFlag = flag.Bool("v", false, "show verbose progress messages")

func main() {
	flag.Parse()

	paths := flag.Args()
	var roots []*rootDir
	if len(paths) == 0 {
		roots = []*rootDir{&rootDir{".", make(chan int64), 0, 0}}
	} else {
		for _, p := range paths {
			roots = append(roots, &rootDir{p, make(chan int64), 0, 0})
		}
	}

	var nDirs sync.WaitGroup
	for _, root := range roots {
		nDirs.Add(1)
		go walkDir(root.path, &nDirs, root.size)
		go func(r *rootDir) {
			for {
				s, open := <-r.size
				if !open {
					break
				}
				r.nFiles++
				r.nByes += s
			}
		}(root)
	}

	done := make(chan struct{})
	go func() {
		nDirs.Wait()
		for _, root := range roots {
			close(root.size)
		}
		close(done)
	}()

	if *vFlag {
		fmt.Printf("%s", strings.Repeat("\n", len(roots)))
		printForEach(roots)
		var tick <-chan time.Time
		tick = time.Tick(500 * time.Millisecond)
	loop:
		for {
			select {
			case <-done:
				break loop
			case <-tick:
				printForEach(roots)
			}
		}
	} else {
		<-done
	}

	printTotal(roots)
}

func printTotal(roots []*rootDir) {
	var nFiles, nByes int64
	for _, r := range roots {
		nFiles += r.nFiles
		nByes += r.nByes
	}
	fmt.Printf("%d files  %.9f GB\n", nFiles, float64(nByes)/1e9)
}

func printForEach(roots []*rootDir) {
	n := len(roots)
	var sb strings.Builder
	for _, r := range roots {
		fmt.Fprintf(&sb, "\033[K%s: %.3f MB (%d files)\n", r.path, float64(r.nByes)/1e6, r.nFiles)
	}
	fmt.Printf("\033[%dA%s", n, sb.String())
}

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

var sema = make(chan struct{}, 20)

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}
	defer func() { <-sema }()

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
