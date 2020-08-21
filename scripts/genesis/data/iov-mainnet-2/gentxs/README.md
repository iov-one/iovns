# Add your validator to the iov-mainnet-2 genesis file #

Follow the procedure [here](https://docs.iov.one/for-validators/mainnet).  Then, find your star1 account in the [iov-mainnet-2 genesis file](https://gist.github.com/davepuchyr/4fe7e002061c537ddb116fee7a2f8e47/raw/genesis.json) so that you know the number of `uiov` tokens that are available to you.  Finally,

```bash
# pick-up env vars
su ${USER_IOV}
set -o allexport ; source /etc/systemd/system/starname.env ; set +o allexport

# clone iovns
cd ~ && git clone https://github.com/iov-one/iovns.git && cd iovns

# set the amount of your delegation in uiov; 1 IOV equals 1 million uiov
export AMOUNT=1000000 # denominated in uiov

# create your genesis transaction (do `iovnsd gentx --help` for the list of available flags)
${DIR_IOVNS}/iovnsd gentx \
  --amount ${AMOUNT}uiov \
  --pubkey $(${DIR_IOVNS}/iovnsd tendermint show-validator --home ${DIR_WORK}) \
  --home ${DIR_WORK} \
  --name ${SIGNER} \
  --output-document "./scripts/genesis/data/iov-mainnet-2/gentxs/${MONIKER}.json"

# create a PR
git checkout -b "${MONIKER}"
git add .
git commit -m "Add '${MONIKER}' gentx"
git push origin "${MONIKER}"
```

Once your PR is verified then it will be merged and your gentx will be written into the genesis file.

Have your **iov-mainnet-2** validator online before the `genesis_time`. :)
