module github.com/iov-one/iovns

go 1.14

require (
	github.com/cosmos/cosmos-sdk v0.38.5
	github.com/fatih/structs v1.1.0
	github.com/golang/mock v1.3.1 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.2
	github.com/influxdata/influxdb1-client v0.0.0-20191209144304-8bf82d3c094d // indirect
	github.com/jinzhu/gorm v1.9.14
	github.com/lib/pq v1.1.1
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/common v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.6
	github.com/tendermint/tm-db v0.5.1
)

// patch bech32 decoding, TODO this should not be the fix for the multisig isue
replace github.com/btcsuite/btcutil => github.com/kava-labs/btcutil v0.0.0-20200522184203-886d33430f06
