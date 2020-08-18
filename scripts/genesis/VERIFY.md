# Verify The iov-mainnet-2 Genesis File

This document assumes that you followed the procedure [here](https://github.com/iov-one/docs/blob/master/docs/iov-name-service/validator/01-mainnet.md) when you joined **iov-mainnet** as a validator.  It will deliniate the procedure necessary to verify the genesis file for **iov-mainnet-2**.  It requires that `diff`, `git`, `go v1.14.2+`, `jq`, `make`, `node`, `sed`, and `yarn` are installed on your system, and user `$USER_IOV` exists.

The verfication process consists of three parts:
1. Pulling the state and iov-mainnet-2 genesis file from IOV's repo.
1. Dumping the state from your node and comparing it to IOV's state.
1. Generating the iov-mainnet-2 genesis file from the dumped state and comparing it to IOV's genesis file.

The following procedure can be run on your validator or sentry nodes.  You can even do a practice run before the official decommissioning of the **iov-mainnet** chain takes place.  Note, however, in that case, you might see differences in the final `git diff` since the genesis file in IOV's repo is only an infrequent snapshot of the current state until the legacy chain is halted.

That said, let's go!

```bash
# stop iovns.service
sudo systemctl stop iovns.service

# start the genesis file verification process
su - ${USER_IOV}
set -o allexport ; source /etc/systemd/system/iovns.env ; set +o allexport # pick-up env vars

# pull the exported state from IOV and build iovnsd
cd ~ \
&& git clone --recursive https://github.com/iov-one/iovns.git \
&& cd iovns \
&& git submodule foreach git checkout master \
&& export HEIGHT=$(jq -r .height ./scripts/genesis/data/dump/dump.json) \
&& make build \
&& export PATH=~/iovns/build:$PATH

# dump local state and compare it with IOV's state
cd ~ \
&& git clone https://github.com/iov-one/weave.git \
&& cd weave/cmd/dumpstate \
&& make test \
&& make \
&& ./dumpstate -db ${DIR_WORK} -height ${HEIGHT} -out ${HEIGHT}.json \
&& sed --in-place "s/\"escrow\"/\"height\":${HEIGHT},\"escrow\"/" ${HEIGHT}.json \
&& jq --sort-keys . ${HEIGHT}.json > dump.json \
&& diff dump.json ~/iovns/scripts/genesis/data/dump/dump.json \
|| echo 'BAD state!'

# generate the iov-mainnet-2 genesis file and compare it with IOV's genesis file
cd ~/iovns/scripts/genesis \
&& yarn \
&& yarn test \
&& node -r esm genesis.js iov-mainnet-2 \
&& cd data/iov-mainnet-2/config \
&& git diff \
|| echo 'BAD genesis!'

exit # ${USER_IOV}
```

If the above is executed successfully then the genesis file that you generated in `.../iovns/scripts/genesis/data/iov-mainnet-2/config/genesis.json` is identical to this [gist](https://gist.githubusercontent.com/davepuchyr/4fe7e002061c537ddb116fee7a2f8e47/raw/genesis.json).  You're ready to follow the procedure [here](https://docs.iov.one/for-validators/mainnet) and wait for `genesis_time`.

`node -r esm genesis.js iov-mainnet-2` is where the conversion from **iov-mainnet**'s state to **iov-mainnet-2**'s genesis file takes place.  It simply does the transformations to fork from weave to cosmos-sdk.  [Check it out](genesis.js).
