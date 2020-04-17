set -o nounset

rm -rf "$HOME/.appcli"
rm -rf "$HOME/.appd"
# init config files
./appd init "$MONIKER" --chain-id iovns
# configure cli
./appcli config chain-id iovns
./appcli config output json
./appcli config trust-node true

# use keyring backend TODO learn what is a keyring backend
./appcli config keyring-backend test

# create accounts
./appcli keys add fd
./appcli keys add dp

# give the accounts some money
./appd add-genesis-account $(./appcli keys show fd -a) 1000iov,10000000000stake
./appd add-genesis-account $(./appcli keys show dp -a) 1000iov,10000000000stake

# save configs for the daemon
./appd gentx --name fd --keyring-backend test

# input genTx to the genesis file
./appd collect-gentxs
# verify genesis file is fine
./appd validate-genesis
