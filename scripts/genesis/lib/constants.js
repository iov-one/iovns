"use strict";


export const multisigs = {
   iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n: {
      "//name": "reward fund",
      address: "cond:gov/rule/0000000000000002",
      star1: "reward fund star1", // TODO
   },
   iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq: {
      "//name": "IOV SAS",
      address: "cond:multisig/usage/0000000000000001",
      star1: "IOV SAS star1", // TODO
   },
   iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0: {
      "//name": "IOV SAS employee bonus pool/colloboration appropriation pool",
      address: "cond:multisig/usage/0000000000000002",
      star1: "IOV SAS employee bonus pool/colloboration appropriation pool star1", // TODO
   },
   iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz: {
      "//name": "IOV SAS pending deals pocket; close deal or burn",
      address: "cond:multisig/usage/0000000000000003",
      star1: "IOV SAS pending deals pocket; close deal or burn star1", // TODO
   },
   iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu: {
      "//name": "IOV SAS bounty fund",
      address: "cond:multisig/usage/0000000000000004",
      star1: "IOV SAS bounty fund star1", // TODO
   },
   iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw: {
      "//name": "Unconfirmed contributors/co-founders",
      address: "cond:multisig/usage/0000000000000005",
      star1: "Unconfirmed contributors/co-founders star1", // TODO
   },
   iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc: {
      "//name": "Custodian of missing star1 accounts",
      address: "cond:multisig/usage/0000000000000006",
      star1: "custodial star1", // TODO
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
      star1: "IOV SAS multisig star1_TBD_isabella*iov", // TODO
   },
   iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph: {
      "//id": "escrow kadima*iov",
      star1: "IOV SAS multisig star1_TBD_kadima*iov", // TODO
   },
   iov149cn0rauw2773lfdp34njyejg3cfz2d56c0m5t: {
      "//id": "escrow joghurt*iov",
      star1: "IOV SAS multisig star1_TBD_joghurt*iov", // TODO
   },
   iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u: {
      "//id": "vaildator guaranteed reward fund",
      star1: "IOV SAS multisig star1_TBD_guaranteed", // TODO
   },
};
