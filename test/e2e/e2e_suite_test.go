package e2e

import (
	"testing"

	"github.com/api7/kong-to-apisix/pkg/apisix"
	"github.com/api7/kong-to-apisix/pkg/utils"
	e2eutils "github.com/api7/kong-to-apisix/test/e2e/utils"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestMigrate(t *testing.T) {
	var apisixConfig []utils.YamlItem
	if err := apisix.EnableAPISIXStandalone(&apisixConfig); err != nil {
		panic(err)
	}
	if err := utils.AppendToConfigYaml(&apisixConfig, e2eutils.ApisixConfigYamlPath); err != nil {
		panic(err)
	}
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "migrate suite")
}
