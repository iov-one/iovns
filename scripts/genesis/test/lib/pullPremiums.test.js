import pullPremiums from "../../lib/pullPremiums";

"use strict";


describe( "Tests ../../lib/pullPremiums.js.", () => {
   it( `Should get the premium starnames that were pushed from firebase to a gist.`, async () => {
      const premiums = await pullPremiums();

      expect( Object.keys( premiums ).length ).toBeGreaterThanOrEqual( 105 );

      const starnames = Object.values( premiums ).reduce( ( accumulator, array ) => {
         return accumulator.concat( array );
      }, [] );

      expect( starnames.length ).toBeGreaterThanOrEqual( 270 );
   } );
} );
