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
 * Adds an "//id" property to multisig accounts
 * @param {Object} dumped - the state of the weave-based chain
 * @param {Object} multisigs - a map of iov1 addresses to multisig account data
 */
export const labelMultisigs = ( dumped, multisigs ) => {
   Object.keys( multisigs ).forEach( iov1 => {
      const index = dumped.cash.findIndex( wallet => wallet.address == iov1 );

      if ( index == -1 ) throw new Error( `Couldn't find ${iov1} in dumped.cash.` );

      dumped.cash[index]["//id"] = multisigs[iov1]["//name"];
   } );
}

/**
 * Adds an "//id" property to ordinary accounts
 * @param {Object} dumped - the state of the weave-based chain
 * @param {Object} osaka - the original genesis data
 */
export const labelAccounts = ( dumped, osaka ) => {
   osaka.app_state.cash.forEach( wallet => {
      const account = dumped.cash.find( account => wallet.address.indexOf( account.address ) != -1 );

      if ( account ) account["//id"] = wallet["//id"];
   } );
}

/**
 * Performs all the necessary transformations to migrate from the weave-based chain to a cosmos-sdk-based chain.
 * @param {Object} args - various objects required for the transformation
 */
export const migrate = args => {
   const dumped = args.dumped;
   const genesis = args.genesis;
   const multisigs = args.multisigs;
   const osaka = args.osaka;

   burnTokens( dumped );
   labelAccounts( dumped, osaka );
   labelMultisigs( dumped, multisigs );
};
