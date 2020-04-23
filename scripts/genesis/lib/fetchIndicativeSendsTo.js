import fetchSendsTo from "./fetchSendsTo";

"use strict";


const fetchIndicativeSendsTo = async ( recipient, reMemo ) => {
   const sends = await fetchSendsTo( recipient ).catch( e => { throw e } );
   const filtered = sends.filter( tx => reMemo.test( tx.message.details.memo ) );

   return filtered;
}


export default fetchIndicativeSendsTo;
