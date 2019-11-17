package github
import (
    "encoding/json"
    "strings"
    "fmt"
    "time"
    "net/http"
    "net/url"
)

func GetIssues(PersonalAccessToken, Userid string) (*[]Issues, error) {
    url := IssuesURL + Userid + "/" + Repo + "/" + "issues?state=all"
    req, err := http.NewRequest(MethodGet, url,  nil)
    req.Header.Add("Authorization", "Bearer " + PersonalAccessToken)
    //req.Header.Add("User-Agent","Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36")
    //req.Header.Add("Content-Type","application/json")
    client := &http.Client{}
    //Send req using http Client
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        resp.Body.Close()
        return nil, fmt.Errorf("search query failed: %s", resp.Status)
    }

    var result []Issues

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        resp.Body.Close()
        return nil, err
    }

    resp.Body.Close()
    return &result, nil
}

