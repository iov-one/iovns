"use strict";

/**
 * Compares expected and got, recursively if necessary.
 * @param {*} expected
 * @param {*} got
 * @throws an expect()'s error with a "depth" property that indicates the level of recursive of the error
 */
const compareObjects = ( expected, got, depth = [] ) => {
   try {
      expect( typeof got ).toBe( typeof expected );

      if ( expected === null || expected === undefined ) {
         expect( got ).toBe( expected );
      } else if ( expected instanceof Array ) {
         expect( got.length ).toBe( expected.length );

         expected.forEach( ( el, i ) => {
            depth.push( i );
            compareObjects( el, got[ i ], depth );
            depth.pop();
         } );
      } else if ( expected instanceof Object ) {
         expect( Object.keys( got ).length ).toBe( Object.keys( expected ).length );

         Object.keys( expected ).forEach( key => {
            depth.push( key );
            compareObjects( expected[key], got[key], depth );
            depth.pop();
         } );
      } else if ( typeof expected == "number" && !Number.isInteger( expected ) ) {
         expect( got ).toBeCloseTo( expected, 8 ); // HARD-CODED 8
      } else {
         expect( got ).toBe( expected );
      }
   } catch ( e ) {
      if ( !e.depth ) {
         e.depth = depth;
         e.message += `\nDepth:\n  ${depth.join( "." )}`;
      }

      throw e;
   }
}

export default compareObjects;
