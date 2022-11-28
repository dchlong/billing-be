package config

import (
	"os"
	"path"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/dchlong/billing-be/pkg/infra"
)

type AppConfig struct {
	HTTPAddr                string                `json:"http_addr"`
	ReadHeaderTimeout       int                   `json:"read_header_timeout"`
	DatabaseConfig          *infra.DatabaseConfig `json:"database_config"`
	NumberOfSecondsInABlock int64                 `json:"number_of_seconds_in_a_block"`
}

func ProviderDatabaseConfig(appConfig *AppConfig) *infra.DatabaseConfig {
	return appConfig.DatabaseConfig
}

func ProviderAppConfig() (*AppConfig, error) {
	configPath := "configs"
	configFile := os.Getenv("CONFIG_FILE")
	if configFile != "" {
		configPath = path.Dir(configFile)
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.AddConfigPath(configPath)
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	appCfg := &AppConfig{}
	err = viper.Unmarshal(appCfg, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "json"
	})
	if err != nil {
		return nil, err
	}

	return appCfg, nil
}
