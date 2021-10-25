package kong

import (
	"github.com/api7/kong-to-apisix/pkg/apisix"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestMigratePlugins(t *testing.T) {
	var err error
	var kongConfig Config
	var apisixConfig apisix.Config

	serviceId := uuid.NewV4().String()
	routeId := uuid.NewV4().String()
	kongConfig.Services = Services{
		{
			ID:   serviceId,
			Name: "svc",
		},
	}
	err = MigrateService(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	apisixConfig.Routes = apisix.Routes{
		{
			ID:   routeId,
			Name: "route",
		},
	}

	consumerId := uuid.NewV4().String()
	kongConfig.Consumers = Consumers{
		{
			ID:       consumerId,
			Username: "consumer",
		},
	}
	err = MigrateConsumer(&kongConfig, &apisixConfig)
	assert.NoError(t, err)

	kongPluginKeyAuthConfig := make(map[string]interface{})
	kongPluginKeyAuthConfig["key_names"] = make([]interface{}, 1)
	kongPluginKeyAuthConfigKeyNames := make([]interface{}, 1)
	for i, v := range []string{"apikey"} {
		kongPluginKeyAuthConfigKeyNames[i] = v
	}
	kongPluginKeyAuthConfig["key_names"] = kongPluginKeyAuthConfigKeyNames
	kongConfig.Plugins = Plugins{
		{
			ID:        uuid.NewV4().String(),
			Name:      PluginKeyAuth,
			Config:    kongPluginKeyAuthConfig,
			ServiceID: serviceId,
			Enabled:   true,
		},
		{
			ID:      uuid.NewV4().String(),
			Name:    PluginRateLimiting,
			RouteID: routeId,
			Config: map[string]interface{}{
				"second": 5,
				"policy": PluginRateLimitingPolicyLocal,
			},
			Enabled: true,
		},
	}

	kongConfig.KeyAuthCredentials = KeyAuthCredentials{
		{
			ConsumerID: consumerId,
			Key:        "test",
		},
	}
	kongConfig.BasicAuthCredentials = BasicAuthCredentials{
		{
			Username:   "test",
			Password:   "123456",
			ConsumerID: consumerId,
		},
	}
	kongConfig.HmacAuthCredentials = HmacAuthCredentials{
		{
			Username:   "test",
			Secret:     "123456",
			ConsumerID: consumerId,
		},
	}
	kongConfig.JwtSecrets = JwtSecrets{
		{
			Key:        "test",
			Secret:     "123456",
			Algorithm:  "HS256",
			ConsumerID: consumerId,
		},
	}

	err = MigratePlugins(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	//assert.Equal(t, apisixConfig.Services[0].Plugins.KeyAuth, &apisix.KeyAuth{})
	assert.Equal(t, apisixConfig.Routes[0].Plugins.LimitCount.Count, kongConfig.Plugins[1].Config["second"])
	assert.Equal(t, apisixConfig.Routes[0].Plugins.LimitCount.TimeWindow, 1)
	assert.Equal(t, apisixConfig.Routes[0].Plugins.LimitCount.Policy, kongConfig.Plugins[1].Config["policy"])
	assert.Equal(t, apisixConfig.Consumers[0].Plugins.KeyAuth.Key, kongConfig.KeyAuthCredentials[0].Key)
	assert.Equal(t, apisixConfig.Consumers[0].Plugins.BasicAuth.Username, kongConfig.BasicAuthCredentials[0].Username)
	assert.Equal(t, apisixConfig.Consumers[0].Plugins.BasicAuth.Password, kongConfig.BasicAuthCredentials[0].Password)
	assert.Equal(t, apisixConfig.Consumers[0].Plugins.HmacAuth.AccessKey, kongConfig.HmacAuthCredentials[0].Username)
	assert.Equal(t, apisixConfig.Consumers[0].Plugins.HmacAuth.SecretKey, kongConfig.HmacAuthCredentials[0].Secret)
	assert.Equal(t, apisixConfig.Consumers[0].Plugins.JwtAuth.Key, kongConfig.JwtSecrets[0].Key)
	assert.Equal(t, apisixConfig.Consumers[0].Plugins.JwtAuth.Secret, kongConfig.JwtSecrets[0].Secret)
	assert.Equal(t, apisixConfig.Consumers[0].Plugins.JwtAuth.Algorithm, kongConfig.JwtSecrets[0].Algorithm)
}
