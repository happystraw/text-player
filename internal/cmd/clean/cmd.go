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

package clean

import (
	"fmt"
	"github.com/happystraw/text-player/internal/cmd/config"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var Cmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete generated files",
	RunE:  run,
}

func run(*cobra.Command, []string) error {
	if err := clean(); err != nil {
		return fmt.Errorf("clean error: %s", err)
	}
	fmt.Println("clean done")
	return nil
}

func clean() error {
	cfg := config.GetConfig()
	path := cfg.GetCachePath()
	dir, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !dir.IsDir() {
		return fmt.Errorf("error: cache file [%s] must be a directory, not a file", path)
	}

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return os.Remove(path)
	})
}
