package kong

import (
	"github.com/api7/kong-to-apisix/pkg/apisix"
	uuid "github.com/satori/go.uuid"
)

func MigrateConsumer(kongConfig *Config, apisixConfig *apisix.Config) error {
	kongConsumers := kongConfig.Consumers

	for index, consumer := range kongConsumers {
		consumerId := consumer.ID
		if len(consumerId) <= 0 {
			consumerId = uuid.NewV4().String()
			kongConfig.Consumers[index].ID = consumerId
		}
		username := consumer.Username
		if len(username) <= 0 {
			username = consumer.CustomID
		}

		var apisixConsumer apisix.Consumer
		apisixConsumer.ID = consumerId
		apisixConsumer.Username = username

		if len(consumer.KeyAuthCredentials) > 0 {
			KTAUpdateApisixConsumerPlugin(&apisixConsumer, &consumer.KeyAuthCredentials[0])
		}

		if len(consumer.BasicAuthCredentials) > 0 {
			KTAUpdateApisixConsumerPlugin(&apisixConsumer, &consumer.BasicAuthCredentials[0])
		}

		if len(consumer.HmacAuthCredentials) > 0 {
			KTAUpdateApisixConsumerPlugin(&apisixConsumer, &consumer.HmacAuthCredentials[0])
		}

		if len(consumer.JwtSecrets) > 0 {
			KTAUpdateApisixConsumerPlugin(&apisixConsumer, &consumer.JwtSecrets[0])
		}

		apisixConfig.Consumers = append(apisixConfig.Consumers, apisixConsumer)
	}

	return nil
}
