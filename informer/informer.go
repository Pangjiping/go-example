package informer

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
)

// client-go中的informer使用

// getK8sClientByKubeConfig 从本地kubeconfig文件中初始化k8sClient [grpc]
// kubeconfig地址: ~/.kube/config
func getK8sClientByKubeConfig() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return k8sClient, nil
}

func ExecInformer() error {
	k8sClient, err := getK8sClientByKubeConfig()
	if err != nil {
		return err
	}

	sharedInformerFactory := informers.NewSharedInformerFactory(k8sClient, time.Minute*10)

	stopChan := make(chan struct{})
	sharedInformerFactory.Start(stopChan)

	podLister := sharedInformerFactory.Core().V1().Pods().Lister()
	pods, err := podLister.List(labels.Nothing())
	if err != nil {
		return err
	}
	fmt.Println(pods)

	//pod, err := podLister.Pods("kube-system").Get("kube-dns")
	//if err != nil {
	//	return err
	//}
	//fmt.Println(pod)

	list, err := podLister.Pods("kube-system").List(labels.Nothing())
	if err != nil {
		return err
	}
	fmt.Println(list)

	return nil
}
