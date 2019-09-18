#!/bin/sh
set -e
# Script to install fn cli inside the containers for testing
# and registering the certificates
echo "Running entrypoint"

apk update && apk add curl
apk add --no-cache libc6-compat
curl -LSs https://raw.githubusercontent.com/fnproject/cli/master/install | sh

cat /app/.certificates/*.crt >> /etc/ssl/certs/ca-certificates.crt
FN_LOG_LEVEL=debug ./fnserver
