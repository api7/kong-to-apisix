package kong

import (
	"testing"

	"github.com/api7/kong-to-apisix/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestMigrateConsumer(t *testing.T) {
	configYaml, kongConsumers := &[]utils.YamlItem{}, &Config{
		Consumers: Consumers{
			{
				CustomID: "test-id-01",
				Username: "test-user-01",
				KeyAuthCredentials: []struct {
					Key string `yaml:"key"`
				}{
					{
						Key: "test-key-01",
					},
				},
			},
			{
				CustomID:           "test-id-02",
				Username:           "test-user-02",
				KeyAuthCredentials: nil,
			},
		},
	}

	apisixConsumers, err := MigrateConsumer(kongConsumers, configYaml)
	assert.NoError(t, err)
	assert.Equal(t, len(kongConsumers.Consumers), len(apisixConsumers))
	assert.Equal(t, (apisixConsumers)[0].Username, "test-user-01")
	assert.Equal(t, (apisixConsumers)[1].Username, "test-user-02")
}
