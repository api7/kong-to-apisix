package main

import (
	"github.com/globocom/gokong"
	"github.com/yiyiyimu/kongtoapisix/pkg/apisixcli"
	"github.com/yiyiyimu/kongtoapisix/pkg/consumer"
	"github.com/yiyiyimu/kongtoapisix/pkg/plugin"
	"github.com/yiyiyimu/kongtoapisix/pkg/route"
	"github.com/yiyiyimu/kongtoapisix/pkg/upstream"
)

func main() {
	apisixCli, err := apisixcli.CreateAPISIXCli()
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
