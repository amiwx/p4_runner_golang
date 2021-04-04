package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB *DBConfig
	P4 *P4Config
}

type DBConfig struct {
	Dialect  string `mapstructure:"driver"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Charset  string `mapstructure:"charset"`
}

type P4Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	ParseURL string `mapstructure:"parseUrl"`
	ServeURL string `mapstructure:"serveUrl"`
}

func GetConfig() (config *Config, err error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	var dbConfig DBConfig
	viper.UnmarshalKey("db_config", &dbConfig)

	var p4Config P4Config
	viper.UnmarshalKey("p4_config", &p4Config)

	config = &Config{DB: &dbConfig, P4: &p4Config}

	return
}
