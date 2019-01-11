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

package cmd

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const runningInConsole = 1

const runningInServer = 2

const logo = ` ______        __  ___  __
/_  __/____ __/ /_/ _ \/ /__ ___ _____ ____
 / / / -_) \ / __/ ___/ / _ ` + "`" + `/ // / -_) __/
/_/  \__/_\_\\__/_/  /_/\_,_/\_, /\__/_/
                            /___/
`

var runningEnv = runningInConsole

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.0.1",
	Use:     "text-player",
	Short:   "Text Player is a text-to-speech and play-to-speech gadget.",
	Long:    logo + "Text Player is a text-to-speech and play-to-speech gadget.",

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		record(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SetVersionTemplate(
		`Text Player {{printf "version %s" .Version}}
`,
	)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.text-player.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			record(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".text-player" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".text-player")

		// Init default config
		initDefaultConfig()

		// Create config if file not exists.
		createDefaultConfigIfNotExists(home + "/.text-player.yaml")

		// Create cache directory if path not exists.
		createDefaultCachePathIfNotExists(home + "/.text-player")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && cfgFile != "" {
		record("Using config file:", viper.ConfigFileUsed())
	}

	// Bind flags value to configuration.
	bindFlagsValueToConfig()
}

func isRunningInConsole() bool {
	return runningEnv == runningInConsole
}

func isRunningInServer() bool {
	return runningEnv == runningInServer
}

func record(v ...interface{}) {
	if isRunningInServer() {
		log.Println(v...)
	} else {
		fmt.Println(v...)
	}
}

func recordf(format string, a ...interface{}) {
	if isRunningInServer() {
		log.Printf(format, a...)
	} else {
		fmt.Printf(format, a...)
	}
}
