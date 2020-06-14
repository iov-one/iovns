#!/bin/bash

set -o nounset

echo -n "This script will remove all existing iovns configurations, wallets etc. Do you want to proceed? [Y/n] "
read YN
if ! echo $YN | grep -v -q -i "n"; then exit; fi

rm -rf "$HOME/.iovnscli"
rm -rf "$HOME/.iovnsd"
# init config files
iovnsd init "$(hostname)" --chain-id iovns
# configure cli
iovnscli config chain-id iovns
iovnscli config output json
iovnscli config trust-node true

# use keyring backend
iovnscli config keyring-backend test

# create accounts
iovnscli keys add fd
iovnscli keys add dp
iovnscli keys add ok

# give the accounts some money
iovnsd add-genesis-account $(iovnscli keys show fd -a) 1000iov,10000000000stake
iovnsd add-genesis-account $(iovnscli keys show dp -a) 1000iov,10000000000stake
iovnsd add-genesis-account $(iovnscli keys show ok -a) 1000iov,10000000000stake

# save configs for the daemon
iovnsd gentx --name fd --keyring-backend test

# input genTx to the genesis file
iovnsd collect-gentxs
# verify genesis file is fine
iovnsd validate-genesis
