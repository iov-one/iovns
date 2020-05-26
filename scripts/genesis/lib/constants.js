"use strict";


export const multisigs = {
   iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n: {
      "//name": "reward fund",
      address: "cond:gov/rule/0000000000000002",
      star1: "TBD", // TODO
   },
   iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq: {
      "//name": "IOV SAS",
      address: "cond:multisig/usage/0000000000000001",
      star1: "TBD", // TODO
   },
   iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0: {
      "//name": "IOV SAS employee bonus pool/colloboration appropriation pool",
      address: "cond:multisig/usage/0000000000000002",
      star1: "TBD", // TODO
   },
   iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz: {
      "//name": "IOV SAS pending deals pocket; close deal or burn",
      address: "cond:multisig/usage/0000000000000003",
      star1: "TBD", // TODO
   },
   iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu: {
      "//name": "IOV SAS bounty fund",
      address: "cond:multisig/usage/0000000000000004",
      star1: "TBD", // TODO
   },
   iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw: {
      "//name": "Unconfirmed contributors/co-founders",
      address: "cond:multisig/usage/0000000000000005",
      star1: "TBD", // TODO
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
