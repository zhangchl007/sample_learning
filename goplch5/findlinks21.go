package main
import (
    "fmt"
    "os"
    "golang.org/x/net/html"
)

func main() {
    doc, err := html.Parse(os.Stdin)
    if err != nil {
        fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
        os.Exit(1)
    }
    for k := range visit(nil, doc) {
        fmt.Println(k)
    }
}

/* func visit(links []string, n *html.Node)[]string {
    if n.Type == html.ElementNode && n.Data == "a" {
        for _, a := range n.Attr {
            if a.Key == "href" {
                links = append(links, a.Val)
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        links = visit(links, c)
    }
    return links
} */

func visit(elems map[string]int, n *html.Node ) map[string]int {
    elems = make(map[string]int)
    if n.Type == html.ElementNode && n.Data == "a" {
        for _, a := range n.Attr {
            fmt.Println(a.Key)
            elems[a.Key]++
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
          elems = visit(elems, c)
    }
    return elems
}
