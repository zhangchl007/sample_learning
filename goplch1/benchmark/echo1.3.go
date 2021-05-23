package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	StrCount()
	StrJoin()
	fmt.Println("works")
}
func StrCount() {
	var s, sep = "", ""
	for j := 1; j < len(os.Args); j++ {
		s += sep + os.Args[j]
		sep = ""
	}
	fmt.Println(s)
}

func StrJoin() {
	var s string
	for _, arg := range os.Args[1:] {
		s += arg
	}
	fmt.Println(strings.Join(os.Args[1:], ""))
}
