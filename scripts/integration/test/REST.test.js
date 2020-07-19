import { chain, fetchObject, gasPrices, iovnscli, memo, msig1, msig1SignTx, postTx, signAndPost, signer, signTx, txUpdateConfigArgs, urlRest, w1, w2, writeTmpJson, } from "./common";
import { Base64 } from "js-base64";
import compareObjects from "./compareObjects";

"use strict";


describe.skip( "Tests the REST API.", () => {
   it( `Should get node_info.`, async () => {
      const fetched = await fetchObject( `${urlRest}/node_info` );

      expect( fetched.node_info.network ).toEqual( chain );
      expect( fetched.application_version.name ).toEqual( "iovns" );
   } );


   it( `Should get syncing and it should be false.`, async () => {
      const fetched = await fetchObject( `${urlRest}/syncing` );

      expect( fetched.syncing ).toEqual( false );
   } );


   it( `Should get configuration.`, async () => {
      const fetched = await fetchObject( `${urlRest}/configuration/query/configuration`, { method: "POST" } );
      const keys = [
         "configurer",
         "valid_domain_name",
         "valid_account_name",
         "valid_uri",
         "valid_resource",
         "domain_renew_period",
         "domain_renew_count_max",
         "domain_grace_period",
         "account_renew_period",
         "account_renew_count_max",
         "account_grace_period",
         "resources_max",
         "certificate_size_max",
         "certificate_count_max",
         "metadata_size_max",
      ];

      keys.forEach( key => expect( fetched.result.configuration.hasOwnProperty( key ) ).toEqual( true ) );
   } );


   it( `Should get fees.`, async () => {
      const fetched = await fetchObject( `${urlRest}/configuration/query/fees`, { method: "POST" } );
      const keys = [
         "fee_coin_denom",
         "fee_coin_price",
         "fee_default",
         "register_account_closed",
         "register_account_open",
         "transfer_account_closed",
         "transfer_account_open",
         "replace_account_resources",
         "add_account_certificate",
         "del_account_certificate",
         "set_account_metadata",
         "register_domain_1",
         "register_domain_2",
         "register_domain_3",
         "register_domain_4",
         "register_domain_5",
         "register_domain_default",
         "register_open_domain_multiplier",
         "transfer_domain_closed",
         "transfer_domain_open",
         "renew_domain_open",
      ];

      keys.forEach( key => expect( fetched.result.fees.hasOwnProperty( key ) ).toEqual( true ) );
   } );


   it( `Should send.`, async () => {
      const amount = 1e6;
      const recipient = w1;
      const balance0 = iovnscli( [ "query", "account", recipient ] );
      const unsigned = iovnscli( [ "tx", "send", signer, recipient, `${amount}uvoi`, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const posted = await signAndPost( unsigned );
      const balance = iovnscli( [ "query", "account", recipient ] );

      expect( posted.ok ).toEqual( true );
      expect( +balance.value.coins[0].amount - +balance0.value.coins[0].amount ).toEqual( amount );
   } );


   it( `Should register a domain, query domainInfo, and delete the domain.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const unsigned = iovnscli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const posted = await signAndPost( unsigned );
      const body = { name: domain };
      const domainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( domainInfo.result.domain.name ).toEqual( domain );
      expect( domainInfo.result.domain.admin ).toEqual( signer );
      expect( domainInfo.result.domain.type.toLowerCase() ).toEqual( "closed" );

      const delDomain = iovnscli( [ "tx", "starname", "del-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const deleted = await signAndPost( delDomain );
      const noDomain = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( deleted.ok ).toEqual( true );
      expect( noDomain.error ).toBeTruthy();
   } );


   it( `Should register a domain, account, add resources, and query resolve.`, async () => {
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
      unsigned.value.fee.amount[0].amount = "100000000";
      unsigned.value.fee.gas = "600000";

      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );

      compareObjects( resources, resolved.result.account.resources );
   } );


   it( `Should register a domain, account, add metadata, and query resolve.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const metadata = "Why the uri suffix?";
      const registerDomain = iovnscli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const registerAccount = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const setMetadata = iovnscli( [ "tx", "starname", "set-account-metadata", "--domain", domain, "--name", name, "--metadata", metadata, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const unsigned = JSON.parse( JSON.stringify( registerDomain ) );

      unsigned.value.msg.push( registerAccount.value.msg[0] );
      unsigned.value.msg.push( setMetadata.value.msg[0] );
      unsigned.value.fee.amount[0].amount = "100000000";
      unsigned.value.fee.gas = "600000";

      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.metadata_uri ).toEqual( metadata );
   } );


   it( `Should register and delete an account and query resolve.`, async () => {
      const domain = "iov";
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const unsigned = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );

      const delAccount = iovnscli( [ "tx", "starname", "del-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const deleted = await signAndPost( delAccount );
      const noAccount = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( deleted.ok ).toEqual( true );
      expect( noAccount.error ).toBeTruthy();
   } );


   it( `Should register a domain, account, add base64 certificate, delete the certificate, and query resolve.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const certificate = JSON.stringify( { my: "certificate", as: "base64" } );
      const base64 = Base64.encode( certificate );
      const registerDomain = iovnscli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const registerAccount = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const addCerts = iovnscli( [ "tx", "starname", "add-certs", "--domain", domain, "--name", name, "--cert", base64, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const unsigned = JSON.parse( JSON.stringify( registerDomain ) );

      unsigned.value.msg.push( registerAccount.value.msg[0] );
      unsigned.value.msg.push( addCerts.value.msg[0] );
      unsigned.value.fee.amount[0].amount = "100000000";
      unsigned.value.fee.gas = "600000";

      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.certificates[0] ).toEqual( base64 );
      expect( Base64.decode( resolved.result.account.certificates[0] ) ).toEqual( certificate );

      const delCerts = iovnscli( [ "tx", "starname", "del-certs", "--domain", domain, "--name", name, "--cert", base64, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const deleted = await signAndPost( delCerts );
      const noCerts = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( deleted.ok ).toEqual( true );
      expect( noCerts.result.account.certificates ).toBeNull();
   } );


   it( `Should register a domain, account, add certificate via file, and query resolve.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const certificate = JSON.stringify( { my: "certificate", as: "base64" } );
      const file = writeTmpJson( certificate );
      const registerDomain = iovnscli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const registerAccount = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const addCerts = iovnscli( [ "tx", "starname", "add-certs", "--domain", domain, "--name", name, "--cert-file", file, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const unsigned = JSON.parse( JSON.stringify( registerDomain ) );

      unsigned.value.msg.push( registerAccount.value.msg[0] );
      unsigned.value.msg.push( addCerts.value.msg[0] );
      unsigned.value.fee.amount[0].amount = "100000000";
      unsigned.value.fee.gas = "600000";

      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.certificates.length ).toEqual( 1 );

      compareObjects( JSON.parse( certificate ), JSON.parse( JSON.parse( Base64.decode( resolved.result.account.certificates[0] ) ) ) );
   } );


   it( `Should register a domain, transfer it with reset flag 2 (ResetNone, the default), and query domainInfo.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const unsigned = iovnscli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const posted = await signAndPost( unsigned );
      const body = { name: domain };
      const domainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( domainInfo.result.domain.name ).toEqual( domain );
      expect( domainInfo.result.domain.admin ).toEqual( signer );

      const recipient = w1;
      const transferDomain = iovnscli( [ "tx", "starname", "transfer-domain", "--domain", domain, "--new-owner", recipient, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const transferred = await signAndPost( transferDomain );
      const newDomainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( transferred.ok ).toEqual( true );
      expect( newDomainInfo.result.domain.name ).toEqual( domain );
      expect( newDomainInfo.result.domain.admin ).toEqual( recipient );
   } );


   it( `Should register a domain, register an account, transfer the domain with reset flag 0 (TransferFlush), and query domainInfo.`, async () => {
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

      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const starname = { starname: `*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );
      const resolvedEmpty = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starname ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.metadata_uri ).toEqual( metadata );
      expect( resolvedEmpty.result.account.owner ).toEqual( signer );
      expect( resolvedEmpty.result.account.metadata_uri ).toEqual( metadataEmpty );

      const recipient = w1;
      const transferDomain = iovnscli( [ "tx", "starname", "transfer-domain", "--domain", domain, "--new-owner", recipient, "--transfer-flag", transferFlag, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const transferred = await signAndPost( transferDomain );
      const newDomainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( { name: domain } ) } );
      const newResolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );
      const newResolvedEmpty = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starname ) } );

      expect( transferred.ok ).toEqual( true );
      expect( newDomainInfo.result.domain.name ).toEqual( domain );
      expect( newDomainInfo.result.domain.admin ).toEqual( recipient );
      expect( newResolved.error ).toBeTruthy();
      expect( newResolvedEmpty.result.account.owner ).toEqual( recipient );
      expect( newResolvedEmpty.result.account.metadata_uri ).toEqual( "" );
   } );


   it( `Should register a domain, register an account, transfer the domain with reset flag 1 (TransferOwned), and query domainInfo.`, async () => {
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

      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const starname = { starname: `*${domain}` };
      const starnameOther = { starname: `${nameOther}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );
      const resolvedEmpty = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starname ) } );
      const resolvedOther = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starnameOther ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.metadata_uri ).toEqual( metadata );
      expect( resolvedEmpty.result.account.owner ).toEqual( signer );
      expect( resolvedEmpty.result.account.metadata_uri ).toEqual( metadataEmpty );
      expect( resolvedOther.result.account.owner ).toEqual( other );

      const recipient = w1;
      const transferDomain = iovnscli( [ "tx", "starname", "transfer-domain", "--domain", domain, "--new-owner", recipient, "--transfer-flag", transferFlag, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const transferred = await signAndPost( transferDomain );
      const newDomainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( { name: domain } ) } );
      const newResolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );
      const newResolvedEmpty = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starname ) } );
      const newResolvedOther = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starnameOther ) } );

      expect( transferred.ok ).toEqual( true );
      expect( newDomainInfo.result.domain.name ).toEqual( domain );
      expect( newDomainInfo.result.domain.admin ).toEqual( recipient );
      expect( newResolved.result.account.owner ).toEqual( recipient );
      expect( newResolved.result.account.metadata_uri ).toEqual( metadata );
      expect( newResolvedEmpty.result.account.owner ).toEqual( recipient );
      expect( newResolvedEmpty.result.account.metadata_uri ).toEqual( metadataEmpty );
      expect( newResolvedOther.result.account.owner ).toEqual( other );
   } );


   it( `Should register a domain with a broker.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const unsigned = iovnscli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const broker = "star1aj9qqrftdqussgpnq6lqj08gwy6ysppf53c8e9"; // w3

      unsigned.value.msg[0].value.broker = broker;

      const posted = await signAndPost( unsigned );
      const body = { name: domain };
      const domainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( domainInfo.result.domain.broker ).toEqual( broker );
   } );


   it( `Should register an account with a broker.`, async () => {
      const domain = "iov";
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const unsigned = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const broker = "star1aj9qqrftdqussgpnq6lqj08gwy6ysppf53c8e9"; // w3

      unsigned.value.msg[0].value.broker = broker;

      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.broker ).toEqual( broker );
   } );


   it( `Should register an account and transfer it without deleting resources, etc.`, async () => {
      const domain = "iov";
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const broker = "star1aj9qqrftdqussgpnq6lqj08gwy6ysppf53c8e9"; // w3
      const certificate = JSON.stringify( { my: "certificate", as: "base64" } );
      const base64 = Base64.encode( certificate );
      const metadata = "Why the uri suffix?";
      const resources = [
         {
            "uri": "cosmos:iov-mainnet-2",
            "resource": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk"
         }
      ];
      const fileResources = writeTmpJson( resources );
      const registerAccount = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const replaceResources = iovnscli( [ "tx", "starname", "replace-resources", "--domain", domain, "--name", name, "--src", fileResources, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const setMetadata = iovnscli( [ "tx", "starname", "set-account-metadata", "--domain", domain, "--name", name, "--metadata", metadata, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const addCerts = iovnscli( [ "tx", "starname", "add-certs", "--domain", domain, "--name", name, "--cert", base64, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const unsigned = JSON.parse( JSON.stringify( registerAccount ) );

      unsigned.value.msg[0].value.broker = broker;
      unsigned.value.msg.push( replaceResources.value.msg[0] );
      unsigned.value.msg.push( setMetadata.value.msg[0] );
      unsigned.value.msg.push( addCerts.value.msg[0] );
      unsigned.value.fee.amount[0].amount = "100000000";
      unsigned.value.fee.gas = "600000";

      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.metadata_uri ).toEqual( metadata );
      expect( resolved.result.account.certificates[0] ).toEqual( base64 );
      expect( Base64.decode( resolved.result.account.certificates[0] ) ).toEqual( certificate );
      compareObjects( resources, resolved.result.account.resources );

      const recipient = w1;
      const transferAccount = iovnscli( [ "tx", "starname", "transfer-account", "--domain", domain, "--name", name, "--new-owner", recipient, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const transferred = await signAndPost( transferAccount );
      const newResolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( transferred.ok ).toEqual( true );
      expect( newResolved.result.account.domain ).toEqual( domain );
      expect( newResolved.result.account.name ).toEqual( name );
      expect( newResolved.result.account.owner ).toEqual( recipient );
      expect( newResolved.result.account.metadata_uri ).toEqual( metadata );
      expect( newResolved.result.account.certificates[0] ).toEqual( base64 );
      expect( Base64.decode( newResolved.result.account.certificates[0] ) ).toEqual( certificate );
      compareObjects( resources, newResolved.result.account.resources );
   } );


   it( `Should register an account and transfer it with deleted resources, etc.`, async () => {
      const domain = "iov";
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const broker = "star1aj9qqrftdqussgpnq6lqj08gwy6ysppf53c8e9"; // w3
      const certificate = JSON.stringify( { my: "certificate", as: "base64" } );
      const base64 = Base64.encode( certificate );
      const metadata = "Why the uri suffix?";
      const resources = [
         {
            "uri": "cosmos:iov-mainnet-2",
            "resource": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk"
         }
      ];
      const fileResources = writeTmpJson( resources );
      const registerAccount = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const replaceResources = iovnscli( [ "tx", "starname", "replace-resources", "--domain", domain, "--name", name, "--src", fileResources, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const setMetadata = iovnscli( [ "tx", "starname", "set-account-metadata", "--domain", domain, "--name", name, "--metadata", metadata, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const addCerts = iovnscli( [ "tx", "starname", "add-certs", "--domain", domain, "--name", name, "--cert", base64, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const unsigned = JSON.parse( JSON.stringify( registerAccount ) );

      unsigned.value.msg[0].value.broker = broker;
      unsigned.value.msg.push( replaceResources.value.msg[0] );
      unsigned.value.msg.push( setMetadata.value.msg[0] );
      unsigned.value.msg.push( addCerts.value.msg[0] );
      unsigned.value.fee.amount[0].amount = "100000000";
      unsigned.value.fee.gas = "600000";

      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.metadata_uri ).toEqual( metadata );
      expect( resolved.result.account.certificates[0] ).toEqual( base64 );
      expect( Base64.decode( resolved.result.account.certificates[0] ) ).toEqual( certificate );
      compareObjects( resources, resolved.result.account.resources );

      const recipient = w1;
      const transferAccount = iovnscli( [ "tx", "starname", "transfer-account", "--reset", "true", "--domain", domain, "--name", name, "--new-owner", recipient, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const transferred = await signAndPost( transferAccount );
      const newResolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( transferred.ok ).toEqual( true );
      expect( newResolved.result.account.domain ).toEqual( domain );
      expect( newResolved.result.account.name ).toEqual( name );
      expect( newResolved.result.account.owner ).toEqual( recipient );
      expect( newResolved.result.account.certificates ).toBeNull();
      expect( newResolved.result.account.metadata_uri ).toEqual( "" );
      expect( newResolved.result.account.resources ).toBeNull();
   } );


   it( `Should renew a domain.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const unsigned = iovnscli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const posted = await signAndPost( unsigned );
      const body = { name: domain };
      const domainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( domainInfo ).toBeTruthy();

      const configuration = await fetchObject( `${urlRest}/configuration/query/configuration`, { method: "POST" } );
      const renew = iovnscli( [ "tx", "starname", "renew-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const renewed = await signAndPost( renew );
      const newDomainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( renewed.ok ).toEqual( true );
      expect( newDomainInfo.result.domain.valid_until ).toBeGreaterThanOrEqual( domainInfo.result.domain.valid_until + configuration.result.configuration.domain_renew_period / 1e9 );
   } );


   it( `Should renew an account.`, async () => {
      const domain = "iov";
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const unsigned = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved ).toBeTruthy();

      const configuration = await fetchObject( `${urlRest}/configuration/query/configuration`, { method: "POST" } );
      const renew = iovnscli( [ "tx", "starname", "renew-account", "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const renewed = await signAndPost( renew );
      const newResolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( renewed.ok ).toEqual( true );
      expect( newResolved.result.account.valid_until ).toBeGreaterThanOrEqual( resolved.result.account.valid_until + configuration.result.configuration.account_renew_period / 1e9 );
   } );


   it( `Should register a domain with a fee payer.`, async () => { // https://github.com/iov-one/iovns/pull/195#issue-433044931
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const recipient = w1;
      const payer = signer;
      const unsigned = iovnscli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", recipient, "--fee-payer", payer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );

      // give the payer credit for brokering the registration
      unsigned.value.msg[0].value.broker = payer;

      const signedRecipient = await signTx( unsigned, recipient );
      const signedPayer = await signTx( signedRecipient, payer );

      // payer must be first signature
      signedPayer.value.signatures = [ signedPayer.value.signatures[1], signedPayer.value.signatures[0] ];

      const balance0 = iovnscli( [ "query", "account", recipient ] );
      const balance0Payer = iovnscli( [ "query", "account", payer ] );
      const posted = await postTx( signedPayer );
      const balance = iovnscli( [ "query", "account", recipient ] );
      const balancePayer = iovnscli( [ "query", "account", payer ] );
      const body = { name: domain };
      const domainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( domainInfo.result.domain.name ).toEqual( domain );
      expect( domainInfo.result.domain.admin ).toEqual( recipient );
      expect( domainInfo.result.domain.broker ).toEqual( payer );
      expect( +balance.value.coins[0].amount ).toEqual( +balance0.value.coins[0].amount );
      expect( +balancePayer.value.coins[0].amount ).toBeLessThan( +balance0Payer.value.coins[0].amount );
   } );


   it( `Should register an account with a fee payer.`, async () => { // https://github.com/iov-one/iovns/pull/195#issue-433044931
      const domain = "iov";
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const recipient = w1;
      const payer = signer;
      const unsigned = iovnscli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", recipient, "--fee-payer", payer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );

      // give the payer credit for brokering the registration
      unsigned.value.msg[0].value.broker = payer;

      const signedRecipient = await signTx( unsigned, recipient );
      const signedPayer = await signTx( signedRecipient, payer );

      // payer must be first signature
      signedPayer.value.signatures = [ signedPayer.value.signatures[1], signedPayer.value.signatures[0] ];

      const balance0 = iovnscli( [ "query", "account", recipient ] );
      const balance0Payer = iovnscli( [ "query", "account", payer ] );
      const posted = await postTx( signedPayer );
      const balance = iovnscli( [ "query", "account", recipient ] );
      const balancePayer = iovnscli( [ "query", "account", payer ] );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.owner ).toEqual( recipient );
      expect( resolved.result.account.broker ).toEqual( payer );
      expect( +balance.value.coins[0].amount ).toEqual( +balance0.value.coins[0].amount );
      expect( +balancePayer.value.coins[0].amount ).toBeLessThan( +balance0Payer.value.coins[0].amount );
   } );


   it( `Should do a multisig send.`, async () => { // https://github.com/iov-one/iovns/blob/master/docs/cli/MULTISIG.md
      const amount = 1000000;
      const signed = msig1SignTx( [ "tx", "send", msig1, w1, `${amount}uvoi`, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );

      const balance0 = iovnscli( [ "query", "account", w1 ] );
      const balance0Payer = iovnscli( [ "query", "account", msig1 ] );
      const posted = await postTx( signed );
      const balance = iovnscli( [ "query", "account", w1 ] );
      const balancePayer = iovnscli( [ "query", "account", msig1 ] );

      expect( posted.ok ).toEqual( true );
      expect( +balance.value.coins[0].amount ).toEqual( amount + +balance0.value.coins[0].amount );
      expect( +balancePayer.value.coins[0].amount ).toBeLessThan( +balance0Payer.value.coins[0].amount - amount );
   } );


   // TODO: don't skip after https://github.com/iov-one/iovns/issues/235 is closed
   it.skip( `Should update configuration.`, async () => {
      const config0 = await fetchObject( `${urlRest}/configuration/query/configuration`, { method: "POST" } );
      const config = JSON.parse( JSON.stringify( config0.result.configuration ) );

      config.account_grace_period = `${1 + +config.account_grace_period / 1e9}s`;
      config.account_renew_count_max += 1;
      config.account_renew_period = `${1 + +config.account_renew_period / 1e9}s`;
      config.blockchain_target_max += 1;
      config.certificate_count_max += 1;
      config.certificate_size_max = 1 + +config.certificate_size_max;
      config.domain_grace_period = `${1 + +config.domain_grace_period / 1e9}s`;
      config.domain_renew_count_max += 1;
      config.domain_renew_period = `${1 + +config.domain_renew_period / 1e9}s`;
      config.metadata_size_max = 1 + +config.metadata_size_max;
      config.valid_account_name = "^[-_.a-z0-9]{1,63}$";
      config.valid_blockchain_address = "^[a-z0-9A-Z]+$";
      config.valid_blockchain_id = "[-a-z0-9A-Z:]+$";
      config.valid_domain_name = "^[-_a-z0-9]{4,15}$";

      const argsConfig = [ ...txUpdateConfigArgs( config, msig1 ), "--memo", memo() ];
      const signed = msig1SignTx( argsConfig );
      const posted = await postTx( signed );
      const updated = await fetchObject( `${urlRest}/configuration/query/configuration`, { method: "POST" } );

      expect( posted.ok ).toEqual( true );
      compareObjects( config, updated.result.configuration );

      // restore original config
      const argsConfig0 = [ ...txUpdateConfigArgs( config0.result.configuration, msig1 ), "--memo", memo() ];
      const restore = msig1SignTx( argsConfig0 );
      const posted0 = await postTx( restore );
      const restored = await fetchObject( `${urlRest}/configuration/query/configuration`, { method: "POST" } );

      expect( posted0.ok ).toEqual( true );
      compareObjects( config0, restored.result.configuration );
   } );


   it( `Should register an account owned by a multisig account.`, async () => {
      const domain = "iov";
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const signed = msig1SignTx( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", msig1, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const posted = await postTx( signed );
      const starname = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starname ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
   } );
} );
