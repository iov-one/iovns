iovnscli tx domain transfer-domain --domain "${DOMAIN}" --from "${WALLET1}" --new-owner $(iovnscli keys show -a "${WALLET2}") --broadcast-mode block --gas-prices 10.0"${DENOM}" -y --transfer-flag 0
