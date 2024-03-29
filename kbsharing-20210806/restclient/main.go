package main

import (
	"flag"
	"path/filepath"
	"fmt"

    "context"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"

	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	namespace := "default"

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
    config.APIPath = "api"
        config.GroupVersion = &corev1.SchemeGroupVersion
        config.NegotiatedSerializer = scheme.Codecs

        fmt.Println("Init RESTClient.")

        // 定义RestClient，用于与k8s API server进行交互
        restClient, err := rest.RESTClientFor(config)
        if err != nil {
                panic(err)
        }

        fmt.Println("Get Pods in cluster.")

        // 获取pod列表。这里只会从namespace为"kube-system"中获取指定的资源(pods)
        result := &corev1.PodList{}
        if err := restClient.
                Get().
                Namespace(namespace).
                Resource("pods").
                VersionedParams(&metav1.ListOptions{Limit: 500}, scheme.ParameterCodec).
                Do(context.TODO()).
                Into(result); err != nil {
                panic(err)
            }

        fmt.Println("Print all listed pods.")

        // 打印所有获取到的pods资源，输出到标准输出
        for _, d := range result.Items {
                fmt.Printf("NAMESPACE: %v NAME: %v \t STATUS: %v \n", d.Namespace, d.Name, d.Status.Phase)
        }
}
