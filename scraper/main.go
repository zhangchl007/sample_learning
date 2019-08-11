package main

import(
    "os"
    "log"
    "fmt"
    "bufio"
    "strings"
    "strconv"
    "sample_learning/scraper/github"
)


func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Enter the filter number of Month!")
    text, _ := reader.ReadString('\n')
    s:= strings.Split(text,"\n")
    mgap, _ := strconv.Atoi(s[0])
    createdtime := github.Timefilter(mgap)
    created := "created:" + createdtime
    os.Args = append(os.Args, created)
    //fmt.Println(os.Args)
    result, err := github.SearchIssues(os.Args[1:])
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%d issue:\n", result.TotalCount)
    for _, item := range result.Items {
        fmt.Printf("#%-5d %9.7s %9.9s %.55s\n", item.Number,item.CreateAt, item.User.Login, item.Title)
    }
}


