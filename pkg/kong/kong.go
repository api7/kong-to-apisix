package kong

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/api7/kongtoapisix/pkg/apisix"
)

var EnvKongYamlPath = "KONG_YAML_PATH"

func Migrate(kongConfig *KongConfig) (*apisix.Config, error) {
	apisixConfig := &apisix.Config{}

	if upstreams, err := MigrateUpstream(kongConfig); err != nil {
		return nil, err
	} else {
		apisixConfig.Upstreams = upstreams
	}

	if routes, err := MigrateRoute(kongConfig); err != nil {
		return nil, err
	} else {
		apisixConfig.Routes = routes
	}

	if consumers, err := MigrateConsumer(kongConfig); err != nil {
		return nil, err
	} else {
		apisixConfig.Consumers = consumers
	}

	if globalRules, err := MigrateGlobalRules(kongConfig); err != nil {
		return nil, err
	} else {
		apisixConfig.GlobalRules = globalRules
	}
	return apisixConfig, nil
}

func ReadYaml() (*KongConfig, error) {
	yamlPath := "kong.yaml"
	if os.Getenv(EnvKongYamlPath) != "" {
		yamlPath = os.Getenv(EnvKongYamlPath)
	}
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
