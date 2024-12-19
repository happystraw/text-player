package config

import (
	"github.com/ebitengine/oto/v3"
)

type Engine struct {
	SampleRate   int        `yaml:"sample-rate" mapstructure:"sample-rate"`
	ChannelCount int        `yaml:"channel-count" mapstructure:"channel-count"`
	Format       oto.Format `yaml:"format" mapstructure:"format"`
	BufferSize   int64      `yaml:"buffer-size" mapstructure:"buffer-size"`
}
