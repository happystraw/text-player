// Copyright © 2019 Yutao Fang <fangyutao1993@hotmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"

	Cfg "github.com/happystraw/text-player/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var version = "0.3.0"

var cfg *Cfg.Config

type options struct {
	// Global Persistent flags
	configFile string
	noCache    bool
	appId      string
	apiKey     string
	apiSecret  string
	// Local flags
	reset bool
	list  bool
}

var opts = options{}

var Cmd = &cobra.Command{
	Version: version,
	Use:     "config",
	Short:   "Save global opts configuration to file",
	RunE:    run,
}

func init() {
	Cmd.Flags().BoolVarP(&opts.reset, "reset", "", false, "reset all to default settings")
	Cmd.Flags().BoolVarP(&opts.list, "list", "", false, "list all settings")
}

func run(cmd *cobra.Command, _ []string) error {
	if opts.list && opts.reset {
		return fmt.Errorf("--list and --reset cannot be set at the same time")
	}

	if opts.reset {
		if err := reset(cmd.Version); err != nil {
			return fmt.Errorf("reset config error: %s", err)
		}
		fmt.Println("Config reset to default")
	} else if opts.list {
		fmt.Printf("Current config: %s\n", opts.configFile)
		fmt.Println("---")
		printConfigToYaml()
	} else {
		if opts.appId == "" || opts.apiKey == "" || opts.apiSecret == "" {
			return fmt.Errorf("--appid, --apikey and --apisecret must be set")
		}
		if err := save(cmd.Version, cfg); err != nil {
			return fmt.Errorf("save config error: %s", err)
		}
		fmt.Println("Config saved to", viper.ConfigFileUsed())
	}

	return nil
}

func reset(v string) error {
	return save(v, Cfg.GetDefaultConfig())
}

func save(v string, cfg *Cfg.Config) error {
	viper.Reset()
	viper.SetConfigFile(opts.configFile)
	viper.Set("version", v)
	viper.Set("config", cfg)
	return viper.WriteConfig()
}

func printConfigToYaml() {
	content, _ := yaml.Marshal(cfg)
	fmt.Println(string(content))
}
