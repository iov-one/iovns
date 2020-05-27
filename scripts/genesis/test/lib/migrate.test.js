import { burnTokens, consolidateEscrows, fixChainIds, labelAccounts, labelMultisigs, migrate } from "../../lib/migrate";
import { chainIds, source2multisig } from "../../lib/constants";

"use strict";


describe( "Tests ../../lib/migrate.js.", () => {
   const dumped = {
      "cash": [
         {
            "address": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "coins": [ { "ticker": "IOV", "whole": 35384615 } ]
         },
         {
            "address": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un",
            "coins": [ { "fractional": 510000000, "ticker": "IOV", "whole": 416 } ]
         },
         {
            "address": "iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n",
            "coins": [ { "ticker": "IOV", "whole": 37 } ]
         },
         {
            "address": "iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0",
            "coins": [ { "ticker": "IOV", "whole": 3570582 } ]
         },
         {
            "address": "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph",
            "coins": [ { "ticker": "IOV", "whole": 1628971 }
            ]
         },
         {
            "address": "iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98",
            "coins": [ { "ticker": "IOV", "whole": 3234710 } ]
         },
         {
            "address": "iov170qvwm0tscn5mza3vmaerkzqllvwc3kykkt7kj",
            "coins": [ { "ticker": "IOV", "whole": 2347987 } ]
         },
         {
            "address": "iov105465l8l3yn06a56h7tqwwvnqq22e8j4nvgf02",
            "coins": [ { "ticker": "IOV", "whole": 269559 } ]
         },
         {
            "address": "iov17gdpegksje9dlh8h0g6ehgk6d4anz9pkfskunr",
            "coins": [ { "ticker": "IOV", "whole": 269559 } ]
         },
         {
            "address": "iov1ppxx0vwx42p47p4pkztzl4d57zh2ctnwsz4fdu",
            "coins": [ { "ticker": "IOV", "whole": 269559 } ]
         },
         {
            "address": "iov1k4dpknrrf4dfm07avau0mmjkrsm6pu863d30us",
            "coins": [ { "ticker": "IOV", "whole": 89853 } ]
         },
         {
            "address": "iov1dfurgye70k7f2gxptztfym697g5t832pp9m94g",
            "coins": [ { "ticker": "IOV", "whole": 89853 } ]
         },
         {
            "address": "iov1497txu54lnwujzl8xhc59y6cmuw82d68udn4l3",
            "coins": [ { "ticker": "IOV", "whole": 89853 } ]
         },
      ],
      "escrow": [
         {
            "address": "iov170qvwm0tscn5mza3vmaerkzqllvwc3kykkt7kj",
            "amount": [ { "ticker": "IOV", "whole": 2347987 } ],
            "arbiter": "iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n",
            "destination": "iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n",
            "source": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "timeout": 1888963200
         },
         {
            "address": "iov105465l8l3yn06a56h7tqwwvnqq22e8j4nvgf02",
            "amount": [ { "ticker": "IOV", "whole": 269559 } ],
            "arbiter": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "destination": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "source": "iov1w2suyhrfcrv5h4wmq3rk3v4x95cxtu0a03gy6x",
            "timeout": 1575979200
         },
         {
            "address": "iov17gdpegksje9dlh8h0g6ehgk6d4anz9pkfskunr",
            "amount": [ { "ticker": "IOV", "whole": 269559 } ],
            "arbiter": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "destination": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "source": "iov1w2suyhrfcrv5h4wmq3rk3v4x95cxtu0a03gy6x",
            "timeout": 1578657600
         },
         {
            "address": "iov1ppxx0vwx42p47p4pkztzl4d57zh2ctnwsz4fdu",
            "amount": [ { "ticker": "IOV", "whole": 269559 } ],
            "arbiter": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "destination": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "source": "iov1w2suyhrfcrv5h4wmq3rk3v4x95cxtu0a03gy6x",
            "timeout": 1581336000
         },
         {
            "address": "iov1k4dpknrrf4dfm07avau0mmjkrsm6pu863d30us",
            "amount": [ { "ticker": "IOV", "whole": 89853 } ],
            "arbiter": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "destination": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "source": "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph",
            "timeout": 1591790400
         },
         {
            "address": "iov1dfurgye70k7f2gxptztfym697g5t832pp9m94g",
            "amount": [ { "ticker": "IOV", "whole": 89853 } ],
            "arbiter": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "destination": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "source": "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph",
            "timeout": 1594382400
         },
         {
            "address": "iov1497txu54lnwujzl8xhc59y6cmuw82d68udn4l3",
            "amount": [ { "ticker": "IOV", "whole": 89853 } ],
            "arbiter": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "destination": "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
            "source": "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph",
            "timeout": 1597060800
         },
      ],
      "username": [
         {
            "Owner": "iov16a42lf29n2h2eurxryspue9fz2d2wnlgpyjv8d",
            "Targets": [
               { "address": "iov16a42lf29n2h2eurxryspue9fz2d2wnlgpyjv8d",    "blockchain_id": "iov-mainnet" },
               { "address": "0x52dBf411B22Af67f33425bf3fbb6B8CF8FB302cd",    "blockchain_id": "ethereum-eip155-1" },
               { "address": "cosmos15dafemy5pkaru4kf23s3e6mnugfv6et9kg2uz7", "blockchain_id": "cosmos-cosmoshub-3" },
               { "address": "star1ayxmc4vqshd9j94hj67r55ppg5hsrhqlmy4dvd",   "blockchain_id": "starname-migration" }
            ],
            "Username": "alpha*iov"
         },
         {
            "Owner": "iov10v69k57z2v0pr3yvtr60pp8g2jx8tdd7f55sv6",
            "Targets": [
               { "address": "star1gfdmksf725qpdgl06e98ks4usg9nmkcwc5qzcg", "blockchain_id": "starname-migration" },
               { "address": "0xa223f22664Ee8bfB41FAD93C388826E7aF24060c",  "blockchain_id": "ethereum-eip155-1" },
               { "address": "iov10v69k57z2v0pr3yvtr60pp8g2jx8tdd7f55sv6",  "blockchain_id": "iov-mainnet" },
               { "address": "4341330819731245941L",                        "blockchain_id": "lisk-ed14889723" }
            ],
            "Username": "blini44*iov"
         },
         {
            "Owner": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un",
            "Targets": [
               { "address": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk", "blockchain_id": "starname-migration" },
               { "address": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un",  "blockchain_id": "iov-mainnet" }
            ],
            "Username": "dave*iov"
         },
         {
            "Owner": "iov1ejk0g6p2xk90lamuvtd3r0kf6jcva09hf4xy74",
            "Targets": [
               { "address": "star1yxxmpqca3l7xzhy4783vkpfx843x4zk749h8fs", "blockchain_id": "starname-migration" },
               { "address": "iov1ejk0g6p2xk90lamuvtd3r0kf6jcva09hf4xy74", "blockchain_id": "iov-mainnet" }
            ],
            "Username": "fish_and_chips*iov"
         },
         {
            "Owner": "iov1tlxqvugk9u5u973a6ee6dq4zsgsv6c5ecr0rvn",
            "Targets": [
               { "address": "0x00C60938d954FEC83E70eE98243B24F7E6EabaC8",  "blockchain_id": "ethereum-eip155-1" },
               { "address": "iov1tlxqvugk9u5u973a6ee6dq4zsgsv6c5ecr0rvn",  "blockchain_id": "iov-mainnet" },
               { "address": "13483265462465913551L",                       "blockchain_id": "lisk-ed14889723" },
               { "address": "star1vmt7wysxug30vfenedfh4ay83y3p75tstagn2y", "blockchain_id": "starname-migration" }
            ],
            "Username": "huobi*iov"
         },
         {
            "Owner": "iov16a42lf29n2h2eurxryspue9fz2d2wnlgpyjv8d",
            "Targets": [
               { "address": "0x52dBf411B22Af67f33425bf3fbb6B8CF8FB302cd",  "blockchain_id": "ethereum-eip155-1" },
               { "address": "iov16a42lf29n2h2eurxryspue9fz2d2wnlgpyjv8d",  "blockchain_id": "iov-mainnet" },
               { "address": "star1ayxmc4vqshd9j94hj67r55ppg5hsrhqlmy4dvd", "blockchain_id": "starname-migration" }
            ],
            "Username": "letsdoit*iov"
         },
         {
            "Owner": "iov1eh6yeyel3zsc8vqnh79fqjtfkcxmj5d8nt49gq",
            "Targets": [
               { "address": "0x2cE327b4EB237313F37a72195d64Cb80F7aeAa15",  "blockchain_id": "ethereum-eip155-1" },
               { "address": "iov1eh6yeyel3zsc8vqnh79fqjtfkcxmj5d8nt49gq",  "blockchain_id": "iov-mainnet" },
               { "address": "16192453558792957658L", "blockchain_id":      "lisk-ed14889723" },
               { "address": "star1y86zdqsegxm7uj9qf7l400y29nc6x9ypqxpdcg", "blockchain_id": "starname-migration" }
            ],
            "Username": "nash.io*iov"
         },
         {
            "Owner": "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph",
            "Targets": [
               { "address": "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph", "blockchain_id": "iov-mainnet" }
            ],
            "Username": "kadima*iov"
         },
         {
            "Owner": "iov1ua6tdcyw8jddn5660qcx2ndhjp4skqk4dkurrl",
            "Targets": [
               { "address": "iov1ua6tdcyw8jddn5660qcx2ndhjp4skqk4dkurrl",  "blockchain_id": "alpe-net" }
            ],
            "Username": "alex*iov"
         },
      ],
   };
   const genesis = {
      chain_id: __filename,
      genesis_time: new Date( "2020-04-15T10:00:00Z" ).toISOString(),
      accounts: [],
      app_hash: "",
      app_state: {},
      auth: {},
      consensus_params: {},
      crisis: {},
      genutil: {},
      gov: {},
   };
   const flammable = [ "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u" ];
   const multisigs = {
      iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n: {
         "//name": "reward fund",
         address: "cond:gov/rule/0000000000000002",
         star1: "star1rewards",
      },
      iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0: {
         "//name": "IOV SAS employee bonus pool/colloboration appropriation pool",
         address: "cond:multisig/usage/0000000000000002",
         star1: "star1bonuses",
      },
   };
   const osaka = {
      app_hash: "",
      app_state: {
         cash: [
            {
               address: "bech32:iov15xzzgu5jkltm24hg9r2ykm6gm2d09tzrcqgrr9",
               "//id": 1957,
               coins: [ "126455 IOV" ]
            },
            {
               address: "bech32:iov1tc4vr2756lcme6hqq2xgdn4dycny32cdev9a8g",
               "//id": 1970,
               coins: [ "62500 IOV" ]
            },
            {
               address: "bech32:iov1s3e835efuht3qulf3lrv02dypn036lnpedq275",
               "//id": 1976,
               coins: [ "626325 IOV" ]
            },
            {
               address: "bech32:iov13adwzjxhqhd79l3y5v58vjugtfszv57tthmv0q",
               "//id": 1978,
               coins: [ "470651 IOV" ]
            },
            {
               address: "bech32:iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph",
               "//id": 2096,
               coins: [ "1000000 IOV" ]
            },
            {
               address: "bech32:iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98",
               "//id": 2101,
               coins: [ "3234708 IOV" ]
            },
            {
               address: "bech32:iov14favyxdrkkdk39kl4qsexc99vgscw8dw22g5d3",
               "//id": 2243,
               coins: [ "555555 IOV" ]
            },
            {
               address: "bech32:iov1wvxg0qw8pg9vnkja9mvvdzk74g6lrms7uh7tn8",
               "//id": 2244,
               coins: [ "107824 IOV" ]
            },
            {
               address: "bech32:iov1jukhxtnh58mmag5y65d8xj2e36wq6083w0t69e",
               "//id": 2246,
               coins: [ "77777 IOV" ]
            },
         ],
      },
      chain_id: "iov-mainnet",
      consensus_params: {},
      genesis_time: new Date( "2019-10-10T10:00:00Z" ).toISOString(),
      validators: [],
   };

   it( `Should burn tokens.`, async () => {
      const copied = JSON.parse( JSON.stringify( dumped ) );

      burnTokens( copied, flammable );

      flammable.forEach( iov1 => {
         const index = copied.cash.findIndex( wallet => wallet.address == iov1 );

         expect( index ).toEqual( -1 );
      } );

      expect( copied.cash.length ).toEqual( dumped.cash.length - flammable.length );
   } );

   it( `Should label multisig accounts.`, async () => {
      labelMultisigs( dumped, multisigs );

      Object.keys( multisigs ).forEach( iov1 => {
         const multisig = dumped.cash.find( wallet => wallet.address == iov1 );

         expect( multisig["//id"] ).toEqual( multisigs[iov1]["//name"] );
      } );
   } );

   it( `Should label accounts.`, async () => {
      labelAccounts( dumped, osaka );

      const id2096 = dumped.cash.find( account => account.address == "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph" );
      const id2101 = dumped.cash.find( account => account.address == "iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98" );

      expect( id2096["//id"] ).toEqual( 2096 );
      expect( id2101["//id"] ).toEqual( 2101 );
      expect( id2096["//iov1"] ).toEqual( "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph" );
      expect( id2101["//iov1"] ).toEqual( "iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98" );
   } );

   it( `Should fix chain ids.`, async () => {
      const invalids = Object.keys( chainIds );
      const valids = Object.values( chainIds );
      const unknowns = dumped.username.reduce( ( accumulator, username ) => {
         username.Targets.forEach( target => {
            const id = target.blockchain_id;

            if ( !valids.includes( id ) && !invalids.includes( id ) ) accumulator.push( id );
         } );

         return accumulator;
      }, [] );

      fixChainIds( dumped, chainIds );

      dumped.username.forEach( username => {
         username.Targets.forEach( target => {
            const id = target.blockchain_id;

            if ( !unknowns.includes( id ) ) {
               expect( valids.includes( id ) ).toEqual( true );
               expect( invalids.includes( id ) ).toEqual( false );
            }
         } );
      } );
   } );

   it( `Should consolidate escrows.`, async () => {
      const dumpedCopy = JSON.parse( JSON.stringify( dumped ) );
      const genesisCopy = JSON.parse( JSON.stringify( genesis ) );

      consolidateEscrows( dumpedCopy, source2multisig, genesisCopy );

      expect( genesisCopy.accounts.length ).toEqual( 3 );

      const guaranteed = genesisCopy.accounts.find( account => account.value.address == source2multisig.iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u.multisig );
      const isabella   = genesisCopy.accounts.find( account => account.value.address == source2multisig.iov1w2suyhrfcrv5h4wmq3rk3v4x95cxtu0a03gy6x.multisig );
      const kadima     = genesisCopy.accounts.find( account => account.value.address == source2multisig.iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph.multisig );

      expect( guaranteed ).toBeTruthy();
      expect( isabella ).toBeTruthy();
      expect( kadima ).toBeTruthy();

      expect( guaranteed.value.coins[0].amount ).toEqual( "2347987" );
      expect( isabella.value.coins[0].amount ).toEqual( "808677" );
      expect( kadima.value.coins[0].amount ).toEqual( "269559" );
   } );

   it( `Should migrate.`, async () => {
      migrate( { chainIds, dumped, flammable, genesis, multisigs, osaka, source2multisig } );
   } );
} );
