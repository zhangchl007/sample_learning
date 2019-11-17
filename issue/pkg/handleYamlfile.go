package github
import (
    "encoding/json"
    "fmt"
    "os"
    "io/ioutil"
    "time"
    "net/http"
     y2j "github.com/ghodss/yaml"
    "github.com/spf13/viper"
    yaml "gopkg.in/yaml.v2"
)

func (yamlfile *Issueyamlfile) UpdateIssueyaml(Title, Body, State string, Locked bool, Assignees, Labels *[]string) (*[]byte, string, bool){
    tmpfile := "/tmp/a.txt"
    IssueTemplate :="issue_template"
    IssueyamlPath  :="src/sample_learning/issue/example/"
    TemplateFile := IssueyamlPath + IssueTemplate + ".yaml"
    viper.SetConfigName(IssueTemplate)
    viper.AddConfigPath(IssueyamlPath)
    err := viper.ReadInConfig()
    if err != nil {
        log.Fatal(err)
    }

    err = viper.Unmarshal(&yamlfile)
    if err != nil{
       log.Fatalf("unable to decode into struct, %v", err)
    }
   //viper.Set(yamlfile.Title, Title)
   yamlfile.Title = Title
   yamlfile.Body  = Body
   yamlfile.State = State
   yamlfile.Locked = Locked
   yamlfile.Assignees = *Assignees
   yamlfile.Labels = *Labels

   //viper.WriteConfig()

   // encode yaml again
   d, err := yaml.Marshal(&yamlfile)
   if err != nil {
      log.Fatalf("error: %v", err)
   }
   WriteToFile(tmpfile, d)
   MoveFile(tmpfile, TemplateFile)
   fmt.Println("The yaml file of issue template had been created succesfully!")
   //yaml to json for issue creation
   y2 := []byte(string(d))
   j2, err := y2j.YAMLToJSON(y2)
   if err != nil {
       fmt.Println("err: %v\n", err)
    }
   //fmt.Println(string(j2))

   return &j2, State, Locked
}

// generate the tempfile for yamlfile
func WriteToFile(f string ,d []byte)  {
    err := ioutil.WriteFile(f, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

//replace yamlfile file
func MoveFile(from,to string) {
	err := os.Rename(from, to)
	if err != nil {
		log.Fatal(err)
	}
}

