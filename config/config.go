package config

import (
    "github.com/spf13/viper"
    "log"
)

var Conf *viper.Viper

func LoadConfig() {
    Conf = viper.New()
    Conf.SetConfigFile(".env")
    Conf.AutomaticEnv()
    if err := Conf.ReadInConfig(); err != nil {
        log.Fatalf("Failed to read config: %v", err)
    }
}
