package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var intArgs []int
	for _, a := range os.Args[1:] {
		i, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		intArgs = append(intArgs, i)
	}
	var ints IntSet
	ints.AddAll(intArgs...)
	fmt.Println(ints)
	fmt.Println(ints.String())
}

type IntSet struct {
	words []uint64
}

// beg: ch07
func (s *IntSet) AddAll(ints ...int) {
	for _, x := range ints {
		s.Add(x)
	}
}

// end: ch07

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
