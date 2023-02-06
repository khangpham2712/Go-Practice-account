package config

import "github.com/spf13/viper"

type Config struct {
	Port       string `mapstructure:"PORT"`
	DBDriver   string `mapstructure:"DRIVER"`
	DBUsername string `mapstructure:"USERNAME"`
	DBPassword string `mapstructure:"PASSWORD"`
	Source     string `mapstructure:"SOURCE"`
	DBPort     string `mapstructure:"MYSQL_PORT"`
	DBName     string `mapstructure:"DATABASE_NAME"`
}

func ReadFromConfigFile(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&config); err != nil {
		return
	}
	return
}
