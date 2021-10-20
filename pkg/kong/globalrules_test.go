package kong

import (
	"testing"

	"github.com/api7/kong-to-apisix/pkg/apisix"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type RateLimitingExpectResult struct {
	Count      int
	TimeWindow int
	Policy     string
}

type RateLimitingTest struct {
	RateLimitingPluginConfig Plugin
	RateLimitingExpectResult RateLimitingExpectResult
}

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

	kongPluginRateLimitingTests := []RateLimitingTest{
		{
			RateLimitingPluginConfig: Plugin{
				ID:      uuid.NewV4().String(),
				Name:    PluginRateLimiting,
				Enabled: true,
				Config: map[string]interface{}{
					"second": 5,
					"policy": PluginRateLimitingPolicyLocal,
				},
			},
			RateLimitingExpectResult: RateLimitingExpectResult{
				Count:      5,
				TimeWindow: 1,
				Policy:     PluginRateLimitingPolicyLocal,
			},
		},
		{
			RateLimitingPluginConfig: Plugin{
				ID:      uuid.NewV4().String(),
				Name:    PluginRateLimiting,
				Enabled: true,
				Config: map[string]interface{}{
					"minute": 10,
					"policy": PluginRateLimitingPolicyLocal,
				},
			},
			RateLimitingExpectResult: RateLimitingExpectResult{
				Count:      10,
				TimeWindow: 1 * 60,
				Policy:     PluginRateLimitingPolicyLocal,
			},
		},
		{
			RateLimitingPluginConfig: Plugin{
				ID:      uuid.NewV4().String(),
				Name:    PluginRateLimiting,
				Enabled: true,
				Config: map[string]interface{}{
					"hour":   15,
					"policy": PluginRateLimitingPolicyLocal,
				},
			},
			RateLimitingExpectResult: RateLimitingExpectResult{
				Count:      15,
				TimeWindow: 1 * 60 * 60,
				Policy:     PluginRateLimitingPolicyLocal,
			},
		},
		{
			RateLimitingPluginConfig: Plugin{
				ID:      uuid.NewV4().String(),
				Name:    PluginRateLimiting,
				Enabled: true,
				Config: map[string]interface{}{
					"day":    20,
					"policy": PluginRateLimitingPolicyLocal,
				},
			},
			RateLimitingExpectResult: RateLimitingExpectResult{
				Count:      20,
				TimeWindow: 1 * 60 * 60 * 24,
				Policy:     PluginRateLimitingPolicyLocal,
			},
		},
		{
			RateLimitingPluginConfig: Plugin{
				ID:      uuid.NewV4().String(),
				Name:    PluginRateLimiting,
				Enabled: true,
				Config: map[string]interface{}{
					"month":  25,
					"policy": PluginRateLimitingPolicyLocal,
				},
			},
			RateLimitingExpectResult: RateLimitingExpectResult{
				Count:      25,
				TimeWindow: 1 * 60 * 60 * 24 * 30,
				Policy:     PluginRateLimitingPolicyLocal,
			},
		},
		{
			RateLimitingPluginConfig: Plugin{
				ID:      uuid.NewV4().String(),
				Name:    PluginRateLimiting,
				Enabled: true,
				Config: map[string]interface{}{
					"year":   5,
					"policy": PluginRateLimitingPolicyLocal,
				},
			},
			RateLimitingExpectResult: RateLimitingExpectResult{
				Count:      5,
				TimeWindow: 1 * 60 * 60 * 24 * 30 * 365,
				Policy:     PluginRateLimitingPolicyLocal,
			},
		},
	}

	for _, kongPluginRateLimitingTest := range kongPluginRateLimitingTests {
		kongConfig = Config{}
		apisixConfig = apisix.Config{}
		kongConfig.Plugins = append(kongConfig.Plugins, kongPluginRateLimitingTest.RateLimitingPluginConfig)
		err = MigrateGlobalRules(&kongConfig, &apisixConfig)
		assert.NoError(t, err)
		assert.Equal(t, kongPluginRateLimitingTest.RateLimitingPluginConfig.ID, apisixConfig.GlobalRules[0].ID)
		assert.Equal(t, kongPluginRateLimitingTest.RateLimitingExpectResult.Policy,
			apisixConfig.GlobalRules[0].Plugins.LimitCount.Policy)
		assert.Equal(t, kongPluginRateLimitingTest.RateLimitingExpectResult.Count,
			apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
		assert.Equal(t, kongPluginRateLimitingTest.RateLimitingExpectResult.TimeWindow,
			apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)
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
