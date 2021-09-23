package kong

import (
	"fmt"
	"github.com/api7/kong-to-apisix/pkg/apisix"
	"github.com/api7/kong-to-apisix/pkg/utils"
)

// TODO: need to take care of plugin precedence
// https://docs.konghq.com/gateway-oss/2.4.x/admin-api/#precedence
func MigrateGlobalRules(kongConfig *KongConfig, configYamlAll *[]utils.YamlItem) (*[]apisix.GlobalRule, error) {
	kongGlobalPlugins := kongConfig.Plugins
	var apisixGlobalRules []apisix.GlobalRule

	for _, p := range kongGlobalPlugins {
		//fmt.Printf("got plugin: %#v\n", p)
		if f, ok := pluginMap[p.Name]; ok {
			if p.Enabled {
				if apisixPlugin, configYaml, err := f(p); err != nil {
					return nil, err
				} else {
					apisixGlobalRule := apisix.GlobalRule{
						Plugins: apisix.Plugins(apisixPlugin),
					}
					apisixGlobalRules = append(apisixGlobalRules, apisixGlobalRule)
					for _, c := range configYaml {
						*configYamlAll = append(*configYamlAll, c)
					}
				}
			}
		} else {
			fmt.Printf("Plugin %s not supported by apisix yet\n", p.Name)
		}
	}
	return &apisixGlobalRules, nil
}
