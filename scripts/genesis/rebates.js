import { atob } from "Base64";
import { Client } from "pg";
import fetch from "node-fetch";
import fs from "fs";
import YAML from "yaml";


const main = async () => {
   // upgrade recompense was announced on May 25, so pull accounts that existed as of that date
   const dump = await fetch( "https://gist.githubusercontent.com/davepuchyr/bf9ab1d2d9ca70326cf58c7c70376554/raw/8c07caa095777ad820342ab3ba3edc186d407019/dump.json" );
   const dumped = await dump.json();
   const eligible = dumped.cash.map( wallet => wallet.address ).sort();

   // connect to the block-metrics db
   const yaml = fs.readFileSync( "/run/user/500/keybase/kbfs/team/iov_one.blockchain/credentials/testnet-common/block-metrics-db-settings.yaml", "utf8" )
   const secrets = YAML.parse( yaml );
   const client = new Client( {
      user: atob( secrets.data.POSTGRES_USER ),
      host: "68.183.242.211",
      database: atob( secrets.data.POSTGRES_DB ),
      password: atob( secrets.data.POSTGRES_PASSWORD ),
      port: 5432,
   } );

   client.connect()

   // pull potential rebate candidates that used change_token_targets
   const changes = ( await client.query( `
      SELECT *  FROM public.transactions
      WHERE message::text LIKE '%change_token_targets%star1%'
      ORDER BY block_id asc
   ` ) ).rows;
   // pull potential rebate candidates that used register_token
   const registers = ( await client.query( `
      SELECT *  FROM public.transactions
      WHERE message::text LIKE '%register%star1%'
      ORDER BY block_id asc
   ` ) ).rows;
   // pull potential rebate candidates that used send (Ledger users)
   const sends = ( await client.query( `
      SELECT *  FROM public.transactions
      WHERE message -> 'details' ->> 'destination' = 'iov10v69k57z2v0pr3yvtr60pp8g2jx8tdd7f55sv6'
      ORDER BY block_id asc
   ` ) ).rows;
   // pull previously paid; TODO: fix LIKE
   const paid = ( await client.query( `
      SELECT *  FROM public.transactions
      WHERE message -> 'details' ->> 'source' = 'iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un'
      AND   message -> 'details' ->> 'memo' LIKE '%ick%'
      ORDER BY block_id asc
   ` ) ).rows.map( row => row.message.details.destination );

   const sendRebate = ( iov1, amount ) => {
      if ( iov1 == "iov1gt83tnjgjg92md2yk25ca5hty5te9yd6vz8fw0" ) debugger; // dave*iov on May 25 somehow
      if ( !eligible.includes( iov1 ) || paid.includes( iov1 ) ) return;

      console.log( `${iov1} ${amount}` );

      paid.push( iov1 );
   };

   changes.forEach( row => {
      const target = row.message.details.new_targets.find( target => target.blockchain_id == "iov-mainnet" );

      if ( target ) sendRebate( target.address, 1 );
   } );
   registers.forEach( row => {
      const target = row.message.details.targets.find( target => target.blockchain_id == "iov-mainnet" );

      if ( target ) sendRebate( target.address, 10 );
   } );
   sends.forEach( row => {
      const iov1 = row.message.details.source;
      const amount = ( row.message.details.amount.whole || 0 ) + ( row.message.details.amount.fractional || 0 );

      sendRebate( iov1, amount + 0.5 ); // 0.5 for the anti-spam fee
   } );

   console.log( `changes == ${changes.length}; registers == ${registers.length}; sends == ${sends.length}; `);
}


main().then( () => {
   process.exit( 0 );
} ).catch( e => {
   console.error( e.stack || e.message || e );
   process.exit( -1 );
} );
