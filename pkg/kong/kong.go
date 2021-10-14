package kong

import (
	"io/ioutil"

	"github.com/api7/kong-to-apisix/pkg/apisix"
	"github.com/api7/kong-to-apisix/pkg/utils"

	"gopkg.in/yaml.v2"
)

func Migrate(kongConfig *Config) (*apisix.Config, *[]utils.YamlItem, error) {
	apisixConfig := &apisix.Config{}
	var configYamlAll []utils.YamlItem

	if err := MigrateService(kongConfig, apisixConfig); err != nil {
		return nil, nil, err
	}

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

func ReadYaml(yamlPath string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}
	var kongConfig *Config
	err = yaml.Unmarshal(yamlFile, &kongConfig)
	if err != nil {
		return nil, err
	}

	return kongConfig, nil
}
