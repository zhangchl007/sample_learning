package main

import (
    "flag"
    "strings"
    "context"
    "path/filepath"
    "os"
    "fmt"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

// Default redisconfig
var defaultredisConf = `save 900 1
save 300 10
save 60 10000
appendonly yes
cluster-config-file /data/nodes.conf
appendfilename "appendonly.aof"
cluster-node-timeout 5000
cluster-migration-barrier 1
cluster-enabled yes
cluster-require-full-coverage no
dir /data`

var kubeconfig *string
var initredisConf string

func init() {

    if home := homeDir(); home != "" {
    kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    } else {
         kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    }

   flag.Parse()
}


func main() {

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
   // generate redisConf
    generateredisConf()
    configMapData := make(map[string]string, 0)
    configMapData["redisconf.yaml"] = fmt.Sprintf(`%s`, initredisConf)
    configMap := corev1.ConfigMap{
    TypeMeta: metav1.TypeMeta{
    Kind:       "ConfigMap",
    APIVersion: "v1",
    },
    ObjectMeta: metav1.ObjectMeta{
      Name:      "redisconfig",
      Namespace: "redis-community1",
    },
    Data: configMapData,
   }

   var cm *corev1.ConfigMap
   if _, err := clientset.CoreV1().ConfigMaps("redis-community1").Get(context.TODO(), "redisconfig",  metav1.GetOptions{}); errors.IsNotFound(err) {
       cm, _ = clientset.CoreV1().ConfigMaps("redis-community1").Create(context.TODO(), &configMap, metav1.CreateOptions{})
   } else {
       cm, _ = clientset.CoreV1().ConfigMaps("redis-community1").Update(context.TODO(), &configMap, metav1.UpdateOptions{})
   }

   fmt.Println(cm)
   //defaultredisConf := generateConf()

}

func homeDir() string {
    if h := os.Getenv("HOME"); h != "" {
        return h
    }
    return os.Getenv("USERPROFILE") // windows
}
func generateredisConf() {
      var initConf = make(map[string][]string, 0)
      var defConf = make([]string,0)
      initConf["cluster-enabled"] = []string{"yes"}
      initConf["cluster-node-timeout"] = []string{"5000"}
      initConf["cluster-require-full-coverage"] = []string{"no"}
      initConf["cluster-migration-barrier"] = []string{"1"}
      initConf["save"] = []string{"900 1", "300 10", "60 10000"}
      initConf["appendonly"] = []string{"yes"}
      initConf["cluster-config-file"] = []string{"/data/nodes.conf"}
      initConf["dir"] = []string{"/data"}
      initConf["appendfilename"] = []string{"\"appendonly.aof\""}

     for k,v := range initConf{
         if len(v) > 1 {
             for i := range v{
                 line1 := k + " " + v[i]
                 defConf = append(defConf,line1)
                 }
             }else {
                 line2 := k + " " + v[0]
                 defConf = append(defConf,line2)
         }
     }
     fmt.Println(initConf)
     initredisConf = strings.Join(defConf, "\n")
//     fmt.Println(initredisConf)
}
