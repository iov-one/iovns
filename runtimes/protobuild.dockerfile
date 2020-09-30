FROM alpine as protoc

ARG PROTOC_VERSION="3.13.0"
ARG PROTOC_ZIP=protoc-${PROTOC_VERSION}-linux-aarch_64.zip

RUN apk --no-cache add curl
WORKDIR /protoc
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}
RUN unzip -o $PROTOC_ZIP
RUN ls -al
RUN ls -al ./include

FROM golang:1.15-alpine
# add g++ required for ledger
RUN apk add g++
# add make
RUN apk add make
# install protobuf
RUN apk add protobuf
# copy protobuf includes
COPY --from=protoc /protoc/include /protobuf/include
# include protobuf in path
ENV PATH=$PATH:/protobuf:/go/bin
# INSTALL GIT
RUN apk --no-cache add git
# INSTALL gogoproto
RUN go get github.com/gogo/protobuf/protoc-gen-gofast
CMD ["sh"]