import { spawnSync } from "child_process";


const iovnscli = args => {
   const cli = spawnSync( "iovnscli", [ ...args, "--output", "json" ] );

   if ( cli.status ) throw cli.error ? cli.error : new Error( cli.stderr.length ? cli.stderr : cli.stdout );

   return JSON.parse( cli.stdout );
};


const main = async () => {
   const pubkeys = {
      a: "starpub1addwnpepqww83fu9wdrswn0m2x9uq9df4ml78mtd333u3w9txq2auc7lpgens7259gh", // antoine
      b: "starpub1addwnpepqtzmf3cy284sn7s7rmcu7qpkg6j6pcpjpanpr67nym88ctltvm3n7jwh63a", // benjamin
      d: "starpub1addwnpepqvpmxxwmnyl3ng6sc2qynyv2866spmp6an6zhhe97hjgz873f8w2uncy6ym", // dave
      k: "starpub1addwnpepq2lhq7m62tmmr5acqrwcdnlzu6mgng8cydfmfm9ckuaxvth6wxz3cnp4kln", // karim
      o: "starpub1addwnpepqd5vtc5j8sl0mxmvxa7ydu4flwh5taxnea423q7d9qup548snshpzyp9z58", // olivier
   };
   const users = {
      a: "antoine",
      b: "benjamin",
      d: "dave",
      k: "karim",
      o: "olivier",
   };
   const user = await new Promise( ( resolve, reject ) => {
      console.log( `Who are you? ("a" for "Antoine", "b" for "Benjamin", etc) [abdko]` );

      process.stdin.resume();
      process.stdin.once( "data", abdko => {
         const input = String( abdko ).trim().toLocaleLowerCase();
         const pubkey = pubkeys[input];

         if ( !pubkey ) reject( new Error( `'${input}' is not one of a, b, d, k, or o.` ) );

         resolve( users[input] );
      } );
   } );
   const have = iovnscli( [ "keys", "list" ] ).map( value => value.name );
   const need = Object.keys( pubkeys ).filter( key => !have.includes( key ) );

   need.forEach( key => {
      const name = users[key];

      if ( name == user ) throw new Error( `Key '${user}' should exist already!  Did you do 'iovnscli keys add ${user} --ledger'?` );
   } );
};


main().then( () => {
   process.exit( 0 );
} ).catch( e => {
   console.error( e.stack || e.message || e );
   process.exit( -1 );
} );
