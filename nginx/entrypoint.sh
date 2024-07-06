#!/bin/sh

envsubst '${DOMAIN} ${BACKEND_PORT}' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf

nginx -g "daemon off;"
