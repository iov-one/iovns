"use strict";

/**
 * Burns tokens from the dumped state by deleting their entry in dumped.cash.
 * @param {Object} dumped - the state of the weave-based chain
 */
export const burnTokens = dumped => {
   const hex0x0 = "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u";
   const index = dumped.cash.findIndex( wallet => wallet.address == hex0x0 );

   if ( index == -1 ) throw new Error( `Couldn't find ${hex0x0} in dumped.cash.` );

   dumped.cash.splice( index, 1 );
};

/**
 * Performs all the necessary transformations to migrate from the weave-based chain to a cosmos-sdk-based chain.
 * @param {Object} args - various objects required for the transformation
 */
export const migrate = args => {
   const dumped = args.dumped;
   const genesis = args.genesis;
   const osaka = args.osaka;

   burnTokens( dumped );
};
