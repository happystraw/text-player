package config

type Cache struct {
	Path    string `yaml:"path" mapstructure:"path"`
	Disable bool   `yaml:"disable" mapstructure:"disable"`
}
