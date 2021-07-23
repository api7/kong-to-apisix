package apisix

import (
	"io/ioutil"

	"github.com/api7/kong-to-apisix/pkg/utils"
	"gopkg.in/yaml.v2"
)

var yamlEndFlag = []byte("#END")

func EnableAPISIXStandalone(apisixConfig *[]utils.YamlItem) error {
	*apisixConfig = append(*apisixConfig, utils.YamlItem{
		Value: false,
		Path:  []interface{}{"apisix", "enable_admin"},
	})
	*apisixConfig = append(*apisixConfig, utils.YamlItem{
		Value: "yaml",
		Path:  []interface{}{"apisix", "config_center"},
	})
	return nil
}

func MarshalYaml(apisixConfig *Config) ([]byte, error) {
	apisixYaml, err := yaml.Marshal(apisixConfig)
	if err != nil {
		return nil, err
	}
	apisixYaml = append(apisixYaml, yamlEndFlag...)
	return apisixYaml, nil
}

func WriteToFile(apisixYamlPath string, apisixYaml []byte) error {
	if err := ioutil.WriteFile(apisixYamlPath, apisixYaml, 0644); err != nil {
		return err
	}
	return nil
}
