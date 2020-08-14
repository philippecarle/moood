package config

import "github.com/kelseyhightower/envconfig"

// Config is holding all the configuration variables of this app
type Config struct {
	Mongo struct {
		User     string `envconfig:"MONGO_USERNAME"`
		Password string `envconfig:"MONGO_PASSWORD"`
	}
	Mercure struct {
		URL  string `envconfig:"MERCURE_HUB_URL"`
		Port int    `envconfig:"MERCURE_HUB_PORT"`
	}
	RabbitMQ struct {
		User     string `envconfig:"RABBITMQ_USER"`
		Password string `envconfig:"RABBITMQ_PASSWORD"`
		URL      string `envconfig:"RABBITMQ_URL"`
		Port     int    `envconfig:"RABBITMQ_PORT"`
	}
	JWTPrivateKey string `envconfig:"JWT_PRIVATE_KEY"`
}

// GetConfig returns a Config struct after it processed the env variables
func GetConfig() Config {
	var conf Config

	envconfig.Process("", &conf)

	return conf
}
