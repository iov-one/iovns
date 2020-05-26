"use strict";

/**
 * Burns tokens from the dumped state by deleting their entry in dumped.cash.
 * @param {Object} dumped - the state of the weave-based chain
 */
export const burnTokens = dumped => {
   const hex0x0 = dumped.cash.findIndex( wallet => wallet.address == "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u" );

   dumped.cash.splice( hex0x0, 1 );
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
