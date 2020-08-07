const { multisigs, source2multisig } = require( "./lib/constants" );
const fs = require( "fs" );


const main = async () => {
   if ( !fs.existsSync( process.argv[2] ) ) throw new Error( `usage: node holdings.js /path/to/genesis.json` );

   const json = fs.readFileSync( process.argv[2] );
   const genesis = JSON.parse( json );
   const custodian = "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y";
   const { iov2coin, iov2star, star2iov, nostar2data } = genesis.app_state.auth.accounts.reduce( ( o, account ) => {
      let iov1 = account["//iov1"];
      const star1 = account.value.address;

      if ( !iov1 ) { // escrowed
         if ( star1 == "star15u4kl3lalt8pm2g4m23erlqhylz76rfh50cuv8" ) iov1 = "iov149cn0rauw2773lfdp34njyejg3cfz2d56c0m5t"; // joghurt
         if ( star1 == "star1hjf04872s9rlcdg2wqwvapwttvt3p4gjpp0xmc" ) iov1 = "iov1ytn5f79hlu5y25phgwx79jmeqj78l6d2ey4s35"; // kadima
         if ( star1 == "star1elad203jykd8la6wgfnvk43rzajyqpk0wsme9g" ) iov1 = "iov1cre24vrfn0klc0rc0wy0a9tf4d360d0fc7f3hl"; // isabella
         if ( star1 == "star17w7fjdkr9laphtyj4wxa32rf0evu94xgywxgl4" ) iov1 = "iov170qvwm0tscn5mza3vmaerkzqllvwc3kykkt7kj"; // vaildator guaranteed reward fund
      }

      o.iov2coin[iov1] = {
         IOV: account.value.coins[0]["//IOV"],
         uiov: +account.value.coins[0].amount,
      };
      o.iov2star[iov1] = star1;
      o.star2iov[star1] = iov1;

      if ( star1 == custodian ) {
         o.nostar2data = Object.keys( account ).reduce( ( no, key ) => {
            if ( key.indexOf( "//no star1" ) != -1 ) {
               const id = key.split( " " )[2];

               no[id] = account[key];
            }

            return no;
         }, {} );
      }

      return o;
   }, { iov2star: {}, star2iov: {}, iov2coin: {} } );
   const exchange = "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96"; // https://internetofvalues.slack.com/archives/GPYCU2AJJ/p1596436694013900
   const iov2domain = genesis.app_state.starname.domains.filter( domain => domain.admin != exchange ).reduce( ( o, domain ) => {
      let iov1 = star2iov[domain.admin];

      if ( !iov1 && domain["//iov1"] ) {
         iov1 = domain["//iov1"];
         iov2star[iov1] = domain.admin;
      }

      o[iov1] = o[iov1] ? o[iov1].concat( domain.name ).sort() : [ domain.name ];

      return o;
   }, {} );
   const iov2username = genesis.app_state.starname.accounts.reduce( ( o, account ) => {
      let iov1 = star2iov[account.owner];

      if ( account.owner != custodian && account.owner != exchange && iov1 != account["//iov1"] ) {
         if ( !iov1 ) {
            iov1 = account["//iov1"]; // iov1 doesn't have tokens but does have a username
         } else {
            console.error( `iov1 mismatch on ${account.name}*${account.domain}!  ${iov1} != ${account["//iov1"]}` );
         }
      }

      const starname = `${account.name}*${account.domain}`;

      o[iov1] = o[iov1] ? o[iov1].concat( starname ).sort() : [ starname ];

      return o;
   }, {} );
   const usernamer = ( msigs, id ) => { // hack names/ids to multisig accounts
      Object.keys( msigs ).forEach( iov1 => {
         iov2username[iov1] = [ msigs[iov1][id] ];
      } );
   };

   usernamer( multisigs, "//name" );
   usernamer( source2multisig, "//id" );
   usernamer( {
      iov170qvwm0tscn5mza3vmaerkzqllvwc3kykkt7kj: { id: "vaildator guaranteed reward fund" },
      iov1cre24vrfn0klc0rc0wy0a9tf4d360d0fc7f3hl: { id: "escrow isabella*iov" }, // iov1 is last tranche
      iov1ytn5f79hlu5y25phgwx79jmeqj78l6d2ey4s35: { id: "escrow kadima*iov" }, // iov1 is last tranche
   }, "id" );

   const iov2true = Object.keys( iov2coin ).concat( Object.keys( iov2domain ), Object.keys( iov2username ) ).reduce( ( o, iov1 ) => {
      o[iov1] = true;

      return o;
   }, {} );
   const iov1s = Object.keys( iov2true );
   const header = [ "iov1", "IOV", "uiov", "star1", "starnames" ];

   console.log( header.join( "," ) );
   iov1s.forEach( iov1 => {
      const coins = iov2coin[iov1];
      const iov = coins && coins.IOV ? coins.IOV : "";
      const uiov =  coins && coins.uiov ? coins.uiov : "";
      const usernames = iov2username[iov1] || [];
      const domains = iov2domain[iov1] || [];
      const star1 = iov2star[iov1];
      const line = [ iov1, iov, uiov, star1, ...usernames, ...domains ];

      console.log( line.join( "," ) );
   } );

   Object.keys( nostar2data ).forEach( iov1 => {
      const line = [ iov1 ];
      const data = nostar2data[iov1];
      const type = typeof data;
      const a = type == "object" ? data : [ data ];

      line.push( typeof a[0] == "number" ? a.shift() : "" ); // IOV
      line.push( "" ); // uiov
      line.push( "" ); // star1
      line.push( a ); // usernames and domains

      console.log( line.join( "," ) );
   } );

}


main().then( () => {
   process.exit( 0 );
} ).catch( e => {
   console.error( e.stack || e.message || e );
   process.exit( -1 );
} );
