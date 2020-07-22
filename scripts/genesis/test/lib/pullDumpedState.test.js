import pullDumpedState from "../../lib/pullDumpedState";

"use strict";


describe( "Tests ../../lib/pullDumpedState.js.", () => {
   it( `Should get the weave-based mainnet state.`, async () => {
      const weave = await pullDumpedState();
      const burned = weave.cash.find( wallet => wallet.address == "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u" );
      const dave = weave.username.find( username => username.Username == "dave*iov" );

      expect( burned ).toBeTruthy();
      expect( burned.coins[0].whole ).toBeGreaterThanOrEqual( 35384615 );
      expect( dave ).toBeTruthy();
      expect( dave.Owner ).toEqual( "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" );
      expect( dave.Targets.length ).toBeGreaterThanOrEqual( 1 );
      expect( weave.contract.length ).toEqual( 6 );
      expect( weave.escrow.length ).toBeLessThanOrEqual( 18 );
      expect( weave.height ).toBeGreaterThanOrEqual( 65318 );
   } );
} );
