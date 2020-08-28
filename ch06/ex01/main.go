package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var ints1 IntSet
	for _, a := range os.Args[1:] {
		i, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		ints1.Add(i)
	}
	fmt.Println(ints1)
	fmt.Println(ints1.String())

	fmt.Println("Copy()")
	ints2 := ints1.Copy()
	fmt.Println(ints2.String())

	fmt.Printf("length: %d\n", ints1.Len())

	fmt.Println("Remove(3)")
	ints1.Remove(3)
	fmt.Println(ints1.String())

	fmt.Printf("Remove(122): %d\n", ints1.String())
	fmt.Println(ints1.String())

	fmt.Println("Clear()")
	ints1.Clear()
	fmt.Printf("ints1: %s\n", ints1.String())
	fmt.Printf("ints2: %s\n", ints2.String()) // verify ints2 remained
}

type IntSet struct {
	words []uint64
}

// beg: ch06
func (s *IntSet) Len() int {
	return len(s.words)
}

func (s *IntSet) Remove(x int) {
	if s.Has(x) {
		word, bit := wordAndBit(x)
		s.words[word] = s.words[word] ^ (1 << bit)
	}
}

func (s *IntSet) Clear() {
	for i, _ := range s.words {
		s.words[i] = 0
	}
}

func (s *IntSet) Copy() *IntSet {
	var copied IntSet
	for i, _ := range s.words {
		copied.words = append(copied.words, 0)
		copied.words[i] = s.words[i]
	}
	return &copied
}

func wordAndBit(x int) (int, uint) {
	return x / 64, uint(x % 64)
}

// end: ch06

func (s *IntSet) Has(x int) bool {
	word, bit := wordAndBit(x)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := wordAndBit(x)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}
