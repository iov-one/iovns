import fetchIndicativeSendsTo from "./lib/fetchIndicativeSendsTo";
import fs from "fs";
import stringify from "json-stable-stringify";


const main = async () => {
   const chain_id = "iov-mainnet2";
   const genesis_time = new Date( "2020-04-15T10:00:00Z" ).toISOString();
   const ledgers = await fetchIndicativeSendsTo( "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un", /./ ).catch( e => { throw e } );  // TODO: remove this placeholder
   const genesis = { // TODO: remove placeholders
      chain_id: chain_id,
      genesis_time: genesis_time,
      accounts: [
         ...ledgers,
      ],
      app_hash: "",
      app_state: {
         bank: {
            send_enabled: true
         },
         distribution: {
            fee_pool: {
               community_pool: []
            },
            community_tax: "0.020000000000000000",
            base_proposer_reward: "0.010000000000000000",
            bonus_proposer_reward: "0.040000000000000000",
            withdraw_addr_enabled: true,
            delegator_withdraw_infos: [],
            previous_proposer: "",
            outstanding_rewards: [],
            validator_accumulated_commissions: [],
            validator_historical_rewards: [],
            validator_current_rewards: [],
            delegator_starting_infos: [],
            validator_slash_events: []
         },
         slashing: {
            params: {
               downtime_jail_duration: "600000000000",
               max_evidence_age: "120000000000",
               min_signed_per_window: "0.050000000000000000",
               signed_blocks_window: "5000",
               slash_fraction_double_sign: "0.050000000000000000",
               slash_fraction_downtime: "0.000000000000000000"
            },
            signing_infos: {},
            missed_blocks: {}
         },
         staking: {
            params: {
               unbonding_time: "259200000000000",
               max_validators: 100,
               max_entries: 7,
               bond_denom: "umuon"
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
               inflation: "0.130000000000000000",
               annual_provisions: "0.000000000000000000"
            },
            params: {
               blocks_per_year: "4855015",
               goal_bonded: "0.670000000000000000",
               inflation_max: "0.200000000000000000",
               inflation_min: "0.070000000000000000",
               inflation_rate_change: "0.130000000000000000",
               mint_denom: "umuon"
            }
         },
         params: null,
      },
      auth: {
         params: {
            max_memo_characters: "256",
            sig_verify_cost_ed25519: "590",
            sig_verify_cost_secp256k1: "1000",
            tx_sig_limit: "7",
            tx_size_cost_per_byte: "10"
         }
      },
      consensus_params: {
         block: {
            max_bytes: "500000",
            max_gas: "5000000",
            time_iota_ms: "1000"
         },
         evidence: {
            max_age: "100000"
         },
         validator: {
            pub_key_types: [
               "ed25519"
            ]
         }
      },
      crisis: {
         constant_fee: {
            denom: "muon",
            amount: "100000"
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
                  denom: "umuon",
                  amount: "100000"
               }
            ],
            max_deposit_period: "172800000000000"
         },
         voting_params: {
            voting_period: "172800000000000"
         },
         tally_params: {
            quorum: "0.334000000000000000",
            threshold: "0.500000000000000000",
            veto: "0.334000000000000000"
         }
      },
   }

   fs.writeFileSync( "genesis.json", stringify( genesis, { space: "  " } ) );
}


main().catch( e => {
   console.error( e );
   process.exit( -1 );
} );
