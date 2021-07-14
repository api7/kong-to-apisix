package apisix

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/api7/kongtoapisix/pkg/utils"
	"gopkg.in/yaml.v2"
)

var yamlEndFlag = []byte("#END")

func EnableAPISIXStandalone(filePath string) error {
	if err := utils.AddValueToYaml(filePath, false, "apisix", "enable_admin"); err != nil {
		return err
	}
	if err := utils.AddValueToYaml(filePath, "yaml", "apisix", "config_center"); err != nil {
		return err
	}
	return nil
}

func WriteToFile(apisixConfig *Config) error {
	apisixYaml, err := yaml.Marshal(apisixConfig)
	if err != nil {
		return err
	}
	apisixYaml = append(apisixYaml, yamlEndFlag...)
	exportPath := "apisix.yaml"
	if os.Getenv("EXPORT_PATH") != "" {
		exportPath = filepath.Join(os.Getenv("EXPORT_PATH"), exportPath)
	}
	if err := ioutil.WriteFile(exportPath, apisixYaml, 0644); err != nil {
		return err
	}
	return nil
}
