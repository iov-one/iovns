"use strict";


export const multisigs = {
   iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n: {
      "//name": "reward fund",
      address: "cond:gov/rule/0000000000000002",
      star1: "star1scfumxscrm53s4dd3rl93py5ja2ypxmxlhs938",
   },
   iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq: {
      "//name": "IOV SAS",
      address: "cond:multisig/usage/0000000000000001",
      star1: "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78",
   },
   iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0: {
      "//name": "IOV SAS employee bonus pool/colloboration appropriation pool",
      address: "cond:multisig/usage/0000000000000002",
      star1: "star16tm7scg0c2e04s0exk5rgpmws2wk4xkd84p5md",
   },
   iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz: {
      "//name": "IOV SAS pending deals pocket; close deal or burn",
      address: "cond:multisig/usage/0000000000000003",
      star1: "star1uyny88het6zaha4pmkwrkdyj9gnqkdfe4uqrwq",
   },
   iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu: {
      "//name": "IOV SAS bounty fund",
      address: "cond:multisig/usage/0000000000000004",
      star1: "star1m7jkafh4gmds8r0w79y2wu2kvayqvrwt7cy7rf",
   },
   iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw: {
      "//name": "Unconfirmed contributors/co-founders",
      address: "cond:multisig/usage/0000000000000005",
      star1: "star1p0d75y4vpftsx9z35s93eppkky7kdh220vrk8n",
   },
   iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc: {
      "//name": "Custodian of missing star1 accounts",
      address: "cond:multisig/usage/0000000000000006",
      star1: "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
   },
};

export const conds = Object.keys( multisigs ).reduce( ( accumulator, key ) => {
   const multisig = multisigs[key];

   accumulator[multisig.address] = {
      "//name": multisig["//name"],
      iov1: key,
      star1: multisig.star1,
   }

   return accumulator;
}, {} );

export const names = Object.keys( multisigs ).reduce( ( accumulator, key ) => {
   const multisig = multisigs[key];

   accumulator[multisig["//name"]] = {
      cond: multisig.address,
      iov1: key,
      star1: multisig.star1,
   }

   return accumulator;
}, {} );

export const chainIds = {
   "bip122-tmp-bcash":            "asset:bch",
   "bip122-tmp-bitcoin":          "asset:btc",
   "bip122-tmp-litecoin":         "asset:ltc",
   "cosmos-binance-chain-tigris": "asset:bnb",
   "cosmos-columbus-3":           "asset:luna", // terra
   "cosmos-cosmoshub-3":          "asset:atom",
   "cosmos-emoney-1":             "asset:ngm",
   "cosmos-irishub":              "asset:iris",
   "cosmos-kava-2":               "asset:kava",
   "ethereum-eip155-1":           "asset:eth",
   "iov-mainnet":                 "asset:iov",
   "lisk-ed14889723":             "asset:lsk",
   "starname-migration":          "asset:iov",
   "tezos-tmp-mainnet":           "asset:xtz",
};

export const source2multisig = {
   iov1w2suyhrfcrv5h4wmq3rk3v4x95cxtu0a03gy6x: {
      "//id": "escrow isabella*iov",
      star1: "star1elad203jykd8la6wgfnvk43rzajyqpk0wsme9g",
   },
   iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph: {
      "//id": "escrow kadima*iov",
      star1: "star1hjf04872s9rlcdg2wqwvapwttvt3p4gjpp0xmc",
   },
   iov149cn0rauw2773lfdp34njyejg3cfz2d56c0m5t: {
      "//id": "escrow joghurt*iov",
      star1: "star15u4kl3lalt8pm2g4m23erlqhylz76rfh50cuv8",
   },
   iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u: {
      "//id": "vaildator guaranteed reward fund",
      star1: "star17w7fjdkr9laphtyj4wxa32rf0evu94xgywxgl4",
   },
};
