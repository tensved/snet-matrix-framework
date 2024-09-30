# Installation

## Installation

First of all, download and install [Docker Engine](https://docs.docker.com/engine/install/) and [Compose plugin](https://docs.docker.com/compose/install/linux/)

1.  Clone the repo and pull submodules\


    ```bash
    git clone https://github.com/tensved/snet-matrix-framework.git
    cd snet-matrix-framework
    git submodule update --init --recursive
    ```
2. Create an `.env` file based on `example.env` and add data to it
3. Create an `.env.local` file in `frontend` directory based on `frontend/example.env` and add data to it
4.  Build and run docker containers\


    ```bash
    docker compose up -d --build
    ```
5.  Run certbot to issue a certificate\


    ```bash
    docker compose run --rm  certbot certonly --webroot --webroot-path /var/www/certbot/ -d yourdomain.com
    ```
6. Uncomment all lines in the file `nginx/templates/default.conf.template`
7.  Restart the `nginx` container so that Nginx starts using the new certificate and config\


    ```bash
    docker compose restart nginx
    ```

## Notes

* Make sure your domain has the correct A records configured
* Use `docker compose run --rm certbot renew` to renew certs
* The minimal example of service is located at the path `pkg/lib/examples/snet/main.go`
