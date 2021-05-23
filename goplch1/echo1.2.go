package main

import (
	"fmt"
	"os"
)

func main() {

	for s, arg := range os.Args[1:] {
		//s += sep + arg
		//sep = ""
		fmt.Println(s, arg)
	}
	//fmt.Println(s)
	//fmt.Println(strings.Join(os.Args[1:], ""))
	//fmt.Println(os.Args[0], s)

}
