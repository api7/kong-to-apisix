package kong

import (
	"github.com/api7/kong-to-apisix/pkg/apisix"
)

const PluginKeyAuth = "key-auth"
const PluginRateLimiting = "rate-limiting"
const PluginRateLimitingPolicyLocal = "local"
const PluginRateLimitingPolicyRedis = "redis"
const PluginProxyCache = "proxy-cache"

func KTAIsKongGlobalPlugin(kongPlugin Plugin) bool {
	return len(kongPlugin.ServiceID) <= 0 && len(kongPlugin.RouteID) <= 0 && len(kongPlugin.ConsumerID) <= 0
}

func KTAConversionKongPluginKeyAuth(kongPlugin Plugin) *apisix.KeyAuth {
	var keyName string
	config := kongPlugin.Config
	if config["key_names"] != nil {
		keyName = config["key_names"].([]interface{})[0].(string)
	}
	if len(keyName) > 0 {
		var apisixKeyAuth apisix.KeyAuth
		//apisixKeyAuth.Header = keyName
		//apisixKeyAuth.Query = keyName
		return &apisixKeyAuth
	}
	return nil
}

func KTAConversionKongPluginRateLimiting(kongPlugin Plugin) *apisix.LimitCount {
	var count int
	var timeWindow int
	var policy string
	var redisHost string
	var redisPort int
	var redisPassword string
	var redisTimeout int
	var redisDatabase int

	config := kongPlugin.Config
	if config["second"] != nil {
		count = config["second"].(int)
		timeWindow = 1
	} else if config["minute"] != nil {
		count = config["minute"].(int)
		timeWindow = 1 * 60
	} else if config["hour"] != nil {
		count = config["hour"].(int)
		timeWindow = 1 * 60 * 60
	} else if config["day"] != nil {
		count = config["day"].(int)
		timeWindow = 1 * 60 * 60 * 24
	} else if config["month"] != nil {
		count = config["month"].(int)
		timeWindow = 1 * 60 * 60 * 24 * 30
	} else if config["year"] != nil {
		count = config["year"].(int)
		timeWindow = 1 * 60 * 60 * 24 * 30 * 365
	}

	if config["policy"] != nil {
		policy = config["policy"].(string)
		switch policy {
		case PluginRateLimitingPolicyRedis:
			if config["redis_host"] != nil {
				redisHost = config["redis_host"].(string)
			}
			if config["redis_port"] != nil {
				redisPort = config["redis_port"].(int)
			}
			if config["redis_password"] != nil {
				redisPassword = config["redis_password"].(string)
			}
			if config["redis_timeout"] != nil {
				redisTimeout = config["redis_timeout"].(int)
			}
			if config["redis_database"] != nil {
				redisDatabase = config["redis_database"].(int)
			}
		default:
			// other type reset to local
			policy = PluginRateLimitingPolicyLocal
		}
	}

	if timeWindow > 0 && count > 0 && len(policy) > 0 {
		var apisixLimitCount apisix.LimitCount
		apisixLimitCount.Count = count
		apisixLimitCount.TimeWindow = timeWindow
		apisixLimitCount.Policy = policy
		apisixLimitCount.RejectedCode = 429
		if policy == PluginRateLimitingPolicyRedis {
			apisixLimitCount.RedisHost = redisHost
			apisixLimitCount.RedisPort = redisPort
			apisixLimitCount.RedisPassword = redisPassword
			apisixLimitCount.RedisTimeout = redisTimeout
			apisixLimitCount.RedisDatabase = redisDatabase
		}
		return &apisixLimitCount
	}
	return nil
}

func KTAConversionKongPluginProxyCache(kongPlugin Plugin) *apisix.ProxyCache {
	var cacheMethod []string
	var cacheHttpStatus []int
	config := kongPlugin.Config
	if config["request_method"] != nil {
		for _, method := range config["request_method"].([]interface{}) {
			cacheMethod = append(cacheMethod, method.(string))
		}
	}
	if config["response_code"] != nil {
		for _, code := range config["response_code"].([]interface{}) {
			cacheHttpStatus = append(cacheHttpStatus, code.(int))
		}
	}

	if len(cacheMethod) <= 0 && len(cacheHttpStatus) <= 0 {
		return nil
	}
	var apisixProxyCache apisix.ProxyCache
	if len(cacheMethod) >= 1 {
		apisixProxyCache.CacheMethod = cacheMethod
	}
	if len(cacheHttpStatus) >= 1 {
		apisixProxyCache.CacheHttpStatus = cacheHttpStatus
	}
	return &apisixProxyCache
}
