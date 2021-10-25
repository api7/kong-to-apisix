package kong

import (
	"testing"

	"github.com/api7/kong-to-apisix/pkg/apisix"

	"github.com/stretchr/testify/assert"
)

type TestConsumerExpect struct {
	ID        string
	Username  string
	KeyAuth   *apisix.KeyAuthCredential
	BasicAuth *apisix.BasicAuthCredential
	HmacAuth  *apisix.HmacAuthCredential
	JwtAuth   *apisix.JwtSecrets
}

type TestConsumer struct {
	Consumer Consumer
	Expect   TestConsumerExpect
}

func TestMigrateConsumer(t *testing.T) {
	var kongConfig Config
	var apisixConfig apisix.Config
	testConsumers := []TestConsumer{
		{
			Consumer: Consumer{
				Username: "only-consumer",
			},
			Expect: TestConsumerExpect{
				Username: "only-consumer",
			},
		},
		{
			Consumer: Consumer{
				Username: "key-auth-consumer",
				KeyAuthCredentials: KeyAuthCredentials{
					{
						Key: "test",
					},
				},
			},
			Expect: TestConsumerExpect{
				Username: "key-auth-consumer",
				KeyAuth: &apisix.KeyAuthCredential{
					Key: "test",
				},
			},
		},
		{
			Consumer: Consumer{
				Username: "basic-auth-consumer",
				BasicAuthCredentials: BasicAuthCredentials{
					{
						Username: "test",
						Password: "123456",
					},
				},
			},
			Expect: TestConsumerExpect{
				Username: "basic-auth-consumer",
				BasicAuth: &apisix.BasicAuthCredential{
					Username: "test",
					Password: "123456",
				},
			},
		},
		{
			Consumer: Consumer{
				Username: "hmac-auth-consumer",
				HmacAuthCredentials: HmacAuthCredentials{
					{
						Username: "test",
						Secret:   "123456",
					},
				},
			},
			Expect: TestConsumerExpect{
				Username: "hmac-auth-consumer",
				HmacAuth: &apisix.HmacAuthCredential{
					AccessKey: "test",
					SecretKey: "123456",
				},
			},
		},
		{
			Consumer: Consumer{
				Username: "jwt-auth-consumer",
				JwtSecrets: JwtSecrets{
					{
						Key:       "test",
						Secret:    "123456",
						Algorithm: "HS256",
					},
				},
			},
			Expect: TestConsumerExpect{
				Username: "jwt-auth-consumer",
				JwtAuth: &apisix.JwtSecrets{
					Key:       "test",
					Secret:    "123456",
					Algorithm: "HS256",
				},
			},
		},
	}
	for _, testConsumer := range testConsumers {
		kongConfig = Config{}
		apisixConfig = apisix.Config{}
		kongConfig.Consumers = append(kongConfig.Consumers, testConsumer.Consumer)
		err := MigrateConsumer(&kongConfig, &apisixConfig)
		assert.NoError(t, err)
		assert.Equal(t, apisixConfig.Consumers[0].Username, testConsumer.Expect.Username)
		if testConsumer.Expect.KeyAuth != nil {
			assert.Equal(t, apisixConfig.Consumers[0].Plugins.KeyAuth.Key, testConsumer.Expect.KeyAuth.Key)
		}
		if testConsumer.Expect.BasicAuth != nil {
			assert.Equal(t, apisixConfig.Consumers[0].Plugins.BasicAuth.Username, testConsumer.Expect.BasicAuth.Username)
			assert.Equal(t, apisixConfig.Consumers[0].Plugins.BasicAuth.Password, testConsumer.Expect.BasicAuth.Password)
		}
		if testConsumer.Expect.HmacAuth != nil {
			assert.Equal(t, apisixConfig.Consumers[0].Plugins.HmacAuth.AccessKey, testConsumer.Expect.HmacAuth.AccessKey)
			assert.Equal(t, apisixConfig.Consumers[0].Plugins.HmacAuth.SecretKey, testConsumer.Expect.HmacAuth.SecretKey)
		}
		if testConsumer.Expect.JwtAuth != nil {
			assert.Equal(t, apisixConfig.Consumers[0].Plugins.JwtAuth.Key, testConsumer.Expect.JwtAuth.Key)
			assert.Equal(t, apisixConfig.Consumers[0].Plugins.JwtAuth.Secret, testConsumer.Expect.JwtAuth.Secret)
			assert.Equal(t, apisixConfig.Consumers[0].Plugins.JwtAuth.Algorithm, testConsumer.Expect.JwtAuth.Algorithm)
		}
	}
}
