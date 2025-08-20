package config

import (
	"regexp"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Global configuration variables.
var (
	Postgres   PostgresConfig   // Configuration for PostgreSQL.
	App        AppConfig        // Configuration for the application.
	Matrix     MatrixConfig     // Configuration for Matrix (chat protocol).
	Blockchain BlockchainConfig // Configuration for Blockchain.
	IPFS       IPFSConfig       // Configuration for IPFS (InterPlanetary File System).
)

// PostgresConfig holds the configuration values for connecting to a PostgreSQL database.
type PostgresConfig struct {
	User     string `env:"DB_USER"`     // The username for the PostgreSQL database.
	Password string `env:"DB_PASSWORD"` // The password for the PostgreSQL database.
	Host     string `env:"DB_HOST"`     // The host address for the PostgreSQL database.
	Port     string `env:"DB_PORT"`     // The port number for the PostgreSQL database.
	Name     string `env:"DB_NAME"`     // The name of the PostgreSQL database.
}

// AppConfig holds the configuration values for the application.
type AppConfig struct {
	Port           string `env:"APP_PORT"`        // The port on which the application runs.
	Domain         string `env:"DOMAIN"`          // The domain name of the application.
	IsProduction   bool   `env:"PRODUCTION"`      // Boolean flag indicating if the app is in production mode.
	PaymentTimeout int    `env:"PAYMENT_TIMEOUT"` // Number of minutes during which user can pay for a service call
}

// IPFSConfig holds the configuration values for connecting to an IPFS provider.
type IPFSConfig struct {
	IPFSProviderURL  string         `env:"IPFS_PROVIDER_URL"` // The URL of the IPFS provider.
	Timeout          string         `env:"IPFS_TIMEOUT"`      // The timeout value for IPFS operations.
	HashCutterRegexp *regexp.Regexp // Regexp for remove special character from ipfs hash
}

// BlockchainConfig holds the configuration values for connecting to a blockchain network.
type BlockchainConfig struct {
	AdminPrivateKey    string `env:"ADMIN_PRIVATE_KEY"`    // The private key of the admin account.
	AdminPublicAddress string `env:"ADMIN_PUBLIC_ADDRESS"` // The public address of the admin account.
	EthProviderURL     string `env:"ETH_PROVIDER_URL"`     // The URL of the Ethereum provider.
	EthProviderWSURL   string `env:"ETH_PROVIDER_WS_URL"`  // The WebSocket URL of the Ethereum provider.
	ChainID            string `env:"CHAIN_ID"`             // The chain ID of the blockchain network.
}

// MatrixConfig holds the configuration values for connecting to a Matrix homeserver.
type MatrixConfig struct {
	HomeserverURL string `env:"MATRIX_HOMESERVER_URL"` // The URL of the Matrix homeserver.
	Servername    string `env:"MATRIX_SERVERNAME"`     // The server name of the Matrix homeserver.
	Username      string `env:"MATRIX_BOT_USERNAME"`   // The username of the Matrix bot.
	Password      string `env:"MATRIX_BOT_PASSWORD"`   // The password of the Matrix bot.
	PickleKey     string `env:"MATRIX_PICKLE_KEY"`     // The pickle key for crypto operations.
}

// Init loads environment variables and parses them into the respective configuration structs.
// It first attempts to load environment variables from a .env file, and then parses the loaded variables.
func Init() {
	if err := godotenv.Load(); err != nil {
		log.Error().Err(err)
	} else {
		log.Debug().Msg(".env file loaded successfully") // Log a debug message if .env file is loaded successfully.
	}

	if err := env.Parse(&App); err != nil {
		log.Error().Err(err)
	}

	if err := env.Parse(&Postgres); err != nil {
		log.Error().Err(err)
	}

	if err := env.Parse(&Matrix); err != nil {
		log.Error().Err(err)
	}

	if err := env.Parse(&Blockchain); err != nil {
		log.Error().Err(err)
	}

	if err := env.Parse(&IPFS); err != nil {
		log.Error().Err(err)
	}

	log.Debug().Msg("configuration loading completed")
}
