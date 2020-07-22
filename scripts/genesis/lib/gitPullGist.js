import { spawn } from "child_process";
import fs from "fs";
import path from "path";

"use strict";

/**
 * Changes directory relative to path.dirname( process.argv[1] ) based on parts and then does `git pull`.
 * @param {Array} parts - subdirectory(s) leading to the gist and its file name, eg [ "data", "dump", "dump.json" ]
 */
const gitPullGist = async ( parts ) => {
   const pwd = path.dirname( process.argv[1] );

   process.chdir( path.join( pwd, ...parts.slice( 0, parts.length - 1 ) ) );

   const git = spawn( "git", [ "pull" ] );
   const stdout = [];

   git.once( "exit", code => {
      if ( code ) console.error( stdout.join( "\n" ) );
   } );
   git.once( "error", error => {
      throw new Error( error );
   } );

   for await ( const data of git.stdout ) {
      stdout.push( String( data ) );
   }

   const text = fs.readFileSync( path.join( pwd, ...parts ), "utf-8" );

   return text;
}


export default gitPullGist;
