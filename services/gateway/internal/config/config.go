package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port           int    `env:"PORT" env-default:"8081"`
	AuthAddress    string `env:"AUTH_ADDRESS" env-required:"true"`
	PaymentAddress string `env:"PAYMENT_ADDRESS" env-required:"true"`
	BillingAddress string `env:"BILLING_ADDRESS" env-required:"true"`
	JWTSecret      []byte `env:"JWT_SECRET" env-default:"my_super_secret_key"`
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
