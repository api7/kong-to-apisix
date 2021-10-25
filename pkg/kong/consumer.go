package kong

import (
	"github.com/api7/kong-to-apisix/pkg/apisix"
	uuid "github.com/satori/go.uuid"
)

func MigrateConsumer(kongConfig *Config, apisixConfig *apisix.Config) error {
	kongConsumers := kongConfig.Consumers

	for _, consumer := range kongConsumers {
		username := consumer.Username
		if len(username) <= 0 {
			username = consumer.CustomID
		}
		var apisixConsumer apisix.Consumer
		if len(consumer.ID) > 0 {
			apisixConsumer.ID = consumer.ID
		} else {
			apisixConsumer.ID = uuid.NewV4().String()
		}
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
