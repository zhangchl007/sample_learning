package main
import (
    "fmt"
)

func main() {

    fmt.Printf("%T\n", f) // "func(...int)"
    fmt.Printf("%T\n", g) // "func([]int)")

}
func f(...int) {}
func g([]int) {}


