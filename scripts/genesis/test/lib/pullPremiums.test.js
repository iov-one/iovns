import pullPremiums from "../../lib/pullPremiums";

"use strict";


describe( "Tests ../../lib/pullPremiums.js.", () => {
   it( `Should get the premium starnames that were pushed from firebase to a gist.`, async () => {
      const premiums = await pullPremiums();

      expect( Object.keys( premiums ).length ).toBeGreaterThanOrEqual( 105 );

      const starnames = Object.values( premiums ).reduce( ( accumulator, array ) => {
         return accumulator.concat( array );
      }, [] );

      expect( starnames.length ).toBeGreaterThanOrEqual( 274 );
      expect( premiums.iov127r6ct2mmzr0x7qvju6sna3j8k4hdvkdm2q0c9[0] ).toEqual( "foundation" );
      expect( premiums.iov12gd6weg7py6vs7ujn22h82422arek8cxzhe85p[0] ).toEqual( "adrian" );
      expect( premiums.iov12gd6weg7py6vs7ujn22h82422arek8cxzhe85p[84] ).toEqual( "world" );
      expect( premiums.iov173arhwsd632w7kx2h0hwn9xs4uxpwpw5snklxn[0] ).toEqual( "chris" );
      expect( premiums.iov173arhwsd632w7kx2h0hwn9xs4uxpwpw5snklxn[8] ).toEqual( "vegan" );
      expect( premiums.iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un[0] ).toEqual( "in3s" );
      expect( premiums.iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un[5] ).toEqual( "example" );
      expect( premiums.iov1vssva25f5fec20weadfkewhctg5rycuc9rgxek[0] ).toEqual( "airport" );
      expect( premiums.iov1vssva25f5fec20weadfkewhctg5rycuc9rgxek[11] ).toEqual( "widmer" );
   } );
} );
