import { chainIds, multisigs, source2multisig } from "./lib/constants";
import { migrate, patchGalaxynet, patchMainnet, } from "./lib/migrate";
import fetchIndicativeSendsTo from "./lib/fetchIndicativeSendsTo";
import readOsakaGenesisFile from "./lib/readOsakaGenesisFile";
import filterReserveds from "./lib/filterReserveds";
import path from "path";
import pullDumpedState from "./lib/pullDumpedState";
import pullPremiums from "./lib/pullPremiums";


const main = async () => {
   // network dependent constants
   const mainnet = process.argv[2].indexOf( "mainnet" ) != -1;
   const chain_id = mainnet ? "iov-mainnet-2" : "iovns-galaxynet";
   const home = mainnet ? path.join( __dirname, "data", chain_id ) : path.join( __dirname, "data", "galaxynet" );
   const gentxs = mainnet ? path.join( __dirname, "data", chain_id, "gentxs" ) : path.join( __dirname, "data", "galaxynet", "gentxs" );
   const patch = mainnet ? patchMainnet : patchGalaxynet;

   // genesis file scaffolding
   const genesis = {
      chain_id: chain_id,
      genesis_time: new Date( "2020-08-26T08:00:00Z" ).toISOString(),
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
              "//note duration": "all the durations are in nanoseconds",
              account_grace_period: 30 * 24 * 60 * 60 + "000000000", // (ab)use javascript
              account_renew_count_max: 2,
              account_renew_period: 365.25 * 24 * 60 * 60 + "000000000",
              resources_max: 10, // https://internetofvalues.slack.com/archives/GPYCU2AJJ/p1592563251000500
              certificate_count_max: 5,
              certificate_size_max: "1024",
              configurer: "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78",
              domain_grace_period: 30 * 24 * 60 * 60 + "000000000",
              domain_renew_count_max: 2,
              domain_renew_period: 365.25 * 24 * 60 * 60 + "000000000",
              metadata_size_max: "1024",
              valid_account_name: "^[-_.a-z0-9]{1,64}$",
              valid_resource: "^[a-z0-9A-Z]+$",
              valid_uri: "^[-a-z0-9A-Z:]+$",
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
            }
         },
         crisis: {
            constant_fee: {
               denom: "uiov",
               amount: "1000000000"
            }
         },
         starname: {
            domains: [
               {
                  account_renew: "31557600",
                  admin: multisigs.iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq.star1,
                  broker: null,
                  name: "iov",
                  type: "open",
                  valid_until: "1924991999"
                },
            ],
            accounts: []
         },
         distribution: {
            fee_pool: {
               community_pool: []
            },
            params: {
               community_tax: "0.000000000000000000",
               base_proposer_reward: "0.050000000000000000",
               bonus_proposer_reward: "0.050000000000000000",
               withdraw_addr_enabled: true,
            },
            delegator_withdraw_infos: [],
            previous_proposer: "",
            outstanding_rewards: [],
            validator_accumulated_commissions: [],
            validator_historical_rewards: [],
            validator_current_rewards: [],
            delegator_starting_infos: [],
            validator_slash_events: []
         },
         evidence: {
            evidence: [
            ],
            params: {
               max_evidence_age: "1814400000000000"
            }
         },
         genutil: {
            gentxs: [
            ]
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
         slashing: {
            params: {
               downtime_jail_duration: "600000000000",
               max_evidence_age: "1814400000000000",
               min_signed_per_window: "0.500000000000000000",
               signed_blocks_window: "10000",
               slash_fraction_double_sign: "0.050000000000000000",
               slash_fraction_downtime: "0.010000000000000000"
            },
            signing_infos: {},
            missed_blocks: {}
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
         mint: {
            minter: {
               inflation: "0.12",
               annual_provisions: "0.000000000000000000"
            },
            params: {
               blocks_per_year: "5259600", // assume 6 seconds per block since create_empty_blocks=false is broken
               "//note": "goal_bonded cannot be 0: module=consensus err='division by zero'",
               goal_bonded: "0.8",
               inflation_max: "0.25",
               inflation_min: "0.12",
               inflation_rate_change: "0.13",
               mint_denom: "uiov"
            }
         },
         params: null,
         upgrade: {
         }
      },
      consensus_params: {
         block: {
            "max_bytes": "500000",
            "max_gas": "-1",
            "time_iota_ms": "1000"
         },
         evidence: {
            "max_age_num_blocks": "100000",
            "max_age_duration": "172800000000000"
         },
         validator: {
            pub_key_types: [
               "ed25519"
            ]
         }
      },
   }

   // other data
   const dumped = await pullDumpedState().catch( e => { throw e } );
   const flammable = [ "iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u" ]; // accounts to burn; "pending deals" tokens were effectively burned by sending to this 0x0 hex account
   const indicatives = await fetchIndicativeSendsTo( "iov10v69k57z2v0pr3yvtr60pp8g2jx8tdd7f55sv6", /(star1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{38})/ ).catch( e => { throw e } );
   const osaka = await readOsakaGenesisFile().catch( e => { throw e } );
   const premiums = await pullPremiums().catch( e => { throw e } );
   const reserveds = filterReserveds( genesis.app_state.configuration.config.valid_domain_name );

   // migration
   migrate( { chainIds, dumped, flammable, genesis, gentxs, home, indicatives, multisigs, osaka, patch, premiums, reserveds, source2multisig } );
}


main().catch( e => {
   console.error( e.stack || e.message || e );
   process.exit( -1 );
} );
