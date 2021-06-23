package apisixcli

import (
	"os"
	"strings"

	"github.com/apache/apisix-ingress-controller/pkg/apisix"
)

var (
	apisixBaseURL  = "http://127.0.0.1:9080/apisix/admin"
	apisixAdminKey = "edd1c9f034335f136f87ad84b625c8f1"
)

const (
	EnvAPISIXAdminHostAddress = "APISIX_ADMIN_ADDR"
	EnvAPISIXAdminToken       = "APISIX_ADMIN_TOKEN"
)

func CreateAPISIXCli() (apisix.Cluster, error) {
	apisixCli, err := apisix.NewClient()
	if err != nil {
		return nil, err
	}
	clusterName := "apisix"
	apisixOptions := &apisix.ClusterOptions{
		Name:     clusterName,
		AdminKey: apisixAdminKey,
		BaseURL:  apisixBaseURL,
	}
	if os.Getenv(EnvAPISIXAdminToken) != "" {
		apisixOptions.AdminKey = strings.TrimRight(os.Getenv(EnvAPISIXAdminToken), "/")
	}
	if os.Getenv(EnvAPISIXAdminHostAddress) != "" {
		apisixOptions.BaseURL = strings.TrimRight(os.Getenv(EnvAPISIXAdminHostAddress), "/")
	}
	apisixCli.AddCluster(apisixOptions)
	return apisixCli.Cluster(clusterName), nil
}
