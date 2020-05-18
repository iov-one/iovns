# Multisig wallets guide

## Create keys if you don't them ready

```shell script
iovnscli keys add w1
iovnscli keys add w2
iovnscli keys add w3
```

## Import addresses that does not exist locally

```shell script
iovnscli keys add \
  p1 \
  --pubkey=starpub1addwnpepqv80htam6gc7fudf9jseldx3afy8nu8anvk935qdctek0yr27jcqj4yv044
```

## Generate multisig wallet

Signature:
```shell script
iovnscli keys add --multisig=name1,name2,name3[...] --multisig-threshold=K new_key_name
```
In our example:
```shell script
iovnscli keys add --multisig=w1,w2,w3,p1 --multisig-threshold=3 msig1
```

## Generate transaction
```shell script
iovnscli tx send $(iovnscli keys show -a msig1) $(iovnscli keys show -a w1)  10iov  --generate-only > unsignedTx
.json
```
Note: `$(iovnscli keys show -a msig1)` returns the address of given account

Tx will be saved to unsignedTx.json. Participants of the multisig wallet will sign this json file.

Warning: When you create a wallet locally, it does not have an account number that is assigned by the network. 
To proceed to the next stages, you need to execute a tx with the wallet, most simple is sending some coins to it:
```shell script
iovnscli tx send $(iovnscli keys show w1 -a) $(iovnscli keys show msig1 -a)  10iov
```

## Signing the transaction

Using wallet `w1`:

```shell script
iovnscli tx sign unsignedTx.json --from=$(iovnscli keys show -a w1) --multisig=$(iovnscli keys show -a msig1) \
  --output-document=w1sig.json
```
Repeat this process for other participants.

## Combining the signatures

After required amount of signatures collected, you need to combine the sigs.
```shell script
iovnscli tx multisign unsignedTx.json msig1 w1sig.json w2sig.json w3sig.json > completeTx.json
``` 

## Broadcast the transaction
```shell script
iovnscli tx broadcast completeTx.json
```
