import { chainIds, multisigs, source2multisig } from "./lib/constants";
import { migrate, patchGalaxynet, patchMainnet, } from "./lib/migrate";
import fetchIndicativeSendsTo from "./lib/fetchIndicativeSendsTo";
import fetchOsakaGenesisFile from "./lib/fetchOsakaGenesisFile";
import path from "path";
import pullDumpedState from "./lib/pullDumpedState";
import pullPremiums from "./lib/pullPremiums";


const main = async () => {
   // network dependent constants
   const mainnet = process.argv[2].indexOf( "mainnet" ) != -1;
   const chain_id = mainnet ? "iov-mainnet-2" : "iovns-galaxynet";
   const home = mainnet ? path.join( __dirname, "data", chain_id ) : path.join( __dirname, "data", "galaxynet" );
   const gentxs = mainnet ? path.join( __dirname, "data", "gentxs" ) : undefined;
   const patch = mainnet ? patchMainnet : patchGalaxynet;

   // genesis file scaffolding
   const genesis = {
      chain_id: chain_id,
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
              "//note duration": "all the durations are in nanoseconds",
              account_grace_period: 30 * 24 * 60 * 60 + "000000000", // (ab)use javascript
              account_renew_count_max: 2,
              account_renew_period: 365.25 * 24 * 60 * 60 + "000000000",
              blockchain_target_max: 16,
              certificate_count_max: 16,
              certificate_size_max: "1024",
              configurer: "star1 IOV SAS", // TODO
              domain_grace_period: 30 * 24 * 60 * 60 + "000000000",
              domain_renew_count_max: 2,
              domain_renew_period: 365.25 * 24 * 60 * 60 + "000000000",
              metadata_size_max: "1024",
              valid_account_name: "[-_\\.a-z0-9]{1,64}$",
              valid_blockchain_address: "^[a-z0-9A-Z]+$",
              valid_blockchain_id: "[-a-z0-9A-Z:]+$",
              valid_domain_name: "^[-_a-z0-9]{4,16}$",
            },
            fees: {
               AddAccountCertificate: "100",
               DefaultFee: "0.5", // the fee for messages that don't explicitly have a fee
               DelAccountCertificate: "0.5",
               FeeCoinDenom: "uiov",
               FeeCoinPrice: "0.1",
               RegisterClosedAccount: "0.5",
               RegisterDomain: "0.5", // TODO: drop this when https://github.com/iov-one/iovns/issues/197 is resolved
               RegisterDomain1: "10000", // domain name with 1 char
               RegisterDomain2: "5000",
               RegisterDomain3: "2000",
               RegisterDomain4: "1000",
               RegisterDomain5: "500",
               RegisterDomainDefault: "250", // domain name with 6 or more chars
               RegisterOpenAccount: "10",
               RegisterOpenDomainMultiplier: "5.5",
               RenewOpenDomain: "12345",
               ReplaceAccountTargets: "10",
               SetAccountMetadata: "500",
               TransferClosedAccount: "10",
               TransferDomainClosed: "10",
               TransferOpenAccount: "10",
               TransferDomainOpen: "10",
            }
         },
         crisis: {
            constant_fee: {
               denom: "uiov",
               amount: "1000000000"
            }
         },
         domain: {
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
               max_evidence_age: "120000000000"
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
   const osaka = await fetchOsakaGenesisFile().catch( e => { throw e } );
   const premiums = await pullPremiums().catch( e => { throw e } );
   const reserveds = []; // TODO

   // migration
   migrate( { chainIds, dumped, flammable, genesis, gentxs, home, indicatives, multisigs, osaka, patch, premiums, reserveds, source2multisig } );
}


main().catch( e => {
   console.error( e.stack || e.message || e );
   process.exit( -1 );
} );
