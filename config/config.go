package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"sync"
)

type Config struct {
	HTTP struct {
		Host string `env:"HTTP_HOST" env-default:"localhost"`
		Port string `env:"HTTP_PORT" env-default:"8080"`
	}

	PSQL struct {
		Host     string `env:"DB_HOST" env-default:"localhost"`
		Port     string `env:"DB_PORT" env-default:"5432"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		DBName   string `env:"DB_NAME"`
		SSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
	}

	DEBUG bool `env:"DEBUG" env-default:"false"`
}

var (
	config    Config
	once      sync.Once
	configErr error
)

func GetConfig() (Config, error) {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			configErr = err
		}

		err = cleanenv.ReadEnv(&config)
		if err != nil {
			configErr = err
			return
		}

	})
	return config, configErr
}
