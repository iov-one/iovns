import { spawn } from "child_process";
import fs from "fs";
import path from "path";

"use strict";


const pullPremiums = async () => {
   const pwd = path.dirname( process.argv[1].indexOf( "build" ) !== -1 ? process.argv[1].replace( "build", "" ) : process.argv[1] );

   process.chdir( path.join( pwd, "data", "premium" ) );

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

   const text = fs.readFileSync( path.join( pwd, "data", "premium", "premium.csv" ), "utf-8" );
   const premiums = text.split( "\n" ).reduce( ( accumulator, line ) => {
      const columns = line.split( /,/g );
      accumulator[columns.shift()] = columns;
      return accumulator;
   }, {} );

   return premiums;
}


export default pullPremiums;
