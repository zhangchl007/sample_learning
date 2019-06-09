package main
import (
    "fmt"
)

func main() {
    s := []int{5, 6, 7, 8, 9}
    i :=len(s)
    fmt.Println(i)
    fmt.Println(remove(s, 2))
    fmt.Println(rm_dup(s, 2))

}

func remove(slice []int, i int) []int {
    slice[i] = slice[len(slice)-1]
    return slice[:len(slice)-1]
}
func rm_dup(s []int, i int) []int {
    copy(s[i:],s[i+1:])
    return s[:len(s)-1]
}
