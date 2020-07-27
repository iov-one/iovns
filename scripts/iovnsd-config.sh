#!/bin/bash
set -e

GENESIS=https://gist.githubusercontent.com/davepuchyr/e1482e63cb81443cde1616f353c4779f/raw/f6f69eb18e81b03c1303552491d61c11165941f9/genesis.json
PATH=~/src/iov/iovns/build:$PATH

which iovnsd
iovnsd version --long
rm -f ~/.iovnsd/config/genesis.json ~/.iovnsd/config/gentx/*
iovnsd unsafe-reset-all

curl ${GENESIS} > genesis.json
CHAIN_ID=$(jq -r .chain_id genesis.json)
iovnsd init ${CHAIN_ID} --chain-id ${CHAIN_ID}
sed --in-place 's/skip_timeout_commit = false/skip_timeout_commit = true/' ~/.iovnsd/config/config.toml
mv genesis.json ~/.iovnsd/config
iovnsd add-genesis-account $(iovnscli keys show ${CHAIN_ID} -a) 1112111000000uvoi
iovnsd gentx --name ${CHAIN_ID} --keyring-backend test --amount 1111111000000uvoi
iovnsd collect-gentxs > /dev/null 2> /dev/null
iovnsd validate-genesis
echo -e "Do \e[32miovnsd start --minimum-gas-prices 10.0uvoi &\e[0m or start iovnsd in your IDE"
echo -e "Bonus points: \e[32miovnscli rest-server --chain-id ${CHAIN_ID} --node http://localhost:26657 --trust-node true &\e[0m"
