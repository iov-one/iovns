#!/bin/bash

src=iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un

amount=$1
dst=$2
tm=$3
shift 3
memo=$@

bnscli send-tokens -amount "$amount IOV" -src "bech32:$src" -dst "bech32:$dst" -memo "$memo" \
   | bnscli with-fee -amount "0.5 IOV" \
   | bnscli sign -tm $tm \
   | bnscli submit -tm $tm
