#!/bin/bash

src=$1
amount=$2
dst=$3
tm=$4
shift 4
memo=$@

bnscli send-tokens -amount "$amount IOV" -src "bech32:$src" -dst "bech32:$dst" -memo "$memo" \
   | bnscli with-fee -amount "0.5 IOV" \
   | bnscli sign -tm $tm \
   | bnscli submit -tm $tm
