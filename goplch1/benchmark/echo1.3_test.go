package main

import (
	"testing"
)

//var args = []string{"hi", "there", "buddy", "boy", "5", "6", "7", "8", "9"}

func BenchmarkStrCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StrCount()
	}
}

func BenchmarkStrJoinStrJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StrJoin()
	}
}
