import { Base64 } from "js-base64";
import { fetchObject, gasPrices, iovnscli, memo, msig1, msig1SignTx, signAndPost, signer, urlRest, w1, writeTmpJson } from "./common";
import forge from "node-forge";

"use strict";


describe( "Tests the CLI.", () => {
   const validCertificate = `{"cert": {"certifier": {"id": "WeStart", "public_key": "b'344a77619d8d6a90d0fbc092880d89607117a9f6fee00ebbf7d3ffa47015fe01'", "URL": "https://www.westart.co/publickey"}, "entity": {"entity_type": "for profit", "registered_name": "IOV SAS", "registred_country": "FR", "VAT_number": "FR31813849017", "URL": "iov.one", "registration_date": "01/03/2018", "address": "55 rue la Boetie", "registered_city": "Paris"}, "starname": {"starname_owner_address": "hjkwbdkj", "starname": "*bestname"}}, "signature": "b'aeef538a01b2ca99a46cd119c9a33a3db1ed7aac15ae890dfe5e29efe329f9dfb7ce179fb4bd4b0ff7424a5981cb9f9408ebcbc8ea998d8478f9bc1276080e0a'"}`;
   const validator = iovnscli( [ "query", "staking", "validators" ] ).find( validator => !validator.jailed ).operator_address;


   it( `Should do a multisig delegate.`, async () => {
      let delegated0 = 0;
      try {
         const delegation = iovnscli( [ "query", "staking", "delegation", msig1, validator ] );
         delegated0 = +delegation.balance.amount;
      } catch ( e ) {
         // no-op on no delegations yet
      }

      const amount = 1e9;
      const signed = msig1SignTx( [ "tx", "staking", "delegate", validator, `${amount}uvoi`, "--from", msig1, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const signedTmp = writeTmpJson( signed );

      const broadcasted = iovnscli( [ "tx", "broadcast", signedTmp, "--broadcast-mode", "block", "--gas-prices", gasPrices ] );
      const delegated = iovnscli( [ "query", "staking", "delegation", msig1, validator ] );

      expect( broadcasted.gas_used ).toBeDefined();
      expect( +delegated.balance.amount ).toEqual( delegated0 + amount );
   } );


   it( `Should do a multisig send.`, async () => {
      const amount = 1000000;
      const signed = msig1SignTx( [ "tx", "send", msig1, w1, `${amount}uvoi`, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const signedTmp = writeTmpJson( signed );

      const balance0 = iovnscli( [ "query", "account", w1 ] );
      const balance0Payer = iovnscli( [ "query", "account", msig1 ] );
      const broadcasted = iovnscli( [ "tx", "broadcast", signedTmp, "--broadcast-mode", "block", "--gas-prices", gasPrices ] );
      const balance = iovnscli( [ "query", "account", w1 ] );
      const balancePayer = iovnscli( [ "query", "account", msig1 ] );

      expect( broadcasted.gas_used ).toBeDefined();
      expect( +balance.value.coins[0].amount ).toEqual( amount + +balance0.value.coins[0].amount );
      expect( +balancePayer.value.coins[0].amount ).toBeLessThan( +balance0Payer.value.coins[0].amount - amount );
   } );


   it( `Should update fees.`, async () => {
      const fees0 = await fetchObject( `${urlRest}/configuration/query/fees`, { method: "POST" } );
      const fees = JSON.parse( JSON.stringify( fees0.result.fees ) );

      Object.keys( fees ).forEach( key => {
         if ( isFinite( parseFloat( fees[key] ) ) ) {
            fees[key] = String( 1.01 * fees[key] ).substring( 0, 17 ); // max precision is 18
         }
      } );

      const feesTmp = writeTmpJson( fees );
      const signed = msig1SignTx( [ "tx", "configuration", "update-fees", "--from", msig1, "--fees-file", feesTmp, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const signedTmp = writeTmpJson( signed );
      const broadcasted = iovnscli( [ "tx", "broadcast", signedTmp, "--broadcast-mode", "block", "--gas-prices", gasPrices ] );
      const updated = await fetchObject( `${urlRest}/configuration/query/fees`, { method: "POST" } );
      const compare = ( had, got ) => {
         Object.keys( had ).forEach( key => {
            if ( isFinite( parseFloat( had[key] ) ) ) {
               expect( 1. * had[key] ).toEqual( 1. * got[key] );
            } else {
               expect( had[key] ).toEqual( got[key] );
            }
         } );
      };

      expect( broadcasted.gas_used ).toBeDefined();
      compare( fees, updated.result.fees );

      // restore original fees
      const fees0Tmp = writeTmpJson( fees0.result.fees );
      const signed0 = msig1SignTx( [ "tx", "configuration", "update-fees", "--from", msig1, "--fees-file", fees0Tmp, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const signed0Tmp = writeTmpJson( signed0 );
      const restore = iovnscli( [ "tx", "broadcast", signed0Tmp, "--broadcast-mode", "block", "--gas-prices", gasPrices ] );
      const restored = await fetchObject( `${urlRest}/configuration/query/fees`, { method: "POST" } );

      expect( restore.gas_used ).toBeDefined();
      compare( fees0.result.fees, restored.result.fees );
   } );


   it( `Should verify the fidelity of a certificate.`, async () => {
      const domain = "iov";
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const certificate0 = validCertificate;
      const base64 = Base64.encode( certificate0 );
      const registerAccount = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const addCerts = iovnscli( [ "tx", "starname", "add-certs", "--domain", domain, "--name", name, "--cert", base64, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const unsigned = JSON.parse( JSON.stringify( registerAccount ) );

      unsigned.value.msg.push( addCerts.value.msg[0] );
      unsigned.value.fee.amount[0].amount = "100000000";
      unsigned.value.fee.gas = "400000";

      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.certificates[0] ).toEqual( base64 );
      expect( Base64.decode( resolved.result.account.certificates[0] ) ).toEqual( certificate0 );

      // verify signature
      const decoded = Base64.decode( resolved.result.account.certificates[0] );
      const message = decoded.match( /{"cert": (.*), "signature"/ )[1]; // fragile!
      const certificate = JSON.parse( decoded );
      const verified = forge.ed25519.verify( {
         message: message,
         encoding: "utf8",
         signature: forge.util.hexToBytes( certificate.signature.slice( 2, -1 ) ),
         publicKey: forge.util.hexToBytes( certificate.cert.certifier.public_key.slice( 2, -1 ) ),
      } );

      expect( verified ).toEqual( true );
   } )


   it( `Should puke on an invalid certificate.`, async () => {
      const domain = "iov";
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const invalidity = "scammer";
      const invalid = validCertificate.replace( "hjkwbdkj", invalidity ); // invalidate the certificate
      const base64 = Base64.encode( invalid );
      const registerAccount = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const addCerts = iovnscli( [ "tx", "starname", "add-certs", "--domain", domain, "--name", name, "--cert", base64, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const unsigned = JSON.parse( JSON.stringify( registerAccount ) );

      unsigned.value.msg.push( addCerts.value.msg[0] );
      unsigned.value.fee.amount[0].amount = "100000000";
      unsigned.value.fee.gas = "400000";

      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.certificates[0] ).toEqual( base64 );
      expect( Base64.decode( resolved.result.account.certificates[0] ) ).toEqual( invalid );

      // verify signature
      const decoded = Base64.decode( resolved.result.account.certificates[0] );
      const message = decoded.match( /{"cert": (.*), "signature"/ )[1]; // fragile!
      const certificate = JSON.parse( decoded );
      const verified = forge.ed25519.verify( {
         message: message,
         encoding: "utf8",
         signature: forge.util.hexToBytes( certificate.signature.slice( 2, -1 ) ),
         publicKey: forge.util.hexToBytes( certificate.cert.certifier.public_key.slice( 2, -1 ) ),
      } );

      expect( certificate.cert.starname.starname_owner_address ).toEqual( invalidity );
      expect( verified ).toEqual( false );
   } )


   it.only( `Should register a domain with a broker.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const broker = w1;

      const registered = iovnscli( [ "tx", "starname", "register-domain", "--yes", "--broadcast-mode", "block", "--domain", domain, "--broker", broker, "--from", signer, "--gas-prices", gasPrices, "--memo", memo() ] );
      const domainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( { name: domain } ) } );

      expect( registered.txhash ).toBeDefined();
      expect( domainInfo.result.domain.broker ).toEqual( broker );
   } );


   it( `Should register an account with a broker.`, async () => {
      const domain = "iov";
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const broker = w1;

      const registered = iovnscli( [ "tx", "starname", "register-account", "--yes", "--broadcast-mode", "block", "--domain", domain, "--name", name, "--broker", broker, "--from", signer, "--gas-prices", gasPrices, "--memo", memo() ] );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( registered.txhash ).toBeDefined();
      expect( resolved.result.account.broker ).toEqual( broker );
   } );


   it( `Should do a multisig reward withdrawl.`, async () => {
      const signed = msig1SignTx( [ "tx", "distribution", "withdraw-rewards", validator, "--from", msig1, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const signedTmp = writeTmpJson( signed );

      const balance0 = iovnscli( [ "query", "account", msig1 ] );
      const broadcasted = iovnscli( [ "tx", "broadcast", signedTmp, "--broadcast-mode", "block", "--gas-prices", gasPrices ] );
      const balance = iovnscli( [ "query", "account", msig1 ] );

      expect( broadcasted.gas_used ).toBeDefined();
      expect( +balance.value.coins[0].amount + parseFloat( gasPrices ) * broadcasted.gas_wanted ).toBeGreaterThan( +balance0.value.coins[0].amount );
   } );
} );
