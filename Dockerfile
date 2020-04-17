FROM golang:1.14.1-alpine3.11
ENV MONIKER=idk
ENV HOME=/iovnsd
WORKDIR /iovnsd
# create build dir
RUN mkdir /source
# copy build files to build dir
COPY ./ /source
# cd to source
WORKDIR /source
# install modules
RUN go mod download
# build all
RUN go build ./cmd/iovnsd
RUN go build ./cmd/iovnscli
# move binaries to iovnsd
RUN mv iovnsd /iovnsd/iovnsd && mv iovnscli /iovnsd/iovnscli
# change to working dir
WORKDIR /iovnsd
# delete build dir
RUN rm -rf /source
# copy utility scripts
COPY ./scripts .