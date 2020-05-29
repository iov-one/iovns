"use strict";

/**
 * Burns tokens from the dumped state by deleting their entry in dumped.cash.
 * @param {Object} dumped - the state of the weave-based chain
 * @param {Array} iov1s - addresses to burn
 */
export const burnTokens = ( dumped, iov1s ) => {
   iov1s.forEach( iov1 => {
      const index = dumped.cash.findIndex( wallet => wallet.address == iov1 );

      if ( index == -1 ) throw new Error( `Couldn't find ${iov1} in dumped.cash.` );

      dumped.cash.splice( index, 1 );
   } );
};

/**
 * Adds "//id" and "//iov1" properties to multisig accounts.
 * @param {Object} dumped - the state of the weave-based chain
 * @param {Object} multisigs - a map of iov1 addresses to multisig account data
 */
export const labelMultisigs = ( dumped, multisigs ) => {
   Object.keys( multisigs ).forEach( iov1 => {
      const index = dumped.cash.findIndex( wallet => wallet.address == iov1 );

      if ( index == -1 ) throw new Error( `Couldn't find ${iov1} in dumped.cash.` );

      dumped.cash[index]["//id"] = multisigs[iov1]["//name"];
      dumped.cash[index]["//iov1"] = iov1;
   } );
}

/**
 * Adds an "//id" property to ordinary accounts and "//iov1" property to all accounts.
 * @param {Object} dumped - the state of the weave-based chain
 * @param {Object} osaka - the original genesis data
 */
export const labelAccounts = ( dumped, osaka ) => {
   osaka.app_state.cash.forEach( wallet => {
      const account = dumped.cash.find( account => wallet.address.indexOf( account.address ) != -1 );

      if ( account ) account["//id"] = wallet["//id"];
   } );

   dumped.cash.forEach( account => account["//iov1"] = account.address );
}

/**
 * Fixes chain ids.
 * @param {Object} dumped - the state of the weave-based chain
 * @param {Object} ids - a map of old chain ids to new chain ids
 */
export const fixChainIds = ( dumped, ids ) => {
   dumped.username.forEach( username => {
      username.Targets.forEach( target => {
         if ( ids[target.blockchain_id] ) target.blockchain_id = ids[target.blockchain_id];
      } );
   } );
}

/**
 * Creates a cosmos-sdk account object.
 * @param {Object} args - optional address, denom, and/or amount
 */
export const createAccount = ( args = {} ) => {
   const template = {
      "type": "cosmos-sdk/Account",
      "value": {
         "address": args.address || "",
         "coins": [
            {
               "denom": args.denom ? args.denom : "iov",
               "amount": args.amount ? String( args.amount ) : "0",
            }
         ],
         "public_key": "",
         "account_number": 0,
         "sequence": 0
      }
   };

   if ( args.id ) template["//id"] = args.id;
   if ( args.iov1 ) template["//iov1"] = args.iov1;

   return template;
};

/**
 * Consolidates a given sources' escrows into an IOV SAS controlled multisig account.
 * @param {Object} dumped - the state of the weave-based chain
 * @param {Object} source2multisig - a map of escrow sources to multisig accounts
 */
export const consolidateEscrows = ( dumped, source2multisig ) => {
   const sources = Object.keys( source2multisig );
   const escrows = dumped.escrow.reduce( ( accumulator, escrow ) => {
      if ( !sources.includes( escrow.source ) ) throw new Error( `Unknown escrow source ${escrow.source} in escrow ${JSON.stringify( escrow )}` );

      const existing = accumulator[escrow.source] || [];

      accumulator[escrow.source] = existing.concat( escrow );

      return accumulator;
   }, {} );
   const multisigs = Object.keys( escrows ).reduce( ( accumulator, source ) => {
      const flammable = escrows[source].map( escrow => escrow.address );

      // burn the tokens before...
      burnTokens( dumped, flammable );

      // ...adding them to multisigs
      escrows[source].forEach( escrow => {
         const account = accumulator[source] || createAccount();
         const value = account.value;

         account["//id"] = `consolidated escrows with source ${source} (${source2multisig[source]["//id"]})`;
         account[`//timeout ${new Date( escrow.timeout * 1000 ).toISOString()}`] = `${escrow.address} yields ${escrow.amount[0].whole} ${value.coins[0].denom}`;
         value.address = source2multisig[source].star1;
         value.coins[0].amount = `${parseInt( value.coins[0].amount ) + escrow.amount[0].whole}`; // no fractionals; must be a string

         accumulator[source] = account;
      } );

      return accumulator;
   }, {} );

   return multisigs;
}

/**
 * Maps iov1 addresses to star1 addresses.
 * @param {Object} dumped - the state of the weave-based chain
 * @param {Object} multisigs - a map of iov1 addresses to multisig account data
 */
export const mapIovToStar = ( dumped, multisigs ) => {
   const iov2star = {};

   dumped.username.forEach( username => {
      const target = username.Targets.find( target => target.address.indexOf( "star1" ) == 0 );

      iov2star[username.Owner] = target ? target.address : false;
   } );
   Object.keys( multisigs ).forEach( iov1 => iov2star[iov1] = multisigs[iov1].star1 );

   return iov2star;
}

/**
 * Converts weave wallets into cosmos-sdk accounts.
 * @param {Object} dumped - the state of the weave-based chain
 * @param {Object} multisigs - a map of iov1 addresses to multisig account data
 * @param {Object} iov2star - a map of iov1 address to star1 addresses
 */
export const convertAccounts = ( dumped, iov2star, multisigs ) => {
   const accounts = [];
   const getAmount = wallet => {
      const coins0 = wallet.coins[0];
      const amount = ( coins0.whole || 0 ) + ( coins0.fractional / 1e9 || 0 );

      return amount;
   };

   Object.keys( multisigs ).forEach( iov1 => {
      const index = dumped.cash.findIndex( wallet => wallet.address == iov1 );
      const wallet = dumped.cash[index];
      const address = multisigs[iov1].star1;
      const amount = getAmount( wallet );
      const id = multisigs[iov1]["//name"];
      const account = createAccount( { address, amount, id, iov1 } );

      account["//alias"] = multisigs[iov1].address;

      // remove multisig account from dumped.cash before...
      burnTokens( dumped, [ iov1 ] );
      // ...adding it to accounts since we're soon to loop on dumped.cash
      accounts.push( account );
   } );

   const custodian = accounts.find( account => account["//iov1"] == "iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc" );
   const copied = [].concat( dumped.cash ); // copy because the burnTokens() below splices dumped.cash

   copied.sort( ( a, b ) => getAmount( b ) - getAmount( a ) ).forEach( wallet => {
      const iov1 = wallet.address;
      const address = iov2star[iov1];
      const amount = getAmount( wallet );
      const id = wallet["//id"];
      const star1 = iov2star[iov1];

      if ( star1 ) {
         const account = createAccount( { amount, address, id, iov1 } );

         accounts.push( account );
      } else {
         // burn before...
         burnTokens( dumped, [ iov1 ] );
         // ...adding to the custodial account
         custodian[`//no star1 ${iov1}`] = amount;
         custodian.value.coins[0].amount = String( +custodian.value.coins[0].amount + +amount )
      }
   } );

   return accounts;
}

/**
 * Performs all the necessary transformations to migrate from the weave-based chain to a cosmos-sdk-based chain.
 * @param {Object} args - various objects required for the transformation
 */
export const migrate = args => {
   const chainIds = args.chainIds;
   const dumped = args.dumped;
   const flammable = args.flammable;
   const genesis = args.genesis;
   const multisigs = args.multisigs;
   const osaka = args.osaka;
   const source2multisig = args.source2multisig;

   // massage inputs...
   burnTokens( dumped, flammable );
   labelAccounts( dumped, osaka );
   labelMultisigs( dumped, multisigs );
   fixChainIds( dumped, chainIds );

   // ...transform (order matters)...
   const iov2star = mapIovToStar( dumped, multisigs, source2multisig );
   const escrows = consolidateEscrows( dumped, source2multisig );
   const accounts = convertAccounts( dumped, iov2star, multisigs );

   // ...mutate genesis
   genesis.accounts.push( ...Object.values( accounts ) );
   genesis.accounts.push( ...Object.values( escrows ) );
};
