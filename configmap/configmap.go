package main

import (
    "flag"
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
const (
defaultredisConf = `bind=“0.0.0.0”
protected-mode=“yes”
tcp-backlog="511"
timeout="0"
tcp-keepalive="300"
daemonize="no"
supervised="no"
pidfile="/var/run/redis.pid"
cluster-enabled="yes"
cluster-node-timeout="5000"
cluster-require-full-coverage="no"
cluster-migration-barrier="1"
cluster-config-file="/data/nodes.conf"
save="900 1"
save="300 10"
save="60 10000"
appendonly="yes"
appendfilename="appendonly.aof"
dir="/data"`
)

var kubeconfig *string

// intailizing variable
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

    //configmap data
    configMapData := make(map[string]string, 0)
    configMapData["redisconf.yaml"] = defaultredisConf
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
   if _, err := clientset.CoreV1().ConfigMaps("redis-community1").Get(context.TODO(), "redisconfig1",  metav1.GetOptions{}); errors.IsNotFound(err) {
       cm, _ = clientset.CoreV1().ConfigMaps("redis-community1").Create(context.TODO(), &configMap, metav1.CreateOptions{})
   } else {
       cm, _ = clientset.CoreV1().ConfigMaps("redis-community1").Update(context.TODO(), &configMap, metav1.UpdateOptions{})
   }

   fmt.Println(cm)

}
//homeDir for get kubeconfig default profile
func homeDir() string {
    if h := os.Getenv("HOME"); h != "" {
        return h
    }
    return os.Getenv("USERPROFILE") // windows
}
