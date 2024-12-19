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
	"github.com/happystraw/text-player/internal/cmd/clean"
	"github.com/happystraw/text-player/internal/cmd/config"
	"github.com/happystraw/text-player/internal/cmd/play"
	"github.com/happystraw/text-player/internal/cmd/serve"
	"github.com/spf13/cobra"
	"os"
)

const logo = ` ______        __  ___  __
/_  __/____ __/ /_/ _ \/ /__ ___ _____ ____
 / / / -_) \ / __/ ___/ / _ ` + "`" + `/ // / -_) __/
/_/  \__/_\_\\__/_/  /_/\_,_/\_, /\__/_/
                            /___/
`

// cmd represents the base command when called without any subcommands
var cmd = &cobra.Command{
	Version: config.Cmd.Version,
	Use:     "text-player",
	Short:   "Text Player is a small tool for converting text to speech and playing speech.",
	Long:    logo + "Text Player is a small tool for converting text to speech and playing speech.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	cmd.SetVersionTemplate(
		`Text Player {{printf "version %s" .Version}}
`,
	)

	// Add commands
	cmd.AddCommand(play.Cmd)
	cmd.AddCommand(serve.Cmd)
	cmd.AddCommand(config.Cmd)
	cmd.AddCommand(clean.Cmd)

	// Persistent flags for all commands.
	config.ApplyPersistentFlags(cmd)

	// Init configuration.
	cobra.OnInitialize(config.InitConfig)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the cmd.
func Execute() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
