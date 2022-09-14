package e2e

import (
	"fmt"
	roa "github.com/alibabacloud-go/tea-roa/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"net/http"
)

const (
	SDK_VERSION                = "2015-12-15"
	AUTH_TYPE                  = "AK"
	SDK_COMMON_RESPONSE_FROMAT = "[OPENAPI CALL SUCCESS!!!]\nResponse: %v\n[SDK alibabacloud-go]"
)

type E2eClient struct {
	roaClient  *roa.Client
	clusterId  string
	nodepoolId string
}

func NewE2eClient(accessKeyId string, accessKeySecret string) (*E2eClient, error) {
	if accessKeyId == "" || accessKeySecret == "" {
		return nil, fmt.Errorf("missing ak")
	}

	csClient, err := createClient(tea.String(accessKeyId), tea.String(accessKeySecret))
	if err != nil {
		return nil, err
	}
	return &E2eClient{roaClient: csClient}, nil
}

func createClient(accessKeyId *string, accessKeySecret *string) (*roa.Client, error) {
	config := &roa.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}

	config.Endpoint = tea.String("cs.aliyuncs.com")
	return roa.NewClient(config)
}

func (c *E2eClient) RegisterCluster(clusterId string, nodepoolId string) {
	if clusterId != "" {
		c.clusterId = clusterId
	}
	if nodepoolId != "" {
		c.nodepoolId = nodepoolId
	}
}

func (c *E2eClient) RemoveInstanceForNodePool(instances []string) error {
	action := "RemoveNodes"
	resp, err := c.roaClient.DoRequestWithAction(
		tea.String(action),
		tea.String(SDK_VERSION),
		nil,
		tea.String(http.MethodDelete),
		tea.String(AUTH_TYPE),
		tea.String(fmt.Sprintf("/clusters/%s/nodepools/%s/nodes", c.clusterId, c.nodepoolId)),
		nil,
		nil,
		instances,
		&util.RuntimeOptions{},
	)
	if err != nil {
		return err
	}

	fmt.Printf(SDK_COMMON_RESPONSE_FROMAT, resp)
	return nil
}
