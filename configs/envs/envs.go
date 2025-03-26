package envs

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MongoURL              string `envconfig:"MONGO_URL" required:"true"`
	MongoDB               string `envconfig:"MONGO_DATABASE" required:"true"`
	RedisURL              string `envconfig:"REDIS_URL" required:"true"`
	InsuranceProviderURL  string `envconfig:"INSURANCE_PROVIDER_URL" required:"true"`
	InsuranceProvideToken string `envconfig:"INSURANCE_PROVIDER_TOKEN" required:"true"`
}

var AppConfig Config

func LoadEnvs() {
	err := envconfig.Process("", &AppConfig)
	if err != nil {
		log.Fatalf("Erro ao carregar as variáveis de ambiente: %v", err)
	}
}
