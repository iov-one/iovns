import fs from "fs";
import path from "path";

"use strict";


const readOsakaGenesisFile = async () => {
   const result = fs.readFileSync( path.join( "data", "osaka.json" ), "utf-8" );
   const json = JSON.parse( result );
   const genesis = json.result.genesis;

   return genesis;
}


export default readOsakaGenesisFile;
