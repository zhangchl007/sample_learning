package main
import (
    "fmt"
    "os"
    "log"
    "samplelearning.com/goplch5/links"
)

func main() {
    //crawl the web breadthFirst,
    //starting from the command-line arguments.
    breadthFirst(crawl, os.Args[1:])
}

    //breadthFirst calls f for each itme in worklist.
    // Any items returned by f are added to the worklist.
    //f is called at most once for each item.
func breadthFirst(f func(item string) []string, worklist []string) {
     seen := make(map[string]bool)
     for len(worklist) > 0 {
         items := worklist
         worklist = nil
         for _, item := range items {
             if !seen[item] {
                 seen[item] = true
                 worklist = append(worklist, f(item)...)
             }
         }

     }
}

func crawl(url string) []string {
    fmt.Println(url)
    list, err := links.Extract(url)
    if err != nil {
        log.Print(err)
    }
    return list
}
