import { spawnSync } from "child_process";
import fs from "fs";
import path from "path";
import stringify from "json-stable-stringify";

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
      username.Targets && username.Targets.forEach( target => {
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
               "denom": args.denom ? args.denom : "uiov",
               "amount": args.amount ? String( args.amount ) : "0",
            }
         ],
         "public_key": null,
         "account_number": "0",
         "sequence": "0"
      }
   };

   if ( args.id ) template["//id"] = args.id;
   if ( args.iov1 ) template["//iov1"] = args.iov1;
   if ( isFinite( args.iov ) ) template.value.coins[0]["//IOV"] = args.iov;

   return template;
};

/**
 * Creates a starname object.
 * @param {Object} args - optional address, name
 */
export const createStarname = ( args = {} ) => {
   const resources = args.targets && args.targets.length ? args.targets : null;

   resources && resources.forEach( target => { // convert target to resource
      target.uri = target.blockchain_id;
      target.resource = target.address;
      delete( target.blockchain_id );
      delete( target.address );
   } );

   const template = {
      "broker": null,
      "certificates": null,
      "domain": args.domain || "",
      "metadata_uri": "",
      "name": args.name || "",
      "owner": args.address || "",
      "resources": resources,
      "valid_until": String( new Date( "2020-10-01T00:00:00Z" ).getTime() / 1000 ), // just after listing date
   };

   if ( args.iov1 ) template["//iov1"] = args.iov1;

   return template;
};

/**
 * Creates a domain object.
 * @param {Object} args - optional address, name
 */
export const createDomain = ( args = {} ) => {
   const template = {
      "account_renew": String( 10 * 365.25 * 24 * 60 * 60 ), // 10 years in seconds
      "admin": args.address,
      "broker": null,
      "name": args.domain,
      "type": "closed",
      "valid_until": String( args.valid_until || new Date( "2020-10-01T00:00:00Z" ).getTime() / 1000 ), // just after listing date
   };

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
         const account = accumulator[source] || createAccount( { iov:0 } );
         const value = account.value;
         const iov = ( escrow.amount[0].whole || 0 ) + ( escrow.amount[0].fractional / 1e9 || 0 );

         account["//id"] = source2multisig[source]["//id"];
         account["//note"] = `consolidated escrows with source ${source}`;
         account[`//timeout ${new Date( escrow.timeout * 1000 ).toISOString()}`] = `${escrow.address} yields ${escrow.amount[0].whole} ${escrow.amount[0].ticker}`;
         value.address = source2multisig[source].star1;
         value.coins[0]["//IOV"] += iov;
         value.coins[0].amount = `${parseInt( value.coins[0].amount ) + 1e6 * iov}`; // must be a string

         accumulator[source] = account;
      } );

      return accumulator;
   }, {} );

   return multisigs;
}

/**
 * Fixes errors that arose during the migration.
 * @param {Object} dumped - the state of the weave-based chain
 * @param {Array} indicatives - an array of weave txs stemming from sends to star1*iov
 */
export const fixErrors = ( dumped, indicatives ) => {
   // iov1yhk8qqp3wsdg7tefd8u457n9zqsny4nqzp6960 both "upgraded" via Neuma and sent to star1*iov, so drop the star1*iov data as requested: https://internetofvalues.slack.com/archives/CPNRVHG94/p1591714233003600
   const index = indicatives.findIndex( indicative => indicative.message.details.source == "iov1yhk8qqp3wsdg7tefd8u457n9zqsny4nqzp6960" );

   indicatives.splice( index, 1 );

   // iov1fpezwaxfnmef8tyyg4t7avz9a2d9gqh3yh8d8n upgraded Ledger firmware
   dumped.username.find( username => username.Owner == "iov1fpezwaxfnmef8tyyg4t7avz9a2d9gqh3yh8d8n" ).Owner = "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un";
}

/**
 * Maps iov1 addresses to star1 addresses.
 * @param {Object} dumped - the state of the weave-based chain
 * @param {Object} multisigs - a map of iov1 addresses to multisig account data
 * @param {Array} indicatives - an array of weave txs stemming from sends to star1*iov
 * @param {Object} premiums - a map of iov1 addresses to { star1, starnames }
 */
export const mapIovToStar = ( dumped, multisigs, indicatives, premiums ) => {
   const iov2star = {};
   const reMemo = /(star1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{38})/;

   dumped.username.forEach( username => {
      if ( !iov2star[username.Owner] ) {
         const target = username.Targets && username.Targets.find( target => target.address.indexOf( "star1" ) == 0 );

         iov2star[username.Owner] = target ? target.address : false;
      }
   } );
   Object.keys( multisigs ).forEach( iov1 => iov2star[iov1] = multisigs[iov1].star1 );
   indicatives.forEach( tx => {
      const iov1 = tx.message.details.source;
      const star1 = tx.message.details.memo.match( reMemo )[0];

      if ( iov2star[iov1] && iov2star[iov1] != star1 ) throw new Error( `star1 mismatch on ${iov1}!  ${iov2star[iov1]} != ${star1}!` );

      iov2star[iov1] = star1;
   } );
   Object.keys( premiums ).forEach( iov1 => {
      if ( !iov2star[iov1] ) {
         if ( reMemo.test( premiums[iov1].star1 ) ) iov2star[iov1] = premiums[iov1].star1;
      } else {
         if ( premiums[iov1].star1.length && iov2star[iov1] != premiums[iov1].star1 ) console.warn( `star1 mismatch on ${iov1}!  ${iov2star[iov1]} != ${premiums[iov1].star1}!` );
      }
    } );

   return iov2star;
}

/**
 * Converts weave wallets and usernames into cosmos-sdk accounts and accounts (starnames).
 * @param {Object} dumped - the state of the weave-based chain
 * @param {Object} iov2star - a map of iov1 address to star1 addresses
 * @param {Object} multisigs - a map of iov1 addresses to multisig account data
 * @param {Object} premiums - a map of iov1 addresses to arrays of domains
 * @param {Array} reserveds - domains reserved by IOV SAS
 */
export const convertToCosmosSdk = ( dumped, iov2star, multisigs, premiums, reserveds ) => {
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
      const iov = getAmount( wallet );
      const amount = 1e6 * iov;
      const id = multisigs[iov1]["//name"];
      const account = createAccount( { address, amount, id, iov, iov1 } );

      account["//alias"] = multisigs[iov1].address;

      // remove multisig account from dumped.cash before...
      burnTokens( dumped, [ iov1 ] );
      // ...adding it to accounts since we're soon to loop on dumped.cash
      accounts.push( account );
   } );

   const custodian = accounts.find( account => account["//iov1"] == "iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc" );
   const copied = [].concat( dumped.cash ); // copy because the burnTokens() below splices dumped.cash
   const safeguarded = copied.sort( ( a, b ) => getAmount( b ) - getAmount( a ) ).reduce( ( previous, wallet ) => {
      const iov1 = wallet.address;
      const address = iov2star[iov1];
      const iov = getAmount( wallet );
      const amount = 1e6 * iov; // convert to uiov;
      const id = wallet["//id"];

      if ( address ) {
         const account = createAccount( { amount, address, id, iov, iov1 } );

         accounts.push( account );
      } else {
         // burn before...
         burnTokens( dumped, [ iov1 ] );
         // ...adding to the custodial account...
         custodian[`//no star1 ${iov1}`] = iov;
         previous += amount; // ...after reduction
      }

      return previous;
   }, Math.floor( custodian.value.coins[0].amount ) );
   custodian.value.coins[0].amount = String( safeguarded );

   const starnames = dumped.username.sort( ( a, b ) => a.Username.localeCompare( b.Username ) ).map( username => {
      const iov1 = username.Owner;
      const address = iov2star[iov1] || custodian.value.address; // add to the custodial account if needed
      const starname = username.Username;
      const [ name, domain ] = starname.split( "*" );
      const targets = username.Targets && username.Targets.filter( target => target.address != iov1 ) || []; // drop the legacy IOV target

      if ( address == custodian.value.address ) {
         const previous = custodian[`//no star1 ${iov1}`];
         const current = !previous ? starname : ( typeof previous == "object" ? previous.concat( starname ) : [ previous, starname ] );

         custodian[`//no star1 ${iov1}`] = current;
      }

      return createStarname( { address, iov1, domain, name, targets } );
   } );

   const domains = [];

   Object.keys( premiums ).forEach( iov1 => {
      const address = iov2star[iov1] || custodian.value.address; // add to the custodial account if needed

      premiums[iov1].starnames.forEach( domain => {
         if ( address == custodian.value.address ) {
            const previous = custodian[`//no star1 ${iov1}`];
            const current = !previous ? domain : ( typeof previous == "object" ? previous.concat( domain ) : [ previous, domain ] );

            custodian[`//no star1 ${iov1}`] = current;
         }

         domains.push( createDomain( { address, iov1, domain } ) );
      } );
   } );

   // reserve domains
   const address = "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96"; // https://internetofvalues.slack.com/archives/GPYCU2AJJ/p1596436694013900
   const releases = [ // give 8 months to sell
      new Date( "2020-09-16T10:00:00Z" ),
      new Date( "2020-10-14T10:00:00Z" ),
      new Date( "2020-11-18T10:00:00Z" ),
      new Date( "2020-12-16T10:00:00Z" ),
      new Date( "2021-01-20T10:00:00Z" ),
      new Date( "2021-02-17T10:00:00Z" ),
      new Date( "2021-03-17T10:00:00Z" ),
      new Date( "2021-04-14T10:00:00Z" ),
   ];
   reserveds.forEach( ( domain, i ) => {
      if ( !domains.find( existing => existing.name == domain ) ) { // don't allow duplicates
         const valid_until = releases[i % releases.length].getTime() / 1000;

         domains.push( createDomain( { address, domain, valid_until } ) );
      }
   } );

   domains.sort( ( a, b ) => a.name.localeCompare( b.name ) );

   return { accounts, starnames, domains };
}

/**
 * Add gentxs to the genesis.json file in home/config.
 * @param {string} gentxs - the value of the --gentx-dir flag for `iovnsd collect-gentxs`
 * @param {string} home - the value of the --home flag for `iovnsd collect-gentxs`
 */
export const addGentxs = ( gentxs, home ) => {
   const iovnsd = spawnSync( "iovnsd", [ "collect-gentxs", "--gentx-dir", gentxs, "--home", home, "--trace" ] );

   if ( iovnsd.stdout.length ) {
      console.log( `${iovnsd.stdout}` );
   }

   if ( iovnsd.stderr.length ) {
      const error = `${iovnsd.stderr}`;

      if ( error.indexOf( "ERROR" ) != -1 || error.indexOf( "panic" ) != -1 ) throw new Error( error );
   };

   // clean-up automagically generated files so that subsequent runs don't fail
   [ "app.toml", "config.toml", "node_key.json", "priv_validator_key.json" ].forEach( file => fs.unlinkSync( path.join( home, "config", file ) ) );
   fs.unlinkSync( path.join( home, "data", "priv_validator_state.json" ) );
}

/**
 * Patches the jestnet genesis object.
 * @param {Object} genesis - the jestnet genesis object
 */
export const patchJestnet = genesis => {
   if ( genesis.chain_id != "jestnet" ) throw new Error( `Wrong chain_id: ${genesis.chain_id} != jestnet.` );

   genesis.app_state.starname.domains[0].account_renew = "3600";
}

/**
 * Patches the iovns-galaxynet genesis object.
 * @param {Object} genesis - the iovns-galaxynet genesis object
 */
export const patchGalaxynet = genesis => {
   if ( genesis.chain_id != "iovns-galaxynet" ) throw new Error( `Wrong chain_id: ${genesis.chain_id} != iovns-galaxynet.` );

   // make dave and bojack rich for testing
   const dave = genesis.app_state.auth.accounts.find( account => account.value.address == "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
   const bojack = genesis.app_state.auth.accounts.find( account => account.value.address == "star1z6rhjmdh2e9s6lvfzfwrh8a3kjuuy58y74l29t" );

   if ( dave ) dave.value.coins[0].amount = "1000000000000";
   if ( bojack ) bojack.value.coins[0].amount = "1000000000000";

   // add other test accounts
   const accounts = [
      {
         "//name": "Cosmostation",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star186lx23hw4vgc3xzs6eh85y0a294wrva7cznafs",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "100000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
      {
         "//name": "faucet",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star13hestkc5egttc2d7v4f0kcpxzlr5j0zhyq2jxh",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
      {
         "//name": "antoine",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star1ttf6p8ek3s28luqhnhsxjjh6f7r7t6af5u4895",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
      {
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star10lalxx8ml63hs86j64nk76kucf72dsucluexz8",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
      {
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star1936a62ple4uayhsynvzkx5zzz8jv4z2n8x09fu",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
      {
         "//name": "msig1",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
      {
         "//name": "w1",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star19jj4wc3lxd54hkzl42m7ze73rzy3dd3wry2f3q",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
      {
         "//name": "w2",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star1l4mvu36chkj9lczjhy9anshptdfm497fune6la",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
      {
         "//name": "w3",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star1aj9qqrftdqussgpnq6lqj08gwy6ysppf53c8e9",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
   ];

   genesis.app_state.auth.accounts.push( ...accounts );

   // set the configuration owner and parameters
   const config = genesis.app_state.configuration.config;

   config["//note"] = "msig1 multisig address from w1,w2,w3,p1 in iovns/docs/cli, threshold 3";
   config.account_grace_period = 1 * 60 + "000000000"; // (ab)use javascript
   config.account_renew_count_max = 2;
   config.account_renew_period = 3 * 60 + "000000000";
   config.resources_max = 10;
   config.certificate_count_max = 3;
   config.certificate_size_max = "1000";
   config.configurer = "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg"; // intentionally not a mainnet multisig
   config.domain_grace_period = 1 * 60 + "000000000";
   config.domain_renew_count_max = 2;
   config.domain_renew_period = 5 * 60 + "000000000";
   config.metadata_size_max = "1000";

   // use uvoi as the token denomination
   genesis.app_state.auth.accounts.forEach( account => account.value.coins[0].denom = "uvoi" );
   genesis.app_state.mint.params.mint_denom = "uvoi";
   genesis.app_state.staking.params.bond_denom = "uvoi";
   genesis.app_state.crisis.constant_fee.denom = "uvoi";
   genesis.app_state.gov.deposit_params.min_deposit[0].denom = "uvoi";
   genesis.app_state.configuration.fees = { // https://internetofvalues.slack.com/archives/GPYCU2AJJ/p1593018862011300?thread_ts=1593017152.004100&cid=GPYCU2AJJ
      "fee_coin_denom": "uvoi",
      "fee_coin_price": "0.0000001",
      "fee_default": "0.500000000000000000",
      "register_account_closed": "0.500000000000000000",
      "register_account_open": "0.500000000000000000",
      "transfer_account_closed": "0.500000000000000000",
      "transfer_account_open": "10.000000000000000000",
      "replace_account_resources": "1.000000000000000000",
      "add_account_certificate": "50.000000000000000000",
      "del_account_certificate": "10.000000000000000000",
      "set_account_metadata": "15.000000000000000000",
      "register_domain_1": "1000.000000000000000000",
      "register_domain_2": "500.000000000000000000",
      "register_domain_3": "200.000000000000000000",
      "register_domain_4": "100.000000000000000000",
      "register_domain_5": "50.000000000000000000",
      "register_domain_default": "25.000000000000000000",
      "register_open_domain_multiplier": "10.00000000000000000",
      "transfer_domain_closed": "12.500000000000000000",
      "transfer_domain_open": "125.000000000000000000",
      "renew_domain_open": "12345.000000000000000000",
   };

   // convert URIs to testnet
   genesis.app_state.starname.accounts.forEach( account => {
      const resource = account.resources ? account.resources.find( resource => resource.uri == "asset:iov" ) : null;

      if ( resource ) resource.uri = "asset-testnet:iov"; // https://internetofvalues.slack.com/archives/CPNRVHG94/p1595965860011800
   } );

   // attempt a decentralized launch
   genesis.genesis_time = new Date( "2020-08-13T12:00:00Z" ).toISOString();
}

/**
 * Patches the iov-mainnet-2 genesis object.
 * @param {Object} genesis - the iov-mainnet-2 genesis object
 */
export const patchMainnet = genesis => {
   if ( genesis.chain_id != "iov-mainnet-2" ) throw new Error( `Wrong chain_id: ${genesis.chain_id} != iov-mainnet-2.` );

   const custodian = genesis.app_state.auth.accounts.find( account => account["//id"] == "Custodian of missing star1 accounts" );
   const lostKeysInCustody = { // lost keys/ledger firmware upgraded
      iov1jq8z8xl9tqdwjsp44gtkd2c5rpq33e556kg0ft: {
         star1: "star1k9ktkefsdxtydga262re596agdklwjmrf9et90",
         id: 2033,
      },
      iov153n95ekuw9rxfhzspgarqjdwnadmvdt0chcjs4: {
         star1: "star1keaxspy5rgw84azg5w640pp8zdla72ra0n5xh2",
         id: 2024,
      },
      iov14qk7zrz2ewhdmy7cjj68sk6jn3rst4vd7u930y: {
         star1: "star1lgh6ekcnkufs4742qr5znvtlz4vglul9g2p6xl",
         id: 2046,
      },
   };

   Object.keys( lostKeysInCustody ).forEach( iov1 => {
      const recover = custodian[`//no star1 ${iov1}`];
      const iov = recover[0];
      const amount = 1.e6 * iov;
      const address = lostKeysInCustody[iov1].star1;
      const id = lostKeysInCustody[iov1].id;
      const [ name, domain ] = recover[1].split( "*" );

      // remove custody of tokens
      delete( custodian[`//no star1 ${iov1}`] );
      custodian.value.coins[0].amount = String( +custodian.value.coins[0].amount - amount );

      // remove custody of starname
      const starname = genesis.app_state.starname.accounts.find( account => account.domain == domain && account.name == name );
      if ( !starname ) throw new Error( `Starname doesn't exist for ${recover[1]}!` );
      starname.owner = address;

      // create and add account
      if ( genesis.app_state.auth.accounts.find( account => account["//iov1"] == iov1 ) ) throw new Error( `Account for ${iov1} already exists!` );
      const account = createAccount( { address, amount, id, iov, iov1 } );
      genesis.app_state.auth.accounts.push( account );
   } );

   const lostKeysWithStar1 = { // lost keys after star1 address generation
      iov1lfjspe4x5u404sskmv5md4q7u9jcz96zya8krw: {
         star1: "star1lsk9ckth2s870kjqcyl6x5af7gazj6eg7msluq",
         id: 2191,
      },
      iov1ja0syy203qncn28cqmz5zh9kh2xl0xxt36m4qx: {
         star1: "star1f2jpr2guzq3y5yjv667axr26pl6qzyn2hzthfa",
         id: 2192,
      },
      iov1axxtqae3x9jtvv7wavg6fnjgpc27dx7a9jlp9r: {
         star1: "star1xnzwj34e8zefm7g7vtgnphfj6x2qgnq723rq0j",
         id: 2193,
      },
   };

   Object.keys( lostKeysWithStar1 ).forEach( iov1 => {
      const account = genesis.app_state.auth.accounts.find( account => account["//iov1"] == iov1 );
      const starname = genesis.app_state.starname.accounts.find( account => account["//iov1"] == iov1 );
      const resource = starname.resources.find( resource => resource.uri.indexOf( ":iov" ) != -1 );
      const star1 = lostKeysWithStar1[iov1].star1;

      account.value.address = star1;
      starname.owner = star1;
      resource.resource = star1;
   } );

   const getAmount = account => {
      return +account.value.coins[0];
   };

   genesis.app_state.auth.accounts = genesis.app_state.auth.accounts.sort( ( a, b ) => getAmount( b ) - getAmount( a ) );
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
   const gentxs = args.gentxs;
   const home = args.home;
   const indicatives = args.indicatives;
   const multisigs = args.multisigs;
   const osaka = args.osaka;
   const patch = args.patch;
   const premiums = args.premiums;
   const reserveds = args.reserveds;
   const source2multisig = args.source2multisig;

   // massage inputs...
   burnTokens( dumped, flammable );
   labelAccounts( dumped, osaka );
   labelMultisigs( dumped, multisigs );
   fixChainIds( dumped, chainIds );
   fixErrors( dumped, indicatives );

   // ...transform (order matters)...
   const iov2star = mapIovToStar( dumped, multisigs, indicatives, premiums );
   const escrows = consolidateEscrows( dumped, source2multisig );
   const { accounts, starnames, domains } = convertToCosmosSdk( dumped, iov2star, multisigs, premiums, reserveds );

   // ...mutate genesis
   genesis.app_state.auth.accounts.push( ...Object.values( accounts ) );
   genesis.app_state.auth.accounts.push( ...Object.values( escrows ) );
   genesis.app_state.starname.accounts.push( ...starnames );
   genesis.app_state.starname.domains.push( ...domains );

   if ( patch ) patch( genesis );

   // write genesis.json before...
   const config = path.join( home, "config" );
   const file = path.join( config, "genesis.json" );

   if ( !fs.existsSync( config ) ) fs.mkdirSync( config );
   fs.writeFileSync( file, stringify( genesis, { space: "  " } ), "utf-8" );

   // ...incorporating gentxs
   if ( gentxs && fs.readdirSync( gentxs ).length > 1 ) { // account for README
      addGentxs( gentxs, home );

      const unformatted = JSON.parse( fs.readFileSync( file, "utf-8" ) );
      fs.writeFileSync( file, stringify( unformatted, { space: "  " } ), "utf-8" );
   }
};
