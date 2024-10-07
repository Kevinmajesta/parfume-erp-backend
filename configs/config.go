package configs

import (
	"errors"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string         `env:"ENV" envDefault:"dev"`
	Port     string         `env:"PORT" envDefault:"8080"`
	Postgres PostgresConfig `envPrefix:"POSTGRES_"`
	Redis    RedisConfig    `envPrefix:"REDIS_"`
	JWT      JwtConfig      `envPrefix:"JWT_"`
	Encrypt  EncryptConfig  `envPrefix:"ENCRYPT_"`
	SMTP     SMTPConfig     `envPrefix:"SMTP_"`
}

type SMTPConfig struct {
	Host     string `env:"HOST" envDefault:"smtp.gmail.com"`
	Port     string `env:"PORT" envDefault:"587"`
	Password string `env:"PASSWORD" envDefault:"qloscitkfltxwjkb"`
}

type PostgresConfig struct {
	Host     string `env:"HOST" `
	Port     string `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
}

type JwtConfig struct {
	SecretKey string `env:"SECRET_KEY"`
}

type RedisConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"6379"`
	Password string `env:"PASSWORD" envDefault:""`
}

type EncryptConfig struct {
	SecretKey string `env:"SECRET_KEY"`
	IV        string `env:"IV"`
}

func NewConfig(envPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		return nil, errors.New("failed to load .env file")
	}

	cfg := new(Config)

	err = env.Parse(cfg)
	if err != nil {
		return nil, errors.New("failed to parse config file")
	}

	return cfg, nil
}
