package kong

import (
	"fmt"
	"strconv"

	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
	"github.com/api7/kongtoapisix/pkg/utils"
)

func MigrateRoute(kongConfig *KongConfig, configYamlAll *[]utils.YamlItem) (*[]v1.Route, error) {
	kongServices := kongConfig.Services

	var apisixRoutes []v1.Route
	i := 0
	for _, s := range kongServices {
		for _, r := range s.Routes {
			i++
			apisixRoute := &v1.Route{
				Metadata: v1.Metadata{
					ID: strconv.Itoa(i),
				},
				UpstreamId: s.ID,
				// TODO: need to check if it's the same
				Priority: r.RegexPriority,
				Plugins:  v1.Plugins{},
			}

			if r.Name != "" {
				apisixRoute.Metadata.Name = r.Name
			}

			// TODO: need to tweak between different rules of apisix and kong later
			if len(r.Paths) == 1 {
				apisixRoute.Uri = r.Paths[0] + "*"
			} else {
				var uris []string
				for _, p := range r.Paths {
					uris = append(uris, p+"*")
				}
				apisixRoute.Uris = uris
			}

			if len(r.Hosts) == 1 {
				apisixRoute.Host = r.Hosts[0]
			} else {
				var hosts []string
				for _, h := range r.Hosts {
					hosts = append(hosts, h)
				}
				apisixRoute.Hosts = hosts
			}

			// since proxy-rewrite could only support one line for regex-uri
			// need to split it to several routes if match uris
			if r.StripPath && apisixRoute.Uri != "" {
				err := addProxyRewrite(apisixRoute)
				if err != nil {
					return nil, err
				}
			}

			var methods []string
			for _, m := range r.Methods {
				methods = append(methods, m)
			}
			apisixRoute.Methods = methods

			plugins := apisixRoute.Plugins
			for _, p := range r.Plugins {
				if f, ok := pluginMap[p.Name]; ok {
					if p.Enabled {
						if apisixPlugin, configYaml, err := f(p); err != nil {
							return nil, err
						} else {
							for k, v := range apisixPlugin {
								plugins[k] = v
							}
							for _, c := range configYaml {
								*configYamlAll = append(*configYamlAll, c)
							}
						}
					}
				}
			}
			apisixRoute.Plugins = plugins

			apisixRoutes = append(apisixRoutes, *apisixRoute)
		}
	}

	return &apisixRoutes, nil
}

func addProxyRewrite(route *v1.Route) error {
	pluginConfig := make(map[string]interface{})
	pluginConfig["regex_uri"] = []string{fmt.Sprintf(`^%s/?(.*)`, route.Uri[:len(route.Uri)-1]), "/$1"}

	route.Plugins["proxy-rewrite"] = pluginConfig

	return nil
}
