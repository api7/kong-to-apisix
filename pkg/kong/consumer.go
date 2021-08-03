package kong

import (
	"fmt"

	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
	"github.com/api7/kong-to-apisix/pkg/utils"
)

func MigrateConsumer(kongConfig *KongConfig, configYamlAll *[]utils.YamlItem) (*[]v1.Consumer, error) {
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
		if c.KeyAuthCredentials != nil {
			pluginConfig := make(map[string]interface{})
			pluginConfig["key"] = c.KeyAuthCredentials[0].Key

			apisixConsumer.Plugins = v1.Plugins{"key-auth": pluginConfig}
		if c.JWTCredentials != nil {
			cred := c.JWTCredentials[0]
			pluginConfig := make(map[string]interface{})
			pluginConfig["key"] = cred.Key
			pluginConfig["secret"] = cred.Secret
			algorithm := cred.Algorithm
			if algorithm == "HS384" || algorithm == "ES256" {
				fmt.Printf("APISIX JWT have not support %s yet\n", algorithm)
			} else {
				pluginConfig["algorithm"] = algorithm
				if cred.Secret != "" {
					pluginConfig["secret"] = cred.Secret
				}
				if cred.RSA_Public_Key != "" {
					pluginConfig["rsa_public_key"] = cred.RSA_Public_Key
				}
			}
		}
		apisixConsumers = append(apisixConsumers, *apisixConsumer)
	}

	return &apisixConsumers, nil
}
