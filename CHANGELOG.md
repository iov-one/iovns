# Changelog 

## HEAD
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