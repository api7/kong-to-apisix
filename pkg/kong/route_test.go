package kong

import (
	"testing"

	"github.com/api7/kong-to-apisix/pkg/apisix"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestMigrateRoute(t *testing.T) {
	var kongConfig Config
	var apisixConfig apisix.Config
	var kongService Service
	kongService.ID = uuid.NewV4().String()
	kongService.Name = "svc"
	kongService.Path = "/svc"
	kongService.Host = "example.com"
	kongService.Port = 80
	kongService.Protocol = "http"
	kongService.Retries = 5
	kongService.ConnectTimeout = 60000
	kongService.ReadTimeout = 60000
	kongService.WriteTimeout = 60000
	kongConfig.Services = append(kongConfig.Services, kongService)

	var kongRoute1 Route
	kongRoute1.ID = uuid.NewV4().String()
	kongRoute1.Name = "route-01"
	kongRoute1.Hosts = []string{"route01.test.com"}
	kongRoute1.Paths = []string{"/route01/foo", "/route01/bar"}
	kongRoute1.Service = kongService.ID
	kongRoute1.Methods = []string{"GET", "POST"}
	kongRoute1.Headers = map[string][]string{
		"R-01-HEADER": []string{"foo", "bar"},
	}
	kongConfig.Routes = append(kongConfig.Routes, kongRoute1)

	var kongRoute2 Route
	kongRoute2.ID = uuid.NewV4().String()
	kongRoute2.Name = "route-02"
	kongRoute2.Hosts = []string{"route02.test.com"}
	kongRoute2.Paths = []string{"/route02/foo", "/route02/bar"}
	kongRoute2.Service = kongService.ID
	kongRoute2.Methods = []string{"PUT", "DELETE"}
	kongRoute2.Headers = map[string][]string{
		"R-02-HEADER": []string{"foo", "bar"},
	}
	kongConfig.Routes = append(kongConfig.Routes, kongRoute2)

	err := MigrateRoute(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	assert.Equal(t, len(apisixConfig.Routes), 4)

	assert.Equal(t, apisixConfig.Routes[0].ID, kongRoute1.ID+"-0")
	assert.Equal(t, apisixConfig.Routes[0].Name, kongRoute1.Name+"-0")
	assert.Equal(t, apisixConfig.Routes[0].URI, kongRoute1.Paths[0]+"*")
	assert.Equal(t, apisixConfig.Routes[0].Methods, kongRoute1.Methods)
	assert.Equal(t, apisixConfig.Routes[0].Hosts, kongRoute1.Hosts)
	assert.Equal(t, apisixConfig.Routes[0].ServiceID, kongService.ID)
	assert.Equal(t, len(apisixConfig.Routes[0].Plugins), 1)

	assert.Equal(t, apisixConfig.Routes[1].ID, kongRoute1.ID+"-1")
	assert.Equal(t, apisixConfig.Routes[1].URI, kongRoute1.Paths[1]+"*")
	assert.Equal(t, apisixConfig.Routes[1].Methods, kongRoute1.Methods)
	assert.Equal(t, apisixConfig.Routes[1].Hosts, kongRoute1.Hosts)
	assert.Equal(t, apisixConfig.Routes[1].ServiceID, kongService.ID)
	assert.Equal(t, len(apisixConfig.Routes[1].Plugins), 1)

	assert.Equal(t, apisixConfig.Routes[2].ID, kongRoute2.ID+"-0")
	assert.Equal(t, apisixConfig.Routes[2].URI, kongRoute2.Paths[0]+"*")
	assert.Equal(t, apisixConfig.Routes[2].Methods, kongRoute2.Methods)
	assert.Equal(t, apisixConfig.Routes[2].Hosts, kongRoute2.Hosts)
	assert.Equal(t, apisixConfig.Routes[2].ServiceID, kongService.ID)
	assert.Equal(t, len(apisixConfig.Routes[2].Plugins), 1)

	assert.Equal(t, apisixConfig.Routes[3].ID, kongRoute2.ID+"-1")
	assert.Equal(t, apisixConfig.Routes[3].URI, kongRoute2.Paths[1]+"*")
	assert.Equal(t, apisixConfig.Routes[3].Methods, kongRoute2.Methods)
	assert.Equal(t, apisixConfig.Routes[3].Hosts, kongRoute2.Hosts)
	assert.Equal(t, apisixConfig.Routes[3].ServiceID, kongService.ID)
	assert.Equal(t, len(apisixConfig.Routes[3].Plugins), 1)

}

func TestGenerateProxyRewritePluginConfig(t *testing.T) {
	var config map[string]interface{}
	//https://docs.konghq.com/gateway-oss/2.4.x/admin-api/#path-handling-algorithms
	config = GenerateProxyRewritePluginConfig("/s", "/fv0", false, "v0")
	assert.Equal(t, config["regex_uri"].([]string)[0], "^/?(.*)")
	assert.Equal(t, config["regex_uri"].([]string)[1], "/s/$1")
	config = GenerateProxyRewritePluginConfig("/s", "/fv1", false, "v1")
	assert.Equal(t, config["regex_uri"].([]string)[0], "^/?(.*)")
	assert.Equal(t, config["regex_uri"].([]string)[1], "/s$1")
	config = GenerateProxyRewritePluginConfig("/s", "/tv0", true, "v0")
	assert.Equal(t, config["regex_uri"].([]string)[0], "^/tv0/?(.*)")
	assert.Equal(t, config["regex_uri"].([]string)[1], "/s/$1")
	config = GenerateProxyRewritePluginConfig("/s", "/tv1", true, "v1")
	assert.Equal(t, config["regex_uri"].([]string)[0], "^/tv1/?(.*)")
	assert.Equal(t, config["regex_uri"].([]string)[1], "/s$1")
	config = GenerateProxyRewritePluginConfig("/s", "/fv0/", false, "v0")
	assert.Equal(t, config["regex_uri"].([]string)[0], "^/?(.*)")
	assert.Equal(t, config["regex_uri"].([]string)[1], "/s/$1")
	config = GenerateProxyRewritePluginConfig("/s", "/fv1/", false, "v1")
	assert.Equal(t, config["regex_uri"].([]string)[0], "^/?(.*)")
	assert.Equal(t, config["regex_uri"].([]string)[1], "/s$1")
	config = GenerateProxyRewritePluginConfig("/s", "/tv0/", true, "v0")
	assert.Equal(t, config["regex_uri"].([]string)[0], "^/tv0/?(.*)")
	assert.Equal(t, config["regex_uri"].([]string)[1], "/s/$1")
	config = GenerateProxyRewritePluginConfig("/s", "/tv1/", true, "v1")
	assert.Equal(t, config["regex_uri"].([]string)[0], "^/tv1/?(.*)")
	assert.Equal(t, config["regex_uri"].([]string)[1], "/s$1")
}
