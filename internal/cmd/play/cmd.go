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

package play

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/happystraw/text-player/internal/cmd/config"
	"github.com/happystraw/text-player/internal/player"
	"github.com/happystraw/text-player/internal/tts"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path"
	"sync"
)

const extension = "pcm"

var mkCacheDirOnce = sync.Once{}

var Cmd = &cobra.Command{
	Use:   "play <message>",
	Short: "Text to speech and play it",
	RunE:  run,
}

func run(_ *cobra.Command, args []string) error {
	if len(args) < 1 || len(args[0]) < 1 {
		return fmt.Errorf("message is required")
	}

	if err := Play(args[0]); err != nil {
		return fmt.Errorf("play error: %s", err)
	}

	fmt.Println("play done")
	return nil
}

func Play(msg string) error {
	// Initialize player
	cfg := config.GetConfig()
	player.InitPlayer(cfg.Engine)

	// Check cache
	filename, _ := getCacheFilename(msg)
	if filename != "" {
		// Play from cache
		return playLocal(filename)
	}

	// Play from remote
	return playRemote(msg)
}

func getCacheFilename(msg string) (string, error) {
	cfg := config.GetConfig()
	if cfg.Cache.Disable {
		return "", nil
	}

	dir := cfg.GetCachePath()
	mkCacheDirOnce.Do(func() {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				panic(fmt.Sprintf("error: create cache path [%s] failed: %s", dir, err))
			}
		}
	})

	filename, err := getCacheFileName(msg)
	if err != nil {
		return "", err
	}

	if file, err := os.Stat(filename); os.IsNotExist(err) {
		return "", nil
	} else if file.IsDir() {
		return "", fmt.Errorf("error: cache file [%s] must be a directory, not a file", filename)
	} else if file.Size() == 0 {
		return "", fmt.Errorf("error: cache file [%s] is empty", filename)
	}

	return filename, nil
}

func playLocal(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	return player.GetPlayer().Play(file)
}

func playRemote(msg string) error {
	cfg := config.GetConfig()
	t := tts.New(cfg.Tts)
	raw, err := t.Create(msg)
	if err != nil {
		return err
	}
	if err := cache(msg, raw); err != nil {
		return err
	}
	return player.GetPlayer().Play(io.NopCloser(bytes.NewReader(raw)))
}

func cache(msg string, raw []byte) error {
	cfg := config.GetConfig()
	if cfg.Cache.Disable {
		return nil
	}

	dir := cfg.GetCachePath()
	mkCacheDirOnce.Do(func() {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				panic(fmt.Sprintf("error: create cache path [%s] failed: %s", dir, err))
			}
		}
	})

	filename, err := getCacheFileName(msg)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, raw, 0666)
}

func getCacheFileName(msg string) (string, error) {
	cfg := config.GetConfig()
	dir := cfg.GetCachePath()
	hash := md5.Sum([]byte(msg))
	filename := hex.EncodeToString(hash[:]) + "." + extension
	return path.Join(dir, string(os.PathSeparator), filename), nil
}
