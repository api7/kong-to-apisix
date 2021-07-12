package plugin

import (
	"context"
	"fmt"

	"github.com/apache/apisix-ingress-controller/pkg/apisix"
	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
	"github.com/api7/kongtoapisix/pkg/utils"
	"github.com/globocom/gokong"
)

var funcMap = map[string]func(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient, p *gokong.Plugin) error{
	"proxy-cache":   proxyCache,
	"key-auth":      keyAuth,
	"rate-limiting": rateLimiting,
}

// TODO: need to take care of plugin precedence
// https://docs.konghq.com/gateway-oss/2.4.x/admin-api/#precedence
func MigratePlugin(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient) error {
	plugins, err := kongCli.Plugins().List(&gokong.PluginQueryString{})
	if err != nil {
		return err
	}
	for _, p := range plugins {
		//fmt.Printf("got plugin: %#v\n", p)
		if f, ok := funcMap[p.Name]; ok {
			if p.Enabled {
				err := f(apisixCli, kongCli, p)
				if err != nil {
					return err
				}
				fmt.Printf("migrate plugin %s succeeds\n", p.Name)
			} else {
				fmt.Printf("Plugin %s not enabled\n", p.Name)
			}
		} else {
			fmt.Printf("Plugin %s not supported by apisix yet\n", p.Name)
		}
	}

	return nil
}

// TODO: some configuration need to be configured in config.yaml
//       including cache_ttl
func proxyCache(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient, p *gokong.Plugin) error {
	if p.ConsumerId == nil && p.ServiceId == nil && p.RouteId == nil {
		routes, err := apisixCli.Route().List(context.Background())
		if err != nil {
			return err
		}
		for _, r := range routes {
			pluginConfig := make(map[string]interface{})
			pluginConfig["cache_method"] = p.Config["request_method"]
			pluginConfig["cache_http_status"] = p.Config["response_code"]
			if cacheTTL, ok := p.Config["cache_ttl"].(float64); ok {
				err := utils.AddValueToYaml(fmt.Sprintf("%v", int(cacheTTL))+"s", "apisix", "proxy_cache", "cache_ttl")
				if err != nil {
					return err
				}
			}

			if r.Plugins == nil {
				r.Plugins = v1.Plugins{
					"proxy-cache": pluginConfig,
				}
			} else {
				r.Plugins["proxy-cache"] = pluginConfig
			}

			_, err := apisixCli.Route().Update(context.Background(), r)
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Println("not support non-global plugins for now")
	}
	return nil
}

// TODO: kong could configure rate limiting in different time range
//       for now fetch the value in minimum time range
func rateLimiting(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient, p *gokong.Plugin) error {
	if p.ConsumerId == nil && p.ServiceId == nil && p.RouteId == nil {
		routes, err := apisixCli.Route().List(context.Background())
		if err != nil {
			return err
		}
		for _, r := range routes {
			pluginConfig := make(map[string]interface{})

			if p.Config["second"] != nil {
				pluginConfig["count"] = p.Config["second"]
				pluginConfig["time_window"] = 1
			} else if p.Config["minute"] != nil {
				pluginConfig["count"] = p.Config["minute"]
				pluginConfig["time_window"] = 1 * 60
			} else if p.Config["hour"] != nil {
				pluginConfig["count"] = p.Config["hour"]
				pluginConfig["time_window"] = 1 * 60 * 60
			} else if p.Config["day"] != nil {
				pluginConfig["count"] = p.Config["day"]
				pluginConfig["time_window"] = 1 * 60 * 60 * 24
			} else if p.Config["month"] != nil {
				pluginConfig["count"] = p.Config["day"]
				pluginConfig["time_window"] = 1 * 60 * 60 * 24 * 30
			} else if p.Config["year"] != nil {
				pluginConfig["count"] = p.Config["day"]
				pluginConfig["time_window"] = 1 * 60 * 60 * 24 * 30 * 365
			}

			switch p.Config["policy"] {
			case "local":
				pluginConfig["policy"] = "local"
			case "cluster":
				fmt.Println(`Convert rate limit policy from 'kong cluster' to 'local'\n
					suggest to use redis to achieve global rate limiting.`)
				pluginConfig["policy"] = "local"
			case "redis":
				pluginConfig["policy"] = "redis"
				pluginConfig["redis_host"] = p.Config["redis_host"]
				pluginConfig["redis_port"] = p.Config["redis_port"]
				pluginConfig["redis_password"] = p.Config["redis_password"]
				pluginConfig["redis_timeout"] = p.Config["redis_timeout"]
				pluginConfig["redis_database"] = p.Config["redis_database"]
			}

			pluginConfig["rejected_code"] = 429
			if r.Plugins == nil {
				r.Plugins = v1.Plugins{
					"limit-count": pluginConfig,
				}
			} else {
				r.Plugins["limit-count"] = pluginConfig
			}

			_, err := apisixCli.Route().Update(context.Background(), r)
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Println("not support non-global plugins for now")
	}
	return nil
}

func keyAuth(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient, p *gokong.Plugin) error {
	err := updateKeyAuthForConsumers(apisixCli, kongCli, p)
	if err != nil {
		return err
	}

	if p.ConsumerId == nil && p.ServiceId == nil && p.RouteId != nil {

		r, err := apisixCli.Route().Get(context.Background(), string(*p.RouteId))
		if err != nil {
			return err
		}

		pluginConfig := make(map[string]interface{})
		pluginConfig["key"] = p.Config["key_names"]

		emptyMap := make(map[string]interface{})

		if r.Plugins == nil {
			r.Plugins = v1.Plugins{
				"key-auth": emptyMap,
			}
		} else {
			r.Plugins["key-auth"] = emptyMap
		}

		_, err = apisixCli.Route().Update(context.Background(), r)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("not support non-global plugins for now")
	}
	return nil
}

func updateKeyAuthForConsumers(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient, p *gokong.Plugin) error {
	consumers, err := kongCli.Consumers().List(&gokong.ConsumerQueryString{})
	if err != nil {
		return err
	}

	for _, c := range consumers {
		pcs, err := kongCli.Consumers().GetPluginConfigs(c.Id, "key-auth")
		if err != nil {
			return err
		}
		// TODO: not sure one consumer could bind to several key-auth keys
		pluginConfig := make(map[string]interface{})
		pluginConfig["key"] = pcs[0]["key"]

		plugins := v1.Plugins{
			"key-auth": pluginConfig,
		}

		apisixConsumer := &v1.Consumer{
			Username: c.Username,
			Plugins:  plugins,
		}
		_, err = apisixCli.Consumer().Update(context.Background(), apisixConsumer)
		if err != nil {
			return err
		}
	}
	return nil
}

func getKongRouteNameWithId(kongCli gokong.KongAdminClient, id string) (string, error) {
	route, err := kongCli.Routes().GetById(id)
	if err != nil {
		return "", err
	}
	return *route.Name, nil
}
