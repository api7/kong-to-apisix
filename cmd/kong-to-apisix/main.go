package main

import (
	"io/ioutil"

	"github.com/api7/kongtoapisix/pkg/apisix"
	"github.com/api7/kongtoapisix/pkg/kong"
	"gopkg.in/yaml.v2"
)

func main() {
	apisixConfig := &apisix.APISIX{}
	if err := kong.Migrate(apisixConfig); err != nil {
		panic("migrate failed: " + err.Error())
	}
	apisixYaml, err := yaml.Marshal(apisixConfig)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("apisix.yaml", apisixYaml, 0644); err != nil {
		panic(err)
	}
}
