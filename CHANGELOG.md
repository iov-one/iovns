# Changelog 

## HEAD

- allow empty account name on msg.Validate()

Breaking changes

- change naming of some json keys in genesis.json

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