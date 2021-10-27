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

	if err := MigrateUpstream(kongConfig, apisixConfig); err != nil {
		return nil, nil, err
	}

	if routes, err := MigrateRoute(kongConfig, &configYamlAll); err != nil {
		return nil, nil, err
	} else {
		apisixConfig.Routes = routes
	}

	if err := MigrateConsumer(kongConfig, apisixConfig); err != nil {
		return nil, nil, err
	}

	if err := MigratePlugins(kongConfig, apisixConfig); err != nil {
		return nil, nil, err
	}

	if err := MigrateGlobalRules(kongConfig, apisixConfig); err != nil {
		return nil, nil, err
	}

	if err := AfterMigrate(apisixConfig); err != nil {
		return nil, nil, err
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
