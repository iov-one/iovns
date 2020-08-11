import fs from "fs";
import path from "path";

"use strict";


const readOsakaGenesisFile = async () => {
   const pwd = path.dirname( process.argv[1] );
   const result = fs.readFileSync( path.join( pwd, "data", "osaka.json" ), "utf-8" );
   const json = JSON.parse( result );
   const genesis = json.result.genesis;

   return genesis;
}


export default readOsakaGenesisFile;
