package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
)

var (
	Postgres   PostgresConfig
	App        AppConfig
	Matrix     MatrixConfig
	Blockchain BlockchainConfig
	IPFS       IPFSConfig
)

type PostgresConfig struct {
	URL string `env:"DB_URL" envDefault:"postgresql://postgres:postgres@localhost:5432/postgres"`
}

type AppConfig struct {
	Port         string `env:"PORT" envDefault:"3000"`
	IsProduction bool   `env:"PRODUCTION"`
}

type IPFSConfig struct {
	IPFSProviderURL string `env:"IPFS_PROVIDER_URL"`
	Timeout         string `env:"IPFS_TIMEOUT"`
}

type BlockchainConfig struct {
	PrivateKey     string `env:"PRIVATE_KEY"`
	EthProviderURL string `env:"ETH_PROVIDER_URL"`
	ChainID        string `env:"CHAIN_ID"`
}

type MatrixConfig struct {
	HomeserverURL string `env:"HOMESERVER_URL"`
	Servername    string `env:"SERVERNAME"`
	Username      string `env:"BOT_USERNAME"`
	Password      string `env:"BOT_PASSWORD"`
}

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	if err := env.Parse(&App); err != nil {
		log.Printf("%+v\n", err)
	}
	log.Printf("%+v\n", App)

	if err := env.Parse(&Postgres); err != nil {
		log.Printf("%+v\n", err)
	}

	if err := env.Parse(&Matrix); err != nil {
		log.Printf("%+v\n", err)
	}

	if err := env.Parse(&Blockchain); err != nil {
		log.Printf("%+v\n", err)
	}

	log.Printf("%+v\n", Blockchain)

	if err := env.Parse(&IPFS); err != nil {
		log.Printf("%+v\n", err)
	}
}
