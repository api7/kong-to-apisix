package kong

import (
	"testing"

	"github.com/api7/kong-to-apisix/pkg/apisix"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestMigrateGlobalRules(t *testing.T) {
	var err error
	var kongConfig Config
	var apisixConfig apisix.Config

	// KeyAuth Global Plugin Test
	kongPluginKeyAuthConfig := make(map[string]interface{})
	kongPluginKeyAuthConfig["key_names"] = make([]interface{}, 1)
	kongPluginKeyAuthConfigKeyNames := make([]interface{}, 1)
	for i, v := range []string{"apikey"} {
		kongPluginKeyAuthConfigKeyNames[i] = v
	}
	kongPluginKeyAuthConfig["key_names"] = kongPluginKeyAuthConfigKeyNames
	kongPluginKeyAuths := Plugins{
		{
			ID:      uuid.NewV4().String(),
			Name:    PluginKeyAuth,
			Enabled: true,
			Config:  kongPluginKeyAuthConfig,
		},
	}
	for _, kongPluginKeyAuth := range kongPluginKeyAuths {
		kongConfig = Config{}
		apisixConfig = apisix.Config{}
		kongConfig.Plugins = append(kongConfig.Plugins, kongPluginKeyAuth)
		err = MigrateGlobalRules(&kongConfig, &apisixConfig)
		assert.NoError(t, err)
		assert.Equal(t, kongPluginKeyAuth.ID, apisixConfig.GlobalRules[0].ID)
	}

	// RateLimiting Global Plugin Test
	kongPluginRateLimitings := Plugins{
		{
			ID:      uuid.NewV4().String(),
			Name:    PluginRateLimiting,
			Enabled: true,
			Config: map[string]interface{}{
				"second": 5,
				"policy": PluginRateLimitingPolicyLocal,
			},
		},
		{
			ID:      uuid.NewV4().String(),
			Name:    PluginRateLimiting,
			Enabled: true,
			Config: map[string]interface{}{
				"minute": 5,
				"policy": PluginRateLimitingPolicyLocal,
			},
		},
		{
			ID:      uuid.NewV4().String(),
			Name:    PluginRateLimiting,
			Enabled: true,
			Config: map[string]interface{}{
				"hour":   5,
				"policy": PluginRateLimitingPolicyLocal,
			},
		},
		{
			ID:      uuid.NewV4().String(),
			Name:    PluginRateLimiting,
			Enabled: true,
			Config: map[string]interface{}{
				"day":    5,
				"policy": PluginRateLimitingPolicyLocal,
			},
		},
		{
			ID:      uuid.NewV4().String(),
			Name:    PluginRateLimiting,
			Enabled: true,
			Config: map[string]interface{}{
				"month":  5,
				"policy": PluginRateLimitingPolicyLocal,
			},
		},
		{
			ID:      uuid.NewV4().String(),
			Name:    PluginRateLimiting,
			Enabled: true,
			Config: map[string]interface{}{
				"year":   5,
				"policy": PluginRateLimitingPolicyLocal,
			},
		},
	}

	for _, kongPluginRateLimiting := range kongPluginRateLimitings {
		kongConfig = Config{}
		apisixConfig = apisix.Config{}
		kongConfig.Plugins = append(kongConfig.Plugins, kongPluginRateLimiting)
		err = MigrateGlobalRules(&kongConfig, &apisixConfig)
		assert.NoError(t, err)
		assert.Equal(t, kongPluginRateLimiting.ID, apisixConfig.GlobalRules[0].ID)
		if kongPluginRateLimiting.Config["second"] != nil {
			assert.Equal(t, kongPluginRateLimiting.Config["second"].(int),
				apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
			assert.Equal(t, 1, apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)
		}
		if kongPluginRateLimiting.Config["minute"] != nil {
			assert.Equal(t, kongPluginRateLimiting.Config["minute"].(int),
				apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
			assert.Equal(t, 1*60, apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)
		}
		if kongPluginRateLimiting.Config["hour"] != nil {
			assert.Equal(t, kongPluginRateLimiting.Config["hour"].(int),
				apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
			assert.Equal(t, 1*60*60, apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)
		}
		if kongPluginRateLimiting.Config["day"] != nil {
			assert.Equal(t, kongPluginRateLimiting.Config["day"].(int),
				apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
			assert.Equal(t, 1*60*60*24, apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)
		}
		if kongPluginRateLimiting.Config["month"] != nil {
			assert.Equal(t, kongPluginRateLimiting.Config["month"].(int),
				apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
			assert.Equal(t, 1*60*60*24*30, apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)
		}
		if kongPluginRateLimiting.Config["year"] != nil {
			assert.Equal(t, kongPluginRateLimiting.Config["year"].(int),
				apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
			assert.Equal(t, 1*60*60*24*30*365, apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)
		}
	}

	// ProxyCache Global Plugin Test
	kongPluginProxyCacheConfig := make(map[string]interface{})
	kongPluginProxyCacheConfigRequestMethod := make([]interface{}, 4)
	kongPluginProxyCacheConfigRequestMethods := []string{"GET", "POST", "PUT", "DELETE"}
	for i, v := range kongPluginProxyCacheConfigRequestMethods {
		kongPluginProxyCacheConfigRequestMethod[i] = v
	}
	kongPluginProxyCacheConfig["request_method"] = kongPluginProxyCacheConfigRequestMethod
	kongPluginProxyCacheConfigResponseCode := make([]interface{}, 3)
	kongPluginProxyCacheConfigResponseCodes := []int{200, 204, 302}
	for i, v := range kongPluginProxyCacheConfigResponseCodes {
		kongPluginProxyCacheConfigResponseCode[i] = v
	}
	kongPluginProxyCacheConfig["response_code"] = kongPluginProxyCacheConfigResponseCode

	kongPluginProxyCaches := Plugins{
		{
			ID:      uuid.NewV4().String(),
			Name:    PluginProxyCache,
			Enabled: true,
			Config:  kongPluginProxyCacheConfig,
		},
	}

	for _, kongPluginProxyCache := range kongPluginProxyCaches {
		kongConfig = Config{}
		apisixConfig = apisix.Config{}
		kongConfig.Plugins = append(kongConfig.Plugins, kongPluginProxyCache)
		err = MigrateGlobalRules(&kongConfig, &apisixConfig)
		assert.NoError(t, err)
		assert.Equal(t, kongPluginProxyCache.ID, apisixConfig.GlobalRules[0].ID)
		assert.Equal(t, kongPluginProxyCacheConfigRequestMethods,
			apisixConfig.GlobalRules[0].Plugins.ProxyCache.CacheMethod)
		assert.Equal(t, kongPluginProxyCacheConfigResponseCodes,
			apisixConfig.GlobalRules[0].Plugins.ProxyCache.CacheHttpStatus)
	}
}
