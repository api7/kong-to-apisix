package main

import (
	"fmt"
	"os"

	"github.com/api7/kongtoapisix/pkg/kong"
)

func main() {
	filePath := "kong.yaml"
	if os.Getenv("KONG_YAML_DUMP_PATH") != "" {
		filePath = os.Getenv("KONG_YAML_DUMP_PATH")
	}
	if err := kong.DumpKong(filePath); err != nil {
		panic(err)
	}
	fmt.Println("generate kong configuration file at", filePath)
}