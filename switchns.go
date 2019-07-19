package main

import (
    "flag"
    "bufio"
    "regexp"
    "fmt"
    "os"
    "path/filepath"
    "log"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

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
   // fmt.Printf("The current Namespaces is ")

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
     return cns
}

func Setcontextns(filename, name, tmpfile string) {
    //tmpfile := "/tmp/a.txt"
    file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	for _, line := range txtlines {
        re := regexp.MustCompile(`namespace: (\w+)`)
        if fields :=re.FindStringSubmatch(line); fields != nil {
            str := "namespace: " + name
            line = re.ReplaceAllString(line, str)
        }
    //  fmt.Println(line)
      WriteToFile(tmpfile, line)
    }
    move(tmpfile, filename)
    fmt.Printf("Successfully switch the namespaces: %s!\n",name)
}

func WriteToFile(filename, line string)  {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
     	log.Println(err)
    }
    defer f.Close()
    //fmt.Println(line)
    if _, err := f.WriteString(line + "\n"); err != nil {
    	log.Println(err)
    }
}

func move(from,to string) {
	err := os.Rename(from, to)
	if err != nil {
		log.Fatal(err)
	}
}
