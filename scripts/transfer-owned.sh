#!/bin/bash
set -e

PATH=~/src/iov/iovns/build:$PATH

NOW=$(date +"%s")
DOMAIN="domain${NOW}"
FROM=star1z6rhjmdh2e9s6lvfzfwrh8a3kjuuy58y74l29t # bojack /run/user/500/keybase/kbfs/team/iov_one.blockchain/credentials/test-wallets/bojack.mne.txt
NEW_OWNER=star19jj4wc3lxd54hkzl42m7ze73rzy3dd3wry2f3q # w1 https://github.com/iov-one/iovns/blob/master/docs/cli/MULTISIG.md#create-keys-if-you-dont-have-them-ready
OTHER=star1l4mvu36chkj9lczjhy9anshptdfm497fune6la # w2 https://github.com/iov-one/iovns/blob/master/docs/cli/MULTISIG.md#create-keys-if-you-dont-have-them-ready
#NODE=https://rpc.cluster-galaxynet.iov.one:443
NODE=http://localhost:26657
CHAIN=iovns-galaxynet
FLAGS_PULL="--chain-id ${CHAIN} --node ${NODE} --output json"
FLAGS_PUSH="$FLAGS_PULL --gas-prices 10.0uvoi --keyring-backend test --broadcast-mode block"

iovnscli tx starname register-domain --yes --domain ${DOMAIN} --from ${FROM} ${FLAGS_PUSH} | jq

iovnscli tx starname register-account --yes --domain ${DOMAIN} --name self                   --from ${FROM} ${FLAGS_PUSH} | jq
iovnscli tx starname register-account --yes --domain ${DOMAIN} --name other --owner ${OTHER} --from ${FROM} ${FLAGS_PUSH} | jq

iovnscli query starname resolve --starname      *${DOMAIN} ${FLAGS_PULL} | jq
iovnscli query starname resolve --starname  self*${DOMAIN} ${FLAGS_PULL} | jq
iovnscli query starname resolve --starname other*${DOMAIN} ${FLAGS_PULL} | jq

iovnscli tx starname transfer-domain --yes --domain ${DOMAIN} --new-owner ${NEW_OWNER} --transfer-flag 1 --from ${FROM} ${FLAGS_PUSH} | jq

iovnscli query starname domain-info --domain ${DOMAIN} ${FLAGS_PULL} | jq
iovnscli query starname resolve --starname      *${DOMAIN} ${FLAGS_PULL} | jq
iovnscli query starname resolve --starname  self*${DOMAIN} ${FLAGS_PULL} | jq
iovnscli query starname resolve --starname other*${DOMAIN} ${FLAGS_PULL} | jq
