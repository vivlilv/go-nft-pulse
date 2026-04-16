package events_opensea

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	OPENSEA_API_KEY string `envconfig:"OPENSEA_API_KEY" required:"true"`
}

func NewConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found: %v", err)
		// Continue, as env vars might be set in the OS
	}
	var config Config

	if err := envconfig.Process("", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get events config: %w", err)
		panic(err)
	}

	return config
}
