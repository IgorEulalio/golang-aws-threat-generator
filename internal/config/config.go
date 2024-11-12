package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port      int
	AwsConfig AwsConfig
}

type AwsConfig struct {
	Region     string
	AccessKey  string
	SecretKey  string
	RoleArn    string
	ExternalID string
}

func LoadConfig() *Config {
	port := 8080 // Default port

	if val, exists := os.LookupEnv("PORT"); exists {
		if p, err := strconv.Atoi(val); err == nil {
			port = p
		}
	}

	return &Config{
		Port: port,
		AwsConfig: AwsConfig{
			Region: "sa-east-1",
		},
	}
}
