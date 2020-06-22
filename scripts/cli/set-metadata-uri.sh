#!/bin/sh

set -e

iovnscli tx domain set-account-metadata --domain ${DOMAIN} --name ${ACCOUNT} --metadata https://iov.one \
    --from ${WALLET1} --gas-prices 10.0uvoi --broadcast-mode block -y

