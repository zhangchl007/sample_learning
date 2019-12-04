package main

import (
	"bufio"
	"fmt"
	"os"
)
func main(){

    count := 0
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		count++
    }
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	fmt.Printf("%d\n", count)
}
