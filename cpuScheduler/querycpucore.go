package main

import (
    "fmt"
    "runtime"
)

func main() {

    // NumCPU 返回当前可用的逻辑处理核心的数量
    fmt.Println(runtime.NumCPU())
}
