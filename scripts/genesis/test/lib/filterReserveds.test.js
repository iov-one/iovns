import filterReserveds from "../../lib/filterReserveds";

"use strict";


describe( "Tests ../../lib/filterReserveds.js.", () => {
   it( `Should filter domains with non-numeric characters.`, async () => {
      const reserveds = filterReserveds( "^[0-9]+$" ).sort();

      expect( reserveds.length ).toEqual( 24 );
      expect( reserveds[0] ).toEqual( "101" );
      expect( reserveds[reserveds.length - 1] ).toEqual( "9111" );
   } );


   it( `Should filter domains with non-alphabetic characters.`, async () => {
      const reserveds = filterReserveds( "^[a-z]+$" ).sort();

      expect( reserveds.length ).toEqual( 28254 );
      expect( reserveds[0] ).toEqual( "a" );
      expect( reserveds[reserveds.length - 1] ).toEqual( "zzztube" );
   } );


   it( `Should filter domains with non-alphanumeric characters.`, async () => {
      const reserveds = filterReserveds( "^[a-z0-9]+$" ).sort();

      expect( reserveds.length ).toEqual( 29424 );
      expect( reserveds[0] ).toEqual( "000webhost" );
      expect( reserveds[reserveds.length - 1] ).toEqual( "zzztube" );
   } );


   it( `Should filter invalid domains.`, async () => {
      const reserveds = filterReserveds( "^[-_a-z0-9]{4,16}$" ).sort();

      expect( reserveds.length ).toEqual( 24988 );
      expect( reserveds[0] ).toEqual( "000webhost" );
      expect( reserveds[reserveds.length - 1] ).toEqual( "zzztube" );
   } );
} );
