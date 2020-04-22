# Simple usage with a mounted data directory:
# > docker build -t iovns .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.iovnsd:/root/.iovnsd -v ~/.iovnscli:/root/.iovnscli iovns iovnsd init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.iovnsd:/root/.iovnsd -v ~/.iovnscli:/root/.iovnscli iovns iovnsd start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python

# Set working directory for the build
WORKDIR /go/src/github.com/iov-one/iovns

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/iovnsd /usr/bin/iovnsd
COPY --from=build-env /go/bin/iovnsd /usr/bin/iovnsd

# Run iovnsd by default, omit entrypoint to ease using container with iovnscli
CMD ["iovnsd"]
