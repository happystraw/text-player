package config

import (
	"fmt"
	C "github.com/happystraw/text-player/internal/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

func init() {
	cfg = C.GetDefaultConfig()
}

func applyOptionsValueToGlobalConfig() {
	if opts.noCache {
		cfg.Cache.Disable = true
	}
	if opts.appId != "" {
		cfg.Tts.AppId = opts.appId
	}
	if opts.apiKey != "" {
		cfg.Tts.ApiKey = opts.apiKey
	}
	if opts.apiSecret != "" {
		cfg.Tts.ApiSecret = opts.apiSecret
	}
}

func InitConfig() {
	// Fall back to default config file if not specified.
	fallback := false
	if opts.configFile == "" {
		fallback = true
		opts.configFile = getDefaultConfigFile()
	}

	if opts.reset {
		return
	}

	if file, err := os.Stat(opts.configFile); err == nil && file.IsDir() {
		// Config file must be a file, not a directory.
		fmt.Printf("Config file [%s] must be a file, not a directory\n", opts.configFile)
		os.Exit(1)
	} else if err != nil && os.IsNotExist(err) && fallback {
		// Create default config file if not exists.
		viper.SetConfigFile(opts.configFile)
		if err := reset(Cmd.Version); err != nil {
			fmt.Printf("Failed to create default config file: %s\n", err)
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Printf("Failed to get config file info: %s\n", err)
		os.Exit(1)
	} else {
		// If a config file is found, read it in.
		viper.SetConfigFile(opts.configFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Failed to read config file: %s\n", err)
			os.Exit(1)
		}
	}

	// Unmarshal configuration.
	if err := viper.UnmarshalKey("config", cfg); err != nil {
		fmt.Println("Try use text-player config --reset to reset the config file")
		fmt.Printf("Config file \"%s\" is invalid\n", opts.configFile)
		os.Exit(1)
	}

	// Apply opts value to global configuration.
	applyOptionsValueToGlobalConfig()
}

func getDefaultConfigFile() string {
	home, _ := homedir.Dir()
	return path.Join(home, ".text-player.yaml")
}

func GetConfig() *C.Config {
	return cfg
}

func ApplyPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&opts.configFile, "config", "c", "", "config file")
	cmd.PersistentFlags().BoolVarP(&opts.noCache, "no-cache", "n", false, "disable cache")
	cmd.PersistentFlags().StringVarP(&opts.appId, "appid", "", "", "xunfei tts api auth appid")
	cmd.PersistentFlags().StringVarP(&opts.apiKey, "apikey", "", "", "xunfei tts api auth apikey")
	cmd.PersistentFlags().StringVarP(&opts.apiSecret, "apisecret", "", "", "xunfei tts api auth apisecret")
}
