import { burnTokens, consolidateEscrows, convertToCosmosSdk, fixChainIds, labelAccounts, labelMultisigs, mapIovToStar, migrate } from "../../lib/migrate";
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
         {
            "address": "iov1m7qjqjuv4ynhzu40xranun4u0r47d4waxc4wh9",
            "coins": [ { "fractional": 500000000, "ticker": "IOV", "whole": 26 } ]
         },
         {
            "address": "iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc",
            "coins": [ { "fractional": 123000, "ticker": "IOV", "whole": 1 } ]
         },
         {
            "address": "iov1q8zjkzk3f2yzfrkh9wswlf9qtmdgel84nnlgs9",
            "coins": [ { "fractional": 657145000, "ticker": "IOV", "whole": 8920 } ]
         },
         {
            "address": "iov1q40tvnph5xy7cjyj3tmqzghukeheykudq246d6",
            "coins": [ { "ticker": "IOV", "whole": 22171 } ]
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
         {
            "Owner": "iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98",
            "Targets": [
               { "address": "iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98", "blockchain_id": "iov-mainnet" }
            ],
            "Username": "confio*iov"
         },
         {
            "Owner": "iov1m7qjqjuv4ynhzu40xranun4u0r47d4waxc4wh9",
            "Targets": [
               { "address": "0x5e415520beb66aa39e00d43cae889f2c5cba7017", "blockchain_id": "ethereum-eip155-1" }
            ],
            "Username": "corentin*iov"
         },
      ],
   };
   const genesis = {
      chain_id: __filename,
      genesis_time: new Date( "2020-04-15T10:00:00Z" ).toISOString(),
      app_hash: "",
      app_state: {
         auth: {
            accounts: [],
         },
         domain: {
            domains: [
               {
                  name: "iov",
                  "//note": "msig1",
                  admin: "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg",
                  valid_until: "1689380911",
                  has_super_user: false,
                  account_renew: "3000",
                  broker: null,
               }
            ],
            accounts: [],
         },
      },
      consensus_params: {},
      crisis: {},
      genutil: {},
      gov: {},
   };
   const flammable = [ "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u" ];
   const indicatives = [
      { "hash": "e0d65bc5377e0806de18f76e07c3234632fad570a799c1063df1f69809bf4337", "block_height": 65609, "message": { "path": "cash/send", "details": { "memo": "star1cnywewxct2p4d5j2fapgkse6yxgh7ecnj4uwpu", "amount": { "whole": 1, "ticker": "IOV" }, "source": "iov1yhk8qqp3wsdg7tefd8u457n9zqsny4nqzp6960", "destination": "iov10v69k57z2v0pr3yvtr60pp8g2jx8tdd7f55sv6" } } },
      { "hash": "20894f0429901e402bb0520d117da9b64dacce2a97b647c66645bf6436af17d7", "block_height": 67029, "message": { "path": "cash/send", "details": { "memo": "star19m9ufykj5ur67l822fpxvz49p535wp3j0m5v3h", "amount": { "ticker": "IOV", "fractional": 1 }, "source": "iov1a9duw7yyxdfh8mrjxmuc0slu8a48muvxkcxvg8", "destination": "iov10v69k57z2v0pr3yvtr60pp8g2jx8tdd7f55sv6" } } }
   ];
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
      iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc: {
         "//name": "Custodian of missing star1 accounts",
         address: "cond:multisig/usage/0000000000000006",
         star1: "star1 custodian", // TODO
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
   const premiums = {
      iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un: [ "in3s", "huth", "tachyon", "sentient" ],
      iov1tlxqvugk9u5u973a6ee6dq4zsgsv6c5ecr0rvn: [ "hell", "hash", "hold" ],
      iov1y63fp8pncpuke7mrc2huqefud59t3munnh0k32: [ "multiverse" ],
      iov1ylw3cnluf3zayfths0ezgjp5cwf6ddvsvwa7l4: [ "lovely" ],
      iov1zr9epgrzysr6zc5s8ucd3qlxkhgj9fwj2a2mkx: [ "gianna", "nodeateam", "tyler", "michael" ],
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
         expect( multisig["//iov1"] ).toEqual( iov1 );
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
      const iov2escrow = consolidateEscrows( dumpedCopy, source2multisig );
      const escrows = Object.values( iov2escrow );

      expect( escrows.length ).toEqual( 3 );

      const guaranteed = escrows.find( account => account.value.address == source2multisig.iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u.star1 );
      const isabella   = escrows.find( account => account.value.address == source2multisig.iov1w2suyhrfcrv5h4wmq3rk3v4x95cxtu0a03gy6x.star1 );
      const kadima     = escrows.find( account => account.value.address == source2multisig.iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph.star1 );

      expect( guaranteed ).toBeTruthy();
      expect( isabella ).toBeTruthy();
      expect( kadima ).toBeTruthy();

      expect( guaranteed.value.coins[0].amount ).toEqual( "2347987" );
      expect( isabella.value.coins[0].amount ).toEqual( "808677" );
      expect( kadima.value.coins[0].amount ).toEqual( "269559" );
   } );

   it( `Should map iov1 addresses to star1 addresses.`, async () => {
      const iov2star = mapIovToStar( dumped, multisigs, indicatives );
      const reMemo = /(star1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{38})/;

      expect( iov2star.iov1ua6tdcyw8jddn5660qcx2ndhjp4skqk4dkurrl ).toEqual( false ); // alex
      expect( iov2star.iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98 ).toEqual( false ); // ethan
      expect( iov2star.iov1m7qjqjuv4ynhzu40xranun4u0r47d4waxc4wh9 ).toEqual( false );
      expect( iov2star.iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un ).toEqual( "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
      expect( iov2star.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n ).toEqual( multisigs.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n.star1 );
      expect( iov2star[indicatives[0].message.details.source] ).toEqual( indicatives[0].message.details.memo.match( reMemo )[0] );
      expect( iov2star[indicatives[1].message.details.source] ).toEqual( indicatives[1].message.details.memo.match( reMemo )[1] );
   } );

   it( `Should convert genesis objects from weave to cosmos-sdk.`, async () => {
      const dumpedCopy = JSON.parse( JSON.stringify( dumped ) );

      burnTokens( dumpedCopy, flammable );
      labelAccounts( dumpedCopy, osaka );
      labelMultisigs( dumpedCopy, multisigs );
      fixChainIds( dumpedCopy, chainIds );

      const iov2star = mapIovToStar( dumpedCopy, multisigs, indicatives );
      const { accounts, starnames, domains } = convertToCosmosSdk( dumpedCopy, iov2star, multisigs, premiums );
      const custodian = accounts.find( account => account["//iov1"] == "iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc" );
      const rewards = accounts.find( account => account["//iov1"] == "iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n" );
      const bonus = accounts.find( account => account["//iov1"] == "iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0" );
      const dave = accounts.find( account => account["//iov1"] == "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" );

      expect( custodian.value.address ).toEqual( multisigs.iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc.star1 );
      expect( custodian.value.coins[0].amount ).toEqual( "8321023.157268" );
      expect( custodian["//no star1 iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98"][0] ).toEqual( 3234710 );
      expect( custodian["//no star1 iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98"][1] ).toEqual( "confio*iov" );
      expect( custodian["//no star1 iov1m7qjqjuv4ynhzu40xranun4u0r47d4waxc4wh9"][0] ).toEqual( 26.5 );
      expect( custodian["//no star1 iov1m7qjqjuv4ynhzu40xranun4u0r47d4waxc4wh9"][1] ).toEqual( "corentin*iov" );
      expect( custodian["//no star1 iov1q8zjkzk3f2yzfrkh9wswlf9qtmdgel84nnlgs9"] ).toEqual( 8920.657145 );
      expect( custodian["//no star1 iov1q40tvnph5xy7cjyj3tmqzghukeheykudq246d6"] ).toEqual( 22171 );
      expect( custodian["//no star1 iov1ua6tdcyw8jddn5660qcx2ndhjp4skqk4dkurrl"] ).toEqual( "alex*iov" );
      expect( custodian["//no star1 iov1ylw3cnluf3zayfths0ezgjp5cwf6ddvsvwa7l4"] ).toEqual( "lovely" );
      expect( custodian["//no star1 iov1y63fp8pncpuke7mrc2huqefud59t3munnh0k32"] ).toEqual( "multiverse" );
      expect( custodian["//no star1 iov1zr9epgrzysr6zc5s8ucd3qlxkhgj9fwj2a2mkx"][0] ).toEqual( "gianna" );
      expect( custodian["//no star1 iov1zr9epgrzysr6zc5s8ucd3qlxkhgj9fwj2a2mkx"][1] ).toEqual( "nodeateam" );
      expect( custodian["//no star1 iov1zr9epgrzysr6zc5s8ucd3qlxkhgj9fwj2a2mkx"][2] ).toEqual( "tyler" );
      expect( custodian["//no star1 iov1zr9epgrzysr6zc5s8ucd3qlxkhgj9fwj2a2mkx"][3] ).toEqual( "michael" );

      expect( rewards.value.address ).toEqual( multisigs.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n.star1 );
      expect( rewards.value.coins[0].amount ).toEqual( "37" );

      expect( bonus.value.address ).toEqual( multisigs.iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0.star1 );
      expect( bonus.value.coins[0].amount ).toEqual( "3570582" );

      expect( dave.value.address ).toEqual( "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
      expect( dave.value.coins[0].amount ).toEqual( "416.51" );

      const alphaiov = starnames.find( starname => starname["//iov1"] == "iov16a42lf29n2h2eurxryspue9fz2d2wnlgpyjv8d" );
      const daveiov = starnames.find( starname => starname["//iov1"] == "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" );
      const huobiiov = starnames.find( starname => starname["//iov1"] == "iov1tlxqvugk9u5u973a6ee6dq4zsgsv6c5ecr0rvn" );

      expect( alphaiov.address ).toEqual( "star1ayxmc4vqshd9j94hj67r55ppg5hsrhqlmy4dvd" );
      expect( alphaiov.starname ).toEqual( "alpha*iov" );

      expect( daveiov.address ).toEqual( "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
      expect( daveiov.starname ).toEqual( "dave*iov" );

      expect( huobiiov.address ).toEqual( "star1vmt7wysxug30vfenedfh4ay83y3p75tstagn2y" );
      expect( huobiiov.starname ).toEqual( "huobi*iov" );

      const alexiov = starnames.find( starname => starname["//iov1"] == "iov1ua6tdcyw8jddn5660qcx2ndhjp4skqk4dkurrl" );
      const confioiov = starnames.find( starname => starname["//iov1"] == "iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98" );
      const kadimaiov = starnames.find( starname => starname["//iov1"] == "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph" );

      expect( alexiov.address ).toEqual( custodian.value.address );
      expect( alexiov.starname ).toEqual( "alex*iov" );

      expect( confioiov.address ).toEqual( custodian.value.address );
      expect( confioiov.starname ).toEqual( "confio*iov" );

      expect( kadimaiov.address ).toEqual( custodian.value.address );
      expect( kadimaiov.starname ).toEqual( "kadima*iov" );

      const gianna = domains.find( domain => domain.name == "gianna" );
      const lovely = domains.find( domain => domain.name == "lovely" );
      const michael = domains.find( domain => domain.name == "michael" );
      const multiverse = domains.find( domain => domain.name == "multiverse" );

      expect( gianna.admin ).toEqual( custodian.value.address );
      expect( lovely.admin ).toEqual( custodian.value.address );
      expect( michael.admin ).toEqual( custodian.value.address );
      expect( multiverse.admin ).toEqual( custodian.value.address );

      const hash = domains.find( domain => domain.name == "hash" );
      const huth = domains.find( domain => domain.name == "huth" );

      expect( hash.admin ).toEqual( "star1vmt7wysxug30vfenedfh4ay83y3p75tstagn2y" );
      expect( huth.admin ).toEqual( "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
   } );

   it( `Should migrate.`, async () => {
      migrate( { chainIds, dumped, flammable, genesis, indicatives, multisigs, osaka, premiums, source2multisig } );
   } );
} );
