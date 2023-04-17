package config

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
)

const fileConfig = ".env"

type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	AppStage    string `mapstructure:"APP_STAGE"`
	AppDev      bool   `mapstructure:"APP_DEV"`
	HTTPPort    string `mapstructure:"HTTP_PORT"`
	MetricHost  string `mapstructure:"OTEL_EXPORTER_METRIC_ENDPOINT"`
	TraceHost   string `mapstructure:"OTEL_EXPORTER_TRACE_ENDPOINT"`
}

// Load the config from file or env to the Config struct
func Load(_ context.Context) (*Config, error) {
	viper.SetConfigFile(fileConfig)
	viper.AutomaticEnv()

	viper.SetDefault("HTTP_PORT", "8080")
	viper.SetDefault("APP_STAGE", "DEV")

	var cfg Config

	if err := loadMappedEnvVariables(&cfg); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Println("unable find or read configuration file")
	}

	if err := viper.UnmarshalExact(&cfg); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &cfg, nil
}

func loadMappedEnvVariables(cfg *Config) error {
	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(cfg, &envKeysMap); err != nil {
		return fmt.Errorf("%v", err)
	}

	for k := range *envKeysMap {
		if err := viper.BindEnv(k); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	return nil
}
