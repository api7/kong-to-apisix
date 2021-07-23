package e2e

import (
	"testing"

	"github.com/api7/kongtoapisix/pkg/apisix"
	"github.com/api7/kongtoapisix/pkg/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMigrate(t *testing.T) {
	var apisixConfig []utils.YamlItem
	if err := apisix.EnableAPISIXStandalone(&apisixConfig); err != nil {
		panic(err)
	}
	if err := utils.AppendToConfigYaml(&apisixConfig, apisixConfigYamlPath); err != nil {
		panic(err)
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "migrate suite")
}
