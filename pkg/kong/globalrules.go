package kong

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/api7/kong-to-apisix/pkg/apisix"
)

// MigrateGlobalRules attention to plugin priority
// https://docs.konghq.com/gateway-oss/2.4.x/admin-api/#precedence
func MigrateGlobalRules(kongConfig *Config, apisixConfig *apisix.Config) error {
	kongPlugins := kongConfig.Plugins
	for _, plugin := range kongPlugins {
		if !KTAIsKongGlobalPlugin(plugin) && plugin.Enabled {
			continue
		}
		var apisixGlobalRule apisix.GlobalRule
		if len(plugin.ID) > 0 {
			apisixGlobalRule.ID = plugin.ID
		} else {
			apisixGlobalRule.ID = uuid.NewV4().String()
		}

		switch plugin.Name {
		case PluginKeyAuth:
			apisixGlobalRule.Plugins.KeyAuth = KTAConversionKongPluginKeyAuth(plugin)
			if apisixGlobalRule.Plugins.KeyAuth == nil {
				fmt.Printf("Kong global plugin %s [ %s ] configuration invalid\n", plugin.Name,
					apisixGlobalRule.ID)
				continue
			}
		case PluginRateLimiting:
			apisixGlobalRule.Plugins.LimitCount = KTAConversionKongPluginRateLimiting(plugin)
			if apisixGlobalRule.Plugins.LimitCount == nil {
				fmt.Printf("Kong global plugin %s [ %s ] configuration invalid\n", plugin.Name,
					apisixGlobalRule.ID)
				continue
			}
		case PluginProxyCache:
			apisixGlobalRule.Plugins.ProxyCache = KTAConversionKongPluginProxyCache(plugin)
			if apisixGlobalRule.Plugins.ProxyCache == nil {
				fmt.Printf("Kong global plugin %s [ %s ] configuration invalid\n", plugin.Name,
					apisixGlobalRule.ID)
				continue
			}
		default:
			fmt.Printf("Kong global plugin %s [ %s ] not supported by apisix yet\n", plugin.Name,
				apisixGlobalRule.ID)
			continue
		}
		apisixConfig.GlobalRules = append(apisixConfig.GlobalRules, apisixGlobalRule)
		fmt.Printf("Kong global plugin %s [ %s ] to APISIX conversion completed\n", plugin.Name,
			apisixGlobalRule.ID)
	}
	fmt.Println("Kong to APISIX global plugins configuration conversion completed")
	return nil
}
