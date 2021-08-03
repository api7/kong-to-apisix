package kong

import (
	"fmt"

	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
	"github.com/api7/kongtoapisix/pkg/utils"
)

var pluginMap = map[string]func(p Plugin) (v1.Plugins, error){
	"proxy-cache":   proxyCache,
	"key-auth":      keyAuth,
	"rate-limiting": rateLimiting,
	"jwt":           jwt,
}

// TODO: some configuration need to be configured in config.yaml
//       including cache_ttl
func proxyCache(p Plugin) (v1.Plugins, error) {
	pluginConfig := make(map[string]interface{})
	pluginConfig["cache_method"] = p.Config["request_method"]
	pluginConfig["cache_http_status"] = p.Config["response_code"]
	if cacheTTL, ok := p.Config["cache_ttl"].(int); ok {
		err := utils.AddValueToYaml(utils.ConfigFilePath, fmt.Sprintf("%v", cacheTTL)+"s", "apisix", "proxy_cache", "cache_ttl")
		if err != nil {
			return nil, err
		}
	}

	return v1.Plugins{"proxy-cache": pluginConfig}, nil
}

// TODO: kong could configure rate limiting in different time range
//       for now fetch the value in minimum time range
func rateLimiting(p Plugin) (v1.Plugins, error) {
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
	return v1.Plugins{"limit-count": pluginConfig}, nil
}

func keyAuth(p Plugin) (v1.Plugins, error) {
	emptyMap := make(map[string]interface{})

	return v1.Plugins{"key-auth": emptyMap}, nil
}

func jwt(p Plugin) (v1.Plugins, error) {
	pluginConfig := make(map[string]interface{})
	pluginConfig["base64_secret"] = p.Config["secret_is_base64"]
	pluginConfig["exp"] = p.Config["maximum_expiration"]

	return v1.Plugins{"jwt": pluginConfig}, nil
}
