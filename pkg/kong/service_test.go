package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"

	uuid "github.com/satori/go.uuid"

	"github.com/api7/kong-to-apisix/pkg/apisix"
)

func TestMigrateService(t *testing.T) {
	var kongConfig Config
	var apisixConfig apisix.Config
	var kongConfigService Service
	kongConfigService.ID = uuid.NewV4().String()
	kongConfigService.Name = "svc"
	kongConfigService.Path = "/svc"
	kongConfigService.Host = "example.com"
	kongConfigService.Port = 80
	kongConfigService.Protocol = "http"
	kongConfigService.Retries = 5
	kongConfigService.ConnectTimeout = 1500
	kongConfigService.ReadTimeout = 150
	kongConfigService.WriteTimeout = 15
	kongConfig.Services = append(kongConfig.Services, kongConfigService)
	err := MigrateService(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	assert.Equal(t, apisixConfig.Services[0].ID, kongConfigService.ID)
	assert.Equal(t, apisixConfig.Services[0].Name, kongConfigService.Name)
	assert.Equal(t, apisixConfig.Services[0].UpstreamID, "svc-"+kongConfigService.ID+"-ups")
	assert.Equal(t, apisixConfig.Upstreams[0].Nodes[0].Host, kongConfigService.Host)
	assert.Equal(t, apisixConfig.Upstreams[0].Nodes[0].Port, kongConfigService.Port)
	assert.Equal(t, apisixConfig.Upstreams[0].Nodes[0].Weight, 100)
	assert.Equal(t, apisixConfig.Upstreams[0].Timeout.Connect, float32(1.5))
	assert.Equal(t, apisixConfig.Upstreams[0].Timeout.Read, float32(0.15))
	assert.Equal(t, apisixConfig.Upstreams[0].Timeout.Send, float32(0.015))
	assert.Equal(t, apisixConfig.Upstreams[0].Scheme, kongConfigService.Protocol)
	assert.Equal(t, apisixConfig.Upstreams[0].Retries, kongConfigService.Retries)
}

func TestGenerateApisixServiceUpstream(t *testing.T) {
	var apisixConfig apisix.Config
	var kongConfigService Service
	kongConfigService.ID = uuid.NewV4().String()
	kongConfigService.Host = "example.com"
	kongConfigService.Port = 80
	kongConfigService.Protocol = "http"
	kongConfigService.Retries = 5
	kongConfigService.ConnectTimeout = 60000
	kongConfigService.ReadTimeout = 60000
	kongConfigService.WriteTimeout = 60000
	upstreamID := GenerateApisixServiceUpstream(kongConfigService, &apisixConfig)
	assert.Equal(t, upstreamID, apisixConfig.Upstreams[0].ID)
	assert.Equal(t, apisixConfig.Upstreams[0].Nodes[0].Host, kongConfigService.Host)
	assert.Equal(t, apisixConfig.Upstreams[0].Nodes[0].Port, kongConfigService.Port)
	assert.Equal(t, apisixConfig.Upstreams[0].Nodes[0].Weight, 100)
	assert.Equal(t, apisixConfig.Upstreams[0].Timeout.Connect, float32(60))
	assert.Equal(t, apisixConfig.Upstreams[0].Timeout.Read, float32(60))
	assert.Equal(t, apisixConfig.Upstreams[0].Timeout.Send, float32(60))
	assert.Equal(t, apisixConfig.Upstreams[0].Scheme, kongConfigService.Protocol)
	assert.Equal(t, apisixConfig.Upstreams[0].Retries, kongConfigService.Retries)
}

func TestFindKongServiceByID(t *testing.T) {
	var kongConfig Config
	var kongConfigService01 Service
	var kongConfigService02 Service
	serviceID1 := uuid.NewV4().String()
	serviceID2 := uuid.NewV4().String()
	assert.NotEqual(t, serviceID1, serviceID2)
	kongConfigService01.ID = serviceID1
	kongConfigService02.ID = serviceID2
	kongConfigService01.Name = "svc01"
	kongConfigService02.Name = "svc02"
	kongConfig.Services = append(kongConfig.Services, kongConfigService01)
	kongConfig.Services = append(kongConfig.Services, kongConfigService02)
	kongConfigService, err := FindKongServiceByID(&kongConfig.Services, serviceID1)
	assert.Nil(t, err)
	assert.Equal(t, kongConfigService.ID, kongConfigService01.ID)
	assert.Equal(t, kongConfigService.Name, kongConfigService01.Name)
	kongConfigService, err = FindKongServiceByID(&kongConfig.Services, serviceID2)
	assert.Nil(t, err)
	assert.Equal(t, kongConfigService.ID, kongConfigService02.ID)
	assert.Equal(t, kongConfigService.Name, kongConfigService02.Name)
}
