package config

import (
	"currency-checker/internal/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Listen struct {
		BindIP string `yaml:"bind_ip" env-default:"0.0.0.0"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`

	EmailServer struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		Credentials struct {
			Email    string `yaml:"email"`
			Password string `yaml:"password"`
		} `yaml:"credentials"`
	} `yaml:"email_server"`

	StoragePath string `yaml:"storage_path" env-default:"emails.txt"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log := logger.New()
		log.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Info(help)
			log.Fatal(err)
		}
	})
	return instance
}
