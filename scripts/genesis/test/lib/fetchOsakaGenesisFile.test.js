import fetchOsakaGenesisFile from "../../lib/fetchOsakaGenesisFile";

"use strict";


describe( "Tests ../../lib/fetchOsakaGenesisFile.js.", () => {
   it( `Should get the weave-based mainnet genesis file.`, async () => {
      const osaka = await fetchOsakaGenesisFile();

      expect( osaka.chain_id ).toEqual( "iov-mainnet" );
      expect( osaka.validators[3].address ).toEqual( "058078082E8ED2431EA61E69657BE27F0D7456FA" );
      expect( osaka.app_state.cash[6]["//id"] ).toEqual( 1957 );
      expect( osaka.app_state.multisig[0]["//name"] ).toEqual( "IOV SAS" );
      expect( osaka.app_state.escrow[25].source ).toEqual( "bech32:iov149cn0rauw2773lfdp34njyejg3cfz2d56c0m5t" );
      expect( osaka.app_state.username[0].owner ).toEqual( "bech32:iov16yrd6qhyd4kxpcklu344ly4f2fay0s9rpz46fm" );
      expect( osaka.app_state.username[0]["//id"] ).toEqual( -26 );
      expect( osaka.app_state.username[0].targets[0].blockchain_id ).toEqual( "iov-mainnet" );
   } );
} );
