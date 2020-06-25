## How to export private key

```bash
iovnscli keys export faucet
```
- Input a passphrase. Remember you will use this pass phrase as env variable.
- Copy armor and use it in env variable:
```
-----BEGIN TENDERMINT PRIVATE KEY-----
kdf: bcrypt
salt: 93EFF493A3EA0A6AB71C00D69176AF19
type: secp256k1

l2zXyG3OOCXzzUxzmYv7Td1OFsc+vnCf7BckhUic8Y11KGCEm76fvRtdzlSwW0A5
fcz4CbdxSMEYktjtW5zyE+nLveB/UoJ3YK8Sbr4=
=Vhcg
-----END TENDERMINT PRIVATE KEY-----
```
## Enviroment variables
- GAS_PRICES
- GAS_ADJUST
- SEND_AMOUNT
- TENDERMINT_RPC
- PORT
- CHAIN_ID
- COIN_DENOM
- ARMOR
- PASSPHRASE

## How to use
`http://localhost:8080/credit?address=<bech32addr>`
