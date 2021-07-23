package kong

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/api7/kongtoapisix/pkg/apisix"
	"github.com/api7/kongtoapisix/pkg/utils"
)

func Migrate(kongConfig *KongConfig) (*apisix.Config, *[]utils.YamlItem, error) {
	apisixConfig := &apisix.Config{}
	var configYamlAll []utils.YamlItem

	if upstreams, err := MigrateUpstream(kongConfig, &configYamlAll); err != nil {
		return nil, nil, err
	} else {
		apisixConfig.Upstreams = upstreams
	}

	if routes, err := MigrateRoute(kongConfig, &configYamlAll); err != nil {
		return nil, nil, err
	} else {
		apisixConfig.Routes = routes
	}

	if consumers, err := MigrateConsumer(kongConfig, &configYamlAll); err != nil {
		return nil, nil, err
	} else {
		apisixConfig.Consumers = consumers
	}

	if globalRules, err := MigrateGlobalRules(kongConfig, &configYamlAll); err != nil {
		return nil, nil, err
	} else {
		apisixConfig.GlobalRules = globalRules
	}
	return apisixConfig, &configYamlAll, nil
}

func ReadYaml(yamlPath string) (*KongConfig, error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}
	var kongConfig *KongConfig
	err = yaml.Unmarshal(yamlFile, &kongConfig)
	if err != nil {
		return nil, err
	}

	return kongConfig, nil
}
