# IOV Name Service (IOVNS)
This repo contains the IOV Name Service daemon (iovnsd) and command-line interface (iovnscli) built on top of the cosmos-sdk.

## Running via docker image
```shell script
# make the apps
make build
# build docker script
docker build -t iovns .
# initialize chain
bash scripts/init.sh
# run docker iovnsd
docker run -it -p 127.0.0.1:46657:26657 -p 127.0.0.1:46656:26656 \
  -v ~/.iovnsd:/app/.iovnsd -v ~/.iovnscli:/app/.iovnscli \
  iovns iovnsd start
```
