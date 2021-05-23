package main

import (
        "fmt"

        "k8s.io/apimachinery/pkg/labels"
        appsclientset "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/client-go/tools/clientcmd"
)

func main() {

         var (
        dcLabelKey       = "app"
        dcLabelValue      = "s2idemo"
        )
        dclabel := labels.SelectorFromSet(labels.Set(map[string]string{dcLabelKey: dcLabelValue}))
        listdcOptions := metav1.ListOptions{ LabelSelector: dclabel.String(),}
        // Instantiate loader for kubeconfig file.
        kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
                clientcmd.NewDefaultClientConfigLoadingRules(),
                &clientcmd.ConfigOverrides{},
        )

        // Determine the Namespace referenced by the current context in the
        // kubeconfig file.
        namespace, _, err := kubeconfig.Namespace()
        if err != nil {
                panic(err)
        }

        // Get a rest.Config from the kubeconfig file.  This will be passed into all
        // the client objects we create.
        restconfig, err := kubeconfig.ClientConfig()
        if err != nil {
                panic(err)
        }
        //create a openshift deploymentconfig client
        appsClient, err := appsclientset.NewForConfig(restconfig)
       if err != nil {
                panic(err)
       }

        err = appsClient.DeploymentConfigs(namespace).DeleteCollection(&metav1.DeleteOptions{}, listdcOptions)
        if err != nil {
                panic(err)
        }else {
               fmt.Printf("The DeploymentConfig: %s: had been deleted in namespace %s:\n", dcLabelValue,namespace)
        }

}
