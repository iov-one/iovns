import fetch from "node-fetch";

"use strict";


const fetchSendsTo = async recipient => {
   const txs = await fetch( "https://explorer-api.cluster-mainnet.iov.one/api/txs/query", { "credentials": "omit", "headers": { "accept": "*/*", "accept-language": "en-US,en;q=0.9", "content-type": "application/x-www-form-urlencoded", "sec-fetch-dest": "empty", "sec-fetch-mode": "cors", "sec-fetch-site": "cross-site" }, "referrer": __filename, "referrerPolicy": "no-referrer-when-downgrade", "body": `Destination=${recipient}`, "method": "POST", "mode": "cors" } ).catch( e => { throw e } );
   const json = await txs.json().catch( e => { throw e } );
   const sends = json
      .filter( tx => tx.message.details.destination == recipient )
      .sort( ( a, b ) => a.block_height - b.block_height )
   ;

   return sends;
}


export default fetchSendsTo;
