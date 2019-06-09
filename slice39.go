package main
import (
    "fmt"
)

func main() {
    a :=[...]int{0, 1, 1, 2, 2, 3, 2, 4, 5, 6, 7, 8}
    //p := &a
    fmt.Println(rm_repeat(a[:]))

}
func rm_repeat(s []int)[]int {
    result := make([]int,0,len(s))
    temp := map[int]struct{}{}
    for _, v := range s {
        if _,ok := temp[v]; !ok {
        temp[v]= struct{}{}
        result =append(result, v)
       }
    }
    return result
}

