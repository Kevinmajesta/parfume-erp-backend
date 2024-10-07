package entity

type Config struct {
	Namespace string
	Redis     RedisConfig
	SMTP      SMTPConfig
}

type SMTPConfig struct {
	Host     string `env:"HOST" envDefault:"smtp.gmail.com"`
	Port     string `env:"PORT" envDefault:"587"`
	Password string `env:"PASSWORD" envDefault:"qloscitkfltxwjkb"`
}

type RedisConfig struct {
	Host string
	Port string
}
