package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var (
	Postgres   PostgresConfig
	App        AppConfig
	Matrix     MatrixConfig
	Blockchain BlockchainConfig
	IPFS       IPFSConfig
)

type PostgresConfig struct {
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Name     string `env:"DB_NAME"`
}

type AppConfig struct {
	Port         string `env:"APP_PORT"`
	Domain       string `env:"DOMAIN"`
	IsProduction bool   `env:"PRODUCTION"`
}

type IPFSConfig struct {
	IPFSProviderURL string `env:"IPFS_PROVIDER_URL"`
	Timeout         string `env:"IPFS_TIMEOUT"`
}

type BlockchainConfig struct {
	AdminPrivateKey    string `env:"ADMIN_PRIVATE_KEY"`
	AdminPublicAddress string `env:"ADMIN_PUBLIC_ADDRESS"`
	EthProviderURL     string `env:"ETH_PROVIDER_URL"`
	EthProviderWSURL   string `env:"ETH_PROVIDER_WS_URL"`
	ChainID            string `env:"CHAIN_ID"`
	TokenAddress       string `env:"TOKEN_ADDRESS"`
}

type MatrixConfig struct {
	HomeserverURL string `env:"MATRIX_HOMESERVER_URL"`
	Servername    string `env:"MATRIX_SERVERNAME"`
	Username      string `env:"MATRIX_BOT_USERNAME"`
	Password      string `env:"MATRIX_BOT_PASSWORD"`
}

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Debug().Msgf("Error loading .env file: %v", err)
	} else {
		log.Debug().Msg(".env file loaded successfully")
	}
	if err := env.Parse(&App); err != nil {
		log.Debug().Msgf("%+v\n", err)
	}

	if err := env.Parse(&Postgres); err != nil {
		log.Debug().Msgf("%+v\n", err)
	}

	if err := env.Parse(&Matrix); err != nil {
		log.Debug().Msgf("%+v\n", err)
	}

	if err := env.Parse(&Blockchain); err != nil {
		log.Debug().Msgf("%+v\n", err)
	}

	if err := env.Parse(&IPFS); err != nil {
		log.Debug().Msgf("%+v\n", err)
	}
}
