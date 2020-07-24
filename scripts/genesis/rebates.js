import { atob } from "Base64";
import { Client } from "pg";
import fetch from "node-fetch";
import fs from "fs";
import YAML from "yaml";


const main = async () => {
   // upgrade recompense was announced on May 25, so pull accounts that existed as of that date
   const dump = await fetch( "https://gist.githubusercontent.com/davepuchyr/bf9ab1d2d9ca70326cf58c7c70376554/raw/2c2c7b5694cfeab9dfab479726cefcafa83520c5/dump.json" );
   const dumped = await dump.json();
   const eligible = dumped.cash;

   // connect to the block-metrics db
   const yaml = fs.readFileSync( "/run/user/500/keybase/kbfs/team/iov_one.blockchain/credentials/testnet-common/block-metrics-db-settings.yaml", "utf8" )
   const secrets = YAML.parse( yaml );
   const client = new Client( {
      user: atob( secrets.data.POSTGRES_USER ),
      host: "68.183.242.211",
      database: atob( secrets.data.POSTGRES_DB ),
      password: atob( secrets.data.POSTGRES_PASSWORD ),
      port: 5432,
   } )

   client.connect()

   const res = await client.query( "SELECT NOW()" );
   console.log( res );
}


main().catch( e => {
   console.error( e.stack || e.message || e );
   process.exit( -1 );
} );
