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
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Erase generated files",
	Run: func(cmd *cobra.Command, args []string) {
		if err := clean(); err != nil {
			record(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func clean() error {
	path, err := getCachePath()
	if err != nil {
		return err
	}
	dir, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !dir.IsDir() {
		return fmt.Errorf("Error: cache file [%s] must be a directory, not a file", path)
	}

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		name := info.Name()
		dotIndex := strings.LastIndex(name, ".")
		if dotIndex < 2 {
			return nil
		}
		if ext := name[dotIndex+1:]; ext != extension {
			return nil
		}
		return os.Remove(path)
	})
}
