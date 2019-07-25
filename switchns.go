package main

import (
    "flag"
    "fmt"
    "os"
    "io/ioutil"
    "path/filepath"
    "log"
    yaml "gopkg.in/yaml.v2"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

// kubeconfig struct
type Data struct {
  Apiv1 string `yaml:"apiVersion"`
  Clusters []Clusters `yaml:"clusters"`
  Contexts []Contexts `yaml:"contexts"`
  Currentctx string `yaml:"current-context"`
  Kind string `yaml:"kind"`
  Preferences Preferences `yaml:"preferences"`
  Users []Users `yaml:"users"`
}
type Clusters struct {
  Cluster Cluster `yaml:"cluster"`
  Name  string `yaml:"name"`
}
type Cluster struct {
  CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
  Server  string `yaml:"server"`
}
type Users struct {
  Name  string `yaml:"name"`
  Utoken Utoken `yaml:"user"`
}
type Utoken struct {
  Token string `yaml:"token"`
  Cscertdata string `yaml:"client-certificate-data"`
  Cskeydata string `yaml:"client-key-data"`
}
type Contexts struct {
  Context Context `yaml:"context"`
  Name    string  `yaml:"name"`
}
type Context struct {
  Cluster string `yaml:"cluster"`
  Namespace  string `yaml:"namespace"`
  User  string `yaml:"user"`
}
type Preferences struct {
    Colors bool `yaml:"colors,omitempty"`
}

func main() {
    var tmpfile string
    var s []string
    tmpfile = "/tmp/a.txt"

    if len(os.Args) == 1 {
        cns :=Getcontextns()
        fmt.Printf("The Current Namespace is %s\n",cns)
	}else if len(os.Args) == 2 {
        fmt.Println("switchns               # Get the current namespace!")
        fmt.Println("switchns -n namespace  # switch the namespace in k8s!")
    }else {
        s = Checkallns(os.Args[2], tmpfile, s)
        //fmt.Printf("All namespaces list below\n")
        //for _, ns := range s {
        //    fmt.Println(ns)
        //}
   }
}

func homeDir() string {
    if h := os.Getenv("HOME"); h != "" {
        return h
    }
    return os.Getenv("USERPROFILE") // windows
}

func Checkallns(name, tmpfile string,s []string)[]string{
    var kubeconfig *string
    var fpath string
    i := 1

    flag.StringVar(&name, "n", "", "change the namespace in current context")

    if home := homeDir(); home != "" {
    kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    fpath = filepath.Join(home, ".kube", "config")
    } else {
         kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    }

   flag.Parse()
   if name == "" {
      name = os.Args[2]
   }
    // uses the current context in kubeconfig
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

    if err != nil {
        panic(err.Error())
    }
    // creates the clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }
    allns, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
    if err != nil {
        log.Fatalln("failed to get Namespaces:", err)
    }
   // Check if the namespace exists in Kubernetes
    for _, ns := range allns.Items {
        m := ns.GetName()
        s =append(s, m)
        if m == name {
            i--
        }
    }
    if i == 0 {
            Setcontextns(fpath, name, tmpfile)
    }else {
        fmt.Println("The kubernetes namespaces doesn't exists!")
    }
    return s
}

func Getcontextns() string{
   // Instantiate loader for kubeconfig file.
   kubecfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
       clientcmd.NewDefaultClientConfigLoadingRules(),
       &clientcmd.ConfigOverrides{},
   )
   cns, _, err := kubecfg.Namespace()
      if err != nil {
         panic(err)
      }
    // restconfig := kubecfg.ConfigAccess()
    //fmt.Println(restconfig)
     return cns
}

func Setcontextns(filename, name, tmpfile string) {

   txtlines, err := ioutil.ReadFile(filename)
   if err != nil {
     panic(err)
   }

// struc Unmasrshal
   var t Data
   err = yaml.Unmarshal(txtlines, &t)
   if err != nil {
     panic(err)
   }
// decode yaml file , then change namespace in the current context
   cls :=  t.Currentctx
   allcls := t.Contexts
   for k, v := range allcls {
       if v.Name == cls {
           p := &v.Context
           p.Setns(name)
           allcls[k] = v
       }
   }
// encode yaml again
   d, err := yaml.Marshal(&t)
   if err != nil {
      log.Fatalf("error: %v", err)
   }
// write yaml to tmpfile and overwrite the kubeconfig file
   WriteToFile(tmpfile, d)
   move(tmpfile, filename)
   fmt.Printf("The current namespace is:  %s\n", name)
}
// generate the tempfile for kubeconfig
func WriteToFile(f string ,d []byte)  {
    err := ioutil.WriteFile(f, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

//replace kubeconfig file
func move(from,to string) {
	err := os.Rename(from, to)
	if err != nil {
		log.Fatal(err)
	}
}

//change the namespace property for Context struct
func (ctx *Context) Setns(name string) string{
    if len(ctx.Namespace) > 0 {
        ctx.Namespace = name
    }else {
        fmt.Println("namespaces is ")
    }
   return ctx.Namespace
}
