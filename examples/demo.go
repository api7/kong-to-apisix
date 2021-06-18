package main

import (
	"github.com/apache/apisix-ingress-controller/pkg/apisix"
	"github.com/globocom/gokong"
	"github.com/yiyiyimu/kongtoapisix/pkg/consumer"
	"github.com/yiyiyimu/kongtoapisix/pkg/plugin"
	"github.com/yiyiyimu/kongtoapisix/pkg/route"
	"github.com/yiyiyimu/kongtoapisix/pkg/upstream"
)

var (
	apisixBaseURL  = "http://127.0.0.1:9080/apisix/admin"
	apisixAdminKey = "edd1c9f034335f136f87ad84b625c8f1"
)

func main() {
	apisixCli, err := createAPISIXCli()
	if err != nil {
		panic(err)
	}
	kongCli := gokong.NewClient(gokong.NewDefaultConfig())

	err = upstream.MigrateUpstream(apisixCli, kongCli)
	if err != nil {
		panic("migrate service failed: " + err.Error())
	}

	err = route.MigrateRoute(apisixCli, kongCli)
	if err != nil {
		panic("migrate route failed: " + err.Error())
	}
	err = consumer.MigrateConsumer(apisixCli, kongCli)
	if err != nil {
		panic("migrate consumer failed: " + err.Error())
	}

	err = plugin.MigratePlugin(apisixCli, kongCli)
	if err != nil {
		panic("migrate plugin failed: " + err.Error())
	}
}

func createAPISIXCli() (apisix.Cluster, error) {
	apisixCli, err := apisix.NewClient()
	if err != nil {
		return nil, err
	}
	clusterName := "apisix"
	apisixCli.AddCluster(&apisix.ClusterOptions{
		Name:     clusterName,
		AdminKey: apisixAdminKey,
		BaseURL:  apisixBaseURL,
	})
	return apisixCli.Cluster(clusterName), nil
}
