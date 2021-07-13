package kong

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/api7/kongtoapisix/pkg/apisix"
)

var EnvKongYamlPath = "KONG_YAML_PATH"

func Migrate(APISIXConfig *apisix.APISIX) error {
	kongConfig, err := readYaml()
	if err != nil {
		return err
	}

	if upstreams, err := MigrateUpstream(kongConfig); err != nil {
		return err
	} else {
		APISIXConfig.Upstreams = upstreams
	}

	if routes, err := MigrateRoute(kongConfig); err != nil {
		return err
	} else {
		APISIXConfig.Routes = routes
	}

	if consumers, err := MigrateConsumer(kongConfig); err != nil {
		return err
	} else {
		APISIXConfig.Consumers = consumers
	}

	if globalRules, err := MigrateGlobalRules(kongConfig); err != nil {
		return err
	} else {
		APISIXConfig.GlobalRules = globalRules
	}
	return nil
}

func readYaml() (*KongConfig, error) {
	yamlPath := "/Users/shuyangwu/yiyiyimu/kong-to-apisix/kong.yaml"
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
