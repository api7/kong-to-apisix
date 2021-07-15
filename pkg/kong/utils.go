package kong

import (
	"context"
	"fmt"
	"os"

	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	"github.com/kong/deck/state"
	"github.com/kong/deck/utils"
)

func DumpKong(fileName string) error {
	kongAddr := "http://localhost:8001"
	if os.Getenv("KONG_ADMIN_ADDR") != "" {
		kongAddr = os.Getenv("KONG_ADMIN_ADDR")
	}
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
