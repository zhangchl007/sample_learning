package main
import (
    "fmt"
)

func main() {
    a :=[...]int{0, 1, 1, 2, 2, 3, 2, 4, 5, 6, 7, 8}
    //p := &a
    fmt.Println(rm_repeat(a[:]))
    //fmt.Println(a)

}
func rm_repeat(s []int)[]int {
    for i,j :=0,1 ; i<len(s)-1; i,j =i+1,j+1{
        if s[i] == s[j] {
            //s = s[:len(s)-1]
            s = append(s[:i], s[j:]...)
        }
    }
    return s
}

