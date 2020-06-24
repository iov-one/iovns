#!/bin/sh

set -e

# shellcheck disable=SC2086
iovnscli tx domain set-account-metadata --domain "${DOMAIN}" --name ${ACCOUNT} --metadata https://iov.one \
    --from "${WALLET1}" --gas-prices 10.0"${DENOM}" --broadcast-mode block -y

