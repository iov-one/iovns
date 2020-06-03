FROM alpine:3.11

ENV NSDHOME /app
ENV NSCLIHOME /app

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq && \
    addgroup iovnsuser && \
    adduser -S -G iovnsuser iovnsuser -h "$NSDHOME" -h "$NSCLIHOME"

# Run the container with iovnsuser by default. (UID=100, GID=1000)
USER iovnsuser

# p2p, rpc and prometheus port
EXPOSE 46656 46657 46660

ARG NSDBINARY=build/iovnsd
ARG NSCLIBINARY=build/iovnscli

COPY $NSDBINARY /usr/bin/iovnsd
COPY $NSCLIBINARY /usr/bin/iovnscli

WORKDIR /app

# Run iovnsd by default, omit entrypoint to ease using container with iovnscli
CMD ["iovnsd"]
STOPSIGNAL SIGTERM

VOLUME $NSDHOME $NSCLIHOME
