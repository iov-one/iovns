import fetch from "node-fetch";

"use strict";


const fetchOsakaGenesisFile = async () => {
   const result = await fetch( "https://rpc-private-a-vip-mainnet.iov.one/genesis", ).catch( e => { throw e } );
   const json = await result.json().catch( e => { throw e } );
   const genesis = json.result.genesis;

   return genesis;
}


export default fetchOsakaGenesisFile;
