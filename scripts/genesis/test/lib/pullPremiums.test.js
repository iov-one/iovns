import pullPremiums from "../../lib/pullPremiums";

"use strict";


describe( "Tests ../../lib/pullPremiums.js.", () => {
   it( `Should get the premium starnames that were pushed from firebase to a gist.`, async () => {
      const premiums = await pullPremiums();

      expect( Object.keys( premiums ).length ).toBeGreaterThanOrEqual( 194 );

      const starnames = Object.keys( premiums ).reduce( ( accumulator, iov1 ) => {
         return accumulator.concat( premiums[iov1].starnames );
      }, [] );

      expect( starnames.length ).toBeGreaterThanOrEqual( 434 );
      expect( premiums.iov127r6ct2mmzr0x7qvju6sna3j8k4hdvkdm2q0c9.starnames[0] ).toEqual( "foundation" );
      expect( premiums.iov12gd6weg7py6vs7ujn22h82422arek8cxzhe85p.starnames[0] ).toEqual( "adrian" );
      expect( premiums.iov12gd6weg7py6vs7ujn22h82422arek8cxzhe85p.starnames[84] ).toEqual( "world" );
      expect( premiums.iov173arhwsd632w7kx2h0hwn9xs4uxpwpw5snklxn.starnames[0] ).toEqual( "chris" );
      expect( premiums.iov173arhwsd632w7kx2h0hwn9xs4uxpwpw5snklxn.starnames[8] ).toEqual( "vegan" );
      expect( premiums.iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un.starnames[0] ).toEqual( "in3s" );
      expect( premiums.iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un.starnames[5] ).toEqual( "example" );
      expect( premiums.iov1vssva25f5fec20weadfkewhctg5rycuc9rgxek.starnames[0] ).toEqual( "airport" );
      expect( premiums.iov1vssva25f5fec20weadfkewhctg5rycuc9rgxek.starnames[11] ).toEqual( "widmer" );
   } );
} );
