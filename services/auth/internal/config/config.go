package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port        int    `env:"PORT" env-default:"8081"`
	DBConn      string `env:"DB_CONN" env-required:"true"`
	KafkaBroker string `env:"KAFKA_BROKER" env-default:"localhost:9092"`
	JWTSecret   []byte `env:"JWT_SECRET" env-default:"my_super_secret_key"`
}

func MustLoad() *Config {
	var cfg Config

	if _, err := os.Stat(".env"); err == nil {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			panic("cannot read config")
		}
	} else {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			panic("cannot read env")
		}
	}

	return &cfg
}
