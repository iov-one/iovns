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
         "public_key": "",
         "account_number": 0,
         "sequence": 0
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
   const template = {
      "broker": null,
      "certificates": null,
      "domain": args.domain || "",
      "metadata_uri": "",
      "name": args.name || "",
      "owner": args.address || "",
      "targets": args.targets && args.targets.length ? args.targets : null,
      "valid_until": String( Math.ceil( Date.now() / 1000 ) + 365.25 * 24 * 60 * 60 ), // 1 year from now
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
      "valid_until": String( Math.ceil( Date.now() / 1000 ) + 365.25 * 24 * 60 * 60 ), // 1 year from now
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
         const iov = parseInt( escrow.amount[0].whole ); // escrows don't have fractional as of 2020.06.07

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
   // iov1yhk8qqp3wsdg7tefd8u457n9zqsny4nqzp6960 both "upgraded" via Neuma and sent to star1*iov, so drop the star1*iov data
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
 */
export const mapIovToStar = ( dumped, multisigs, indicatives ) => {
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
      }

      return previous + iov; // ...after reduction
   }, 0 );
   custodian.value.coins[0].amount = String( Math.ceil( safeguarded ) );

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

      premiums[iov1].forEach( domain => {
         if ( address == custodian.value.address ) {
            const previous = custodian[`//no star1 ${iov1}`];
            const current = !previous ? domain : ( typeof previous == "object" ? previous.concat( domain ) : [ previous, domain ] );

            custodian[`//no star1 ${iov1}`] = current;
         }

         domains.push( createDomain( { address, iov1, domain } ) );
      } );
   } );

   const iov1 = "iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq";
   const address = multisigs[iov1].star1;
   reserveds.forEach( domain => {
      domains.push( createDomain( { address, iov1, domain } ) );
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
}

/**
 * Patches the jestnet genesis object.
 * @param {Object} genesis - the jestnet genesis object
 */
export const patchJestnet = genesis => {
   if ( genesis.chain_id != "jestnet" ) throw new Error( `Wrong chain_id: ${genesis.chain_id} != jestnet.` );

   genesis.app_state.domain.domains[0].account_renew = "3600";
}

/**
 * Patches the iovns-galaxynet genesis object.
 * @param {Object} genesis - the iovns-galaxynet genesis object
 */
export const patchGalaxynet = genesis => {
   if ( genesis.chain_id != "iovns-galaxynet" ) throw new Error( `Wrong chain_id: ${genesis.chain_id} != iovns-galaxynet.` );

   // make dave rich for testing
   const dave = genesis.app_state.auth.accounts.find( account => account.value.address == "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );

   if ( dave ) dave.value.coins[0].amount = "1000000000000";

   // add other test accounts
   const accounts = [
      {
         "//name": "faucet",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star13dq838nu0wmzvx8ge6z5upvu7uze3xlusnts5c",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000"
               }
            ],
            "public_key": "",
            "account_number": 0,
            "sequence": 0
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
            "public_key": "",
            "account_number": 0,
            "sequence": 0
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
            "public_key": "",
            "account_number": 0,
            "sequence": 0
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
            "public_key": "",
            "account_number": 0,
            "sequence": 0
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
            "public_key": "",
            "account_number": 0,
            "sequence": 0
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
            "public_key": "",
            "account_number": 0,
            "sequence": 0
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
            "public_key": "",
            "account_number": 0,
            "sequence": 0
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
            "public_key": "",
            "account_number": 0,
            "sequence": 0
         }
      },
   ];

   genesis.app_state.auth.accounts.push( ...accounts );

   // hack multisig accounts since pubkeys from others are still pending; TODO: delete when possible
   const hackMultisig = {
      "reward fund":                                                  "star1rad8f5rm39ak03h3ev0q4lrshywjdn3v9fn6w3",
      "IOV SAS":                                                      "star12d063hg3ypass56a52fhap25tfgxyaluu6w02r",
      "IOV SAS employee bonus pool/colloboration appropriation pool": "star1v6na4q8kqljynwkh3gt4katlsrqzsk3ewxv6aw",
      "IOV SAS pending deals pocket; close deal or burn":             "star1vhkg66j3xvzqf4smy9qup5ra8euyjwlpdkdyn4",
      "IOV SAS bounty fund":                                          "star1gxchcu6wycentu6fs977hygqx67kv5n7x25w4g",
      "Unconfirmed contributors/co-founders":                         "star1f27zp27q6d8xqeq768r0gffg7ux34ml69dt67j",
      "escrow isabella*iov":                                          "star1uzn9lxhmw0q2vfgy6d5meh2n7m43fqse6ryru6",
      "escrow kadima*iov":                                            "star1hkeufxdyypclg876kc4u9nxjqudkgh2uecrpm7",
      "escrow guaranteed reward fund":                                "star1v875jc00cqh26k5505p5mt4q8w0ylwypsca3jr",
      "vaildator guaranteed reward fund":                             "star1n0et7nukw4htc56lkuqer67heppfjpdhs525ua",
      "Custodian of missing star1 accounts":                          "star1xc7tn8szhtvcat2k29t6072235gsqcrujd60wy",
      "vaildator guaranteed reward fund":                             "star13c7s0xkmpu9uykn56scwwnkjl07svm69j0jm29",
      "escrow isabella*iov":                                          "star1wywlg9ddad2l5zw7zqgcytwx838x00t7t2qqag",
      "escrow kadima*iov":                                            "star1s7dy7pmhzj8t0s48xnvt0ceug873zn9ue4qqma",
      "escrow joghurt*iov":                                           "star1wy4kze7hanky9kpmvrygad5ar8j37wur4e5e3g",
   };
   const hackMultisigKeys = Object.keys( hackMultisig );
   const hackCustodianStar1 = hackMultisig["Custodian of missing star1 accounts"];

   genesis.app_state.auth.accounts.forEach( account => {
      if ( hackMultisigKeys.findIndex( key => key == account["//id"] ) != -1 ) {
         account.value.address = hackMultisig[account["//id"]];
      }
   } );
   genesis.app_state.domain.domains.forEach( domain => {
      if ( domain.admin.toLowerCase().indexOf( "custodia" ) != -1 ) {
         domain.admin = hackCustodianStar1;
      } else if ( domain.type == "open" ) {
         domain.admin = hackMultisig["IOV SAS"];
      }
   } );
   genesis.app_state.domain.accounts.forEach( account => {
      if ( account.owner.toLowerCase().indexOf( "custodia" ) != -1 ) {
         account.owner = hackCustodianStar1;
      }
   } );

   // set the configuration owner and parameters
   const config = genesis.app_state.configuration.config;

   config["//note"] = "msig1 multisig address from w1,w2,w3,p1 in iovns/docs/cli, threshold 3";
   config.account_grace_period = "1800";
   config.account_renew_count_max = 2;
   config.account_renew_period = "1800";
   config.blockchain_target_max = 3;
   config.certificate_count_max = 3;
   config.certificate_size_max = "1000";
   config.configurer = "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg";
   config.domain_grace_period = "400";
   config.domain_renew_count_max = 2;
   config.domain_renew_period = "1800";
   config.metadata_size_max = "1000";

   // set the incorrect fee implementation TODO: delete after https://github.com/iov-one/iovns/issues/22 is closed
   const fees = genesis.app_state.configuration.fees;

   fees.default_fees = {
      "domain/add_certificates_account": {
         "amount": "10",
         "denom": "uiov"
      },
      "domain/delete_account": {
         "amount": "10",
         "denom": "uiov"
      },
      "domain/delete_certificate_account": {
         "amount": "10",
         "denom": "uiov"
      },
      "domain/delete_domain": {
         "amount": "10",
         "denom": "uiov"
      },
      "domain/register_account": {
         "amount": "10",
         "denom": "uiov"
      },
      "domain/register_domain": {
         "amount": "10",
         "denom": "uiov"
      },
      "domain/renew_account": {
         "amount": "10",
         "denom": "uiov"
      },
      "domain/renew_domain": {
         "amount": "10",
         "denom": "uiov"
      },
      "domain/replace_account_targets": {
         "amount": "10",
         "denom": "uiov"
      },
      "domain/set_account_metadata": {
         "amount": "10",
         "denom": "uiov"
      },
      "domain/transfer_account": {
         "amount": "10",
         "denom": "uiov"
      },
      "domain/transfer_domain": {
         "amount": "10",
         "denom": "uiov"
      }
   };
   fees.level_fees = {
      "domain/register_domain": {
         "1": {
            "amount": "10000",
            "denom": "uiov"
         },
         "2": {
            "amount": "5000",
            "denom": "uiov"
         },
         "3": {
            "amount": "2000",
            "denom": "uiov"
         },
         "4": {
            "amount": "1000",
            "denom": "uiov"
         },
         "5": {
            "amount": "500",
            "denom": "uiov"
         }
      }
   };

   // stabilize valid_untils
   const validUntil = 1609415999;
   const fixTransients = hasValidUntils => {
      hasValidUntils.forEach( hasValidUntil => hasValidUntil.valid_until = String( validUntil ) );
   };

   fixTransients( genesis.app_state.domain.domains );
   fixTransients( genesis.app_state.domain.accounts );
}

/**
 * Patches the iov-mainnet-2 genesis object.
 * @param {Object} genesis - the iov-mainnet-2 genesis object
 */
export const patchMainnet = genesis => {
   if ( genesis.chain_id != "iov-mainnet-2" ) throw new Error( `Wrong chain_id: ${genesis.chain_id} != iov-mainnet-2.` );

   // TODO
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
   const iov2star = mapIovToStar( dumped, multisigs, indicatives );
   const escrows = consolidateEscrows( dumped, source2multisig );
   const { accounts, starnames, domains } = convertToCosmosSdk( dumped, iov2star, multisigs, premiums, reserveds );

   // ...mutate genesis
   genesis.app_state.auth.accounts.push( ...Object.values( accounts ) );
   genesis.app_state.auth.accounts.push( ...Object.values( escrows ) );
   genesis.app_state.domain.accounts.push( ...starnames );
   genesis.app_state.domain.domains.push( ...domains );

   if ( patch ) patch( genesis );

   // write genesis.json before...
   const config = path.join( home, "config" );

   if ( !fs.existsSync( config ) ) fs.mkdirSync( config );
   fs.writeFileSync( path.join( config, "genesis.json" ), stringify( genesis, { space: "  " } ), "utf-8" );

   // ...incorporating gentxs
   if ( gentxs ) addGentxs( gentxs, home );
};
