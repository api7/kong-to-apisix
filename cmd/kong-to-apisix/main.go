/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"os"

	"github.com/api7/kongtoapisix/pkg/apisix"
	"github.com/api7/kongtoapisix/pkg/kong"
	"github.com/api7/kongtoapisix/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newMigrateCommand() *cobra.Command {
	var apisixYamlPath string
	var apisixConfigPath string
	var kongYamlPath string
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "migrate from kong.yaml to apisix.yaml",
		RunE: func(cmd *cobra.Command, _ []string) error {
			kongConfig, err := kong.ReadYaml(kongYamlPath)
			if err != nil {
				return errors.Wrap(err, "read kong yaml failed")
			}
			apisixDecl, apisixConfig, err := kong.Migrate(kongConfig)
			if err != nil {
				return errors.Wrap(err, "migrate failed")
			}

			apisixDeclYaml, err := apisix.MarshalYaml(apisixDecl)
			if err != nil {
				return err
			}

			if apisixYamlPath == "" {
				fmt.Printf("apisix.yaml:\n%s", string(apisixDeclYaml))

				apisixConfigYaml, err := utils.ShowConfigYaml(apisixConfig)
				if err != nil {
					return err
				}
				fmt.Printf("\n\nconfig.yaml:\n%s", string(apisixConfigYaml))
			} else {
				if err := apisix.WriteToFile(apisixYamlPath, apisixDeclYaml); err != nil {
					return err
				}
				if apisixConfigPath != "" {
					if err := apisix.EnableAPISIXStandalone(apisixConfig); err != nil {
						return err
					}
					if err := utils.AppendToConfigYaml(apisixConfig, apisixConfigPath); err != nil {
						return err
					}
				}
				fmt.Println("migrate succeed")
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&kongYamlPath, "input", "i", "kong.yaml", "path of input kong declarative configuration file")
	cmd.PersistentFlags().StringVarP(&apisixYamlPath, "output", "o", "", "path of output apisix declarative configuration file")
	cmd.PersistentFlags().StringVarP(&apisixConfigPath, "config", "c", "", "path of output apisix config file")
	return cmd
}

func newDumpCommand() *cobra.Command {
	var kongYamlPath string
	var kongAddr string
	cmd := &cobra.Command{
		Use:   "dump [-k kong-yaml-path] [-a kong-address]",
		Short: "dump kong configuration with bare kong cluster",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := kong.DumpKong(kongAddr, kongYamlPath); err != nil {
				return errors.Wrap(err, "generate kong yaml failed")
			}
			fmt.Println("generated kong configuration file at", kongYamlPath)
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&kongYamlPath, "output", "o", "kong.yaml", "path of kong declarative configuration file")
	cmd.PersistentFlags().StringVarP(&kongAddr, "kong-admin-api", "a", "http://localhost:8001", "address of kong admin API")
	return cmd
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kong-to-apisix [command]",
		Short: "Tool to help you migrate from kong to APISIX",
	}

	cmd.AddCommand(newMigrateCommand())
	cmd.AddCommand(newDumpCommand())
	return cmd
}

func main() {
	root := NewCommand()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
