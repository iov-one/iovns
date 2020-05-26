import { burnTokens, labelMultisigs, migrate } from "../../lib/migrate";

"use strict";


describe( "Tests ../../lib/migrate.js.", () => {
   const dumped = {
      "cash": [
         {
            "address": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "coins": [ { "ticker": "IOV", "whole": 35384615 } ]
         },
         {
            "address": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un",
            "coins": [ { "fractional": 510000000, "ticker": "IOV", "whole": 416 } ]
         },
         {
            "address": "iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n",
            "coins": [ { "ticker": "IOV", "whole": 37 } ]
         },
         {
            "address": "iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0",
            "coins": [ { "ticker": "IOV", "whole": 3570582 } ]
         },
      ],
   };
   const genesis = {
      chain_id: __filename,
      genesis_time: new Date( "2020-04-15T10:00:00Z" ).toISOString(),
      accounts: [],
      app_hash: "",
      app_state: {},
      auth: {},
      consensus_params: {},
      crisis: {},
      genutil: {},
      gov: {},
   };
   const multisigs = {
      iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n: {
         "//name": "reward fund",
         address: "cond:gov/rule/0000000000000002",
         star1: "star1rewards",
      },
      iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0: {
         "//name": "IOV SAS employee bonus pool/colloboration appropriation pool",
         address: "cond:multisig/usage/0000000000000002",
         star1: "star1bonuses",
      },
   };

   it( `Should burn tokens.`, async () => {
      const copied = JSON.parse( JSON.stringify( dumped ) );

      burnTokens( copied );

      const hex0x0 = copied.cash.findIndex( wallet => wallet.address == "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u" );

      expect( hex0x0 ).toEqual( -1 );
      expect( copied.cash.length ).toEqual( dumped.cash.length - 1 );
   } );

   it( `Should label multisig accounts.`, async () => {
      labelMultisigs( dumped, multisigs );

      Object.keys( multisigs ).forEach( iov1 => {
         const multisig = dumped.cash.find( wallet => wallet.address == iov1 );

         expect( multisig["//id"] ).toEqual( multisigs[iov1]["//name"] );
      } );
   } );

   it( `Should migrate.`, async () => {
      migrate( { dumped, genesis, multisigs } );
   } );
} );
