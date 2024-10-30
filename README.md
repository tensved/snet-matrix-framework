<div align="center">

<h3 align="center">snet-matrix-framework</h3>

  <p align="center">
    ai-bots based on the Matrix protocol and snet ecosystem
    <br />
    <a href="https://snet-matrix-framework.gitbook.io/snet-matrix-framework"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/tensved/snet-matrix-framework">View Demo</a>
    ·
    <a href="https://github.com/tensved/snet-matrix-framework/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    ·
    <a href="https://github.com/tensved/snet-matrix-framework/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>
</div>

## About The Project

The snet-matrix-framework is designed to create bots that will connect the messenger user on the Matrix protocol and AI services in the snet ecosystem

## Getting Started

### Installation
1. Install [Docker Engine](https://docs.docker.com/engine/install/) and [Compose plugin](https://docs.docker.com/compose/install/linux/)
2. Clone the repo
   ```sh
   git clone https://github.com/tensved/snet-matrix-framework.git
   cd snet-matrix-framework
   git submodule update --init --recursive
   ```
3. Create an `.env` file based on `example.env` and add data to it
4. Create an `.env.local` file in `frontend` directory based on `frontend/example.env.local` and add data to it
5. Build and run docker containers
   ```sh
   docker compose up -d --build
   ```
6. Run certbot to issue a certificate
   ```sh
   docker compose run --rm  certbot certonly --webroot --webroot-path /var/www/certbot/ -d yourdomain.com
   ```
7. Uncomment all lines in the file `nginx/templates/default.conf.template`
8. Restart the `nginx` container so that Nginx starts using the new certificate
   ```sh
   docker compose restart nginx
   ```
   
### Explanation of variables in .env

#### Database
- DB_USER – username for connecting to the database
- DB_PASSWORD – password for the database user
- DB_HOST – hostname or IP address of the database server
- DB_PORT – port number on which the database server is listening
- DB_NAME – name of the database to connect to

#### App
- APP_PORT – port number on which the application will run
- PRODUCTION – flag indicating whether the application is running in production mode
- DOMAIN – your domain for the client application that provides the payment gateway
- PAYMENT_TIMEOUT – time to pay for the service call

#### Matrix
- MATRIX_HOMESERVER_URL – URL of the Matrix homeserver
- MATRIX_BOT_USERNAME – username for the Matrix bot that will provide access to snet services
- MATRIX_BOT_PASSWORD – password for the Matrix bot that will provide access to snet services
- MATRIX_SERVERNAME – server name for your Matrix

#### Ethereum
- IPFS_PROVIDER_URL – URL of the IPFS provider
- ETH_PROVIDER_URL – HTTP URL for the Ethereum provider (e.g., Infura)
- ETH_PROVIDER_WS_URL – WebSocket URL for the Ethereum provider (e.g., Infura)
- CHAIN_ID – chain ID of the Ethereum network

#### Admin
- ADMIN_PUBLIC_ADDRESS – public address of the admin on the blockchain
- ADMIN_PRIVATE_KEY – admin private key

### Notes

- Make sure your domain has the correct A records configured
- Use `docker compose run --rm certbot renew` to renew certs
- The minimal example of service is located at the path `cmd/main.go`
