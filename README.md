# IOVNS
This repo contains bnsd built on top of the cosmos-sdk

Everything is still very very very experimental, and I expect the design to dramatically change over the course of the coming days.

## Running docker image
```shell script
# build docker script
docker build -t iovns .
# initialize chain
bash scripts/init.sh
# run docker iovnsd
docker run -it -p 127.0.0.1:26657:26657 -p 127.0.0.1:26656:26656 \
  -v ~/.iovnsd:/app/.iovnsd -v ~/.iovnscli:/app/.iovnscli \ 
  iovns iovnsd start
```