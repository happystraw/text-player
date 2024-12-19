package config

type Tts struct {
	Host      string         `yaml:"host" mapstructure:"host"`
	AppId     string         `yaml:"app-id" mapstructure:"app-id"`
	ApiKey    string         `yaml:"api-key" mapstructure:"api-key"`
	ApiSecret string         `yaml:"api-secret" mapstructure:"api-secret"`
	Params    map[string]any `yaml:"params" mapstructure:"params"`
}
