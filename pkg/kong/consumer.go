package kong

import (
	"github.com/api7/kong-to-apisix/pkg/apisix"
	"github.com/api7/kong-to-apisix/pkg/utils"
)

func MigrateConsumer(kongConfig *KongConfig, configYamlAll *[]utils.YamlItem) (*[]apisix.Consumer, error) {
	kongConsumers := kongConfig.Consumers

	var apisixConsumers []apisix.Consumer
	for _, c := range kongConsumers {
		username := c.Username
		if username == "" {
			username = c.CustomId
		}
		apisixConsumer := &apisix.Consumer{
			Username: username,
		}

		// TODO: need to test then it got multiple key
		if len(c.KeyAuthCredentials) > 0 && c.KeyAuthCredentials[0].Key != "" {
			pluginConfig := make(map[string]interface{})
			pluginConfig["key"] = c.KeyAuthCredentials[0].Key

			apisixConsumer.Plugins = apisix.Plugins{"key-auth": pluginConfig}
		}
		apisixConsumers = append(apisixConsumers, *apisixConsumer)
	}

	return &apisixConsumers, nil
}
