import { Base64 } from "js-base64";
import { gasPrices, iovnscli, memo, msig1, msig1SignTx, signAndBroadcastTx, signer, w1, w2, writeTmpJson } from "./common";
import compareObjects from "./compareObjects";
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
      const fees0 = iovnscli( [ "query", "configuration", "get-fees" ] );
      const fees = JSON.parse( JSON.stringify( fees0.fees ) );

      Object.keys( fees ).forEach( key => {
         if ( isFinite( parseFloat( fees[key] ) ) ) {
            fees[key] = String( 1.01 * fees[key] ).substring( 0, 17 ); // max precision is 18
         }
      } );

      const feesTmp = writeTmpJson( fees );
      const signed = msig1SignTx( [ "tx", "configuration", "update-fees", "--from", msig1, "--fees-file", feesTmp, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const signedTmp = writeTmpJson( signed );
      const broadcasted = iovnscli( [ "tx", "broadcast", signedTmp, "--broadcast-mode", "block", "--gas-prices", gasPrices ] );
      const updated = iovnscli( [ "query", "configuration", "get-fees" ] );
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
      compare( fees, updated.fees );

      // restore original fees
      const fees0Tmp = writeTmpJson( fees0.fees );
      const signed0 = msig1SignTx( [ "tx", "configuration", "update-fees", "--from", msig1, "--fees-file", fees0Tmp, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const signed0Tmp = writeTmpJson( signed0 );
      const restore = iovnscli( [ "tx", "broadcast", signed0Tmp, "--broadcast-mode", "block", "--gas-prices", gasPrices ] );
      const restored = iovnscli( [ "query", "configuration", "get-fees" ] );

      expect( restore.gas_used ).toBeDefined();
      compare( fees0.fees, restored.fees );
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

      const broadcasted = signAndBroadcastTx( unsigned );
      const resolved = iovnscli( [ "query", "starname", "resolve", "--starname", `${name}*${domain}` ] );

      expect( broadcasted.gas_used ).toBeDefined();
      if ( !broadcasted.logs ) throw new Error( broadcasted.raw_log );

      expect( resolved.account.domain ).toEqual( domain );
      expect( resolved.account.name ).toEqual( name );
      expect( resolved.account.owner ).toEqual( signer );
      expect( resolved.account.certificates[0] ).toEqual( base64 );
      expect( Base64.decode( resolved.account.certificates[0] ) ).toEqual( certificate0 );

      // verify signature
      const decoded = Base64.decode( resolved.account.certificates[0] );
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

      const broadcasted = signAndBroadcastTx( unsigned );
      const resolved = iovnscli( [ "query", "starname", "resolve", "--starname", `${name}*${domain}` ] );

      expect( broadcasted.gas_used ).toBeDefined();
      if ( !broadcasted.logs ) throw new Error( broadcasted.raw_log );

      expect( resolved.account.domain ).toEqual( domain );
      expect( resolved.account.name ).toEqual( name );
      expect( resolved.account.owner ).toEqual( signer );
      expect( resolved.account.certificates[0] ).toEqual( base64 );
      expect( Base64.decode( resolved.account.certificates[0] ) ).toEqual( invalid );

      // verify signature
      const decoded = Base64.decode( resolved.account.certificates[0] );
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


   it( `Should register a domain with a broker.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const broker = w1;

      const registered = iovnscli( [ "tx", "starname", "register-domain", "--yes", "--broadcast-mode", "block", "--domain", domain, "--broker", broker, "--from", signer, "--gas-prices", gasPrices, "--memo", memo() ] );

      expect( registered.txhash ).toBeDefined();
      if ( !registered.logs ) throw new Error( registered.raw_log );

      const domainInfo = iovnscli( [ "query", "starname", "domain-info", "--domain", domain ] );

      expect( domainInfo.domain.name ).toEqual( domain );
      expect( domainInfo.domain.broker ).toEqual( broker );
   } );


   it( `Should register an account with a broker.`, async () => {
      const domain = "iov";
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const broker = w1;

      const registered = iovnscli( [ "tx", "starname", "register-account", "--yes", "--broadcast-mode", "block", "--domain", domain, "--name", name, "--broker", broker, "--from", signer, "--gas-prices", gasPrices, "--memo", memo() ] );

      expect( registered.txhash ).toBeDefined();
      if ( !registered.logs ) throw new Error( registered.raw_log );

      const resolved = iovnscli( [ "query", "starname", "resolve", "--starname", `${name}*${domain}` ] );

      expect( resolved.account.domain ).toEqual( domain );
      expect( resolved.account.name ).toEqual( name );
      expect( resolved.account.broker ).toEqual( broker );
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


   it( `Should register a domain, register an account, transfer the domain with reset flag 0 (TransferFlush), and query domain-info.`, async () => {
      const transferFlag = "0";
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const metadata = "Why the uri suffix?";
      const metadataEmpty = "top-level corporate info"; // metadata for the empty account
      const registerDomain = iovnscli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const registerAccount = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const setMetadata = iovnscli( [ "tx", "starname", "set-account-metadata", "--domain", domain, "--name", name, "--metadata", metadata, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const setMetadataEmpty = iovnscli( [ "tx", "starname", "set-account-metadata", "--domain", domain, "--name", "", "--metadata", metadataEmpty, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const unsigned = JSON.parse( JSON.stringify( registerDomain ) );

      unsigned.value.msg.push( registerAccount.value.msg[0] );
      unsigned.value.msg.push( setMetadata.value.msg[0] );
      unsigned.value.msg.push( setMetadataEmpty.value.msg[0] );
      unsigned.value.fee.amount[0].amount = "100000000";
      unsigned.value.fee.gas = "600000";

      const broadcasted = signAndBroadcastTx( unsigned );

      expect( broadcasted.gas_used ).toBeDefined();
      if ( !broadcasted.logs ) throw new Error( broadcasted.raw_log );

      const resolved = iovnscli( [ "query", "starname", "resolve", "--starname", `${name}*${domain}` ] );
      const resolvedEmpty = iovnscli( [ "query", "starname", "resolve", "--starname", `*${domain}` ] );

      expect( resolved.account.domain ).toEqual( domain );
      expect( resolved.account.name ).toEqual( name );
      expect( resolved.account.owner ).toEqual( signer );
      expect( resolved.account.metadata_uri ).toEqual( metadata );
      expect( resolvedEmpty.account.owner ).toEqual( signer );
      expect( resolvedEmpty.account.metadata_uri ).toEqual( metadataEmpty );

      const recipient = w1;
      const transferred = iovnscli( [ "tx", "starname", "transfer-domain", "--yes", "--broadcast-mode", "block", "--domain", domain, "--new-owner", recipient, "--transfer-flag", transferFlag, "--from", signer, "--gas-prices", gasPrices, "--memo", memo() ] );

      expect( transferred.gas_used ).toBeDefined();
      if ( !transferred.logs ) throw new Error( transferred.raw_log );

      const newDomainInfo = iovnscli( [ "query", "starname", "domain-info", "--domain", domain ] );
      const newResolvedEmpty = iovnscli( [ "query", "starname", "resolve", "--starname", `*${domain}` ] );

      expect( newDomainInfo.domain.name ).toEqual( domain );
      expect( newDomainInfo.domain.admin ).toEqual( recipient );
      expect( newResolvedEmpty.account.owner ).toEqual( recipient );
      expect( newResolvedEmpty.account.metadata_uri ).toEqual( "" );

      expect( () => {
         iovnscli( [ "query", "starname", "resolve", "--starname", `${name}*${domain}` ] );
      } ).toThrow( `account does not exist: not found in domain ${domain}: ${name}` );
   } );


   it( `Should register a domain, register an account, transfer the domain with reset flag 1 (TransferOwned), and query domain-info.`, async () => {
      const transferFlag = "1";
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const nameOther = `${Math.floor( Math.random() * 1e9 )}`;
      const other = w2; // 3rd party account owner in this case
      const metadata = "Why the uri suffix?";
      const metadataEmpty = "top-level corporate info"; // metadata for the empty account
      const registerDomain = iovnscli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const registerAccount = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const registerAccountOther = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", nameOther, "--owner", other, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const setMetadata = iovnscli( [ "tx", "starname", "set-account-metadata", "--domain", domain, "--name", name, "--metadata", metadata, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const setMetadataEmpty = iovnscli( [ "tx", "starname", "set-account-metadata", "--domain", domain, "--name", "", "--metadata", metadataEmpty, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const unsigned = JSON.parse( JSON.stringify( registerDomain ) );

      unsigned.value.msg.push( registerAccount.value.msg[0] );
      unsigned.value.msg.push( registerAccountOther.value.msg[0] );
      unsigned.value.msg.push( setMetadata.value.msg[0] );
      unsigned.value.msg.push( setMetadataEmpty.value.msg[0] );
      unsigned.value.fee.amount[0].amount = "100000000";
      unsigned.value.fee.gas = "600000";

      const broadcasted = signAndBroadcastTx( unsigned );

      expect( broadcasted.gas_used ).toBeDefined();
      if ( !broadcasted.logs ) throw new Error( broadcasted.raw_log );

      const resolved = iovnscli( [ "query", "starname", "resolve", "--starname", `${name}*${domain}` ] );
      const resolvedEmpty = iovnscli( [ "query", "starname", "resolve", "--starname", `*${domain}` ] );
      const resolvedOther = iovnscli( [ "query", "starname", "resolve", "--starname", `${nameOther}*${domain}` ] );

      expect( resolved.account.domain ).toEqual( domain );
      expect( resolved.account.name ).toEqual( name );
      expect( resolved.account.owner ).toEqual( signer );
      expect( resolved.account.metadata_uri ).toEqual( metadata );
      expect( resolvedEmpty.account.owner ).toEqual( signer );
      expect( resolvedEmpty.account.metadata_uri ).toEqual( metadataEmpty );
      expect( resolvedOther.account.owner ).toEqual( other );

      const recipient = w1;
      const transferred = iovnscli( [ "tx", "starname", "transfer-domain", "--yes", "--broadcast-mode", "block", "--domain", domain, "--new-owner", recipient, "--transfer-flag", transferFlag, "--from", signer, "--gas-prices", gasPrices, "--memo", memo() ] );

      expect( transferred.gas_used ).toBeDefined();
      if ( !transferred.logs ) throw new Error( transferred.raw_log );

      const newDomainInfo = iovnscli( [ "query", "starname", "domain-info", "--domain", domain ] );
      const newResolved = iovnscli( [ "query", "starname", "resolve", "--starname", `${name}*${domain}` ] );
      const newResolvedEmpty = iovnscli( [ "query", "starname", "resolve", "--starname", `*${domain}` ] );
      const newResolvedOther = iovnscli( [ "query", "starname", "resolve", "--starname", `${nameOther}*${domain}` ] );

      expect( newDomainInfo.domain.name ).toEqual( domain );
      expect( newDomainInfo.domain.admin ).toEqual( recipient );
      expect( newResolved.account.owner ).toEqual( recipient );
      expect( newResolved.account.metadata_uri ).toEqual( metadata );
      expect( newResolvedEmpty.account.owner ).toEqual( recipient );
      expect( newResolvedEmpty.account.metadata_uri ).toEqual( metadataEmpty );
      expect( newResolvedOther.account.owner ).toEqual( other );
   } );


   it( `Should register a domain, transfer it with reset flag 2 (ResetNone, the default), and query domain-info.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const registered = iovnscli( [ "tx", "starname", "register-domain", "--yes", "--broadcast-mode", "block", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--memo", memo() ] );

      expect( registered.txhash ).toBeDefined();
      if ( !registered.logs ) throw new Error( registered.raw_log );

      const domainInfo = iovnscli( [ "query", "starname", "domain-info", "--domain", domain ] );

      expect( domainInfo.domain.name ).toEqual( domain );
      expect( domainInfo.domain.admin ).toEqual( signer );

      const recipient = w1;
      const transferred = iovnscli( [ "tx", "starname", "transfer-domain", "--yes", "--broadcast-mode", "block", "--domain", domain, "--new-owner", recipient, "--from", signer, "--gas-prices", gasPrices, "--memo", memo() ] );

      expect( transferred.txhash ).toBeDefined();
      if ( !transferred.logs ) throw new Error( transferred.raw_log );

      const newDomainInfo = iovnscli( [ "query", "starname", "domain-info", "--domain", domain ] );

      expect( newDomainInfo.domain.name ).toEqual( domain );
      expect( newDomainInfo.domain.admin ).toEqual( recipient );
   } );


   it( `Should register an open domain and transfer it.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const registered = iovnscli( [ "tx", "starname", "register-domain", "--yes", "--broadcast-mode", "block", "--type", "open", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--memo", memo() ] );

      expect( registered.txhash ).toBeDefined();
      if ( !registered.logs ) throw new Error( registered.raw_log );

      const domainInfo = iovnscli( [ "query", "starname", "domain-info", "--domain", domain ] );

      expect( domainInfo.domain.name ).toEqual( domain );
      expect( domainInfo.domain.admin ).toEqual( signer );
      expect( domainInfo.domain.type ).toEqual( "open" );

      const recipient = w1;
      const transferred = iovnscli( [ "tx", "starname", "transfer-domain", "--yes", "--broadcast-mode", "block", "--domain", domain, "--new-owner", recipient, "--from", signer, "--gas-prices", gasPrices, "--memo", memo() ] );

      expect( transferred.txhash ).toBeDefined();
      if ( !transferred.logs ) throw new Error( transferred.raw_log );

      const newDomainInfo = iovnscli( [ "query", "starname", "domain-info", "--domain", domain ] );

      expect( newDomainInfo.domain.name ).toEqual( domain );
      expect( newDomainInfo.domain.admin ).toEqual( recipient );
      expect( newDomainInfo.domain.type ).toEqual( "open" );
   } );


   it( `Should register and renew a closed domain.`, async () => {
      // register
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;

      const registered = iovnscli( [ "tx", "starname", "register-domain", "--yes", "--broadcast-mode", "block", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--memo", memo() ] );

      expect( registered.txhash ).toBeDefined();
      if ( !registered.logs ) throw new Error( registered.raw_log );

      const domainInfo = iovnscli( [ "query", "starname", "domain-info", "--domain", domain ] );

      expect( domainInfo.domain.name ).toEqual( domain );

      // renew
      const gas = 1234567;
      const balance0 = iovnscli( [ "query", "account", signer ] );
      const renewed = iovnscli( [ "tx", "starname", "renew-domain", "--yes", "--broadcast-mode", "block", "--domain", domain, "--from", signer, "--gas", gas, "--gas-prices", gasPrices, "--memo", memo() ] );

      expect( renewed.txhash ).toBeDefined();
      if ( !renewed.logs ) throw new Error( renewed.raw_log );

      const newDomainInfo = iovnscli( [ "query", "starname", "domain-info", "--domain", domain ] );
      const starname = iovnscli( [ "query", "starname", "resolve", "--starname", `*${domain}` ] );
      const balance = iovnscli( [ "query", "account", signer ] );
      const fees = iovnscli( [ "query", "configuration", "get-fees" ] ).fees;

      expect( newDomainInfo.domain.name ).toEqual( domain );
      expect( newDomainInfo.domain.valid_until ).toBeGreaterThan( domainInfo.domain.valid_until );
      expect( newDomainInfo.domain.valid_until ).toEqual( starname.account.valid_until );
      // FIXME: BUG: renewal of a closed domain mistakenly only charges for renewing the accounts in the domain; there's only the empty account in this case
      expect( +balance0.value.coins[0].amount - +balance.value.coins[0].amount ).toEqual( gas * parseFloat( gasPrices ) + 1 * fees.register_account_closed / fees.fee_coin_price );
   } );


   it( `Should sign a message, verify it, and fail verification after message alteration.`, async () => {
      const message = "Hello, World!";
      const created = iovnscli( [ "tx", "signutil", "create", "--text", message, "--from", signer, "--memo", memo(), "--generate-only" ] );
      const tmpCreated = writeTmpJson( created );
      const signed = iovnscli( [ "tx", "sign", tmpCreated, "--from", signer, "--offline", "--chain-id", "signed-message-v1", "--account-number", "0", "--sequence", "0" ] );
      const tmpSigned = writeTmpJson( signed );
      const verified = iovnscli( [ "tx", "signutil", "verify", "--file", tmpSigned ] );

      expect( verified.message ).toEqual( message );
      expect( verified.signer ).toEqual( signer );

      // alter the y+NyzKwBpsPJ2xdZMYR4CkFMjhHh004gnRmyXqoWN9J7kqOHxNaevG7TMSvs/NnOT649kbxHUim7koWkvGy8Ew== signature
      signed.value.signatures[0].signature = "z" + signed.value.signatures[0].signature.substr( 1 );

      const tmpAltered = writeTmpJson( signed );

      try {
         iovnscli( [ "tx", "signutil", "verify", "--file", tmpAltered ] );
      } catch ( e ) {
         expect( e.message.indexOf( "ERROR: invalid signature from address found at index 0" ) ).toEqual( 0 );
      }
   } );


   it( `Should do a reverse look-up.`, async () => {
      const uri = "asset:eth";
      const resource = "0x6DF432079347050e0D8dA43C21fa6fe54697AfA7"; // 01node*iov
      const result = iovnscli( [ "query", "starname", "resolve-resource", "--uri", uri, "--resource", resource ] );

      expect( result.accounts.length ).toEqual( 1 );

      const account0 = result.accounts[0];

      expect( account0.name ).toEqual( "01node" );
      expect( account0.domain ).toEqual( "iov" );
      expect( account0.resources.find( r => r.uri == uri && r.resource == resource ) ).toBeDefined();
   } );


   // don't skip once https://github.com/iov-one/iovns/issues/369 is closed
   it.skip( `Should register a domain and account, set resources, and delete resources.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const resources = [
         {
            "uri": "cosmos:iov-mainnet-2",
            "resource": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk"
         }
      ];
      const fileResources = writeTmpJson( resources );
      const registerDomain = iovnscli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const registerAccount = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const replaceResources = iovnscli( [ "tx", "starname", "replace-resources", "--domain", domain, "--name", name, "--src", fileResources, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const unsigned = JSON.parse( JSON.stringify( registerDomain ) );

      unsigned.value.msg.push( registerAccount.value.msg[0] );
      unsigned.value.msg.push( replaceResources.value.msg[0] );
      unsigned.value.fee.amount[0].amount = "400000";
      unsigned.value.fee.gas = "400000";

      const broadcasted = signAndBroadcastTx( unsigned );

      expect( broadcasted.gas_used ).toBeDefined();
      if ( !broadcasted.logs ) throw new Error( broadcasted.raw_log );

      const resolved = iovnscli( [ "query", "starname", "resolve", "--starname", `${name}*${domain}` ] );

      expect( resolved.account.domain ).toEqual( domain );
      expect( resolved.account.name ).toEqual( name );
      expect( resolved.account.owner ).toEqual( signer );
      compareObjects( resources, resolved.account.resources );

      const emptyResources = null;
      const tmpResources = writeTmpJson( emptyResources );
      const replaceResources1 = iovnscli( [ "tx", "starname", "replace-resources", "--domain", domain, "--name", name, "--src", tmpResources, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const broadcasted1 = signAndBroadcastTx( replaceResources1 );

      expect( broadcasted1.gas_used ).toBeDefined();
      if ( !broadcasted1.logs ) throw new Error( broadcasted.raw_log );

      const resolved1 = iovnscli( [ "query", "starname", "resolve", "--starname", `${name}*${domain}` ] );

      compareObjects( emptyResources, resolved1.account.resources );
   } );


   // don't skip once https://github.com/iov-one/iovns/issues/370 is closed
   it.skip( `Should register a domain and account, set metadata, and delete metadata.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const metadata = "Not empty.";
      const registerDomain = iovnscli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const registerAccount = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const setMetadata = iovnscli( [ "tx", "starname", "set-account-metadata", "--domain", domain, "--name", name, "--metadata", metadata, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const unsigned = JSON.parse( JSON.stringify( registerDomain ) );

      unsigned.value.msg.push( registerAccount.value.msg[0] );
      unsigned.value.msg.push( setMetadata.value.msg[0] );
      unsigned.value.fee.amount[0].amount = "400000";
      unsigned.value.fee.gas = "400000";

      const broadcasted = signAndBroadcastTx( unsigned );

      expect( broadcasted.gas_used ).toBeDefined();
      if ( !broadcasted.logs ) throw new Error( broadcasted.raw_log );

      const resolved = iovnscli( [ "query", "starname", "resolve", "--starname", `${name}*${domain}` ] );

      expect( resolved.account.domain ).toEqual( domain );
      expect( resolved.account.name ).toEqual( name );
      expect( resolved.account.owner ).toEqual( signer );
      expect( resolved.account.metadata_uri ).toEqual( metadata );

      const metadata1 = "";
      const setMetadata1 = iovnscli( [ "tx", "starname", "set-account-metadata", "--domain", domain, "--name", name, "--metadata", metadata1, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const broadcasted1 = signAndBroadcastTx( setMetadata1 );

      expect( broadcasted1.gas_used ).toBeDefined();
      if ( !broadcasted1.logs ) throw new Error( broadcasted.raw_log );

      const resolved1 = iovnscli( [ "query", "starname", "resolve", "--starname", `${name}*${domain}` ] );

      expect( resolved1.account.metadata_uri ).toEqual( metadata1 );
   } );
} );
