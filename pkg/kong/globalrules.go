package kong

import (
	"fmt"

	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
	"github.com/api7/kongtoapisix/pkg/utils"
)

// TODO: need to take care of plugin precedence
// https://docs.konghq.com/gateway-oss/2.4.x/admin-api/#precedence
func MigrateGlobalRules(kongConfig *KongConfig, configYamlAll *[]utils.YamlItem) (*[]v1.GlobalRule, error) {
	kongGlobalPlugins := kongConfig.Plugins
	var apisixGlobalRules []v1.GlobalRule

	for _, p := range kongGlobalPlugins {
		//fmt.Printf("got plugin: %#v\n", p)
		if f, ok := pluginMap[p.Name]; ok {
			if p.Enabled {
				if apisixPlugin, configYaml, err := f(p); err != nil {
					return nil, err
				} else {
					apisixGlobalRule := v1.GlobalRule{
						Plugins: apisixPlugin,
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
