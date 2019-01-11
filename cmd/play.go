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
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/hajimehoshi/oto"
	"github.com/happystraw/text-player/tts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const extension = "wav"

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play <message>",
	Short: "Text to speech and play it",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 || len(args[0]) < 1 {
			record("Error: message is required")
			cmd.Usage()
			os.Exit(1)
		}

		if err := Play(args[0]); err != nil {
			record(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
}

// Play audio generated by text
func Play(msg string) error {
	filename, _ := useCacheFile(msg)

	if filename != "" {
		return playLocal(filename)
	}

	return playRemote(msg)
}

func useCacheFile(msg string) (string, error) {
	if disableCache() {
		return "", nil
	}

	filename, err := getCacheFileName(msg)
	if err != nil {
		return "", err
	}

	if file, err := os.Stat(filename); os.IsNotExist(err) {
		return "", nil
	} else if file.IsDir() {
		return "", fmt.Errorf("Error: cache file [%s] must be a directory, not a file", filename)
	}

	return filename, nil
}

func playLocal(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	return play(file)
}

func playRemote(msg string) error {
	xunfei := tts.New(
		viper.GetString("xunfei.host"),
		viper.GetString("xunfei.appid"),
		viper.GetString("xunfei.apikey"),
	)
	params := viper.GetStringMapString("xunfei.params")
	// got wav / pcm audio
	params["aue"] = "raw"
	params["auf"] = "audio/L16;rate=16000"

	raw, err := xunfei.Create(msg, params)

	if err != nil {
		return err
	}

	if err := cache(msg, raw); err != nil {
		recordf("Warning: cache file fail: %s", err)
	}

	return play(ioutil.NopCloser(bytes.NewReader(raw)))
}

func play(r io.ReadCloser) error {
	player, err := oto.NewPlayer(16000, 1, 2, 1024)
	if err != nil {
		return err
	}
	defer player.Close()

	if _, err := io.Copy(player, r); err != nil {
		return err
	}

	return nil
}

func cache(msg string, raw []byte) error {
	if disableCache() {
		return nil
	}
	filename, err := getCacheFileName(msg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, raw, 0666)
}

func getCacheFileName(msg string) (string, error) {
	path, err := getCachePath()
	if err != nil {
		return "", err
	}
	w := md5.New()
	w.Write([]byte(msg))
	return fmt.Sprintf(
		"%s%s%s.%s",
		path,
		string(os.PathSeparator),
		hex.EncodeToString(w.Sum(nil)),
		extension,
	), nil
}
