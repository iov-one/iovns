# Changelog 

## HEAD

## v0.2.4

- domain grace period is time.duration now 
- refactor configuration module to be used with multisig wallets
- allow empty account name on msg.Validate()
- enable fees for all domain module handlers
- remove certificate indexing
- fix account transfer
- fix account renew timestamp
- fix account store keys that end up reading contents of other accounts

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