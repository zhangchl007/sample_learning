package main
import (
    "fmt"
    "bytes"
    "os"
    "golang.org/x/net/html"
)

func main() {
    doc, err := html.Parse(os.Stdin)
    if err != nil {
        fmt.Fprintf(os.Stderr, "outline1: %v\n", err)
        os.Exit(1)
    }
    output(doc)
}

func collecttxt(buf *bytes.Buffer, n *html.Node) {
    if n.Type == html.TextNode {
        buf.WriteString(n.Data)
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        collecttxt(buf, c)
    }
}

func output(n *html.Node){
    if n.Type == html.ElementNode && n.Data == "a" {
        for _, a := range n.Attr {
            if a.Key == "href" {
                text :=&bytes.Buffer{}
                collecttxt(text, n)
                fmt.Println(text)
            }
        }
    }
    for c := n.FirstChild; c !=nil; c = c.NextSibling {
        output(c)

    }
}
