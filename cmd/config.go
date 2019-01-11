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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var isReset bool

var configOption = struct {
	CachePath    string
	DisableCache bool
	AppID        string
	APIKey       string
}{
	CachePath:    "",
	DisableCache: false,
	AppID:        "",
	APIKey:       "",
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Save global flags configuration to file",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if isReset {
			err = resetDefault()
		} else {
			err = viper.WriteConfig()
		}
		if err != nil {
			record(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	rootCmd.PersistentFlags().StringVarP(&configOption.CachePath, "cache-path", "o", "", "path for cache files(default is $HOME/.text-player)")
	rootCmd.PersistentFlags().BoolVarP(&configOption.DisableCache, "disable-cache", "n", false, "disable cache generated speech files")
	rootCmd.PersistentFlags().StringVarP(&configOption.AppID, "appid", "", "", "xunfei tts api auth appid")
	rootCmd.PersistentFlags().StringVarP(&configOption.APIKey, "apikey", "", "", "xunfei tts api auth apikey")

	configCmd.Flags().BoolVarP(&isReset, "reset", "", false, "reset all to default settings")
}

func initDefaultConfig() {
	viper.SetDefault("version", "0.0.1")
	viper.SetDefault("cache-path", "")
	viper.SetDefault("disable-cache", false)
	viper.SetDefault("xunfei.host", "http://api.xfyun.cn/v1/service/v1/tts")
	viper.SetDefault("xunfei.appid", "")
	viper.SetDefault("xunfei.apikey", "")
	viper.SetDefault("xunfei.params", map[string]string{
		"voice_name": "xiaoyan",
	})

}

func createDefaultConfigIfNotExists(filename string) error {
	if _, err := os.Stat(filename); err == nil || !os.IsNotExist(err) {
		return nil
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	return viper.WriteConfig()
}

func createDefaultCachePathIfNotExists(path string) error {
	if info, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0766)
	} else if !info.IsDir() {
		return fmt.Errorf("cache path [%s] is not a directory : ", path)
	}

	return nil
}

func getCachePath() (string, error) {
	path := configOption.CachePath
	if path == "" {
		path = viper.GetString("cache-path")
	}

	if path == "" {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		path = home + "/.text-player"
	}

	if info, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("cache path [%s] not exists", path)
	} else if !info.IsDir() {
		return "", fmt.Errorf("cache path [%s] is not a directory : ", path)
	}

	return path, nil
}

func disableCache() bool {
	return viper.GetBool("disable-cache")
}

func bindFlagsValueToConfig() {
	if configOption.CachePath != "" {
		viper.Set("cache-path", configOption.CachePath)
	}
	if configOption.DisableCache {
		viper.Set("disable-cache", true)
	}
	if configOption.AppID != "" {
		viper.Set("xunfei.appid", configOption.AppID)
	}
	if configOption.APIKey != "" {
		viper.Set("xunfei.apikey", configOption.APIKey)
	}
}

func resetDefault() error {
	filename := viper.ConfigFileUsed()
	viper.Reset()
	viper.SetConfigFile(filename)
	initDefaultConfig()
	return viper.WriteConfig()
}
