package kong

import (
	"context"
	"fmt"

	"github.com/api7/kong-to-apisix/pkg/apisix"

	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	"github.com/kong/deck/state"
	"github.com/kong/deck/utils"
)

func DumpKong(kongAddr string, fileName string) error {
	rootConfig := utils.KongClientConfig{
		Address: kongAddr,
	}

	// if fileName is empty, print it to console
	if fileName == "" {
		fileName = "-"
	}
	wsClient, err := utils.GetKongClient(rootConfig)
	if err != nil {
		return err
	}
	dumpConfig := dump.Config{}

	rawState, err := dump.Get(context.Background(), wsClient, dumpConfig)
	if err != nil {
		return fmt.Errorf("reading configuration from Kong: %w", err)
	}
	ks, err := state.Get(rawState)
	if err != nil {
		return fmt.Errorf("building state: %w", err)
	}

	return file.KongStateToFile(ks, file.WriteConfig{
		SelectTags: dumpConfig.SelectorTags,
		Workspace:  "",
		Filename:   fileName,
		FileFormat: "YAML",
		WithID:     false,
	})
}

func AfterMigrate(apisixConfig *apisix.Config) error {
	// Processing consumer data
	if len(apisixConfig.Consumers) > 0 {
		for consumerIndex, consumer := range apisixConfig.Consumers {
			if len(consumer.ID) > 0 {
				// If there is a consumer ID, it will cause APISIX verification to fail
				apisixConfig.Consumers[consumerIndex].ID = ""
			}
		}
	}
	return nil
}
