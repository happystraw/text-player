package config

import (
	"github.com/ebitengine/oto/v3"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

type Config struct {
	Cache  *Cache  `yaml:"cache" mapstructure:"cache"`
	Engine *Engine `yaml:"engine" mapstructure:"engine"`
	Tts    *Tts    `yaml:"tts" mapstructure:"tts"`
}

var defaultConfig *Config

func init() {
	defaultConfig = &Config{
		Cache: &Cache{
			Disable: false,
			Path:    getDefaultCachePath(),
		},
		Engine: &Engine{
			SampleRate:   16000,
			ChannelCount: 1,
			Format:       oto.FormatSignedInt16LE,
			BufferSize:   0,
		},
		Tts: &Tts{
			Host: "wss://tts-api.xfyun.cn/v2/tts",
			Params: map[string]any{
				"aue": "raw",
				"vcn": "xiaoyan",
				"auf": "audio/L16;rate=16000",
				"tte": "UTF8",
			},
		},
	}
}

func (c *Config) GetCachePath() string {
	if c.Cache != nil && c.Cache.Path != "" {
		return c.Cache.Path
	}
	return getDefaultCachePath()
}

func GetDefaultConfig() *Config {
	return defaultConfig
}

func getDefaultCachePath() string {
	home, _ := homedir.Dir()
	return path.Join(home, string(os.PathSeparator), ".text-player")
}
