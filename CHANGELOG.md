# Changelog

## HEAD

## v0.9.4

- CLI: Fix output on verify
- REST: Fix swagger documentation

## v0.9.3

- add signature module

## v0.9.2

- update docs

## v0.9.1

- upgrade to cosmos sdk 0.39.1
- fix empty account renewal
- use external crud package
- make controllers API nicer
- enhance executors tests
- allow open domain transfers
- fix filtering when primary key is present
- fix finding the smallest set in the filters
- change enhance crud store api and make it more secure

## v0.9.0
- add genesis file generation scripts
- upgrade to cosmos sdk 0.39
- move all files from root module to respective packages
- fix transfer domain and add tests
- fix empty account name ambiguity
- BREAKING: use cosmos-sdk v0.39 (Launchpad)
- fix export genesis function
- fix domain renewal disallowed after grace period
- BREAKING: change Account.name type to pointer
- BREAKING: refactor domain and account keeper by dividing in two
- BREAKING: remove index package
- add crud store
- add helm chart

## v0.4.5

- CHANGE: bump cosmos-sdk and tendermint version
- FIX: AccountRenewalCountMax and DomainRenewalCountMax bumped at configuration update
- FIX: fix cli tests
- FIX: fix domain renewal
- Implement block-metrics
- disable block metrics CI

## v0.4.4

- fix: iovnscli get config
- REST: rename /domain/ query path to /starname/
- REST: rename FromOwner to WithOwner
- CLI: Add broker field to registerDomain and registerAccount
- FIX: fix fee deduct and improve tests
- FIX: TransferDomain flushes empty account content
- Implement faucet

## v0.4.3

- Add configuration module rest features
- Sync swagger ui with recent changes
- Rename resolve-domain to domain-info, resolve-account to resolve
- Alias iovnscli domain to iovnscli starname
- Change account target blockchain id in genesis to blockhain_id

## v0.4.2

- Add cli tests
- Fix domain query responses
- Resolve account by starname functionality
- Normalize fee parameter names

## v0.4.1

- Enable gitian on travis builds

### Breaking changes
- rename targets to resource
- Implement fee payer functionality

## v0.4.0

- change reconciliate with new fee calculator spec
- fix multisig message length
- add ledger support
- Integrate gitian builds
- Remove account renew field in types.Domain
- Improve json field names in msgs
- Improve iovnscli tx add-certs error handling
- Fix delete domain handler
- Enable empty account queries

## v0.3.0
- fix account controller max renew exceed
- add tests to account handlers in domain
- reconcile domain spec
- treat handler as orchestrator
- extend keeper functionality
- move all errors/authorization checks to handlers
- introduce domain and account controllers
- upgrade cosmos-sdk to v0.38.4
- iovnsd: fix fee colletor address
- iovnscli: certificates accepted in base64 json

### Breaking changes

- change hasSuperUser to DomainType
- Open domain's admin is changed from zero address to normal address
- Recon configuration
- Recon register account handler
- Recon transfer account handler
- Recon delete account handler
- Recon replace account targets handler
- Recon add account certs handler
- Recon delete account cert handler
- Recon replace metadata handler

## v0.2.5

- iovnscli: fix has-superuser bool flag bug
- iovnsd: fix duplicate blockchain targets ID
- remove flush domain feature

## v0.2.4

- domain grace period is time.duration now
- refactor configuration module to be used with multisig wallets
- allow empty account name on msg.Validate()
- enable fees for all domain module handlers
- remove certificate indexing
- fix account transfer
- fix account renew timestamp
- fix account store keys that end up reading contents of other accounts
- add logging on panics

### Breaking changes

- configuration struct signature in genesis file changed

## 0.2.3

- fix add signers in msg renew account and renew domain
- add resolve certificates
- add resolve blockchain targets
- add generalized indexing strategy
- abstract indexing
- iovnscli: accept certificate as file

### Breaking changes

- Change naming of some json keys in genesis.json
- change move blockchain address from iovns types to domain module types

## 0.2.2

- implement iovns lite client swagger

## 0.2.1

- fix properly export genesis state
- fix properly init genesis state from old state
- fix iterate all domains
- add iterate all accounts
- change one store key for domain
- change shorter indexing keys
- add MsgSetAccountURI: handlers, cli tx, rest tx
- remove panic if fees are missing

## 0.2.0

- change path prefix to star
