package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/api7/kongtoapisix/pkg/apisix"
	"github.com/api7/kongtoapisix/pkg/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMigrate(t *testing.T) {
	os.Setenv("EXPORT_PATH", "../../repos/apisix-docker/example/apisix_conf")
	configPath := filepath.Join("../../", utils.ConfigFilePath)
	if err := apisix.EnableAPISIXStandalone(configPath); err != nil {
		panic(err)
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "migrate suite")
}
