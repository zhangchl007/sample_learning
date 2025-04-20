package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

// IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words    []uint64
	wordSize int // bits per word: 32 or 64
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/s.wordSize, uint(x%s.wordSize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/s.wordSize, uint(x%s.wordSize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < s.wordSize; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", s.wordSize*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func (s *IntSet) Len() int {
	count := 0
	for _, word := range s.words {
		for word != 0 {
			count += int(word & 1)
			word >>= 1
		}
	}
	return count
}

// Remove removes x from the set.
func (s *IntSet) Remove(x int) {
	word, bit := x/s.wordSize, uint(x%s.wordSize)
	if word < len(s.words) {
		s.words[word] &^= 1 << bit
	}
}

// Clear clears the set.
func (s *IntSet) Clear() {
	s.words = nil
}

// Copy returns a copy of the set.
func (s *IntSet) Copy() *IntSet {
	copiedSet := &IntSet{wordSize: s.wordSize}
	copiedSet.words = make([]uint64, len(s.words))
	copy(copiedSet.words, s.words)
	return copiedSet
}

// AddSlice adds a slice of integers to the set.
func (s *IntSet) AddSlice(slice []int) {
	for _, x := range slice {
		s.Add(x)
	}
}

func (s *IntSet) IntersectWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &= tword
		} else if i < len(s.words) {
			s.words[i] = 0
		}
	}
}

func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &^= tword
		}
	}
}

func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// Elements returns a slice of all elements in the set.
func (s *IntSet) Elements() []int {
	var elements []int
	for i, word := range s.words {
		for j := 0; j < s.wordSize; j++ {
			if word&(1<<uint(j)) != 0 {
				elements = append(elements, s.wordSize*i+j)
			}
		}
	}
	return elements
}

func getWordSize() int {
	// Use strconv.IntSize for the current platform
	return strconv.IntSize
}

func main() {
	wordSize := getWordSize()
	fmt.Printf("Architecture: %s, int size: %d bits\n", runtime.GOARCH, wordSize)

	var x, y IntSet
	x.wordSize = wordSize
	y.wordSize = wordSize

	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String()) // "{1 9 144}"
	y.Add(9)
	y.Add(42)
	fmt.Println(y.Len())
	fmt.Println(y.String()) // "{9 42}"
	x.UnionWith(&y)
	x.Add(123)
	fmt.Println(x.String())           // "{1 9 42 123 144}"
	fmt.Println(x.Has(9), x.Has(123)) //
	x.Remove(123)
	fmt.Println(x.String()) // "{1 9 42 144}"
	c := x.Copy()
	fmt.Println(c.String()) // "{1 9 42 144}"
	x.Clear()
	fmt.Println(x.String()) // "{}"
	s := []int{1, 2, 3, 4, 5}
	x.AddSlice(s)
	for _, v := range x.Elements() {
		fmt.Println(v)
	}
}
