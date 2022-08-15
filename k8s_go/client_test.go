package k8s_go

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InitClientGo(t *testing.T) {
	testcases := []map[string]interface{}{
		{
			"description":  "grpc client test",
			"type":         "grpc",
			"wantErr":      false,
			"clientNotNil": true,
		},
		{
			"description":  "rest client test",
			"type":         "restful",
			"wantErr":      false,
			"clientNotNil": true,
		},
	}

	var client interface{}
	var err error

	for _, tc := range testcases {
		t.Log(tc["description"])
		if tc["type"].(string) == "grpc" {
			client, err = getK8sClientByKubeConfig()
		} else {
			client, err = getRestClient()
		}
		assert.Equal(t, tc["wantErr"].(bool), err != nil)
		assert.Equal(t, tc["clientNotNil"].(bool), client != nil)
	}
}

func Test_AllOpteration(t *testing.T) {

	// init grpc client
	t.Log("test for init grpc client ...")
	k8sClient, err := getK8sClientByKubeConfig()
	assert.Nil(t, err)
	assert.NotNil(t, k8sClient)

	// init restful client
	t.Log("test for init restful client ...")
	restClient, err := getRestClient()
	assert.Nil(t, err)
	assert.NotNil(t, restClient)

	// get nodes
	t.Log("test for get nodes from apiserver ...")
	nodesOne, err := getNodesFromApiServer()
	assert.Nil(t, err)
	assert.NotNil(t, nodesOne)
	assert.Equal(t, true, len(nodesOne) >= 1)

	// get nodes
	t.Log("test for get nodes from apiserver with optimizing...")
	nodesTwo, err := describeNodesFromApiServer()
	assert.Nil(t, err)
	assert.NotNil(t, nodesTwo)
	assert.Equal(t, true, len(nodesTwo) >= 1)

	// get nodeSize
	t.Log("test for get nodeSize from apiserver ...")
	nodeSize, err := getNodeSizeFromApiServer()
	assert.Nil(t, err)
	assert.NotNil(t, nodeSize)
	assert.Equal(t, true, nodeSize >= 1)

	// get pods
	t.Log("test for get pods from apiserver ...")
	pods, err := getPodsFromApiServer()
	assert.Nil(t, err)
	assert.NotNil(t, pods)
	assert.Equal(t, true, len(pods) >= 1)

	// get podSize
	t.Log("test for get podSize from apiserver ...")
	podSize, err := getPodSizeFromApiServer()
	assert.Nil(t, err)
	assert.NotNil(t, podSize)
	assert.Equal(t, true, podSize)
}

// failed
func Test_getPodFromApiServerByUid(t *testing.T) {
	pod, err := getPodFromApiServerByUid()
	assert.Nil(t, err)
	assert.NotNil(t, pod)

	fmt.Printf("get pod uid: %v\n", pod.GetUID())
}
