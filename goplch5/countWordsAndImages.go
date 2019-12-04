package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)
func main(){
    counts := make(map[rune]int)
    dcounts := make(map[rune]int)
    ocounts := make(map[rune]int)
    invalid :=0
    in := bufio.NewReader(os.Stdin)
    for {
        r, n, err := in.ReadRune()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Fprintf(os.Stderr, "charcount: %v\n",err)
            os.Exit(1)
        }
        if r == unicode.ReplacementChar && n == 1 {
            invalid++
            continue
        }
        if unicode.IsLetter(r) {
            counts[r]++
        }else if unicode.IsNumber(r) {
            dcounts[r]++
        }else {
            ocounts[r]++
        }
    }
    fmt.Printf("letter\tcount\n")
    for c, n := range counts{
        fmt.Printf("%q\t%d\n",c, n)
    }
    fmt.Printf("digit\tcount\n")
    for c, n := range dcounts{
        fmt.Printf("%q\t%d\n",c, n)
    }
    fmt.Printf("other\tcount\n")
    for c, n := range ocounts{
        fmt.Printf("%q\t%d\n",c, n)
    }
    if invalid > 0 {
        fmt.Printf("\n%d invalid UTF-8 chracters\n", invalid)
    }
}


