package main

import (
	"fmt"
)

var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"database":              {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating system"},
	"operating system":      {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	for k,v  := range topoSort(prereqs) {
        fmt.Printf("%s:\t%d\n", k, v)
	}
}

func topoSort(m map[string][]string) map[string]int {
	seen := make(map[string]bool)
	course := make(map[string]int)
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, v := range items {
			if !seen[v] {
				seen[v] = true
				visitAll(m[v])
				if yv, ok := m[v]; ok {
					course[v]++
					for _, val := range yv {
						course[val]++
					}
				}
			}
		}
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	visitAll(keys)
	return course
}
