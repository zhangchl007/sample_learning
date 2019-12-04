package main
import (
    "fmt"
    "time"
    "log"
    "os"
    "net/http"
)

func main() {
    for _, url := range os.Args[1:] {
        if err := WaitForServer(url); err != nil {
            //fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
            log.Fatalf("Site is down: %v\n", err)
            os.Exit(1)
        }
    }
}
func WaitForServer(url string) error{
    const timeout = 1 * time.Minute
    deadline := time.Now().Add(timeout)
    for tries := 0; time.Now().Before(deadline); tries++ {
        _,err := http.Head(url)
        if err == nil {
            fmt.Println("Reach it!")
            return nil //success
        }
        log.Printf("server not reponding (%s); retrying...",err)
        time.Sleep(time.Second << uint(tries)) //exponential back-off
    }
    return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}
