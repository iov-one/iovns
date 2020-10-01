import { spawnSync } from "child_process";
import fetch from "node-fetch";
import fs from "fs";
import tmp from "tmp";

"use strict";

const tmpFiles = [];

export const chain = process.env.CHAIN;
export const echo = process.env.IOVNSCLI_ECHO == "true";
export const gasPrices = process.env.GAS_PRICES;
export const signer = process.env.SIGNER;
export const urlRest = process.env.URL_REST;
export const urlRpc = process.env.URL_RPC;
export const w1 = "star19jj4wc3lxd54hkzl42m7ze73rzy3dd3wry2f3q"; // w1
export const w2 = "star1l4mvu36chkj9lczjhy9anshptdfm497fune6la"; // w2
export const w3 = "star1aj9qqrftdqussgpnq6lqj08gwy6ysppf53c8e9"; // w3
export const msig1 = "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg"; // msig1


export const iovnscli = ( args ) => {
   const maybeWithKeyring = args.find( arg => arg == "query" ) ? args : args.concat( [ "--keyring-backend", "test" ] );
   const maybeWithChainId = args.find( arg => arg == "--offline" || arg == "signutil" ) ? maybeWithKeyring : maybeWithKeyring.concat( [ "--chain-id", chain, "--node", urlRpc ] );
   const cliargs = maybeWithChainId.concat( [ "--output", "json" ] );
   const cli = spawnSync( "iovnscli", cliargs );
   if ( echo ) console.info( `\n\x1b[94miovnscli ${cliargs.join( " " )} | jq\x1b[89m\n` );

   if ( cli.status ) throw cli.error ? cli.error : new Error( cli.stderr.length ? cli.stderr : cli.stdout ) ;

   return JSON.parse( cli.stdout );
};


export const writeTmpJson = o => {
   const tmpname = tmp.tmpNameSync( { template: "REST.test-XXXXXX.json", unsafeCleanup: true } );

   fs.writeFileSync( tmpname, JSON.stringify( o ) );
   tmpFiles.push( tmpname );

   return tmpname;
};


export const signTx = ( tx, from, multisig = "" ) => {
   const tmpname = writeTmpJson( tx );
   const args = [ "tx", "sign", tmpname, "--from", from ];
   if ( multisig != "" ) args.push( "--multisig", multisig );
   const signed = iovnscli( args );

   return signed;
};


export const postTx = async ( signed ) => {
   const tx = { tx: signed.value, mode: "block" };
   const fetched = await fetch( `${urlRest}/txs`, { method: "POST", body: JSON.stringify( tx ) } );

   return fetched;
};


export const signAndPost = async ( unsigned, from = signer ) => {
   const tx = signTx( unsigned, from );
   const posted = await postTx( tx );

   return posted;
};


export const signAndBroadcastTx = ( unsigned, from = signer ) => {
   const unsignedTmp = writeTmpJson( unsigned );
   const args = [ "tx", "sign", unsignedTmp, "--from", from ];
   const signed = iovnscli( args );
   const signedTmp = writeTmpJson( signed );
   const broadcasted = iovnscli( [ "tx", "broadcast", signedTmp, "--broadcast-mode", "block", "--gas-prices", gasPrices ] );

   return broadcasted;
};

export const fetchObject = async ( url, options ) => {
   const fetched = await fetch( url, options );
   const o = await fetched.json();

   return o;
};


/**
 * Determine the file and line number of the caller assuming we're in jest.
 * @returns {string} file:line
 **/
export const memo = () => {
   try {
      throw new Error( "memo" );
   } catch ( e ) {
      const lines = e.stack.split( "\n" );
      const matches = lines[2].match( /.*\/(.*):(\d+):(\d+)/ );
      const file = matches[1];
      const line = matches[2];

      return `${process.env.HOSTNAME}:${file}:${line}`;
   }
}


/**
 * Signs a tx on behalf of msig1.
 * @param {Array} args iovnscli arguments for the tx
 * @returns {object} tx signed by msig1
 * @see https://github.com/iov-one/iovns/blob/master/docs/cli/MULTISIG.md
 **/
export const msig1SignTx = ( args ) => {
   const unsigned = iovnscli( args );
   const w1Signed = signTx( unsigned, w1, msig1 );
   const w2Signed = signTx( unsigned, w2, msig1 );
   const w3Signed = signTx( unsigned, w3, msig1 );
   const unsignedTmp = writeTmpJson( unsigned );
   const w1Tmp = writeTmpJson( w1Signed );
   const w2Tmp = writeTmpJson( w2Signed );
   const w3Tmp = writeTmpJson( w3Signed );
   const signed = iovnscli( [ "tx", "multisign", unsignedTmp, "msig1", w1Tmp, w2Tmp, w3Tmp, "--gas-prices", gasPrices ] );

   return signed;
}


/**
 * Generates the arguments for the update-config command given a configuration object.
 * @param {Object} configuration the configuration
 * @param {string} from the signer
 * @returns {Array} an args array ready for iovnscli
 **/
export const txUpdateConfigArgs = ( configuration, from ) => {
   return [
      "tx", "configuration", "update-config",
      "--signer", from,
      "--account-grace-period", configuration.account_grace_period,
      "--account-renew-count-max", configuration.account_renew_count_max,
      "--account-renew-period", configuration.account_renew_period,
      "--resource-max", configuration.resources_max,
      "--certificate-count-max", configuration.certificate_count_max,
      "--certificate-size-max", configuration.certificate_size_max,
      "--configurer", configuration.configurer,
      "--domain-grace-period", configuration.domain_grace_period,
      "--domain-renew-count-max", configuration.domain_renew_count_max,
      "--domain-renew-period", configuration.domain_renew_period,
      "--metadata-size-max", configuration.metadata_size_max,
      "--valid-account-name", configuration.valid_account_name,
      "--valid-resource", configuration.valid_resource,
      "--valid-uri", configuration.valid_uri,
      "--valid-domain-name", configuration.valid_domain_name,
      "--gas-prices", gasPrices,
      "--generate-only",
   ];
};
