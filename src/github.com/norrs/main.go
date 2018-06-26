package main

import (
        "fmt"
        "time"
        "flag"
        "os"
        "path/filepath"
        "github.com/golang/glog"
        "k8s.io/api/core/v1"
        "k8s.io/client-go/tools/clientcmd"
        "k8s.io/client-go/kubernetes"
        "k8s.io/apimachinery/pkg/fields"

        extensions "k8s.io/api/extensions/v1beta1"
        "k8s.io/client-go/tools/cache"
        // Only required to authenticate against GKE clusters
        _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
    var kubeconfig *string
    if home := homeDir(); home != "" {
	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    } else {
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    }
    flag.Parse()

    // use the current context in kubeconfig
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        glog.Errorln(err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        glog.Errorln(err)
    }
clientset.CoreV1().RESTClient()
    watchlist := cache.NewListWatchFromClient(
        clientset.ExtensionsV1beta1().RESTClient(),
        string(v1.ResourceServices),
        v1.NamespaceAll,
        fields.Everything(),
    )
    _, controller := cache.NewInformer( // also take a look at NewSharedIndexInformer
        watchlist,
        &extensions.Ingress{},
        0, //Duration is int64
        cache.ResourceEventHandlerFuncs{
            AddFunc: func(obj interface{}) {
                fmt.Printf("ingress added: %s \n", obj)
            },
            DeleteFunc: func(obj interface{}) {
                fmt.Printf("ingress deleted: %s \n", obj)
            },
            UpdateFunc: func(oldObj, newObj interface{}) {
                fmt.Printf("ingress changed \n")
            },
         },
     )
         // I found it in k8s scheduler module. Maybe it's help if you 
    // interested in.
     // serviceInformer := 
   // cache.NewSharedIndexInformer(watchlist, 
   //  &v1.Service{}, 0, cache.Indexers{
     //     cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
     // })
     // go serviceInformer.Run(stop)
    stop := make(chan struct{})
    defer close(stop)
    go controller.Run(stop)
    for {
        time.Sleep(time.Second)
    }
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
