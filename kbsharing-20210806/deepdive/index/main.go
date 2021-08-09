package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    meta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/tools/cache"
	"fmt"
)


const (
    NamespaceIndexName = "namepace"

    NodeNameIndexName = "nodeName"
)


func NamespaceIndexFunc(obj interface{}) ([]string, error) {
    m, err := meta.Accessor(obj)
    if err != nil {
        return []string{""}, fmt.Errorf("object has no meta: %v", err)
    }
    return []string{m.GetNamespace()},nil
}

func NodeNameIndexFunc(obj interface{}) ([]string, error) {
    //pod, ok := obj{*corev1.Pod}
    pod, ok := obj.(*corev1.Pod)
    if !ok {
        return []string{},nil
    }
    return []string{pod.Spec.NodeName}, nil
}

func main() {
    index := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{
        NamespaceIndexName: NamespaceIndexFunc,
        NodeNameIndexName: NodeNameIndexFunc,
    })

    pod1 := &corev1.Pod {
        ObjectMeta: metav1.ObjectMeta {
            Name: "test-pod-1",
            Namespace: "default",
        },
    Spec: corev1.PodSpec{NodeName: "node1"},
    }
    pod2 := &corev1.Pod {
        ObjectMeta: metav1.ObjectMeta {
            Name: "test-pod-2",
            Namespace: "default",
        },
    Spec: corev1.PodSpec{NodeName: "node2"},
    }
    pod3 := &corev1.Pod {
        ObjectMeta: metav1.ObjectMeta {
            Name: "test-pod-3",
            Namespace: "kube-system",
        },
    Spec: corev1.PodSpec{NodeName: "node2"},
    }

    _ = index.Add(pod1)
    _ = index.Add(pod2)
    _ = index.Add(pod3)

    keys, err := index.IndexKeys(NamespaceIndexName, "default")
    if err != nil {
        panic(err)
    }
    for _, k := range keys {
        fmt.Println(k)
    }

    pods, err := index.ByIndex(NamespaceIndexName, "default")
    if err != nil {
        panic(err)
    }
    for _, pod := range pods {
        fmt.Println(pod.(*corev1.Pod).Name)
    }
    fmt.Println("=============================================")

    pods, err = index.ByIndex(NodeNameIndexName, "node2")

    if err != nil {
        panic(err)
    }

    for _, pod := range pods {
        fmt.Println(pod.(*corev1.Pod).Name)
    }
}
