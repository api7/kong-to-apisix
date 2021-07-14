package main

import (
	"github.com/api7/kongtoapisix/pkg/apisix"
	"github.com/api7/kongtoapisix/pkg/kong"
)

func main() {
	kongConfig, err := kong.ReadYaml()
	if err != nil {
		panic(err)
	}
	apisixConfig, err := kong.Migrate(kongConfig)
	if err != nil {
		panic("migrate failed: " + err.Error())
	}
	if err := apisix.WriteToFile(apisixConfig); err != nil {
		panic(err)
	}
	if err := apisix.EnableAPISIXStandalone(); err != nil {
		panic(err)
	}
}
