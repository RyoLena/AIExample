package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Port               string `mapstructure:"PORT"`
	PythonAIServiceURL string `mapstructure:"PYTHON_AI_SERVER"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("读取失败", err)
	}
	err = viper.Unmarshal(&config)
	return config, err
}
