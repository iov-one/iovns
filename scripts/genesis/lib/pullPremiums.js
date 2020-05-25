import gitPullGist from "./gitPullGist"

"use strict";


const pullPremiums = async () => {
   const gist = await gitPullGist( [ "data", "premium", "premium.csv" ] ).catch( e => { throw e } );
   const premiums = gist.split( "\n" ).reduce( ( accumulator, line ) => {
      const columns = line.split( /,/g );
      const address = columns.shift();
      const existing = accumulator[address] || [];

      accumulator[address] = existing.concat( columns );

      return accumulator;
   }, {} );

   return premiums;
}


export default pullPremiums;
