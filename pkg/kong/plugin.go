package kong

import (
	"fmt"

	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
	"github.com/api7/kongtoapisix/pkg/utils"
)

var pluginMap = map[string]func(p Plugin) (v1.Plugins, []utils.YamlItem, error){
	"proxy-cache":   proxyCache,
	"key-auth":      keyAuth,
	"rate-limiting": rateLimiting,
}

// TODO: some configuration need to be configured in config.yaml
//       including cache_ttl
func proxyCache(p Plugin) (v1.Plugins, []utils.YamlItem, error) {
	pluginConfig := make(map[string]interface{})
	pluginConfig["cache_method"] = p.Config["request_method"]
	pluginConfig["cache_http_status"] = p.Config["response_code"]

	var configYaml []utils.YamlItem
	if cacheTTL, ok := p.Config["cache_ttl"].(int); ok {
		configYaml = append(configYaml, utils.YamlItem{
			Value: fmt.Sprintf("%v", cacheTTL) + "s",
			Path:  []interface{}{"apisix", "proxy_cache", "cache_ttl"},
		})
	}

	return v1.Plugins{"proxy-cache": pluginConfig}, configYaml, nil
}

// TODO: kong could configure rate limiting in different time range
//       for now fetch the value in minimum time range
func rateLimiting(p Plugin) (v1.Plugins, []utils.YamlItem, error) {
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
	return v1.Plugins{"limit-count": pluginConfig}, nil, nil
}

func keyAuth(p Plugin) (v1.Plugins, []utils.YamlItem, error) {
	pluginConfig := make(map[string]interface{})
	pluginConfig["key"] = p.Config["key_names"]

	emptyMap := make(map[string]interface{})

	return v1.Plugins{"key-auth": emptyMap}, nil, nil
}
