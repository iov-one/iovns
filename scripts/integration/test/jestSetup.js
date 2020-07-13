module.exports = async ( globalConfig ) => {
   console.error( `TODO: init global tmpFiles[] in ${__filename}` );

   const { spawnSync } = await import( "child_process" );

   const log = process.env.CONTINUOUS_INTEGRATION ? console.log : () => {};


   const iovnscli = ( args ) => {
      const cli = spawnSync( "iovnscli", args.concat( [ "--output", "json" ] ) );

      if ( cli.status ) throw new Error( cli.stderr );

      log( `Success on 'iovnscli ${args.join( " " )}'.` );

      return JSON.parse( cli.stdout );
   }


   const iovnscliKeysAdd = ( key, args, mnemonic ) => {
      const cliargs = [ "keys", "add", key ].concat( args, [ "--keyring-backend", "test" ] );
      const cli = spawnSync( "iovnscli", cliargs, { input: `${mnemonic}\n` } );

      if ( cli.status ) throw new Error( cli.stderr );

      log( `Success on 'iovnscli ${cliargs.join( " " )}'.` );

      return cli.stderr; // cosmos-sdk stupidly writes to stderr on success
   }


   // https://github.com/iov-one/iovns/blob/master/docs/cli/MULTISIG.md
   const keysWithMnemonic = {
      "bojack": process.env.MNEMONIC_BOJACK,
      "w1": "salad velvet type bamboo neglect prize guess eternal tornado sadness obvious deliver horn capable apart analyst offer echo noise destroy ocean tumble cricket unable",
      "w2": "salmon post develop tumble funny hobby original vintage history length neglect identify frequent tooth then cluster there gravity bridge grow actress trouble obvious elder",
      "w3": "ahead increase coral dutch visual armed good raw skull blur duty move jazz bundle monster surface stairs error trash day ankle meadow famous universe",
   };
   const have = iovnscli( [ "keys", "list", "--keyring-backend", "test" ] );
   const want = Object.keys( keysWithMnemonic ).concat( [ "p1", "msig1" ] );
   const need = want.reduce( ( previous, key ) => {
      previous[key] = !have.find( o => o.name == key );

      return previous;
   }, {} );

   Object.keys( keysWithMnemonic ).forEach( key => {
      if ( need[key] ) iovnscliKeysAdd( key, [ "--recover" ], keysWithMnemonic[key] )
   } );

   if ( need.p1 )   iovnscliKeysAdd( "p1",    [ "--pubkey=starpub1addwnpepqv80htam6gc7fudf9jseldx3afy8nu8anvk935qdctek0yr27jcqj4yv044" ] );
   if ( need.msig1) iovnscliKeysAdd( "msig1", [ "--multisig=w1,w2,w3,p1", "--multisig-threshold=3" ]);

   log( JSON.stringify( iovnscli( [ "keys", "list", "--keyring-backend", "test" ] ), null, "   " ) );
};
