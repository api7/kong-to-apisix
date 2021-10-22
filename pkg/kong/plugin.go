package kong

import (
	"fmt"
	"reflect"

	"github.com/api7/kong-to-apisix/pkg/apisix"

	uuid "github.com/satori/go.uuid"
)

const PluginKeyAuth = "key-auth"
const PluginRateLimiting = "rate-limiting"
const PluginRateLimitingPolicyLocal = "local"
const PluginRateLimitingPolicyRedis = "redis"
const PluginProxyCache = "proxy-cache"

// MigratePlugins This function is only called when the data exported by kong config is used
func MigratePlugins(kongConfig *Config, apisixConfig *apisix.Config) error {
	kongPlugins := kongConfig.Plugins
	for _, kongPlugin := range kongPlugins {
		if len(kongPlugin.ID) <= 0 {
			kongPlugin.ID = uuid.NewV4().String()
		}

		// If the current plugin is a global plugin, skip
		if KTAIsKongGlobalPlugin(kongPlugin) {
			continue
		}

		if !kongPlugin.Enabled {
			fmt.Printf("Kong plugin %s [ %s ] is disabled\n", kongPlugin.Name,
				kongPlugin.ID)
			continue
		}

		if len(kongPlugin.ServiceID) > 0 {
			for index, apisixService := range apisixConfig.Services {
				if apisixService.ID == kongPlugin.ServiceID {
					KTAUpdateApisixServicePlugin(&apisixService, &kongPlugin)
					apisixConfig.Services[index] = apisixService
					break
				}
			}
		}

		if len(kongPlugin.RouteID) > 0 {
			for index, apisixRoute := range apisixConfig.Routes {
				if apisixRoute.ID == kongPlugin.RouteID {
					KTAUpdateApisixRoutePlugin(&apisixRoute, &kongPlugin)
					apisixConfig.Routes[index] = apisixRoute
					break
				}
			}
		}
	}

	apisixConsumers := apisixConfig.Consumers
	for index, apisixConsumer := range apisixConsumers {
		for _, keyAuthCredential := range kongConfig.KeyAuthCredentials {
			if apisixConsumer.ID == keyAuthCredential.ConsumerID {
				KTAUpdateApisixConsumerPlugin(&apisixConsumer, &keyAuthCredential)
				apisixConfig.Consumers[index] = apisixConsumer
				break
			}
		}

		for _, basicAuthCredential := range kongConfig.BasicAuthCredentials {
			if apisixConsumer.ID == basicAuthCredential.ConsumerID {
				KTAUpdateApisixConsumerPlugin(&apisixConsumer, &basicAuthCredential)
				apisixConfig.Consumers[index] = apisixConsumer
				break
			}
		}

		for _, hmacAuthCredential := range kongConfig.HmacAuthCredentials {
			if apisixConsumer.ID == hmacAuthCredential.ConsumerID {
				KTAUpdateApisixConsumerPlugin(&apisixConsumer, &hmacAuthCredential)
				apisixConfig.Consumers[index] = apisixConsumer
				break
			}
		}

		for _, jwtSecret := range kongConfig.JwtSecrets {
			if apisixConsumer.ID == jwtSecret.ConsumerID {
				KTAUpdateApisixConsumerPlugin(&apisixConsumer, &jwtSecret)
				apisixConfig.Consumers[index] = apisixConsumer
				break
			}
		}
	}

	return nil
}

func KTAUpdateApisixServicePlugin(apisixService *apisix.Service, kongPlugin *Plugin) {
	switch kongPlugin.Name {
	case PluginKeyAuth:
		apisixService.Plugins.KeyAuth = KTAConversionKongPluginKeyAuth(*kongPlugin)
	case PluginProxyCache:
		apisixService.Plugins.ProxyCache = KTAConversionKongPluginProxyCache(*kongPlugin)
	case PluginRateLimiting:
		apisixService.Plugins.LimitCount = KTAConversionKongPluginRateLimiting(*kongPlugin)
	default:
		fmt.Printf("Kong service [%s] plugin %s not supported by apisix yet\n",
			apisixService.ID, kongPlugin.Name)
	}
}

func KTAUpdateApisixRoutePlugin(apisixRoute *apisix.Route, kongPlugin *Plugin) {
	switch kongPlugin.Name {
	case PluginKeyAuth:
		apisixRoute.Plugins.KeyAuth = KTAConversionKongPluginKeyAuth(*kongPlugin)
	case PluginProxyCache:
		apisixRoute.Plugins.ProxyCache = KTAConversionKongPluginProxyCache(*kongPlugin)
	case PluginRateLimiting:
		apisixRoute.Plugins.LimitCount = KTAConversionKongPluginRateLimiting(*kongPlugin)
	default:
		fmt.Printf("Kong route [%s] plugin %s not supported by apisix yet\n",
			apisixRoute.ID, kongPlugin.Name)
	}
}

func KTAUpdateApisixConsumerPlugin(apisixConsumer *apisix.Consumer, kongPlugin interface{}) {
	switch reflect.TypeOf(kongPlugin) {
	case reflect.TypeOf(&KeyAuthCredential{}):
		apisixConsumer.Plugins.KeyAuth =
			KTAConversionKongConsumerPluginKeyAuthCredential(*kongPlugin.(*KeyAuthCredential))
	case reflect.TypeOf(&BasicAuthCredential{}):
		apisixConsumer.Plugins.BasicAuth =
			KTAConversionKongConsumerPluginBasicAuthCredential(*kongPlugin.(*BasicAuthCredential))
	case reflect.TypeOf(&HmacAuthCredential{}):
		apisixConsumer.Plugins.HmacAuth =
			KTAConversionKongConsumerPluginHmacAuthCredential(*kongPlugin.(*HmacAuthCredential))
	case reflect.TypeOf(&JwtSecret{}):
		apisixConsumer.Plugins.JwtAuth =
			KTAConversionKongConsumerPluginJwtSecrets(*kongPlugin.(*JwtSecret))
	default:
		fmt.Printf("Kong consumer route [%s] plugin %s not supported by apisix yet\n",
			apisixConsumer.ID, reflect.TypeOf(kongPlugin).Elem().Name())
	}
}

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
		apisixKeyAuth.Header = keyName
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

func KTAConversionKongConsumerPluginKeyAuthCredential(credential KeyAuthCredential) *apisix.KeyAuthCredential {
	if len(credential.Key) > 0 {
		var apisixKeyAuthCredential apisix.KeyAuthCredential
		apisixKeyAuthCredential.Key = credential.Key
		return &apisixKeyAuthCredential
	}
	return nil
}

func KTAConversionKongConsumerPluginBasicAuthCredential(credential BasicAuthCredential) *apisix.BasicAuthCredential {
	if len(credential.Username) > 0 && len(credential.Password) > 0 {
		var apisixBasicAuthCredential apisix.BasicAuthCredential
		apisixBasicAuthCredential.Username = credential.Username
		apisixBasicAuthCredential.Password = credential.Password
		return &apisixBasicAuthCredential
	}
	return nil
}

func KTAConversionKongConsumerPluginHmacAuthCredential(credential HmacAuthCredential) *apisix.HmacAuthCredential {
	if len(credential.Username) > 0 && len(credential.Secret) > 0 {
		var apisixHmacAuthCredential apisix.HmacAuthCredential
		apisixHmacAuthCredential.AccessKey = credential.Username
		apisixHmacAuthCredential.SecretKey = credential.Secret
		return &apisixHmacAuthCredential
	}
	return nil
}

func KTAConversionKongConsumerPluginJwtSecrets(secret JwtSecret) *apisix.JwtSecrets {
	if len(secret.Key) > 0 && len(secret.Secret) > 0 {
		var apisixJwtSecrets apisix.JwtSecrets
		apisixJwtSecrets.Key = secret.Key
		apisixJwtSecrets.Secret = secret.Secret
		if len(secret.Algorithm) > 0 {
			apisixJwtSecrets.Algorithm = secret.Algorithm
		}
		return &apisixJwtSecrets
	}
	return nil
}
