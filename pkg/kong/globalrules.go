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
	for kgpIndex, kongGlobalPlugin := range kongPlugins {
		if !kongGlobalPlugin.Enabled {
			continue
		}

		if !KTAIsKongGlobalPlugin(kongGlobalPlugin) {
			continue
		}

		kongGlobalPluginId := kongGlobalPlugin.ID
		if len(kongGlobalPlugin.ID) <= 0 {
			kongGlobalPluginId = uuid.NewV4().String()
			kongConfig.Plugins[kgpIndex].ID = kongGlobalPluginId
		}

		var apisixGlobalRule apisix.GlobalRule
		apisixGlobalRule.ID = kongGlobalPluginId
		switch kongGlobalPlugin.Name {
		case PluginKeyAuth:
			apisixGlobalRule.Plugins.KeyAuth = KTAConversionKongPluginKeyAuth(kongGlobalPlugin)
			if apisixGlobalRule.Plugins.KeyAuth == nil {
				fmt.Printf("Kong global plugin %s [ %s ] configuration invalid\n", kongGlobalPlugin.Name,
					apisixGlobalRule.ID)
				continue
			}
		case PluginRateLimiting:
			apisixGlobalRule.Plugins.LimitCount = KTAConversionKongPluginRateLimiting(kongGlobalPlugin)
			if apisixGlobalRule.Plugins.LimitCount == nil {
				fmt.Printf("Kong global plugin %s [ %s ] configuration invalid\n", kongGlobalPlugin.Name,
					apisixGlobalRule.ID)
				continue
			}
		case PluginProxyCache:
			apisixGlobalRule.Plugins.ProxyCache = KTAConversionKongPluginProxyCache(kongGlobalPlugin)
			if apisixGlobalRule.Plugins.ProxyCache == nil {
				fmt.Printf("Kong global plugin %s [ %s ] configuration invalid\n", kongGlobalPlugin.Name,
					apisixGlobalRule.ID)
				continue
			}
		default:
			fmt.Printf("Kong global plugin %s [ %s ] not supported by apisix yet\n", kongGlobalPlugin.Name,
				apisixGlobalRule.ID)
			continue
		}

		apisixConfig.GlobalRules = append(apisixConfig.GlobalRules, apisixGlobalRule)
		fmt.Printf("Kong global plugin %s [ %s ] to APISIX conversion completed\n", kongGlobalPlugin.Name,
			apisixGlobalRule.ID)
	}

	fmt.Println("Kong to APISIX global plugins configuration conversion completed")
	return nil
}
