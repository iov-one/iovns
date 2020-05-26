import { burnTokens, labelAccounts, labelMultisigs, migrate } from "../../lib/migrate";

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
         {
            "address": "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph",
            "coins": [ { "ticker": "IOV", "whole": 1628971 }
            ]
         },
         {
            "address": "iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98",
            "coins": [ { "ticker": "IOV", "whole": 3234710 } ]
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
   const osaka = {
      app_hash: "",
      app_state: {
         cash: [
            {
               address: "bech32:iov15xzzgu5jkltm24hg9r2ykm6gm2d09tzrcqgrr9",
               "//id": 1957,
               coins: [ "126455 IOV" ]
            },
            {
               address: "bech32:iov1tc4vr2756lcme6hqq2xgdn4dycny32cdev9a8g",
               "//id": 1970,
               coins: [ "62500 IOV" ]
            },
            {
               address: "bech32:iov1s3e835efuht3qulf3lrv02dypn036lnpedq275",
               "//id": 1976,
               coins: [ "626325 IOV" ]
            },
            {
               address: "bech32:iov13adwzjxhqhd79l3y5v58vjugtfszv57tthmv0q",
               "//id": 1978,
               coins: [ "470651 IOV" ]
            },
            {
               address: "bech32:iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph",
               "//id": 2096,
               coins: [ "1000000 IOV" ]
            },
            {
               address: "bech32:iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98",
               "//id": 2101,
               coins: [ "3234708 IOV" ]
            },
            {
               address: "bech32:iov14favyxdrkkdk39kl4qsexc99vgscw8dw22g5d3",
               "//id": 2243,
               coins: [ "555555 IOV" ]
            },
            {
               address: "bech32:iov1wvxg0qw8pg9vnkja9mvvdzk74g6lrms7uh7tn8",
               "//id": 2244,
               coins: [ "107824 IOV" ]
            },
            {
               address: "bech32:iov1jukhxtnh58mmag5y65d8xj2e36wq6083w0t69e",
               "//id": 2246,
               coins: [ "77777 IOV" ]
            },
         ],
      },
      chain_id: "iov-mainnet",
      consensus_params: {},
      genesis_time: new Date( "2019-10-10T10:00:00Z" ).toISOString(),
      validators: [],
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

   it( `Should label accounts.`, async () => {
      labelAccounts( dumped, osaka );

      const id2096 = dumped.cash.find( account => account.address == "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph" );
      const id2101 = dumped.cash.find( account => account.address == "iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98" );

      expect( id2096["//id"] ).toEqual( 2096 );
      expect( id2101["//id"] ).toEqual( 2101 );
   } );

   it( `Should migrate.`, async () => {
      migrate( { dumped, genesis, multisigs, osaka } );
   } );
} );
