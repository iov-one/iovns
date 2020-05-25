import gitPullGist from "./gitPullGist"

"use strict";


const pullDumpedState = async () => {
   const gist = await gitPullGist( [ "data", "dump", "dump.json" ] );
   const state = JSON.parse( gist );

   return state;
}


export default pullDumpedState;
