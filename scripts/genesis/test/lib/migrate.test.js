import { burnTokens, migrate } from "../../lib/migrate";

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

   it( `Should burn tokens.`, async () => {
      burnTokens( dumped );

      const hex0x0 = dumped.cash.findIndex( wallet => wallet.address == "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u" );

      expect( hex0x0 ).toEqual( -1 );
      expect( dumped.cash.length ).toEqual( 1 );
   } );

   it( `Should migrate.`, async () => {
      migrate( { dumped, genesis } );
   } );
} );
