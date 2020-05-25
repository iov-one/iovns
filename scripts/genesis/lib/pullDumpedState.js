import { spawn } from "child_process";
import fs from "fs";
import path from "path";

"use strict";


const pullDumpedState = async () => {
   const pwd = path.dirname( process.argv[1].indexOf( "build" ) !== -1 ? process.argv[1].replace( "build", "" ) : process.argv[1] );

   process.chdir( path.join( pwd, "data", "dump" ) );

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

   const state = JSON.parse( fs.readFileSync( path.join( pwd, "data", "dump", "dump.json" ), "utf-8" ) );

   return state;
}


export default pullDumpedState;
