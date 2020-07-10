#!/bin/sh

set -e

# if iovnsd does not exist then init the new network
if [ ! -d ".iovnsd" ]; then
  echo ".iovnsd not found... initting a new chain"
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
  iovnsd add-genesis-account $(iovnscli keys show fd -a) 1000000000tiov
  iovnsd add-genesis-account $(iovnscli keys show dp -a) 1000000000tiov
  iovnsd add-genesis-account $(iovnscli keys show ok -a) 1000000000tiov

  # save configs for the daemon
  iovnsd gentx --name fd --keyring-backend test --amount 10000000tiov

  # input genTx to the genesis file
  iovnsd collect-gentxs
  # verify genesis file is fine
  iovnsd validate-genesis

  sed -i 's/stake/tiov/g' ~/.iovnsd/config/genesis.json
fi

iovnsd start