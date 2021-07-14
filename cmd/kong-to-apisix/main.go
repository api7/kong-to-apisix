package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/api7/kongtoapisix/pkg/apisix"
	"github.com/api7/kongtoapisix/pkg/kong"
	"gopkg.in/yaml.v2"
)

var yamlEndFlag = []byte("#END")

func main() {
	apisixConfig := &apisix.APISIX{}
	if err := kong.Migrate(apisixConfig); err != nil {
		panic("migrate failed: " + err.Error())
	}
	apisixYaml, err := yaml.Marshal(apisixConfig)
	if err != nil {
		panic(err)
	}
	apisixYaml = append(apisixYaml, yamlEndFlag...)
	exportPath := "apisix.yaml"
	if os.Getenv("EXPORT_PATH") != "" {
		exportPath = filepath.Join(os.Getenv("EXPORT_PATH"), exportPath)
	}
	if err := ioutil.WriteFile(exportPath, apisixYaml, 0644); err != nil {
		panic(err)
	}
	if err := enableAPISIXStandalone(); err != nil {
		panic(err)
	}
}

func enableAPISIXStandalone() error {
	if err := kong.AddValueToYaml(false, "apisix", "enable_admin"); err != nil {
		return err
	}
	if err := kong.AddValueToYaml("yaml", "apisix", "config_center"); err != nil {
		return err
	}
	return nil
}
