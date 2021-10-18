package kong

import (
	"github.com/api7/kong-to-apisix/pkg/apisix"
	"github.com/api7/kong-to-apisix/pkg/utils"
)

func MigrateConsumer(kongConfig *Config, configYamlAll *[]utils.YamlItem) (apisix.Consumers, error) {
	kongConsumers := kongConfig.Consumers

	var apisixConsumers apisix.Consumers
	for _, c := range kongConsumers {
		username := c.Username
		if username == "" {
			username = c.CustomID
		}
		apisixConsumer := apisix.Consumer{
			Username: username,
		}

		// TODO: need to test then it got multiple key
		if len(c.KeyAuthCredentials) > 0 && len(c.KeyAuthCredentials[0].Key) > 0 {
			apisixConsumer.Plugins.KeyAuth.Key = c.KeyAuthCredentials[0].Key
		}
		apisixConsumers = append(apisixConsumers, apisixConsumer)
	}

	return apisixConsumers, nil
}
