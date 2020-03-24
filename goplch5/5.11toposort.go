package main
import (
    "fmt"
    "sort"
)

var prereqs = map[string][]string{
    "algorithms": {"data structures"},
    "calculus": {"linear algebra"},
    "linear algebra": {"calculus"},
    "compilers": {
        "data structures",
        "formal languages",
        "computer organization",
    },
    "data structures": {"discrete math"},
    "database": {"data structures"},
    "discrete math": {"intro to programming"},
    "formal languages": {"discrete math"},
    "networks": {"operating system"},
    "operating system": {"data structures", "computer organization"},
    "programming languages": {"data structures", "computer organization"},
}

func index(k string, s []string) {
    for _, v := range s {
        if k == v {
            p := []string{k}
            return &p
        }
    }
}

func main(){
    for i, course := range topoSort(prereqs) {
        fmt.Printf("%d:\t%s\n", i+1, course)
    }
}

func topoSort(m map[string][]string) []string{
    var order []string
    seen := make(map[string]bool)
    var visitAll func index(k string, s []string)
    visitAll = func index(k string, s []string) {

        for _, item := range items {
            if !seen[item] {
                seen[item] = true
                visitAll(m[item])
                order = append(order, item)
            }
        }
    }
    var keys []string
    for key :=range m {
        keys = append(keys,key)
    }
    sort.Strings(keys)
    visitAll(keys)
    return order
}
