package kong

import (
	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
)

func MigrateConsumer(kongConfig *KongConfig) (*[]v1.Consumer, error) {
	kongConsumers := kongConfig.Consumers

	var apisixConsumers []v1.Consumer
	for _, c := range kongConsumers {
		username := c.Username
		if username == "" {
			username = c.CustomId
		}
		apisixConsumer := &v1.Consumer{
			Username: username,
		}

		// TODO: need to test then it got multiple key
		if c.KeyAuthCredentials[0].Key != "" {
			pluginConfig := make(map[string]interface{})
			pluginConfig["key"] = c.KeyAuthCredentials[0].Key

			apisixConsumer.Plugins = v1.Plugins{"key-auth": pluginConfig}
		}
		apisixConsumers = append(apisixConsumers, *apisixConsumer)
	}

	return &apisixConsumers, nil
}
