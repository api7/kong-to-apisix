package kong

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func readYaml() {
	yamlFile, err := ioutil.ReadFile("/Users/shuyangwu/yiyiyimu/kong-to-apisix/kong.yaml")
	if err != nil {
		panic(err)
	}
	var kongConfig *KongConfig
	err = yaml.Unmarshal(yamlFile, &kongConfig)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", kongConfig.Consumers)
}
