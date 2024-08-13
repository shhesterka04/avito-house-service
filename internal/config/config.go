package config

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	prjDir     = "/root/"
	configFile = "config.env"
)

type Config struct {
	HostAddr   string `mapstructure:"HOST_ADDR"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = prjDir + configFile
	}

	filename := filepath.Base(path)
	path = filepath.Dir(path)

	viper.SetDefault("HOST_ADDR", "0.0.0.0:8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "postgres")

	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrapf(err, "failed to read config from %s", path)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	return cfg, nil
}
