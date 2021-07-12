package apisixcli

import (
	"os"
	"strings"

	"github.com/apache/apisix-ingress-controller/pkg/apisix"
)

var (
	apisixBaseURL  = "http://127.0.0.1:9080"
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
	}
	if os.Getenv(EnvAPISIXAdminToken) != "" {
		apisixOptions.AdminKey = strings.TrimRight(os.Getenv(EnvAPISIXAdminToken), "/")
	}
	apisixAdminURL := apisixBaseURL
	if os.Getenv(EnvAPISIXAdminHostAddress) != "" {
		apisixAdminURL = strings.TrimRight(os.Getenv(EnvAPISIXAdminHostAddress), "/")
	}
	apisixOptions.BaseURL = apisixAdminURL + "/apisix/admin"
	apisixCli.AddCluster(apisixOptions)
	return apisixCli.Cluster(clusterName), nil
}
