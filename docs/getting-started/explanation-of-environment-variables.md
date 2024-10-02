# Explanation of environment variables

Explanation of  environment variables required for installation.

## .env variables

### Database

* DB\_USER – username for connecting to the database
* DB\_PASSWORD – password for the database user
* DB\_HOST – hostname or IP address of the database server
* DB\_PORT – port number on which the database server is listening
* DB\_NAME – name of the database to connect to

### App

* APP\_PORT – port number on which the application will run
* PRODUCTION – flag indicating whether the application is running in production mode
* DOMAIN – your domain for the client application that provides the payment gateway

### Matrix

* MATRIX\_HOMESERVER\_URL – URL of the Matrix homeserver
* MATRIX\_BOT\_USERNAME – username for the Matrix bot that will provide access to snet services
* MATRIX\_BOT\_PASSWORD – password for the Matrix bot that will provide access to snet services
* MATRIX\_SERVERNAME – server name for your Matrix

### Ethereum

* IPFS\_PROVIDER\_URL – URL of the IPFS provider
* ETH\_PROVIDER\_URL – HTTP URL for the Ethereum provider (e.g., Infura)
* ETH\_PROVIDER\_WS\_URL – WebSocket URL for the Ethereum provider (e.g., Infura)
* CHAIN\_ID – chain ID of the Ethereum network

### Admin

* ADMIN\_PUBLIC\_ADDRESS – public address of the admin on the blockchain
* ADMIN\_PRIVATE\_KEY – admin private key

## .env.local variables

### App

* VITE\_URL\_PATH\_PREFIX — specific path to host, e.g. for https://your.domain.com/path/to/gateway VITE\_URL\_PATH\_PREFIX will be /path/to/gateway. If not specified, / will be used as the default value

### Backend

* VITE\_BACKEND\_URL — URL of the framework. If not specified, Mock Service Worker will be used as the mocked backend (for tests only)

### Wallet Connect

* VITE\_PUBLIC\_WALLETCONNECT\_PROJECT\_ID — a Project ID from walletconnect.com
* VITE\_MAINNET\_RPC\_URL — URL of a Mainnet RPC
* VITE\_SEPOLIA\_RPC\_URL — URL of a Sepolia RPC. It's a testnet, so it's used for testing only. If not specified, Sepolia will not be supported
