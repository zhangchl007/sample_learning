package github
import (
    "encoding/json"
    "strings"
    "fmt"
    "time"
    "net/http"
    "net/url"
)
func SearchIssues(terms []string) (*IssuesSearchResult,error) {
    q := url.QueryEscape(strings.Join(terms, " "))
    fmt.Println(IssuesURL + "?q" +q)
    resp, err := http.Get(IssuesURL + "?q=" + q)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        resp.Body.Close()
        return nil, fmt.Errorf("search query failed: %s", resp.Status)
    }
    var result IssuesSearchResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        resp.Body.Close()
        return nil, err
    }
    resp.Body.Close()
    return &result, nil
}

func Timefilter(mgap int) (createdtime string) {
    end:= time.Now()
    start := end.AddDate(0, -mgap , 0)
    createdtime = start.Format("2006-01")
    return createdtime
}
