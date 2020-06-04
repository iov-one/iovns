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
         configuration: {
            config: {
              account_grace_period: String( 30 * 24 * 60 * 60 ),
              account_renew_count_max: 2,
              account_renew_period: String( 365.25 * 24 * 60 * 60 ),
              blockchain_target_max: 16,
              certificate_count_max: 16,
              certificate_size_max: "1024",
              configurer: "star1 IOV SAS", // TODO
              domain_grace_period: String( 30 * 24 * 60 * 60 ),
              domain_renew_count_max: 2,
              domain_renew_period: String( 365.25 * 24 * 60 * 60 ),
              metadata_size_max: "1024",
              valid_account_name: "[-_\\.a-z0-9]{0,64}$",
              valid_blockchain_address: "^[a-z0-9A-Z]+$",
              valid_blockchain_id: "[-a-z0-9A-Z:]+$",
              valid_domain_name: "^[-_a-z0-9]{4,16}$",
            },
            fees: {
               AddCertificate: "100 / 1",
               DefaultFee: "1 / 2", // the fee for messages that don't explicitly have a fee
               IovTokenPrice: "4 / 10", // price in euros; manually updated
               RegisterClosedAccount: "1 / 2",
               RegisterDomain_1: "10000 / 1", // domain name with 1 char
               RegisterDomain_2: "5000 / 1",
               RegisterDomain_3: "2000 / 1",
               RegisterDomain_4: "1000 / 1",
               RegisterDomain_5: "500 / 1",
               "RegisterDomain_6+": "250 / 1", // domain name with 6 or more chars
               RegisterOpenAccount: "10 / 1",
               RegisterOpenDomainMultiplier: "10 / 1",
               RenewOpenDomain: "12345 / 1",
               ReplaceTargets: "10 / 1",
               SetMetaData: "500 / 3",
               TransferCloseAccount: "10 / 1",
               TransferClosedDomain: "10 / 1",
               TransferOpenAccount: "10 / 1",
               TransferOpenDomain: "10 / 1",
            }
         },
         domain: {
            domains: [],
            accounts: []
         },
         distribution: {
            fee_pool: {
               community_pool: []
            },
            community_tax: "0.000000000000000000",
            base_proposer_reward: "0.050000000000000000",
            bonus_proposer_reward: "0.050000000000000000",
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
               unbonding_time: "259200000000000",
               max_validators: 16,
               max_entries: 7,
               bond_denom: "iov"
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
               mint_denom: "iov"
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
      crisis: {
         constant_fee: {
            denom: "iov",
            amount: "1000"
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
                  denom: "iov",
                  amount: "1000"
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
   }

   fs.writeFileSync( "genesis.json", stringify( genesis, { space: "  " } ) );
}


main().catch( e => {
   console.error( e );
   process.exit( -1 );
} );
