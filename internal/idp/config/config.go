package config

import "github.com/spf13/viper"

type Config struct {
	BaseURL string `mapstructure:"BASE_URL"`
	Key     string `mapstructure:"KEY"`
	Cert    string `mapstructure:"CERT"`
	PSQLurl string `mapstructure:"PSQL_URL"`
}

func NewConfig() (config *Config, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
