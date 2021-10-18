package kong

import (
	"testing"

	"github.com/api7/kong-to-apisix/pkg/apisix"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type RouteExpect struct {
	ID        string
	Name      string
	Path      string
	Hosts     []string
	Methods   []string
	ServiceId string
	Priority  uint
}

type TestMigrateRouteExpect struct {
	Route        Route
	RouteExpects []RouteExpect
}

func TestMigrateRoute(t *testing.T) {
	var err error
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
	err = MigrateService(&kongConfig, &apisixConfig)
	assert.NoError(t, err)

	testRouteId1 := uuid.NewV4().String()
	testRouteId2 := uuid.NewV4().String()
	testMigrateRouteExpects := []TestMigrateRouteExpect{
		{
			Route: Route{
				ID:            testRouteId1,
				Name:          "test-route-01",
				Hosts:         []string{"01.route.test"},
				Paths:         []string{"/test/route/01/foo", "/test/route/01/bar"},
				Methods:       []string{"GET", "POST"},
				ServiceID:     kongService.ID,
				RegexPriority: 10,
			},
			RouteExpects: []RouteExpect{
				{
					ID:        testRouteId1 + "-1",
					Name:      "test-route-01-1",
					Path:      "/test/route/01/foo*",
					Hosts:     []string{"01.route.test"},
					Methods:   []string{"GET", "POST"},
					ServiceId: kongService.ID,
					Priority:  10,
				},
				{
					ID:        testRouteId1 + "-2",
					Name:      "test-route-01-2",
					Path:      "/test/route/01/bar*",
					Hosts:     []string{"01.route.test"},
					Methods:   []string{"GET", "POST"},
					ServiceId: kongService.ID,
					Priority:  10,
				},
			},
		},
		{
			Route: Route{
				ID:            testRouteId2,
				Name:          "test-route-02",
				Hosts:         []string{"02.route.test"},
				Paths:         []string{"/test/route/02/foo", "/test/route/02/bar"},
				Methods:       []string{"PUT", "DELETE"},
				ServiceID:     kongService.ID,
				RegexPriority: 20,
			},
			RouteExpects: []RouteExpect{
				{
					ID:        testRouteId2 + "-1",
					Name:      "test-route-02-1",
					Path:      "/test/route/02/foo*",
					Hosts:     []string{"02.route.test"},
					Methods:   []string{"PUT", "DELETE"},
					ServiceId: kongService.ID,
					Priority:  20,
				},
				{
					ID:        testRouteId2 + "-2",
					Name:      "test-route-02-2",
					Path:      "/test/route/02/bar*",
					Hosts:     []string{"02.route.test"},
					Methods:   []string{"PUT", "DELETE"},
					ServiceId: kongService.ID,
					Priority:  20,
				},
			},
		},
	}

	for _, testMigrateRouteExpect := range testMigrateRouteExpects {
		kongConfig.Routes = Routes{}
		apisixConfig.Routes = apisix.Routes{}
		kongConfig.Routes = append(kongConfig.Routes, testMigrateRouteExpect.Route)
		err = MigrateRoute(&kongConfig, &apisixConfig)
		assert.NoError(t, err)
		assert.Equal(t, len(apisixConfig.Routes), len(testMigrateRouteExpect.RouteExpects))
		assert.Equal(t, len(apisixConfig.Routes[0].Plugins.ProxyRewrite.RegexURI), 2)
		assert.Equal(t, apisixConfig.Routes[0].ID, testMigrateRouteExpect.RouteExpects[0].ID)
		assert.Equal(t, apisixConfig.Routes[0].Name, testMigrateRouteExpect.RouteExpects[0].Name)
		assert.Equal(t, apisixConfig.Routes[0].URI, testMigrateRouteExpect.RouteExpects[0].Path)
		assert.Equal(t, apisixConfig.Routes[0].Methods, testMigrateRouteExpect.RouteExpects[0].Methods)
		assert.Equal(t, apisixConfig.Routes[0].Hosts, testMigrateRouteExpect.RouteExpects[0].Hosts)
		assert.Equal(t, apisixConfig.Routes[0].Priority, testMigrateRouteExpect.RouteExpects[0].Priority)
		assert.Equal(t, apisixConfig.Routes[0].ServiceID, testMigrateRouteExpect.RouteExpects[0].ServiceId)
		assert.Equal(t, len(apisixConfig.Routes[1].Plugins.ProxyRewrite.RegexURI), 2)
		assert.Equal(t, apisixConfig.Routes[1].ID, testMigrateRouteExpect.RouteExpects[1].ID)
		assert.Equal(t, apisixConfig.Routes[1].Name, testMigrateRouteExpect.RouteExpects[1].Name)
		assert.Equal(t, apisixConfig.Routes[1].URI, testMigrateRouteExpect.RouteExpects[1].Path)
		assert.Equal(t, apisixConfig.Routes[1].Methods, testMigrateRouteExpect.RouteExpects[1].Methods)
		assert.Equal(t, apisixConfig.Routes[1].Hosts, testMigrateRouteExpect.RouteExpects[1].Hosts)
		assert.Equal(t, apisixConfig.Routes[1].Priority, testMigrateRouteExpect.RouteExpects[1].Priority)
		assert.Equal(t, apisixConfig.Routes[1].ServiceID, testMigrateRouteExpect.RouteExpects[1].ServiceId)
	}
<<<<<<< HEAD
=======
	kongConfig.Routes = append(kongConfig.Routes, kongRoute2)

	err := MigrateRoute(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	assert.Equal(t, len(apisixConfig.Routes), 4)

	assert.Equal(t, apisixConfig.Routes[0].ID, kongRoute1.ID+"-1")
	assert.Equal(t, apisixConfig.Routes[0].Name, kongRoute1.Name+"-1")
	assert.Equal(t, apisixConfig.Routes[0].URI, kongRoute1.Paths[0]+"*")
	assert.Equal(t, apisixConfig.Routes[0].Methods, kongRoute1.Methods)
	assert.Equal(t, apisixConfig.Routes[0].Hosts, kongRoute1.Hosts)
	assert.Equal(t, apisixConfig.Routes[0].ServiceID, kongService.ID)
	assert.Equal(t, len(apisixConfig.Routes[0].Plugins), 1)

	assert.Equal(t, apisixConfig.Routes[1].ID, kongRoute1.ID+"-2")
	assert.Equal(t, apisixConfig.Routes[1].Name, kongRoute1.Name+"-2")
	assert.Equal(t, apisixConfig.Routes[1].URI, kongRoute1.Paths[1]+"*")
	assert.Equal(t, apisixConfig.Routes[1].Methods, kongRoute1.Methods)
	assert.Equal(t, apisixConfig.Routes[1].Hosts, kongRoute1.Hosts)
	assert.Equal(t, apisixConfig.Routes[1].ServiceID, kongService.ID)
	assert.Equal(t, len(apisixConfig.Routes[1].Plugins), 1)

	assert.Equal(t, apisixConfig.Routes[2].ID, kongRoute2.ID+"-1")
	assert.Equal(t, apisixConfig.Routes[2].Name, kongRoute2.Name+"-1")
	assert.Equal(t, apisixConfig.Routes[2].URI, kongRoute2.Paths[0]+"*")
	assert.Equal(t, apisixConfig.Routes[2].Methods, kongRoute2.Methods)
	assert.Equal(t, apisixConfig.Routes[2].Hosts, kongRoute2.Hosts)
	assert.Equal(t, apisixConfig.Routes[2].ServiceID, kongService.ID)
	assert.Equal(t, len(apisixConfig.Routes[2].Plugins), 1)

	assert.Equal(t, apisixConfig.Routes[3].ID, kongRoute2.ID+"-2")
	assert.Equal(t, apisixConfig.Routes[3].Name, kongRoute2.Name+"-2")
	assert.Equal(t, apisixConfig.Routes[3].URI, kongRoute2.Paths[1]+"*")
	assert.Equal(t, apisixConfig.Routes[3].Methods, kongRoute2.Methods)
	assert.Equal(t, apisixConfig.Routes[3].Hosts, kongRoute2.Hosts)
	assert.Equal(t, apisixConfig.Routes[3].ServiceID, kongService.ID)
	assert.Equal(t, len(apisixConfig.Routes[3].Plugins), 1)

>>>>>>> 23abdb5 (feat: compatible with kong deck and kong config)
}

func TestGenerateProxyRewritePluginConfig(t *testing.T) {
	var config *apisix.ProxyRewrite
	//https://docs.konghq.com/gateway-oss/2.4.x/admin-api/#path-handling-algorithms
	config = GenerateProxyRewritePluginConfig("/s", "/fv0", false, "v0")
	assert.Equal(t, config.RegexURI[0], "^/?(.*)")
	assert.Equal(t, config.RegexURI[1], "/s/$1")
	config = GenerateProxyRewritePluginConfig("/s", "/fv1", false, "v1")
	assert.Equal(t, config.RegexURI[0], "^/?(.*)")
	assert.Equal(t, config.RegexURI[1], "/s$1")
	config = GenerateProxyRewritePluginConfig("/s", "/tv0", true, "v0")
	assert.Equal(t, config.RegexURI[0], "^/tv0/?(.*)")
	assert.Equal(t, config.RegexURI[1], "/s/$1")
	config = GenerateProxyRewritePluginConfig("/s", "/tv1", true, "v1")
	assert.Equal(t, config.RegexURI[0], "^/tv1/?(.*)")
	assert.Equal(t, config.RegexURI[1], "/s$1")
	config = GenerateProxyRewritePluginConfig("/s", "/fv0/", false, "v0")
	assert.Equal(t, config.RegexURI[0], "^/?(.*)")
	assert.Equal(t, config.RegexURI[1], "/s/$1")
	config = GenerateProxyRewritePluginConfig("/s", "/fv1/", false, "v1")
	assert.Equal(t, config.RegexURI[0], "^/?(.*)")
	assert.Equal(t, config.RegexURI[1], "/s$1")
	config = GenerateProxyRewritePluginConfig("/s", "/tv0/", true, "v0")
	assert.Equal(t, config.RegexURI[0], "^/tv0/?(.*)")
	assert.Equal(t, config.RegexURI[1], "/s/$1")
	config = GenerateProxyRewritePluginConfig("/s", "/tv1/", true, "v1")
	assert.Equal(t, config.RegexURI[0], "^/tv1/?(.*)")
	assert.Equal(t, config.RegexURI[1], "/s$1")
}
