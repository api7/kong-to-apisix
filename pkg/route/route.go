package route

import (
	"context"
	"fmt"

	"github.com/apache/apisix-ingress-controller/pkg/apisix"
	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
	"github.com/globocom/gokong"
)

func MigrateRoute(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient) error {
	routes, err := kongCli.Routes().List(&gokong.RouteQueryString{})
	if err != nil {
		return err
	}

	for _, r := range routes {
		//fmt.Printf("got route: %#v\n", *r.Name)

		apisixRoute := &v1.Route{
			Metadata: v1.Metadata{
				ID: *r.Id,
			},
			UpstreamId: string(*r.Service),
			// TODO: need to check if it's the same
			Priority: *r.RegexPriority,
		}

		if r.Name != nil {
			apisixRoute.Metadata.Name = *r.Name
		}

		// TODO: need to tweak between different rules of apisix and kong later
		if len(r.Paths) == 1 {
			apisixRoute.Uri = *r.Paths[0] + "*"
		} else {
			var uris []string
			for _, p := range r.Paths {
				uris = append(uris, *p+"*")
			}
			apisixRoute.Uris = uris
		}

		if len(r.Hosts) == 1 {
			apisixRoute.Host = *r.Hosts[0]
		} else {
			var hosts []string
			for _, p := range r.Hosts {
				hosts = append(hosts, *p)
			}
			apisixRoute.Hosts = hosts
		}

		// since proxy-rewrite could only support one line for regex-uri
		// need to split it to several routes if match uris
		if *r.StripPath && apisixRoute.Uri != "" {
			err := addProxyRewrite(apisixCli, kongCli, apisixRoute)
			if err != nil {
				return err
			}
		}

		var methods []string
		for _, m := range r.Methods {
			methods = append(methods, *m)
		}
		apisixRoute.Methods = methods

		_, err := apisixCli.Route().Create(context.Background(), apisixRoute)
		if err != nil {
			return err
		}

		var printName string
		if r.Name != nil {
			printName = *r.Name
		} else {
			printName = *r.Id
		}
		fmt.Printf("migrate route %s succeeds\n", printName)
	}

	return nil
}

func addProxyRewrite(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient, route *v1.Route) error {
	pluginConfig := make(map[string]interface{})
	pluginConfig["regex_uri"] = []string{fmt.Sprintf(`^%s/?(.*)`, route.Uri[:len(route.Uri)-1]), "/$1"}

	route.Plugins = v1.Plugins{
		"proxy-rewrite": pluginConfig,
	}

	return nil
}
