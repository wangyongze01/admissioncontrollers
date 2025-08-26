package client

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
	"path/filepath"
	"sync"
)

type Client struct {
	Api        *kubernetes.Clientset
	MetricsApi *versioned.Clientset
}

var K8sClient *Client
var onceK8s sync.Once

// 原生创建client
func NewClientK8s() {
	onceK8s.Do(func() {
		var cfg *rest.Config

		useKubeConfig := os.Getenv("USE_KUBECONFIG")
		kubeConfigFilePath := os.Getenv("KUBECONFIG")

		if len(useKubeConfig) == 0 {
			// default to service account in cluster token
			c, err := rest.InClusterConfig()
			if err != nil {
				panic(err.Error())
			}
			cfg = c
		} else {
			//load from a kube config
			var kubeconfig string

			if kubeConfigFilePath == "" {
				if home := homedir.HomeDir(); home != "" {
					kubeconfig = filepath.Join(home, ".kube", "config")
				}
			} else {
				kubeconfig = kubeConfigFilePath
			}

			fmt.Println("kubeconfig: " + kubeconfig)

			c, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				panic(err.Error())
			}
			cfg = c
		}

		K8sClient = &Client{Api: nil, MetricsApi: nil}
		//K8sClient.Api, err = kubernetes.NewForConfig(cfg)
		config, err := versioned.NewForConfig(cfg)
		if err != nil {
			glog.Errorf("kubernetes client failed")
			panic(err.Error())
		}
		K8sClient.MetricsApi = config
		K8sClient.Api, err = kubernetes.NewForConfig(cfg)
		fmt.Println("k8s service success")
		if err != nil {
			glog.Errorf("kubernetes client failed")
			panic(err.Error())
		}
	})
}
