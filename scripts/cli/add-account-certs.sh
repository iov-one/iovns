#!/bin/sh

set -e

iovnscli tx domain add-certs --domain ${DOMAIN} --cert-file ./cert.json --from ${WALLET1} --broadcast-mode block --gas-prices 10.0uvoi -y
