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
	var kongPluginKeyAuth Plugin
	kongPluginKeyAuth.ID = uuid.NewV4().String()
	kongPluginKeyAuth.Name = PluginKeyAuth
	kongPluginKeyAuth.Enabled = true
	kongPluginKeyAuthConfig := make(map[string]interface{})
	kongPluginKeyAuthConfig["key_names"] = make([]interface{}, 1)
	kongPluginKeyAuthConfigKeyNames := make([]interface{}, 1)
	for i, v := range []string{"apikey"} {
		kongPluginKeyAuthConfigKeyNames[i] = v
	}
	kongPluginKeyAuthConfig["key_names"] = kongPluginKeyAuthConfigKeyNames
	kongPluginKeyAuth.Config = kongPluginKeyAuthConfig
	kongConfig.Plugins = append(kongConfig.Plugins, kongPluginKeyAuth)
	err = MigrateGlobalRules(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	assert.Equal(t, kongPluginKeyAuth.ID, apisixConfig.GlobalRules[0].ID)

	kongConfig = Config{}
	apisixConfig = apisix.Config{}
	var kongPluginRateLimitingSecond Plugin
	kongPluginRateLimitingSecond.ID = uuid.NewV4().String()
	kongPluginRateLimitingSecond.Name = PluginRateLimiting
	kongPluginRateLimitingSecond.Enabled = true
	kongPluginRateLimitingSecondConfig := make(map[string]interface{})
	kongPluginRateLimitingSecondConfig["second"] = 5
	kongPluginRateLimitingSecondConfig["policy"] = PluginRateLimitingPolicyLocal
	kongPluginRateLimitingSecond.Config = kongPluginRateLimitingSecondConfig
	kongConfig.Plugins = append(kongConfig.Plugins, kongPluginRateLimitingSecond)
	err = MigrateGlobalRules(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	assert.Equal(t, kongPluginRateLimitingSecond.ID, apisixConfig.GlobalRules[0].ID)
	assert.Equal(t, kongPluginRateLimitingSecondConfig["second"].(int),
		apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
	assert.Equal(t, 1, apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)

	kongConfig = Config{}
	apisixConfig = apisix.Config{}
	var kongPluginRateLimitingMinute Plugin
	kongPluginRateLimitingMinute.ID = uuid.NewV4().String()
	kongPluginRateLimitingMinute.Name = PluginRateLimiting
	kongPluginRateLimitingMinute.Enabled = true
	kongPluginRateLimitingMinuteConfig := make(map[string]interface{})
	kongPluginRateLimitingMinuteConfig["minute"] = 5
	kongPluginRateLimitingMinuteConfig["policy"] = PluginRateLimitingPolicyLocal
	kongPluginRateLimitingMinute.Config = kongPluginRateLimitingMinuteConfig
	kongConfig.Plugins = append(kongConfig.Plugins, kongPluginRateLimitingMinute)
	err = MigrateGlobalRules(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	assert.Equal(t, kongPluginRateLimitingMinute.ID, apisixConfig.GlobalRules[0].ID)
	assert.Equal(t, kongPluginRateLimitingMinuteConfig["minute"].(int),
		apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
	assert.Equal(t, 1*60, apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)

	kongConfig = Config{}
	apisixConfig = apisix.Config{}
	var kongPluginRateLimitingHour Plugin
	kongPluginRateLimitingHour.ID = uuid.NewV4().String()
	kongPluginRateLimitingHour.Name = PluginRateLimiting
	kongPluginRateLimitingHour.Enabled = true
	kongPluginRateLimitingHourConfig := make(map[string]interface{})
	kongPluginRateLimitingHourConfig["hour"] = 5
	kongPluginRateLimitingHourConfig["policy"] = PluginRateLimitingPolicyLocal
	kongPluginRateLimitingHour.Config = kongPluginRateLimitingHourConfig
	kongConfig.Plugins = append(kongConfig.Plugins, kongPluginRateLimitingHour)
	err = MigrateGlobalRules(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	assert.Equal(t, kongPluginRateLimitingHour.ID, apisixConfig.GlobalRules[0].ID)
	assert.Equal(t, kongPluginRateLimitingHourConfig["hour"].(int),
		apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
	assert.Equal(t, 1*60*60, apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)

	kongConfig = Config{}
	apisixConfig = apisix.Config{}
	var kongPluginRateLimitingDay Plugin
	kongPluginRateLimitingDay.ID = uuid.NewV4().String()
	kongPluginRateLimitingDay.Name = PluginRateLimiting
	kongPluginRateLimitingDay.Enabled = true
	kongPluginRateLimitingDayConfig := make(map[string]interface{})
	kongPluginRateLimitingDayConfig["day"] = 5
	kongPluginRateLimitingDayConfig["policy"] = PluginRateLimitingPolicyLocal
	kongPluginRateLimitingDay.Config = kongPluginRateLimitingDayConfig
	kongConfig.Plugins = append(kongConfig.Plugins, kongPluginRateLimitingDay)
	err = MigrateGlobalRules(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	assert.Equal(t, kongPluginRateLimitingDay.ID, apisixConfig.GlobalRules[0].ID)
	assert.Equal(t, kongPluginRateLimitingDayConfig["day"].(int),
		apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
	assert.Equal(t, 1*60*60*24, apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)

	kongConfig = Config{}
	apisixConfig = apisix.Config{}
	var kongPluginRateLimitingMonth Plugin
	kongPluginRateLimitingMonth.ID = uuid.NewV4().String()
	kongPluginRateLimitingMonth.Name = PluginRateLimiting
	kongPluginRateLimitingMonth.Enabled = true
	kongPluginRateLimitingMonthConfig := make(map[string]interface{})
	kongPluginRateLimitingMonthConfig["month"] = 5
	kongPluginRateLimitingMonthConfig["policy"] = PluginRateLimitingPolicyLocal
	kongPluginRateLimitingMonth.Config = kongPluginRateLimitingMonthConfig
	kongConfig.Plugins = append(kongConfig.Plugins, kongPluginRateLimitingMonth)
	err = MigrateGlobalRules(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	assert.Equal(t, kongPluginRateLimitingMonth.ID, apisixConfig.GlobalRules[0].ID)
	assert.Equal(t, kongPluginRateLimitingMonthConfig["month"].(int),
		apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
	assert.Equal(t, 1*60*60*24*30, apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)

	kongConfig = Config{}
	apisixConfig = apisix.Config{}
	var kongPluginRateLimitingYear Plugin
	kongPluginRateLimitingYear.ID = uuid.NewV4().String()
	kongPluginRateLimitingYear.Name = PluginRateLimiting
	kongPluginRateLimitingYear.Enabled = true
	kongPluginRateLimitingYearConfig := make(map[string]interface{})
	kongPluginRateLimitingYearConfig["year"] = 5
	kongPluginRateLimitingYearConfig["policy"] = PluginRateLimitingPolicyLocal
	kongPluginRateLimitingYear.Config = kongPluginRateLimitingYearConfig
	kongConfig.Plugins = append(kongConfig.Plugins, kongPluginRateLimitingYear)
	err = MigrateGlobalRules(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	assert.Equal(t, kongPluginRateLimitingYear.ID, apisixConfig.GlobalRules[0].ID)
	assert.Equal(t, kongPluginRateLimitingYearConfig["year"].(int),
		apisixConfig.GlobalRules[0].Plugins.LimitCount.Count)
	assert.Equal(t, 1*60*60*24*30*365, apisixConfig.GlobalRules[0].Plugins.LimitCount.TimeWindow)

	kongConfig = Config{}
	apisixConfig = apisix.Config{}
	var kongPluginProxyCache Plugin
	kongPluginProxyCache.ID = uuid.NewV4().String()
	kongPluginProxyCache.Name = PluginProxyCache
	kongPluginProxyCache.Enabled = true
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
	kongPluginProxyCache.Config = kongPluginProxyCacheConfig
	kongConfig.Plugins = append(kongConfig.Plugins, kongPluginProxyCache)
	err = MigrateGlobalRules(&kongConfig, &apisixConfig)
	assert.NoError(t, err)
	assert.Equal(t, kongPluginProxyCache.ID, apisixConfig.GlobalRules[0].ID)
	assert.Equal(t, kongPluginProxyCacheConfigRequestMethods,
		apisixConfig.GlobalRules[0].Plugins.ProxyCache.CacheMethod)
	assert.Equal(t, kongPluginProxyCacheConfigResponseCodes,
		apisixConfig.GlobalRules[0].Plugins.ProxyCache.CacheHttpStatus)
}
