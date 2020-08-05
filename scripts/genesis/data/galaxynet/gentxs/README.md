# Add your validator to the galaxynet genesis file #

Follow the procedure at https://docs.iov.one/for-validators/testnet.  Then,

```sh
# pick-up env vars
su - ${USER_IOV}
set -o allexport ; source /etc/systemd/system/starname.env ; set +o allexport

# clone iovns
cd ~ && git clone https://github.com/iov-one/iovns.git && cd iovns

# create your genesis transaction (do `iovnsd gentx --help` for the list of available flags)
${DIR_IOVNS}/iovnsd gentx \
  --amount 1000000000uvoi \
  --pubkey $(${DIR_WORK}/iovnsd tendermint show-validator) \
  --home ${DIR_WORK} \
  --name ${SIGNER} \
  --output-document ./scripts/genesis/data/galaxynet/gentxs/${MONIKER}.json

# create a PR
git checkout -b ${MONIKER}
git add .
git commit -m "Add ${MONIKER} gentx"
git push origin ${MONIKER}
```

Have your validator online before the waiting-for-cosmos-sdk-v0.39.1 release `genesis_time`. :)
