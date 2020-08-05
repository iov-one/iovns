import {
   addGentxs,
   burnTokens,
   consolidateEscrows,
   convertToCosmosSdk,
   fixChainIds,
   fixErrors,
   labelAccounts,
   labelMultisigs,
   mapIovToStar,
   migrate,
   patchGalaxynet,
   patchJestnet,
   patchMainnet,
} from "../../lib/migrate";
import { chainIds, source2multisig } from "../../lib/constants";
import compareObjects from "../compareObjects";
import fs from "fs";
import path from "path";
import stringify from "json-stable-stringify";
import tmp from "tmp";

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
         {
            "address": "iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq",
            "coins": [ { "fractional": 500000000, "ticker": "IOV", "whole": 13015243 } ]
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
         {
            "Owner": "iov1yhk8qqp3wsdg7tefd8u457n9zqsny4nqzp6960",
            "Targets": [
               { "address": "star1qvpth6t72336fjxlej2xv8eu84hrpxdxf5rgzz", "blockchain_id": "starname-migration" },
               { "address": "iov1yhk8qqp3wsdg7tefd8u457n9zqsny4nqzp6960", "blockchain_id": "iov-mainnet" },
               { "address": "16104600299727948959L", "blockchain_id": "lisk-ed14889723" },
               { "address": "0x40698A9DcE4d6a63E766Dd08b83D03c6727DCB1a", "blockchain_id": "ethereum-eip155-1" }
            ],
            "Username": "btc13*iov"
         },
         {
            "Owner": "iov1fpezwaxfnmef8tyyg4t7avz9a2d9gqh3yh8d8n",
            "Targets": [
               { "address": "iov1fpezwaxfnmef8tyyg4t7avz9a2d9gqh3yh8d8n", "blockchain_id": "iov-mainnet" } ],
            "Username": "ledger*iov"
         },
      ],
   };
   const genesis = {
      chain_id: "jestnet",
      genesis_time: new Date( "2020-04-15T10:00:00Z" ).toISOString(),
      app_hash: "",
      app_state: {
         auth: {
            accounts: [],
            params: {
               max_memo_characters: "256",
               sig_verify_cost_ed25519: "590",
               sig_verify_cost_secp256k1: "1000",
               tx_sig_limit: "7",
               tx_size_cost_per_byte: "10"
            },
         },
         bank: {
            send_enabled: true
         },
         configuration: {
            config: {
               account_grace_period: String( 30 * 24 * 60 * 60 ),
               account_renew_count_max: 2,
               account_renew_period: String( 365.25 * 24 * 60 * 60 ),
               resource_target_max: 16,
               certificate_count_max: 16,
               certificate_size_max: "1024",
               configurer: "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78",
               domain_grace_period: String( 30 * 24 * 60 * 60 ),
               domain_renew_count_max: 2,
               domain_renew_period: String( 365.25 * 24 * 60 * 60 ),
               metadata_size_max: "1024",
               valid_account_name: "[-_\\.a-z0-9]{0,64}$",
               valid_resource: "^[a-z0-9A-Z]+$",
               valid_uri: "[-a-z0-9A-Z:]+$",
               valid_domain_name: "^[-_a-z0-9]{4,16}$",
            },
            fees: {
               add_account_certificate: "100",
               del_account_certificate: "0.5",
               fee_coin_denom: "uiov",
               fee_coin_price: "0.000001",
               fee_default: "0.5", // the fee for messages that don't explicitly have a fee
               register_account_closed: "0.5",
               register_account_open: "10",
               register_domain_1: "10000", // domain name with 1 char
               register_domain_2: "5000",
               register_domain_3: "2000",
               register_domain_4: "1000",
               register_domain_5: "500",
               register_domain_default: "250", // domain name with 6 or more chars
               register_open_domain_multiplier: "5.5",
               renew_domain_open: "12345",
               replace_account_resources: "10",
               set_account_metadata: "500",
               transfer_account_closed: "10",
               transfer_account_open: "10",
               transfer_domain_closed: "10",
               transfer_domain_open: "10",
            },
         },
         crisis: {
            constant_fee: {
               denom: "uiov",
               amount: "1000000000"
            }
         },
         distribution: {
            delegator_starting_infos: [],
            delegator_withdraw_infos: [],
            fee_pool: {
               community_pool: []
            },
            outstanding_rewards: [],
            params: {
               base_proposer_reward: "0.050000000000000000",
               bonus_proposer_reward: "0.050000000000000000",
               community_tax: "0.000000000000000000",
               withdraw_addr_enabled: true
            },
            previous_proposer: "",
            validator_accumulated_commissions: [],
            validator_current_rewards: [],
            validator_historical_rewards: [],
            validator_slash_events: []
         },
         evidence: {
            evidence: [],
            params: {
               max_evidence_age: "1814400000000000"
            },
         },
         starname: {
            domains: [
               {
                  name: "iov",
                  "//note": "msig1",
                  admin: "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg",
                  valid_until: String( Math.ceil( Date.now() / 1000 ) + 365.25 * 24 * 60 * 60 ), // 1 year from now
                  type: "open",
                  account_renew: "3000",
                  broker: null,
               }
            ],
            accounts: [],
         },
         genutil: {
         },
         gov: {
            starting_proposal_id: "1",
            deposits: null,
            votes: null,
            proposals: null,
            deposit_params: {
               min_deposit: [
                  {
                     denom: "uiov",
                     amount: "1000000000"
                  }
               ],
               max_deposit_period: "172800000000000"
            },
            voting_params: {
               voting_period: "345600000000000"
            },
            tally_params: {
               quorum: "0.334000000000000000",
               threshold: "0.500000000000000000",
               veto: "0.334000000000000000"
            }
         },
         mint: {
            minter: {
               inflation: "0.000000000000000000",
               annual_provisions: "0.000000000000000000"
            },
            params: {
               blocks_per_year: "105192",
               "//note": "goal_bonded cannot be 0: module=consensus err='division by zero'",
               goal_bonded: "0.000000000000000001",
               inflation_max: "0.0000000000000000",
               inflation_min: "0.0000000000000000",
               inflation_rate_change: "0.000000000000000000",
               mint_denom: "uiov"
            }
         },
         slashing: {
            missed_blocks: {},
            params: {
               downtime_jail_duration: "600000000000",
               max_evidence_age: "1814400000000000",
               min_signed_per_window: "0.500000000000000000",
               signed_blocks_window: "10000",
               slash_fraction_double_sign: "0.050000000000000000",
               slash_fraction_downtime: "0.010000000000000000"
            },
            signing_infos: {},
         },
         staking: {
            params: {
               historical_entries: 0,
               unbonding_time: "259200000000000",
               max_validators: 16,
               max_entries: 7,
               bond_denom: "uiov"
            },
            last_total_power: "0",
            last_validator_powers: null,
            validators: null,
            delegations: null,
            unbonding_delegations: null,
            redelegations: null,
            exported: false
         },
         supply: {
            supply: []
         },
      },
      consensus_params: {
         block: {
            max_bytes: "500000",
            max_gas: "-1",
            time_iota_ms: "1000"
         },
         evidence: {
            max_age_num_blocks: "100000",
            max_age_duration: "172800000000000"
         },
         validator: {
            pub_key_types: [
               "ed25519"
            ]
         }
      },
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
         star1: "star1scfumxscrm53s4dd3rl93py5ja2ypxmxlhs938",
      },
      iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0: {
         "//name": "IOV SAS employee bonus pool/colloboration appropriation pool",
         address: "cond:multisig/usage/0000000000000002",
         star1: "star16tm7scg0c2e04s0exk5rgpmws2wk4xkd84p5md",
      },
      iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc: {
         "//name": "Custodian of missing star1 accounts",
         address: "cond:multisig/usage/0000000000000006",
         star1: "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
      },
      iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq: {
         "//name": "IOV SAS",
         address: "cond:multisig/usage/0000000000000001",
         star1: "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78",
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
      iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un: { star1: "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk", starnames: [ "in3s", "huth", "tachyon", "sentient" ] },
      iov1tlxqvugk9u5u973a6ee6dq4zsgsv6c5ecr0rvn: { star1: "star1vmt7wysxug30vfenedfh4ay83y3p75tstagn2y", starnames: [ "hell", "hash", "hold" ] },
      iov1y63fp8pncpuke7mrc2huqefud59t3munnh0k32: { star1: "", starnames: [ "multiverse" ] },
      iov1ylw3cnluf3zayfths0ezgjp5cwf6ddvsvwa7l4: { star1: "", starnames: [ "lovely" ] },
      iov1zr9epgrzysr6zc5s8ucd3qlxkhgj9fwj2a2mkx: { star1: "star18awsa7fhwtsevta28p3uw8ymtznvpwtzl3ep5f", starnames: [ "gianna", "nodeateam", "tyler", "michael" ] },
      iov12gd6weg7py6vs7ujn22h82422arek8cxzhe85p: { star1: "star1usl4zpltjesrp5rqae3fdjdyj5dyymakmhq6mt", starnames: [ "adrian", "84more" ] },
      zGLlamFypWMPUeHVVsvo4mXFFOE63: { star1: "", starnames: [ "cosmostation", "ibcwallet", "korea", "mintscan", "seoul", "station" ] },
      zHbPpUYyRguRlhAiC30zimM05hGx2: { star1: "", starnames: [ "jim" ] },
   };
   const reserveds = [
      "goldman",
      "socgen",
      "twitter",
      "youtube",
      "ivanontech",
      "goldman_sachs",
      "fcbarca",
      "china",
      "goldman1",
      "socgen1",
      "twitter1",
      "youtube1",
      "ivanontech1",
      "goldman_sachs1",
      "fcbarca1",
      "china1",
   ];

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

      expect( guaranteed.value.coins[0].amount ).toEqual( "2347987000000" );
      expect( isabella.value.coins[0].amount ).toEqual( "808677000000" );
      expect( kadima.value.coins[0].amount ).toEqual( "269559000000" );
   } );

   it( `Should fix human errors.`, async () => {
      const dumpedCopy = JSON.parse( JSON.stringify( dumped ) );
      const indicativesCopy = JSON.parse( JSON.stringify( indicatives ) );
      const previous = [].concat( indicativesCopy );

      fixErrors( dumpedCopy, indicativesCopy );

      expect( indicativesCopy.length ).toEqual( previous.length - 1 );
      expect( indicativesCopy.findIndex( indicative => indicative.message.details.source == "iov1yhk8qqp3wsdg7tefd8u457n9zqsny4nqzp6960" ) ).toEqual( -1 );

      const ledger = dumpedCopy.username.find( username => username.Username == "ledger*iov" );

      expect( ledger.Owner ).toEqual( "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" );
   } );

   it( `Should map iov1 addresses to star1 addresses.`, async () => {
      const dumpedCopy = JSON.parse( JSON.stringify( dumped ) );
      const indicativesCopy = JSON.parse( JSON.stringify( indicatives ) );

      fixErrors( dumpedCopy, indicativesCopy );

      const iov2star = mapIovToStar( dumped, multisigs, indicativesCopy, premiums );
      const reMemo = /(star1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{38})/;

      expect( iov2star.iov1ua6tdcyw8jddn5660qcx2ndhjp4skqk4dkurrl ).toEqual( false ); // alex
      expect( iov2star.iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98 ).toEqual( false ); // ethan
      expect( iov2star.iov1m7qjqjuv4ynhzu40xranun4u0r47d4waxc4wh9 ).toEqual( false );
      expect( iov2star.iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un ).toEqual( "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
      expect( iov2star.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n ).toEqual( multisigs.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n.star1 );
      expect( iov2star[indicativesCopy[0].message.details.source] ).toEqual( indicativesCopy[0].message.details.memo.match( reMemo )[0] );
   } );

   it( `Should convert genesis objects from weave to cosmos-sdk.`, async () => {
      const dumpedCopy = JSON.parse( JSON.stringify( dumped ) );
      const indicativesCopy = JSON.parse( JSON.stringify( indicatives ) );

      burnTokens( dumpedCopy, flammable );
      labelAccounts( dumpedCopy, osaka );
      labelMultisigs( dumpedCopy, multisigs );
      fixChainIds( dumpedCopy, chainIds );
      fixErrors( dumpedCopy, indicativesCopy );

      const iov2star = mapIovToStar( dumpedCopy, multisigs, indicativesCopy, premiums );
      const { accounts, starnames, domains } = convertToCosmosSdk( dumpedCopy, iov2star, multisigs, premiums, reserveds );
      const custodian = accounts.find( account => account["//iov1"] == "iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc" );
      const iov = accounts.find( account => account["//iov1"] == "iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq" );
      const rewards = accounts.find( account => account["//iov1"] == "iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n" );
      const bonus = accounts.find( account => account["//iov1"] == "iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0" );
      const dave = accounts.find( account => account["//iov1"] == "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" );

      expect( custodian.value.address ).toEqual( multisigs.iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc.star1 );
      expect( custodian.value.coins[0].amount ).toEqual( "8321438667145" );
      expect( custodian["//no star1 iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98"][0] ).toEqual( 3234710 );
      expect( custodian["//no star1 iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98"][1] ).toEqual( "confio*iov" );
      expect( custodian["//no star1 iov1m7qjqjuv4ynhzu40xranun4u0r47d4waxc4wh9"][0] ).toEqual( 26.5 );
      expect( custodian["//no star1 iov1m7qjqjuv4ynhzu40xranun4u0r47d4waxc4wh9"][1] ).toEqual( "corentin*iov" );
      expect( custodian["//no star1 iov1q8zjkzk3f2yzfrkh9wswlf9qtmdgel84nnlgs9"] ).toEqual( 8920.657145 );
      expect( custodian["//no star1 iov1q40tvnph5xy7cjyj3tmqzghukeheykudq246d6"] ).toEqual( 22171 );
      expect( custodian["//no star1 iov1ua6tdcyw8jddn5660qcx2ndhjp4skqk4dkurrl"] ).toEqual( "alex*iov" );
      expect( custodian["//no star1 iov1ylw3cnluf3zayfths0ezgjp5cwf6ddvsvwa7l4"] ).toEqual( "lovely" );
      expect( custodian["//no star1 iov1y63fp8pncpuke7mrc2huqefud59t3munnh0k32"] ).toEqual( "multiverse" );

      expect( rewards.value.address ).toEqual( multisigs.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n.star1 );
      expect( rewards.value.coins[0].amount ).toEqual( "37000000" );

      expect( bonus.value.address ).toEqual( multisigs.iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0.star1 );
      expect( bonus.value.coins[0].amount ).toEqual( "3570582000000" );

      expect( dave.value.address ).toEqual( "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
      expect( dave.value.coins[0].amount ).toEqual( "416510000" );

      const alphaiov = starnames.find( starname => starname["//iov1"] == "iov16a42lf29n2h2eurxryspue9fz2d2wnlgpyjv8d" );
      const daveiov = starnames.find( starname => starname["//iov1"] == "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" );
      const huobiiov = starnames.find( starname => starname["//iov1"] == "iov1tlxqvugk9u5u973a6ee6dq4zsgsv6c5ecr0rvn" );

      expect( alphaiov.owner ).toEqual( "star1ayxmc4vqshd9j94hj67r55ppg5hsrhqlmy4dvd" );
      expect( alphaiov.domain ).toEqual( "iov" );
      expect( alphaiov.name ).toEqual( "alpha" );

      expect( daveiov.owner ).toEqual( "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
      expect( daveiov.domain ).toEqual( "iov" );
      expect( daveiov.name ).toEqual( "dave" );

      expect( huobiiov.owner ).toEqual( "star1vmt7wysxug30vfenedfh4ay83y3p75tstagn2y" );
      expect( huobiiov.domain ).toEqual( "iov" );
      expect( huobiiov.name ).toEqual( "huobi" );

      const alexiov = starnames.find( starname => starname["//iov1"] == "iov1ua6tdcyw8jddn5660qcx2ndhjp4skqk4dkurrl" );
      const confioiov = starnames.find( starname => starname["//iov1"] == "iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98" );
      const kadimaiov = starnames.find( starname => starname["//iov1"] == "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph" );

      expect( alexiov.owner ).toEqual( custodian.value.address );
      expect( alexiov.domain ).toEqual( "iov" );
      expect( alexiov.name ).toEqual( "alex" );

      expect( confioiov.owner ).toEqual( custodian.value.address );
      expect( confioiov.domain ).toEqual( "iov" );
      expect( confioiov.name ).toEqual( "confio" );

      expect( kadimaiov.owner ).toEqual( custodian.value.address );
      expect( kadimaiov.domain ).toEqual( "iov" );
      expect( kadimaiov.name ).toEqual( "kadima" );

      const lovely = domains.find( domain => domain.name == "lovely" );
      const multiverse = domains.find( domain => domain.name == "multiverse" );

      expect( lovely.admin ).toEqual( custodian.value.address );
      expect( multiverse.admin ).toEqual( custodian.value.address );

      const hash = domains.find( domain => domain.name == "hash" );
      const huth = domains.find( domain => domain.name == "huth" );
      const goldman = domains.find( domain => domain.name == "goldman" );
      const socgen = domains.find( domain => domain.name == "socgen" );
      const twitter = domains.find( domain => domain.name == "twitter" );
      const youtube = domains.find( domain => domain.name == "youtube" );
      const externalCustodian = "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96"; // HARD-CODED in conjunction with convertToCosmosSdk().address

      expect( hash.admin ).toEqual( "star1vmt7wysxug30vfenedfh4ay83y3p75tstagn2y" );
      expect( huth.admin ).toEqual( "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
      expect( goldman.admin ).toEqual( externalCustodian );
      expect( socgen.admin ).toEqual( externalCustodian );
      expect( twitter.admin ).toEqual( externalCustodian );
      expect( youtube.admin ).toEqual( externalCustodian );

      const limiteds = domains.filter( domain => domain.valid_until != lovely.valid_until );

      expect( limiteds.length ).toEqual( reserveds.length );

      const buckets = limiteds.reduce( ( accumulator, limited ) => {
         let count = accumulator[limited.valid_until] || 0;

         accumulator[limited.valid_until] = ++count;

         return accumulator;
      }, {} );

      expect( Object.keys( buckets ).length ).toEqual( 8 ); // HARD-CODED in conjunction with convertToCosmosSdk().releases
      Object.keys( buckets ).forEach( valid_until => expect( isFinite( valid_until ) ).toEqual( true ) ); // milliseconds since epoch
   } );

   it( `Should fail to add gentxs due to floating point amount.`, async () => {
      const tmpobj = tmp.dirSync( { template: "migrate-test-gentxs-XXXXXX", unsafeCleanup: true } );
      const home = tmpobj.name;
      const config = path.join( home, "config" );
      const gentx = "gentx-61e1f6d195f022cab0fe18f2ac1a4d33430999eb.json";
      const gentxs = path.join( home, "gentxs" );
      const genesisCopy = JSON.parse( JSON.stringify( genesis ) );

      genesisCopy.app_state.auth.accounts.push( { // add the account used in gentx
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
            "coins": [
               {
                  "denom": "iov",
                  "amount": "416.51" // must be an integer
               }
            ],
            "public_key": "",
            "account_number": 0,
            "sequence": 0
         },
         "//iov1": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un"
      } );

      fs.mkdirSync( config );
      fs.mkdirSync( gentxs );
      fs.copyFileSync( path.join( __dirname, gentx ), path.join( gentxs, gentx ) );
      fs.writeFileSync( path.join( config, "genesis.json" ), stringify( genesisCopy, { space: "  " } ), "utf-8" );

      try {
         addGentxs( gentxs, home );
      } catch ( e ) {
         expect( e.message.indexOf( "416.51" ) ).not.toEqual( -1 );
      }

      tmpobj.removeCallback();
   } );

   it( `Should fail to add gentxs due to missing account.`, async () => {
      const tmpobj = tmp.dirSync( { template: "migrate-test-gentxs-XXXXXX", unsafeCleanup: true } );
      const home = tmpobj.name;
      const config = path.join( home, "config" );
      const gentx = "gentx-61e1f6d195f022cab0fe18f2ac1a4d33430999eb.json";
      const gentxs = path.join( home, "gentxs" );

      fs.mkdirSync( config );
      fs.mkdirSync( gentxs );
      fs.copyFileSync( path.join( __dirname, gentx ), path.join( gentxs, gentx ) );
      fs.writeFileSync( path.join( config, "genesis.json" ), stringify( genesis, { space: "  " } ), "utf-8" );

      try {
         addGentxs( gentxs, home );
      } catch ( e ) {
         expect( e.message.indexOf( "account star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk not in genesis.json" ) ).not.toEqual( -1 );
      }

      tmpobj.removeCallback();
   } );

   it( `Should add gentxs.`, async () => {
      const tmpobj = tmp.dirSync( { template: "migrate-test-gentxs-XXXXXX", unsafeCleanup: true } );
      const home = tmpobj.name;
      const config = path.join( home, "config" );
      const gentx = "gentx-61e1f6d195f022cab0fe18f2ac1a4d33430999eb.json";
      const gentxs = path.join( home, "gentxs" );
      const genesisCopy = JSON.parse( JSON.stringify( genesis ) );

      genesisCopy.app_state.auth.accounts.push( { // add the account used in gentx
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "416"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         },
         "//iov1": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un"
      } );

      fs.mkdirSync( config );
      fs.mkdirSync( gentxs );
      fs.copyFileSync( path.join( __dirname, gentx ), path.join( gentxs, gentx ) );
      fs.writeFileSync( path.join( config, "genesis.json" ), stringify( genesisCopy, { space: "  " } ), "utf-8" );

      addGentxs( gentxs, home );

      const json = fs.readFileSync( path.join( config, "genesis.json" ), "utf-8" );
      const validatored = JSON.parse( json );
      const slim = {
         "type": "cosmos-sdk/StdTx",
         "value": {
            "msg": [
               {
                  "type": "cosmos-sdk/MsgCreateValidator",
                  "value": {
                     "description": {
                        "moniker": "slim.avaritia.com",
                        "identity": "",
                        "website": "",
                        "security_contact": "",
                        "details": ""
                     },
                     "commission": {
                        "rate": "0.100000000000000000",
                        "max_rate": "0.200000000000000000",
                        "max_change_rate": "0.010000000000000000"
                     },
                     "min_self_delegation": "1",
                     "delegator_address": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
                     "validator_address": "starvaloper1478t4fltj689nqu83vsmhz27quk7uggjtjp2gl",
                     "pubkey": "starvalconspub1zcjduepqds57cwz6kgzprcsuermllsyglcwz9w2z85nuar575z82mujtrhws0n4m0g",
                     "value": {
                        "denom": "uiov",
                        "amount": "1"
                     }
                  }
               }
            ],
            "fee": {
               "amount": [],
               "gas": "200000"
            },
            "signatures": [
               {
                  "pub_key": {
                     "type": "tendermint/PubKeySecp256k1",
                     "value": "AwOzGduZPxmjUMKASZGKPrUA7Drs9CvfJfXkgR/RSdyu"
                  },
                  "signature": "6ccfT5yLqEKJ1wB9KRbALLeKWilnWTwl0MfRQQr7YlAQWPSzPSnwQhbskOHPlGYxxD7XRX03MecQ9ZLxbqiOKA=="
               }
            ],
            "memo": "720e88e97dca985e303fc4ce1a8fc355571778cc@192.168.1.46:26656"
         }
      };

      compareObjects( slim, validatored.app_state.genutil.gentxs[0] );

      tmpobj.removeCallback();
   } );

   it( `Should fail to patch wrong-chain_id.`, async () => {
      const genesisCopy = JSON.parse( JSON.stringify( genesis ) );

      genesisCopy.chain_id = "wrong-chain_id";

      expect( () => { patchJestnet( genesisCopy ) } ).toThrow( `Wrong chain_id: ${genesisCopy.chain_id} != jestnet.` );
   } );

   it( `Should patch jestnet.`, async () => {
      const genesisCopy = JSON.parse( JSON.stringify( genesis ) );
      const previous = genesisCopy.app_state.starname.domains[0].account_renew;

      patchJestnet( genesisCopy );

      const current = genesisCopy.app_state.starname.domains[0].account_renew;

      expect( current ).not.toEqual( previous );
      expect( current ).toEqual( "3600" );
   } );

   it( `Should patch iovns-galaxynet.`, async () => {
      const genesisCopy = JSON.parse( JSON.stringify( genesis ) );
      const previous = [].concat( genesisCopy.app_state.auth.accounts );
      const poor =  {
         "//name": "dave*iov",
         "//iov1": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
            "coins": [
               {
                  "denom": "iov",
                  "amount": "416.51"
               }
            ],
            "public_key": "",
            "account_number": 0,
            "sequence": 0,
         },
      };
      const iovsas = {
         "//alias": "cond:multisig/usage/0000000000000001",
         "//id": "IOV SAS",
         "//iov1": "iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq",
         "type": "cosmos-sdk/Account",
         "value": {
            "account_number": 0,
            "address": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
            "coins": [
               {
                  "//IOV": 13015243.5,
                  "amount": "13015243500000",
                  "denom": "uiov"
               }
            ],
            "public_key": "",
            "sequence": 0
         }
      };
      const custodian = {
         "//alias": "cond:multisig/usage/0000000000000006",
         "//id": "Custodian of missing star1 accounts",
         "//iov1": "iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc",
         "type": "cosmos-sdk/Account",
         "value": {
            "account_number": 0,
            "address": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
            "coins": [
               {
                  "//IOV": 1.000123,
                  "amount": "72898628",
                  "denom": "uiov"
               }
            ],
            "public_key": "",
            "sequence": 0
         }
      };
      const isabella = {
         "//id": source2multisig.iov1w2suyhrfcrv5h4wmq3rk3v4x95cxtu0a03gy6x["//id"],
         "type": "cosmos-sdk/Account",
         "value": {
            "account_number": 0,
            "address": "star1elad203jykd8la6wgfnvk43rzajyqpk0wsme9g",
            "coins": [
               {
                  "//IOV": 2965149,
                  "amount": "2965149000000",
                  "denom": "uiov"
               }
            ],
            "public_key": "",
            "sequence": 0
         }
      };
      const accounts = [ poor, iovsas, custodian, isabella ];

      previous.push( ...JSON.parse( JSON.stringify( accounts ) ) );
      genesisCopy.app_state.auth.accounts = [].concat( previous );
      genesisCopy.chain_id = "iovns-galaxynet";
      genesisCopy.app_state.starname.domains = [
         {
            "account_renew": "315576000",
            "admin": multisigs.iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq.star1,
            "broker": null,
            "name": "iov",
            "type": "open",
            "valid_until": "1609415999"
         },
         {
            "//iov1": "zEaAIrHRUZTZF9uEWy0KJZ92J42T2",
            "account_renew": "315576000",
            "admin": multisigs.iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc.star1,
            "broker": null,
            "name": "0000",
            "type": "closed",
            "valid_until": "1609415999"
         },
      ];
      genesisCopy.app_state.starname.accounts = [
         {
            "//iov1": "iov1akhp7t0gtuaq4dwdw6qf0nvv6d2vf4vz8kwyl8",
            "broker": null,
            "certificates": null,
            "domain": "iov",
            "metadata_uri": "",
            "name": "...",
            "owner": multisigs.iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc.star1,
            "targets": null,
            "valid_until": "1609415999"
         },
         {
            "//iov1": "iov16qzp8q9kffdgamwtfcztg6z7puet374mgsxvhr",
            "broker": null,
            "certificates": null,
            "domain": "iov",
            "metadata_uri": "",
            "name": "01node",
            "owner": multisigs.iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc.star1,
            "targets": [
              {
                "address": "992290736603857528L",
                "blockchain_id": "lip9:9ee11e9df416b18b"
              },
              {
                "address": "0x6DF432079347050e0D8dA43C21fa6fe54697AfA7",
                "blockchain_id": "eip155:1"
              }
            ],
            "valid_until": "1609415999"
          },
      ];

      patchGalaxynet( genesisCopy );

      const current = genesisCopy.app_state.auth.accounts;

      expect( current.length ).not.toEqual( previous.length );
      expect( current.length ).toEqual( 13 );

      const antoine = current.find( account => account["//name"] == "antoine" );
      const dave = current.find( account => account["//name"] == "dave*iov" );
      const faucet = current.find( account => account["//name"] == "faucet" );
      const msig1 = current.find( account => account["//name"] == "msig1" );
      const w1 = current.find( account => account["//name"] == "w1" );
      const w2 = current.find( account => account["//name"] == "w2" );
      const w3 = current.find( account => account["//name"] == "w3" );

      expect( antoine ).toBeTruthy();
      expect( dave ).toBeTruthy();
      expect( faucet ).toBeTruthy();
      expect( msig1 ).toBeTruthy();
      expect( w1 ).toBeTruthy();
      expect( w2 ).toBeTruthy();
      expect( w3 ).toBeTruthy();

      expect( dave.value.address ).toEqual( "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
      expect( msig1.value.address ).toEqual( "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg" );
      expect( w1.value.address ).toEqual( "star19jj4wc3lxd54hkzl42m7ze73rzy3dd3wry2f3q" );
      expect( w2.value.address ).toEqual( "star1l4mvu36chkj9lczjhy9anshptdfm497fune6la" );
      expect( w3.value.address ).toEqual( "star1aj9qqrftdqussgpnq6lqj08gwy6ysppf53c8e9" );

      expect( dave.value.coins[0].amount ).not.toEqual( poor.value.coins[0].amount );
      expect( dave.value.coins[0].amount ).toEqual( "1000000000000" );

      const config = genesisCopy.app_state.configuration.config;

      expect( config["//note"] ).toEqual( "msig1 multisig address from w1,w2,w3,p1 in iovns/docs/cli, threshold 3" );
      expect( config.configurer ).toEqual( "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg" );
      expect( config.account_grace_period ).toEqual( "60000000000" );
      expect( config.account_renew_count_max ).toEqual( 2 );
      expect( config.account_renew_period ).toEqual( "180000000000" );
      expect( config.resources_max ).toEqual( 10 );
      expect( config.certificate_count_max ).toEqual( 3 );
      expect( config.certificate_size_max ).toEqual( "1000" );
      expect( config.domain_grace_period ).toEqual( "60000000000" );
      expect( config.domain_renew_count_max ).toEqual( 2 );
      expect( config.domain_renew_period ).toEqual( "300000000000" );
      expect( config.metadata_size_max ).toEqual( "1000" );

      const iov = genesisCopy.app_state.starname.domains.find( domain => domain.name == "iov" );
      const zeros = genesisCopy.app_state.starname.domains.find( domain => domain.name == "0000" );
      const dots = genesisCopy.app_state.starname.accounts.find( account => account.name == "..." );
      const claudiu  = genesisCopy.app_state.starname.accounts.find( account => account.name == "01node" );
      const escrow  = genesisCopy.app_state.auth.accounts.find( account => account["//id"] == "escrow isabella*iov" );

      expect( iov ).toBeTruthy();
      expect( zeros ).toBeTruthy();
      expect( dots ).toBeTruthy();
      expect( claudiu ).toBeTruthy();
      expect( escrow ).toBeTruthy();

      expect( iov.admin ).toEqual( "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78" );
      expect( zeros.admin ).toEqual( "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y" );
      expect( dots.owner ).toEqual( "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y" );
      expect( claudiu.owner ).toEqual( "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y" );
      expect( escrow.value.address ).toEqual( "star1elad203jykd8la6wgfnvk43rzajyqpk0wsme9g" );

      expect( dave.value.coins[0].denom ).toEqual( "uvoi" );
      expect( escrow.value.coins[0].denom ).toEqual( "uvoi" );
      expect( genesisCopy.app_state.mint.params.mint_denom ).toEqual( "uvoi" );
      expect( genesisCopy.app_state.staking.params.bond_denom ).toEqual( "uvoi" );
      expect( genesisCopy.app_state.configuration.fees.fee_coin_denom ).toEqual( "uvoi" );
      expect( genesisCopy.app_state.crisis.constant_fee.denom ).toEqual( "uvoi" );
      expect( genesisCopy.app_state.gov.deposit_params.min_deposit[0].denom ).toEqual( "uvoi" );
   } );

   it( `Should patch iov-mainnet-2.`, async () => {
      const genesisCopy = JSON.parse( JSON.stringify( genesis ) );

      genesisCopy.chain_id = "iov-mainnet-2";
      genesisCopy.app_state.auth.accounts.push( {
         "//alias": "cond:multisig/usage/0000000000000006",
         "//id": "Custodian of missing star1 accounts",
         "//iov1": "iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc",
         "//no star1 iov14qk7zrz2ewhdmy7cjj68sk6jn3rst4vd7u930y": [
            122534,
            "misang*iov"
         ],
         "//no star1 iov1jq8z8xl9tqdwjsp44gtkd2c5rpq33e556kg0ft": [
            700666,
            "charlief*iov"
         ],
         "//no star1 iov153n95ekuw9rxfhzspgarqjdwnadmvdt0chcjs4": [
            1111111,
            "gillesd*iov"
         ],
         "//no star1 iov12gd6weg7py6vs7ujn22h82422arek8cxzhe85p": [
            "adrian",
            "adrianirimia",
            "world"
         ],
         "type": "cosmos-sdk/Account",
         "value": {
            "account_number": 0,
            "address": "star1xc7tn8szhtvcat2k29t6072235gsqcrujd60wy",
            "coins": [
               {
                  "//IOV": 1.000123,
                  "amount": "73262964",
                  "denom": "uvoi"
               }
            ],
            "public_key": "",
            "sequence": 0
         }
      },
      {
         "//id": 2191,
         "//iov1": "iov1lfjspe4x5u404sskmv5md4q7u9jcz96zya8krw",
         "type": "cosmos-sdk/Account",
         "value": {
            "account_number": "0",
            "address": "star1alx5xdyws6hzeynewdzm4x5wg0v6tjl3wcq34g",
            "coins": [
               {
                  "//IOV": 49998,
                  "amount": "49998000000",
                  "denom": "uvoi"
               }
            ],
            "public_key": null,
            "sequence": "0"
         }
      },
      {
         "//id": 2192,
         "//iov1": "iov1ja0syy203qncn28cqmz5zh9kh2xl0xxt36m4qx",
         "type": "cosmos-sdk/Account",
         "value": {
            "account_number": "0",
            "address": "star16yqeeqp26h64rldcql7cpsz0zte7rwx9jc8s6x",
            "coins": [
               {
                  "//IOV": 22220,
                  "amount": "22220000000",
                  "denom": "uvoi"
               }
            ],
            "public_key": null,
            "sequence": "0"
         }
      },
      {
         "//id": 2193,
         "//iov1": "iov1axxtqae3x9jtvv7wavg6fnjgpc27dx7a9jlp9r",
         "type": "cosmos-sdk/Account",
         "value": {
            "account_number": "0",
            "address": "star1myg5g0uvsjdq20atrjv5t49ffflwz2m7d7u677",
            "coins": [
               {
                  "//IOV": 235925.5,
                  "amount": "235925500000",
                  "denom": "uvoi"
               }
            ],
            "public_key": null,
            "sequence": "0"
         }
      } );
      genesisCopy.app_state.starname.accounts.push( {
         "//iov1": "iov14qk7zrz2ewhdmy7cjj68sk6jn3rst4vd7u930y",
         "broker": null,
         "certificates": null,
         "domain": "iov",
         "metadata_uri": "",
         "name": "misang",
         "owner": "star1xc7tn8szhtvcat2k29t6072235gsqcrujd60wy",
         "targets": null,
         "valid_until": "1609415999"
      },
      {
         "//iov1": "iov1jq8z8xl9tqdwjsp44gtkd2c5rpq33e556kg0ft",
         "broker": null,
         "certificates": null,
         "domain": "iov",
         "metadata_uri": "",
         "name": "charlief",
         "owner": "star1xc7tn8szhtvcat2k29t6072235gsqcrujd60wy",
         "targets": null,
         "valid_until": "1609415999"
      },
      {
         "//iov1": "iov1jq8z8xl9tqdwjsp44gtkd2c5rpq33e556kg0ft",
         "broker": null,
         "certificates": null,
         "domain": "iov",
         "metadata_uri": "",
         "name": "gillesd",
         "owner": "star1xc7tn8szhtvcat2k29t6072235gsqcrujd60wy",
         "targets": null,
         "valid_until": "1609415999"
      },
      {
         "//iov1": "iov1axxtqae3x9jtvv7wavg6fnjgpc27dx7a9jlp9r",
         "broker": null,
         "certificates": null,
         "domain": "iov",
         "metadata_uri": "",
         "name": "de26",
         "owner": "star1myg5g0uvsjdq20atrjv5t49ffflwz2m7d7u677",
         "resources": [
            {
               "resource": "star1myg5g0uvsjdq20atrjv5t49ffflwz2m7d7u677",
               "uri": "asset:iov"
            }
         ],
         "valid_until": "1609415999"
      },
      {
         "//iov1": "iov1lfjspe4x5u404sskmv5md4q7u9jcz96zya8krw",
         "broker": null,
         "certificates": null,
         "domain": "iov",
         "metadata_uri": "",
         "name": "jean501",
         "owner": "star1alx5xdyws6hzeynewdzm4x5wg0v6tjl3wcq34g",
         "resources": [
            {
               "resource": "star1alx5xdyws6hzeynewdzm4x5wg0v6tjl3wcq34g",
               "uri": "asset:iov"
            }
         ],
         "valid_until": "1609415999"
      },
      {
         "//iov1": "iov1ja0syy203qncn28cqmz5zh9kh2xl0xxt36m4qx",
         "broker": null,
         "certificates": null,
         "domain": "iov",
         "metadata_uri": "",
         "name": "mam",
         "owner": "star16yqeeqp26h64rldcql7cpsz0zte7rwx9jc8s6x",
         "resources": [
            {
               "resource": "star16yqeeqp26h64rldcql7cpsz0zte7rwx9jc8s6x",
               "uri": "asset:iov"
            }
         ],
         "valid_until": "1609415999"
      } );

      const accounts0 = [].concat( genesisCopy.app_state.auth.accounts );
      const starnames0 = [].concat( genesisCopy.app_state.starname.accounts );
      const custodian0 = JSON.parse( JSON.stringify( accounts0.find( account => account["//id"] == "Custodian of missing star1 accounts" ) ) );

      patchMainnet( genesisCopy );

      expect( genesisCopy.app_state.auth.accounts.length ).toEqual( accounts0.length + 3 ); // charlie, gilles, misang
      expect( genesisCopy.app_state.starname.accounts.length ).toEqual( starnames0.length );

      // lostKeysInCustody
      const custodian = genesisCopy.app_state.auth.accounts.find( account => account["//id"] == "Custodian of missing star1 accounts" );
      const charliefAmount = 700666 * 1e6;
      const charliestar1 = "star1k9ktkefsdxtydga262re596agdklwjmrf9et90";
      const charlief = genesisCopy.app_state.auth.accounts.find( account => account.value.address == charliestar1 );
      const charliefiov = genesisCopy.app_state.starname.accounts.find( account => account.owner == charliestar1 );
      const gillesdAmount = 1111111 * 1e6;
      const gillesdstar1 = "star1keaxspy5rgw84azg5w640pp8zdla72ra0n5xh2";
      const gillesd = genesisCopy.app_state.auth.accounts.find( account => account.value.address == gillesdstar1 );
      const gillesdiov = genesisCopy.app_state.starname.accounts.find( account => account.owner == gillesdstar1 );
      const misangAmount = 122534 * 1e6;
      const misangstar1 = "star1lgh6ekcnkufs4742qr5znvtlz4vglul9g2p6xl";
      const misang = genesisCopy.app_state.auth.accounts.find( account => account.value.address == misangstar1 );
      const misangiov = genesisCopy.app_state.starname.accounts.find( account => account.owner == misangstar1 );

      expect( custodian.value.coins[0].amount ).toEqual( String( +custodian0.value.coins[0].amount - charliefAmount - gillesdAmount - misangAmount ) );
      expect( charlief.value.coins[0].amount ).toEqual( String( charliefAmount ) );
      expect( charliefiov ).toBeTruthy();
      expect( gillesd.value.coins[0].amount ).toEqual( String( gillesdAmount ) );
      expect( gillesdiov ).toBeTruthy();
      expect( misang.value.coins[0].amount ).toEqual( String( misangAmount ) );
      expect( misangiov ).toBeTruthy();

      // lostKeysWithStar1
      const de26star1 = "star1xnzwj34e8zefm7g7vtgnphfj6x2qgnq723rq0j";
      const de26 = genesisCopy.app_state.auth.accounts.find( account => account.value.address == de26star1 );
      const de26iov = genesisCopy.app_state.starname.accounts.find( account => account.owner == de26star1 );
      const jean501star1 = "star1lsk9ckth2s870kjqcyl6x5af7gazj6eg7msluq";
      const jean501 = genesisCopy.app_state.auth.accounts.find( account => account.value.address == jean501star1 );
      const jean501iov = genesisCopy.app_state.starname.accounts.find( account => account.owner == jean501star1 );
      const mamstar1 = "star1f2jpr2guzq3y5yjv667axr26pl6qzyn2hzthfa";
      const mam = genesisCopy.app_state.auth.accounts.find( account => account.value.address == mamstar1 );
      const mamiov = genesisCopy.app_state.starname.accounts.find( account => account.owner == mamstar1 );

      expect( de26 ).toBeTruthy();
      expect( de26iov.name ).toEqual( "de26" );
      expect( de26iov.resources[0].resource ).toEqual( de26star1 );
      expect( jean501 ).toBeTruthy();
      expect( jean501iov.name ).toEqual( "jean501" );
      expect( jean501iov.resources[0].resource ).toEqual( jean501star1 );
      expect( mam ).toBeTruthy();
      expect( mamiov.name ).toEqual( "mam" );
      expect( mamiov.resources[0].resource ).toEqual( mamstar1 );
   } );

   it( `Should migrate.`, async () => {
      const tmpobj = tmp.dirSync( { template: "migrate-test-migrate-XXXXXX", unsafeCleanup: true } );
      const home = tmpobj.name;
      const config = path.join( home, "config" );
      const gentx = "gentx-61e1f6d195f022cab0fe18f2ac1a4d33430999eb.json";
      const gentxs = path.join( home, "gentxs" );

      fs.mkdirSync( config );
      fs.mkdirSync( gentxs );
      fs.copyFileSync( path.join( __dirname, gentx ), path.join( gentxs, gentx ) );

      migrate( { chainIds, dumped, flammable, genesis, gentxs, home, indicatives, multisigs, osaka, premiums, reserveds, source2multisig } );

      const nextGen = {
         "app_hash": "",
         "app_state": {
            "auth": {
               "accounts": [
                  {
                     "//alias": "cond:gov/rule/0000000000000002",
                     "//id": "reward fund",
                     "//iov1": "iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n",
                     "type": "cosmos-sdk/Account",
                     "value": {
                        "account_number": "0",
                        "address": "star1scfumxscrm53s4dd3rl93py5ja2ypxmxlhs938",
                        "coins": [
                           {
                              "//IOV": 37,
                              "amount": "37000000",
                              "denom": "uiov"
                           }
                        ],
                        "public_key": null,
                        "sequence": "0"
                     }
                  },
                  {
                     "//alias": "cond:multisig/usage/0000000000000002",
                     "//id": "IOV SAS employee bonus pool/colloboration appropriation pool",
                     "//iov1": "iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0",
                     "type": "cosmos-sdk/Account",
                     "value": {
                        "account_number": "0",
                        "address": "star16tm7scg0c2e04s0exk5rgpmws2wk4xkd84p5md",
                        "coins": [
                           {
                              "//IOV": 3570582,
                              "amount": "3570582000000",
                              "denom": "uiov"
                           }
                        ],
                        "public_key": null,
                        "sequence": "0"
                     }
                  },
                  {
                     "//alias": "cond:multisig/usage/0000000000000006",
                     "//id": "Custodian of missing star1 accounts",
                     "//iov1": "iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc",
                     "//no star1 iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98": [
                        3234710,
                        "confio*iov"
                     ],
                     "//no star1 iov1m7qjqjuv4ynhzu40xranun4u0r47d4waxc4wh9": [
                        26.5,
                        "corentin*iov"
                     ],
                     "//no star1 iov1q40tvnph5xy7cjyj3tmqzghukeheykudq246d6": 22171,
                     "//no star1 iov1q8zjkzk3f2yzfrkh9wswlf9qtmdgel84nnlgs9": 8920.657145,
                     "//no star1 iov1ua6tdcyw8jddn5660qcx2ndhjp4skqk4dkurrl": "alex*iov",
                     "//no star1 iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph": [
                        1628971,
                        "kadima*iov"
                     ],
                     "//no star1 iov1y63fp8pncpuke7mrc2huqefud59t3munnh0k32": "multiverse",
                     "//no star1 iov1ylw3cnluf3zayfths0ezgjp5cwf6ddvsvwa7l4": "lovely",
                     "//no star1 zGLlamFypWMPUeHVVsvo4mXFFOE63": [
                        "cosmostation",
                        "ibcwallet",
                        "korea",
                        "mintscan",
                        "seoul",
                        "station"
                     ],
                     "//no star1 zHbPpUYyRguRlhAiC30zimM05hGx2": "jim",
                     "type": "cosmos-sdk/Account",
                     "value": {
                        "account_number": "0",
                        "address": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                        "coins": [
                           {
                              "//IOV": 1.000123,
                              "amount": "4895215667145",
                              "denom": "uiov"
                           }
                        ],
                        "public_key": null,
                        "sequence": "0"
                     }
                  },
                  {
                     "//alias": "cond:multisig/usage/0000000000000001",
                     "//id": "IOV SAS",
                     "//iov1": "iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq",
                     "type": "cosmos-sdk/Account",
                     "value": {
                        "account_number": "0",
                        "address": "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78",
                        "coins": [
                           {
                              "//IOV": 13015243.5,
                              "amount": "13015243500000",
                              "denom": "uiov"
                           }
                        ],
                        "public_key": null,
                        "sequence": "0"
                     }
                  },
                  {
                     "//iov1": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un",
                     "type": "cosmos-sdk/Account",
                     "value": {
                        "account_number": "0",
                        "address": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
                        "coins": [
                           {
                              "//IOV": 416.51,
                              "amount": "416510000",
                              "denom": "uiov"
                           }
                        ],
                        "public_key": null,
                        "sequence": "0"
                     }
                  },
                  {
                     "//id": "vaildator guaranteed reward fund",
                     "//note": "consolidated escrows with source iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u",
                     "//timeout 2029-11-10T00:00:00.000Z": "iov170qvwm0tscn5mza3vmaerkzqllvwc3kykkt7kj yields 2347987 IOV",
                     "type": "cosmos-sdk/Account",
                     "value": {
                        "account_number": "0",
                        "address": "star17w7fjdkr9laphtyj4wxa32rf0evu94xgywxgl4",
                        "coins": [
                           {
                              "//IOV": 2347987,
                              "amount": "2347987000000",
                              "denom": "uiov"
                           }
                        ],
                        "public_key": null,
                        "sequence": "0"
                     }
                  },
                  {
                     "//id": "escrow isabella*iov",
                     "//note": "consolidated escrows with source iov1w2suyhrfcrv5h4wmq3rk3v4x95cxtu0a03gy6x",
                     "//timeout 2019-12-10T12:00:00.000Z": "iov105465l8l3yn06a56h7tqwwvnqq22e8j4nvgf02 yields 269559 IOV",
                     "//timeout 2020-01-10T12:00:00.000Z": "iov17gdpegksje9dlh8h0g6ehgk6d4anz9pkfskunr yields 269559 IOV",
                     "//timeout 2020-02-10T12:00:00.000Z": "iov1ppxx0vwx42p47p4pkztzl4d57zh2ctnwsz4fdu yields 269559 IOV",
                     "type": "cosmos-sdk/Account",
                     "value": {
                        "account_number": "0",
                        "address": "star1elad203jykd8la6wgfnvk43rzajyqpk0wsme9g",
                        "coins": [
                           {
                              "//IOV": 808677,
                              "amount": "808677000000",
                              "denom": "uiov"
                           }
                        ],
                        "public_key": null,
                        "sequence": "0"
                     }
                  },
                  {
                     "//id": "escrow kadima*iov",
                     "//note": "consolidated escrows with source iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph",
                     "//timeout 2020-06-10T12:00:00.000Z": "iov1k4dpknrrf4dfm07avau0mmjkrsm6pu863d30us yields 89853 IOV",
                     "//timeout 2020-07-10T12:00:00.000Z": "iov1dfurgye70k7f2gxptztfym697g5t832pp9m94g yields 89853 IOV",
                     "//timeout 2020-08-10T12:00:00.000Z": "iov1497txu54lnwujzl8xhc59y6cmuw82d68udn4l3 yields 89853 IOV",
                     "type": "cosmos-sdk/Account",
                     "value": {
                        "account_number": "0",
                        "address": "star1hjf04872s9rlcdg2wqwvapwttvt3p4gjpp0xmc",
                        "coins": [
                           {
                              "//IOV": 269559,
                              "amount": "269559000000",
                              "denom": "uiov"
                           }
                        ],
                        "public_key": null,
                        "sequence": "0"
                     }
                  }
               ],
               "params": {
                  "max_memo_characters": "256",
                  "sig_verify_cost_ed25519": "590",
                  "sig_verify_cost_secp256k1": "1000",
                  "tx_sig_limit": "7",
                  "tx_size_cost_per_byte": "10"
               }
            },
            "bank": {
               "send_enabled": true
            },
            "configuration": {
               "config": {
                  "account_grace_period": "2592000",
                  "account_renew_count_max": 2,
                  "account_renew_period": "31557600",
                  "certificate_count_max": 16,
                  "certificate_size_max": "1024",
                  "configurer": "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78",
                  "domain_grace_period": "2592000",
                  "domain_renew_count_max": 2,
                  "domain_renew_period": "31557600",
                  "metadata_size_max": "1024",
                  "resource_target_max": 16,
                  "valid_account_name": "[-_\\.a-z0-9]{0,64}$",
                  "valid_domain_name": "^[-_a-z0-9]{4,16}$",
                  "valid_resource": "^[a-z0-9A-Z]+$",
                  "valid_uri": "[-a-z0-9A-Z:]+$"
               },
               "fees": {
                  "add_account_certificate": "100",
                  "del_account_certificate": "0.5",
                  "fee_coin_denom": "uiov",
                  "fee_coin_price": "0.000001",
                  "fee_default": "0.5",
                  "register_account_closed": "0.5",
                  "register_account_open": "10",
                  "register_domain_1": "10000",
                  "register_domain_2": "5000",
                  "register_domain_3": "2000",
                  "register_domain_4": "1000",
                  "register_domain_5": "500",
                  "register_domain_default": "250",
                  "register_open_domain_multiplier": "5.5",
                  "renew_domain_open": "12345",
                  "replace_account_resources": "10",
                  "set_account_metadata": "500",
                  "transfer_account_closed": "10",
                  "transfer_account_open": "10",
                  "transfer_domain_closed": "10",
                  "transfer_domain_open": "10"
               }
            },
            "crisis": {
               "constant_fee": {
                  "amount": "1000000000",
                  "denom": "uiov"
               }
            },
            "distribution": {
               "delegator_starting_infos": [
               ],
               "delegator_withdraw_infos": [
               ],
               "fee_pool": {
                  "community_pool": [
                  ]
               },
               "outstanding_rewards": [
               ],
               "params": {
                  "base_proposer_reward": "0.050000000000000000",
                  "bonus_proposer_reward": "0.050000000000000000",
                  "community_tax": "0.000000000000000000",
                  "withdraw_addr_enabled": true
               },
               "previous_proposer": "",
               "validator_accumulated_commissions": [
               ],
               "validator_current_rewards": [
               ],
               "validator_historical_rewards": [
               ],
               "validator_slash_events": [
               ]
            },
            "evidence": {
               "evidence": [
               ],
               "params": {
                  "max_evidence_age": "1814400000000000"
               }
            },
            "genutil": {
            },
            "gov": {
               "deposit_params": {
                  "max_deposit_period": "172800000000000",
                  "min_deposit": [
                     {
                        "amount": "1000000000",
                        "denom": "uiov"
                     }
                  ]
               },
               "deposits": null,
               "proposals": null,
               "starting_proposal_id": "1",
               "tally_params": {
                  "quorum": "0.334000000000000000",
                  "threshold": "0.500000000000000000",
                  "veto": "0.334000000000000000"
               },
               "votes": null,
               "voting_params": {
                  "voting_period": "345600000000000"
               }
            },
            "mint": {
               "minter": {
                  "annual_provisions": "0.000000000000000000",
                  "inflation": "0.000000000000000000"
               },
               "params": {
                  "//note": "goal_bonded cannot be 0: module=consensus err='division by zero'",
                  "blocks_per_year": "105192",
                  "goal_bonded": "0.000000000000000001",
                  "inflation_max": "0.0000000000000000",
                  "inflation_min": "0.0000000000000000",
                  "inflation_rate_change": "0.000000000000000000",
                  "mint_denom": "uiov"
               }
            },
            "slashing": {
               "missed_blocks": {
               },
               "params": {
                  "downtime_jail_duration": "600000000000",
                  "max_evidence_age": "1814400000000000",
                  "min_signed_per_window": "0.500000000000000000",
                  "signed_blocks_window": "10000",
                  "slash_fraction_double_sign": "0.050000000000000000",
                  "slash_fraction_downtime": "0.010000000000000000"
               },
               "signing_infos": {
               }
            },
            "staking": {
               "delegations": null,
               "exported": false,
               "last_total_power": "0",
               "last_validator_powers": null,
               "params": {
                  "bond_denom": "uiov",
                  "historical_entries": 0,
                  "max_entries": 7,
                  "max_validators": 16,
                  "unbonding_time": "259200000000000"
               },
               "redelegations": null,
               "unbonding_delegations": null,
               "validators": null
            },
            "starname": {
               "accounts": [
                  {
                     "//iov1": "iov1ua6tdcyw8jddn5660qcx2ndhjp4skqk4dkurrl",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "alex",
                     "owner": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "resources": null,
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov16a42lf29n2h2eurxryspue9fz2d2wnlgpyjv8d",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "alpha",
                     "owner": "star1ayxmc4vqshd9j94hj67r55ppg5hsrhqlmy4dvd",
                     "resources": [
                        {
                           "resource": "0x52dBf411B22Af67f33425bf3fbb6B8CF8FB302cd",
                           "uri": "asset:eth"
                        },
                        {
                           "resource": "cosmos15dafemy5pkaru4kf23s3e6mnugfv6et9kg2uz7",
                           "uri": "asset:atom"
                        },
                        {
                           "resource": "star1ayxmc4vqshd9j94hj67r55ppg5hsrhqlmy4dvd",
                           "uri": "asset:iov"
                        }
                     ],
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov10v69k57z2v0pr3yvtr60pp8g2jx8tdd7f55sv6",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "blini44",
                     "owner": "star1gfdmksf725qpdgl06e98ks4usg9nmkcwc5qzcg",
                     "resources": [
                        {
                           "resource": "star1gfdmksf725qpdgl06e98ks4usg9nmkcwc5qzcg",
                           "uri": "asset:iov"
                        },
                        {
                           "resource": "0xa223f22664Ee8bfB41FAD93C388826E7aF24060c",
                           "uri": "asset:eth"
                        },
                        {
                           "resource": "4341330819731245941L",
                           "uri": "asset:lsk"
                        }
                     ],
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1yhk8qqp3wsdg7tefd8u457n9zqsny4nqzp6960",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "btc13",
                     "owner": "star1qvpth6t72336fjxlej2xv8eu84hrpxdxf5rgzz",
                     "resources": [
                        {
                           "resource": "star1qvpth6t72336fjxlej2xv8eu84hrpxdxf5rgzz",
                           "uri": "asset:iov"
                        },
                        {
                           "resource": "16104600299727948959L",
                           "uri": "asset:lsk"
                        },
                        {
                           "resource": "0x40698A9DcE4d6a63E766Dd08b83D03c6727DCB1a",
                           "uri": "asset:eth"
                        }
                     ],
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1j43xew5yq7ap2kesgjnlzru0z22grs94qsyf98",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "confio",
                     "owner": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "resources": null,
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1m7qjqjuv4ynhzu40xranun4u0r47d4waxc4wh9",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "corentin",
                     "owner": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "resources": [
                        {
                           "resource": "0x5e415520beb66aa39e00d43cae889f2c5cba7017",
                           "uri": "asset:eth"
                        }
                     ],
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "dave",
                     "owner": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
                     "resources": [
                        {
                           "resource": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
                           "uri": "asset:iov"
                        }
                     ],
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1ejk0g6p2xk90lamuvtd3r0kf6jcva09hf4xy74",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "fish_and_chips",
                     "owner": "star1yxxmpqca3l7xzhy4783vkpfx843x4zk749h8fs",
                     "resources": [
                        {
                           "resource": "star1yxxmpqca3l7xzhy4783vkpfx843x4zk749h8fs",
                           "uri": "asset:iov"
                        }
                     ],
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1tlxqvugk9u5u973a6ee6dq4zsgsv6c5ecr0rvn",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "huobi",
                     "owner": "star1vmt7wysxug30vfenedfh4ay83y3p75tstagn2y",
                     "resources": [
                        {
                           "resource": "0x00C60938d954FEC83E70eE98243B24F7E6EabaC8",
                           "uri": "asset:eth"
                        },
                        {
                           "resource": "13483265462465913551L",
                           "uri": "asset:lsk"
                        },
                        {
                           "resource": "star1vmt7wysxug30vfenedfh4ay83y3p75tstagn2y",
                           "uri": "asset:iov"
                        }
                     ],
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "kadima",
                     "owner": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "resources": null,
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "ledger",
                     "owner": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
                     "resources": [
                        {
                           "resource": "iov1fpezwaxfnmef8tyyg4t7avz9a2d9gqh3yh8d8n",
                           "uri": "asset:iov"
                        }
                     ],
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov16a42lf29n2h2eurxryspue9fz2d2wnlgpyjv8d",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "letsdoit",
                     "owner": "star1ayxmc4vqshd9j94hj67r55ppg5hsrhqlmy4dvd",
                     "resources": [
                        {
                           "resource": "0x52dBf411B22Af67f33425bf3fbb6B8CF8FB302cd",
                           "uri": "asset:eth"
                        },
                        {
                           "resource": "star1ayxmc4vqshd9j94hj67r55ppg5hsrhqlmy4dvd",
                           "uri": "asset:iov"
                        }
                     ],
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1eh6yeyel3zsc8vqnh79fqjtfkcxmj5d8nt49gq",
                     "broker": null,
                     "certificates": null,
                     "domain": "iov",
                     "metadata_uri": "",
                     "name": "nash.io",
                     "owner": "star1y86zdqsegxm7uj9qf7l400y29nc6x9ypqxpdcg",
                     "resources": [
                        {
                           "resource": "0x2cE327b4EB237313F37a72195d64Cb80F7aeAa15",
                           "uri": "asset:eth"
                        },
                        {
                           "resource": "16192453558792957658L",
                           "uri": "asset:lsk"
                        },
                        {
                           "resource": "star1y86zdqsegxm7uj9qf7l400y29nc6x9ypqxpdcg",
                           "uri": "asset:iov"
                        }
                     ],
                     "valid_until": "1628034048"
                  }
               ],
               "domains": [
                  {
                     "//note": "msig1",
                     "account_renew": "3000",
                     "admin": "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg",
                     "broker": null,
                     "name": "iov",
                     "type": "open",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov12gd6weg7py6vs7ujn22h82422arek8cxzhe85p",
                     "account_renew": "315576000",
                     "admin": "star1usl4zpltjesrp5rqae3fdjdyj5dyymakmhq6mt",
                     "broker": null,
                     "name": "84more",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov12gd6weg7py6vs7ujn22h82422arek8cxzhe85p",
                     "account_renew": "315576000",
                     "admin": "star1usl4zpltjesrp5rqae3fdjdyj5dyymakmhq6mt",
                     "broker": null,
                     "name": "adrian",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "china",
                     "type": "closed",
                     "valid_until": "1621418400"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "china1",
                     "type": "closed",
                     "valid_until": "1621418400"
                  },
                  {
                     "//iov1": "zGLlamFypWMPUeHVVsvo4mXFFOE63",
                     "account_renew": "315576000",
                     "admin": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "broker": null,
                     "name": "cosmostation",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "fcbarca",
                     "type": "closed",
                     "valid_until": "1618961400"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "fcbarca1",
                     "type": "closed",
                     "valid_until": "1618961400"
                  },
                  {
                     "//iov1": "iov1zr9epgrzysr6zc5s8ucd3qlxkhgj9fwj2a2mkx",
                     "account_renew": "315576000",
                     "admin": "star18awsa7fhwtsevta28p3uw8ymtznvpwtzl3ep5f",
                     "broker": null,
                     "name": "gianna",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "goldman",
                     "type": "closed",
                     "valid_until": "1602664200"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "goldman_sachs",
                     "type": "closed",
                     "valid_until": "1615986000"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "goldman_sachs1",
                     "type": "closed",
                     "valid_until": "1615986000"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "goldman1",
                     "type": "closed",
                     "valid_until": "1602664200"
                  },
                  {
                     "//iov1": "iov1tlxqvugk9u5u973a6ee6dq4zsgsv6c5ecr0rvn",
                     "account_renew": "315576000",
                     "admin": "star1vmt7wysxug30vfenedfh4ay83y3p75tstagn2y",
                     "broker": null,
                     "name": "hash",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1tlxqvugk9u5u973a6ee6dq4zsgsv6c5ecr0rvn",
                     "account_renew": "315576000",
                     "admin": "star1vmt7wysxug30vfenedfh4ay83y3p75tstagn2y",
                     "broker": null,
                     "name": "hell",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1tlxqvugk9u5u973a6ee6dq4zsgsv6c5ecr0rvn",
                     "account_renew": "315576000",
                     "admin": "star1vmt7wysxug30vfenedfh4ay83y3p75tstagn2y",
                     "broker": null,
                     "name": "hold",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un",
                     "account_renew": "315576000",
                     "admin": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
                     "broker": null,
                     "name": "huth",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "zGLlamFypWMPUeHVVsvo4mXFFOE63",
                     "account_renew": "315576000",
                     "admin": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "broker": null,
                     "name": "ibcwallet",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un",
                     "account_renew": "315576000",
                     "admin": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
                     "broker": null,
                     "name": "in3s",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "ivanontech",
                     "type": "closed",
                     "valid_until": "1613529000"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "ivanontech1",
                     "type": "closed",
                     "valid_until": "1613529000"
                  },
                  {
                     "//iov1": "zHbPpUYyRguRlhAiC30zimM05hGx2",
                     "account_renew": "315576000",
                     "admin": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "broker": null,
                     "name": "jim",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "zGLlamFypWMPUeHVVsvo4mXFFOE63",
                     "account_renew": "315576000",
                     "admin": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "broker": null,
                     "name": "korea",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1ylw3cnluf3zayfths0ezgjp5cwf6ddvsvwa7l4",
                     "account_renew": "315576000",
                     "admin": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "broker": null,
                     "name": "lovely",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1zr9epgrzysr6zc5s8ucd3qlxkhgj9fwj2a2mkx",
                     "account_renew": "315576000",
                     "admin": "star18awsa7fhwtsevta28p3uw8ymtznvpwtzl3ep5f",
                     "broker": null,
                     "name": "michael",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "zGLlamFypWMPUeHVVsvo4mXFFOE63",
                     "account_renew": "315576000",
                     "admin": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "broker": null,
                     "name": "mintscan",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1y63fp8pncpuke7mrc2huqefud59t3munnh0k32",
                     "account_renew": "315576000",
                     "admin": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "broker": null,
                     "name": "multiverse",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1zr9epgrzysr6zc5s8ucd3qlxkhgj9fwj2a2mkx",
                     "account_renew": "315576000",
                     "admin": "star18awsa7fhwtsevta28p3uw8ymtznvpwtzl3ep5f",
                     "broker": null,
                     "name": "nodeateam",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un",
                     "account_renew": "315576000",
                     "admin": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
                     "broker": null,
                     "name": "sentient",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "zGLlamFypWMPUeHVVsvo4mXFFOE63",
                     "account_renew": "315576000",
                     "admin": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "broker": null,
                     "name": "seoul",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "socgen",
                     "type": "closed",
                     "valid_until": "1605726000"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "socgen1",
                     "type": "closed",
                     "valid_until": "1605726000"
                  },
                  {
                     "//iov1": "zGLlamFypWMPUeHVVsvo4mXFFOE63",
                     "account_renew": "315576000",
                     "admin": "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y",
                     "broker": null,
                     "name": "station",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "//iov1": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un",
                     "account_renew": "315576000",
                     "admin": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk",
                     "broker": null,
                     "name": "tachyon",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "twitter",
                     "type": "closed",
                     "valid_until": "1608096600"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "twitter1",
                     "type": "closed",
                     "valid_until": "1608096600"
                  },
                  {
                     "//iov1": "iov1zr9epgrzysr6zc5s8ucd3qlxkhgj9fwj2a2mkx",
                     "account_renew": "315576000",
                     "admin": "star18awsa7fhwtsevta28p3uw8ymtznvpwtzl3ep5f",
                     "broker": null,
                     "name": "tyler",
                     "type": "closed",
                     "valid_until": "1628034048"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "youtube",
                     "type": "closed",
                     "valid_until": "1610553600"
                  },
                  {
                     "account_renew": "315576000",
                     "admin": "star1v794jm5am4qpc52kvgmxxm2j50kgu9mjszcq96",
                     "broker": null,
                     "name": "youtube1",
                     "type": "closed",
                     "valid_until": "1610553600"
                  }
               ]
            },
            "supply": {
               "supply": [
               ]
            }
         },
         "chain_id": "jestnet",
         "consensus_params": {
            "block": {
               "max_bytes": "500000",
               "max_gas": "-1",
               "time_iota_ms": "1000"
            },
            "evidence": {
               "max_age_duration": "172800000000000",
               "max_age_num_blocks": "100000"
            },
            "validator": {
               "pub_key_types": [
                  "ed25519"
               ]
            }
         },
         "genesis_time": "2020-04-15T10:00:00.000Z"
      };

      // hack around transient values before...
      const fixTransients = ( previous, current ) => {
         for ( let i = 0, n = previous.length; i < n; ++i ) {
            expect( +current[i].valid_until ).toBeGreaterThanOrEqual( +previous[i].valid_until );

            previous[i].valid_until = current[i].valid_until;
         };
      };

      fixTransients( nextGen.app_state.starname.domains, genesis.app_state.starname.domains );
      fixTransients( nextGen.app_state.starname.accounts, genesis.app_state.starname.accounts );

      // ...comparing
      compareObjects( nextGen, genesis );

      tmpobj.removeCallback();
   } );
} );
