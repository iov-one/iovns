# Hard-Fork To iov-mainnet-2

The hard-fork from weave to cosmos-sdk marks a new era for IOV and starnames.  This document aims to outline how to make the fork as smooth as possible.  It is intended for validators on the legacy chain.

Validator nodes on the legacy chain do not maintain the full history of the chain due to pruning.  That, coupled with the fact that blocks are created on-demand via transactions, means that some coordination is required in order for validators to be able to verify the exported state from the legacy chain.  In order to facilitate that coordination, IOV will add a validator to the legacy validator set that has ⅔+ `voting_power` so that when it is stopped then the legacy chain will be halted.  That will allow all validators to export state locally and use it to verify the **iov-mainnet-2** genesis file.

IOV will halt the legacy chain at least an hour before the new chain's `genesis_time`.  That may not seem like a lot of time to verify the new genesis file but, as you'll see, verifying the file only takes a matter of seconds.

The technical procedure for verifying the genesis file for **iov-mainnet-2** is [here](VERIFY.md).  Wait for IOV to announce that it has generated the genesis file before attempting the procedure.  Otherwise you'd be comparing your local genesis file to one that is out-of-date.  Announce on Telegram in the **IOV Validators** channel whether or not you were able to replicate the **iov-mainnet-2** genesis file.  When all validators are good to go then announce on Twitter, too :).  If there's a problem then we will debug and potentially restart the legacy chain and repeat the technical procedure until consensus on the genesis file is achieved.

Once consensus on the **iov-mainnet-2** genesis file is achieved then follow the procedure at https://docs.iov.one/mainnet and wait for `genesis_time`.  Note that https://docs.iov.one/mainnet builds on https://docs.iov.one/for-validators/testnet, so it'd be good to be familiar with that or even join the testnet before the launch of **iov-mainnet-2**.

Once the new chain is started then feel free to release all the resources that the legacy chain used.


## Proof Of Ownership ##

There will undoubtedly be some legacy token holders that failed to provide IOV with a star1 address for the new chain.  All tokens and *iov names owned by them will be in the custody of IOV on the new chain.  In order to claim their property the token/name holders will have to prove their ownership.  In order to do so, IOV will restart its validator on the legacy chain and require that owners map their iov1 address to their star1 address.  Other validators are not required to operate on the legacy chain since the IOV validator will have ⅔+ `voting_power`.

## Questions? ##

If you have any questions/concerns/suggestions then please post them in Telegram in the **IOV Validators** channel.
