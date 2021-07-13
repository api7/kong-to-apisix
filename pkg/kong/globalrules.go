package kong

import (
	"fmt"

	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
)

// TODO: need to take care of plugin precedence
// https://docs.konghq.com/gateway-oss/2.4.x/admin-api/#precedence
func MigrateGlobalRules(kongConfig *KongConfig) (*[]v1.GlobalRule, error) {
	kongGlobalPlugins := kongConfig.Plugins
	var apisixGlobalRules []v1.GlobalRule

	for _, p := range kongGlobalPlugins {
		//fmt.Printf("got plugin: %#v\n", p)
		if f, ok := pluginMap[p.Name]; ok {
			if p.Enabled {
				if apisixPlugin, err := f(p); err != nil {
					return nil, err
				} else {
					apisixGlobalRule := v1.GlobalRule{
						Plugins: apisixPlugin,
					}
					apisixGlobalRules = append(apisixGlobalRules, apisixGlobalRule)
				}
			}
		} else {
			fmt.Printf("Plugin %s not supported by apisix yet\n", p.Name)
		}
	}
	return &apisixGlobalRules, nil
}
