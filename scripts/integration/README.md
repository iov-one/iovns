# Jest driven tests

Jest is a powerful testing framework that can be used to test the APIs `iovnsd`.  The battery of tests can be selectively executed with the [only](https://jestjs.io/docs/en/api#testonlyname-fn-timeout) and [skip](https://jestjs.io/docs/en/api#testskipname-fn) functions, which is very handy.

# Setup

1. `pwd | grep -q 'easy-testnets/scripts' || echo 'Wrong directory; be in easy-testnets/scripts'`
1. `yarn install`
1. Add account `bojack` to your keys; the mnemonic is in keybase under `...team/iov_one.blockchain/credentials/test-wallets`.
1. Make sure that bojack has funds with `iovnscli q account star1z6rhjmdh2e9s6lvfzfwrh8a3kjuuy58y74l29t --output json | jq`.
1. Edit `.env` appropriately, ie point the `URL`s at the testnet that you want to test, etc.
1. `yarn test`

## Visual Studio Code

Running jest through VSCode's debugger makes life easier. A [launch.json](.vscode/launch.json) file exists so it should be possible to run the tests out-of-the box.
