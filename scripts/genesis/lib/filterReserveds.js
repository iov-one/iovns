import fs from "fs";
import path from "path";

"use strict";

/**
 * Reads reserveds.csv and filters domains based on reDomain.
 * @param {string} regex - a string representation of a regular expression for valid domains
 */
const filterReserveds = ( regex ) => {
   const pwd = path.dirname( process.argv[1] );
   const csvs = fs.readFileSync( path.join( pwd, "data", "reserveds.csv" ), "utf-8" ); // UGLY, but so be it
   const reDomain = new RegExp( regex );
   const lines = csvs.split( "\n" );
   const reserveds = [];

   for ( let i = 1, n = lines.length; i < n; ++i ) { // ignore the header
      const domain = lines[i].split( "," )[0].toLocaleLowerCase();

      if ( reDomain.test( domain ) ) reserveds.push( domain );
   }

   return reserveds;
}


export default filterReserveds;
