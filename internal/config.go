package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"addr" env:"HTTP_ADDR" env-default:":8080"`
}

type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true" env-default:"development"`
	StoragePath string     `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	// If not found, use flag
	if configPath == "" {
		var flagPath string
		flag.StringVar(&flagPath, "config", "", "Path to the configuration file")
		flag.Parse()

		if flagPath != "" {
			configPath = flagPath
		} else {
			// Fallback config path
			configPath = "./config/config.yaml"
		}
	}

	// Check file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", configPath)
	}

	// Read the config
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	return &cfg
}
