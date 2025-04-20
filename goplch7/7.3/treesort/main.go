package main

// This program demonstrates the use of the treesort package
import (
	"fmt"

	"github.com/zhangchl007/sample_learning/goplch7/7.3/treesort"
)

func main() {
	// Example usage of treesort
	values := []int{5, 3, 8, 1, 4}
	fmt.Println("Before sorting:", values)
	treesort.Sort(values)
	fmt.Println("After sorting:", values)
}
