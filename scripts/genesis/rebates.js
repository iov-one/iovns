import { atob } from "Base64";
import { Client } from "pg";
import { spawnSync } from "child_process";
import fetch from "node-fetch";
import fs from "fs";
import YAML from "yaml";


const memo = "Thank-you for registering your Starname Network address!"; // HARD-CODED in conjunction with LIKE '%Thank%Starname%Network%'


const main = async () => {
   // upgrade recompense was announced on May 25, so pull accounts that existed just before the first starname-migration registration
   const dump = await fetch( "https://gist.githubusercontent.com/davepuchyr/bf9ab1d2d9ca70326cf58c7c70376554/raw/f3f879ff7fa29c8de8d5b5610ce52d5a38323d31/dump.json" );
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
      WHERE message::text LIKE '%register_token%star1%'
      ORDER BY block_id asc
   ` ) ).rows;
   // pull potential rebate candidates that used send (Ledger users)
   const sends = ( await client.query( `
      SELECT *  FROM public.transactions
      WHERE message -> 'details' ->> 'destination' = 'iov10v69k57z2v0pr3yvtr60pp8g2jx8tdd7f55sv6'
      AND   message -> 'details' ->> 'memo' LIKE '%star1%'
      ORDER BY block_id asc
   ` ) ).rows;
   // pull previously paid; HARD-CODED LIKE in conjunction with memo
   const paid = ( await client.query( `
      SELECT *  FROM public.transactions
      WHERE message -> 'details' ->> 'source' = 'iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un'
      AND   message -> 'details' ->> 'memo' LIKE '%Thank%Starname%Network%'
      ORDER BY block_id asc
   ` ) ).rows.map( row => row.message.details.destination );
   const paid0 = paid.length;
   let payout = 0;
   const recipients = {};

   const filterRecipient = ( iov1, amount, info ) => {
      if ( !eligible.includes( iov1 ) || paid.includes( iov1 ) ) return;

      paid.push( iov1 );
      payout += amount;
      recipients[iov1] = {
         amount,
         info,
         iov1,
      };
   };

   changes.forEach( row => {
      const target = row.message.details.new_targets.find( target => target.blockchain_id == "iov-mainnet" );

      if ( target ) filterRecipient( target.address, 1, row.message.details.username );
   } );
   registers.forEach( row => {
      const target = row.message.details.targets.find( target => target.blockchain_id == "iov-mainnet" );

      if ( target ) filterRecipient( target.address, 10, row.message.details.username );
   } );
   sends.forEach( row => {
      const iov1 = row.message.details.source;
      const amount = ( row.message.details.amount.whole || 0 ) + ( row.message.details.amount.fractional || 0 );

      filterRecipient( iov1, amount + 0.5, row.message.details.memo ); // 0.5 for the anti-spam fee
   } );

   const iov1s = Object.values( recipients ).sort( ( a, b ) => {
      return a.amount == b.amount ? a.info.localeCompare( b.info ) : a.amount - b.amount;
   } ).map( recipient => recipient.iov1 );

   iov1s.forEach( iov1 => console.log( `${iov1} ${recipients[iov1].amount} ${recipients[iov1].info}` ) );
   console.log( `changes == ${changes.length}; registers == ${registers.length}; sends == ${sends.length}; paid0 == ${paid0}; rebates == ${paid.length - paid0}; payout == ${payout};`);

   await new Promise( resolve => {
      console.log( `Send ${paid.length - paid0} rebates totalling ${payout + 0.5 * ( paid.length - paid0 )} including fees via ${process.env.TM}? [yN]` );

      process.stdin.resume();
      process.stdin.once( "data", yn => {
         if ( `${yn}`.toLocaleLowerCase().indexOf( "y" ) == -1 ) process.exit( 0 );

         resolve();
      } );
   } );

   // send rebates
   iov1s.forEach( iov1 => {
      const cli = spawnSync( "bash", [ "-x", "rebate.sh", recipients[iov1].amount, iov1, process.env.TM, memo ] );

      if ( cli.status ) throw cli.error ? cli.error : new Error( cli.stderr.length ? cli.stderr : cli.stdout );
   } );
}


main().then( () => {
   process.exit( 0 );
} ).catch( e => {
   console.error( e.stack || e.message || e );
   process.exit( -1 );
} );
