package k8s_go

import (
	"context"
	"flag"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

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

// getRestClient 从本地kubeconfig文件初始化restfulClient [http]
// kubeconfig地址: ~/.kube/config
func getRestClient() (*rest.RESTClient, error) {
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

	config.APIPath = "api"
	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, err
	}
	return restClient, nil
}

// getNodeFromApiServer 从apiserver获取node信息
// 使用grpc，获取node列表
// 对于大集群，如果不对etcd做优化的话，绝对不应该访问etcd，前提是不要求强一致性
func getNodesFromApiServer() ([]v1.Node, error) {
	k8sClient, err := getK8sClientByKubeConfig()
	if err != nil {
		return nil, err
	}
	nodeList, err := k8sClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		LabelSelector:        "",             // 节点labelSelector,用"key=val"或者"key!=val"和","拼接成的字符串
		FieldSelector:        "",             // 和label同理
		ResourceVersion:      "0",            // 0 表示从apiserver获取最新缓存，不会打到etcd
		ResourceVersionMatch: "NotOlderThan", // Match和resourceVersion的搭配见 https://kubernetes.io/zh-cn/docs/reference/using-api/api-concepts/
		//TimeoutSeconds:       nil,            // 访问超时时间
		//Limit:                1,              // 获取全量nodes没啥用,在获取nodeSize时有大用
	})
	if err != nil {
		return nil, err
	}
	return nodeList.Items, nil
}

// getNodesWithPagination 从apiserver全量list node的一种优化
// 具体使用见 describeNodesFromApiServer
func getNodesWithPagination(nextToken string, maxResult int64) (*v1.NodeList, error) {
	k8sClient, err := getK8sClientByKubeConfig()
	if err != nil {
		return nil, err
	}
	nodeList, err := k8sClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		Continue:             nextToken,
		Limit:                maxResult,
		ResourceVersion:      "0",
		ResourceVersionMatch: "NotOlderThan",
	})
	if err != nil {
		return nil, err
	}
	return nodeList, nil
}

// describeNodesFromApiServer 从apiserver全量list node
// 针对大集群的一种优化手段
func describeNodesFromApiServer() ([]v1.Node, error) {
	var nextToken = ""
	var step int64 = 100
	var totalCount int64
	nodeList := make([]v1.Node, 0)

	for {
		nodes, err := getNodesWithPagination(nextToken, step)
		if err != nil {
			return nil, err
		}
		nodeList = append(nodeList, nodes.Items...)
		totalCount += int64(len(nodes.Items))
		if nodes.Continue == "" {
			break
		}
		nextToken = nodes.Continue
	}
	return nodeList, nil
}

// getNodeSizeFromApiServer 只获取nodeSize信息，不需要元数据
// 使用limit和remainingCount做优化，返回最少的数据
func getNodeSizeFromApiServer() (int64, error) {
	k8sClient, err := getK8sClientByKubeConfig()
	if err != nil {
		return 0, err
	}
	nodeList, err := k8sClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		Limit:                1,
		ResourceVersion:      "0",
		ResourceVersionMatch: "NotOlderThan",
	})
	if err != nil {
		return 0, err
	}

	size := *nodeList.RemainingItemCount + int64(len(nodeList.Items))
	// 注意nodeList.Size()不是node的size, 而是返回的nodeList数据的大小, 多少字节
	return size, nil
}

// getPodsFromApiServer 获取pods列表
// 有很多可选择的优化项, 某个namespace、某个节点、某些标签
// pod不像node, 尽量不要全量list
func getPodsFromApiServer() ([]v1.Pod, error) {
	k8sClient, err := getK8sClientByKubeConfig()
	if err != nil {
		return nil, err
	}

	// 对于labelSelector和fieldSelector如果不知道可以用kubectl看一下具体的字段含义
	podList, err := k8sClient.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{
		LabelSelector:        "",
		FieldSelector:        "spec.nodeName=nodeName", // 这个list某个节点上的
		ResourceVersion:      "0",
		ResourceVersionMatch: "NotOlderThan",
	})
	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}

// getPodSizeFromApiServer 获取podSize
// 使用limit和remainingCount优化
func getPodSizeFromApiServer() (int64, error) {
	k8sClient, err := getK8sClientByKubeConfig()
	if err != nil {
		return 0, err
	}

	podList, err := k8sClient.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{
		ResourceVersion:      "0",
		ResourceVersionMatch: "NotOlderThan",
		Limit:                1,
	})
	if err != nil {
		return 0, err
	}
	return *podList.RemainingItemCount + int64(len(podList.Items)), nil
}

// getPodFromApiServerByUid 匹配单个pod只能通过pod_name来做
func getPodFromApiServerByUid() (*v1.Pod, error) {
	k8sClient, err := getK8sClientByKubeConfig()
	if err != nil {
		return nil, err
	}
	podList, err := k8sClient.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{
		ResourceVersion:      "0",
		ResourceVersionMatch: "NotOlderThan",
		FieldSelector:        fields.OneTermEqualSelector("metadata.uid", "1a52ef95-50aa-41be-bfd1-3dc6b7f54568").String(),
		Limit:                1,
	})
	if err != nil {
		return nil, err
	}
	return &podList.Items[0], nil
}
