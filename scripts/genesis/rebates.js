import { Client } from "pg";
import YAML from "yaml";
import fs from "fs";
import { atob } from "Base64";


const main = async () => {
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
