package kong

import (
	"fmt"
	"strconv"

	"github.com/api7/kong-to-apisix/pkg/apisix"
	"github.com/api7/kong-to-apisix/pkg/utils"
)

func MigrateRoute(kongConfig *Config, configYamlAll *[]utils.YamlItem) (apisix.Routes, error) {
	kongServices := kongConfig.Services

	var apisixRoutes apisix.Routes
	i := 0
	for _, s := range kongServices {
		for _, r := range s.Routes {
			i++
			apisixRoute := &apisix.Route{
				ID:        strconv.Itoa(i),
				ServiceID: s.ID,
				// TODO: need to check if it's the same
				Priority: uint(r.RegexPriority),
				Plugins:  apisix.Plugins{},
			}

			if r.Name != "" {
				apisixRoute.Name = r.Name
			}

			// TODO: need to tweak between different rules of apisix and kong later
			if len(r.Paths) == 1 {
				apisixRoute.URI = r.Paths[0] + "*"
			} else {
				var uris []string
				for _, p := range r.Paths {
					uris = append(uris, p+"*")
				}
				apisixRoute.URIs = uris
			}

			if len(r.Hosts) == 1 {
				apisixRoute.Host = r.Hosts[0]
			} else {
				var hosts []string
				for hostIndex := range r.Hosts {
					hosts = append(hosts, r.Hosts[hostIndex])
				}
				apisixRoute.Hosts = hosts
			}

			// since proxy-rewrite could only support one line for regex-uri
			// need to split it to several routes if match uris
			if r.StripPath && apisixRoute.URI != "" {
				proxyRewrite := addProxyRewrite(apisixRoute)
				if proxyRewrite != nil {
					apisixRoute.Plugins.ProxyRewrite = proxyRewrite
				}
			}

			var methods []string
			for methodIndex := range r.Methods {
				methods = append(methods, r.Methods[methodIndex])
			}
			apisixRoute.Methods = methods

			for _, p := range r.Plugins {
				if !p.Enabled {
					continue
				}
				switch p.Name {
				case PluginKeyAuth:
					apisixRoute.Plugins.KeyAuth = KTAConversionKongPluginKeyAuth(p)
				case PluginProxyCache:
					apisixRoute.Plugins.ProxyCache = KTAConversionKongPluginProxyCache(p)
				case PluginRateLimiting:
					apisixRoute.Plugins.LimitCount = KTAConversionKongPluginRateLimiting(p)
				default:
					fmt.Printf("Kong route [%s] plugin %s not supported by apisix yet\n", r.ID, p.Name)
				}
			}

			apisixRoutes = append(apisixRoutes, *apisixRoute)
		}
	}

	return apisixRoutes, nil
}

func addProxyRewrite(route *apisix.Route) *apisix.ProxyRewrite {
	var proxyRewrite apisix.ProxyRewrite
	proxyRewrite.RegexURI = []string{fmt.Sprintf(`^%s/?(.*)`, route.URI[:len(route.URI)-1]), "/$1"}
	return &proxyRewrite
}
