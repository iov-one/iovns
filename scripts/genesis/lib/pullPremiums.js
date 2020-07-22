import gitPullGist from "./gitPullGist"

"use strict";


const pullPremiums = async () => {
   const gist = await gitPullGist( [ "data", "premium", "premium.csv" ] ).catch( e => { throw e } );
   const premiums = gist.split( "\n" ).reduce( ( accumulator, line ) => {
      const columns = line.split( /,/g );
      const iov1 = columns.shift();
      const star1 = columns.shift();
      const existing = accumulator[iov1] ? accumulator[iov1].starnames : [];

      accumulator[iov1] = {
         star1: star1,
         starnames: existing.concat( columns ),
      };

      return accumulator;
   }, {} );

   return premiums;
}


export default pullPremiums;
