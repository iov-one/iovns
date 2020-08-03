import { spawnSync } from "child_process";
import fs from "fs";


const iovnscli = args => {
   const cli = spawnSync( "iovnscli", [ ...args, "--output", "json" ] );

   if ( cli.status ) throw cli.error ? cli.error : new Error( cli.stderr.length ? cli.stderr : cli.stdout );

   let o = {};

   try {
    o = JSON.parse( cli.stdout );
   } catch ( e ) {
      // no-op on non-json output
   }

   return o;
};


const main = async () => {
   // read pubkeys
   const raw = fs.readFileSync( "multisig_pubkeys.txt", "utf8" );
   const lines = raw.split( "\n" );
   const reKey = /pubkey: (.*)/;
   const reName = /- name: (.*)/;
   const pubkeys = {};

   for ( let i = 0, n = lines.length; i < n; ++i ) {
      const line = lines[i];

      if ( line.indexOf( "- name:" ) != -1 ) {
         const name = line.match( reName )[1].trim();
         const key = lines[++i].match( reKey )[1].trim();

         pubkeys[name] = key;
      }
   }

   // get user
   const user = await new Promise( ( resolve, reject ) => {
      console.log( `Who are you? ("a" for "Antoine", "b" for "Benjamin", etc) [abdko]` );

      process.stdin.resume();
      process.stdin.once( "data", abdko => {
         const input = String( abdko ).trim().toLocaleLowerCase();
         const users = {
            a: "antoine",
            b: "benjamin",
            d: "dave*iov",
            k: "karim",
            o: "olivier",
         }
         const user = users[input.slice( 0, 1 )];

         user && resolve( user ) || reject( new Error( `'${input}' is not one of a, b, d, k, or o.` ) );
      } );
   } );

   // create needed keys
   const nice = { // convert to nice names
      "reward fund": "rewardFund",
      "IOV SAS": "iovSAS",
      "IOV SAS employee bonus pool/colloboration appropriation pool": "employeePool",
      "IOV SAS pending deals pocket; close deal or burn": "dealsPocket",
      "IOV SAS bounty fund": "bountyFund",
      "Unconfirmed contributors/co-founders": "cofounders",
      "escrow isabella*iov": "escrowIsabella",
      "escrow kadima*iov": "escrowKadima",
      "escrow vaildator guaranteed reward fund": "escrowValidators",
   };
   const have = iovnscli( [ "keys", "list" ] ).map( value => value.name );
   const need = Object.keys( pubkeys ).filter( key => !have.includes( nice[key] ? nice[key] : key ) );

   need.forEach( name => {
      if ( name == user ) throw new Error( `Key '${user}' should exist already!  Did you do 'iovnscli keys add ${user} --ledger'?` );

      iovnscli( [ "keys", "add", nice[name] ? nice[name] : name, "--pubkey", pubkeys[name] ] );
   } );
};


main().then( () => {
   process.exit( 0 );
} ).catch( e => {
   console.error( e.stack || e.message || e );
   process.exit( -1 );
} );
