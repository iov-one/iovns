package pkg

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type Configuration struct {
	TendermintRPC string
	Port          string
	ChainID       string
	CoinDenom     string
	Armor         string
	Passphrase    string
	Memo          string
	SendAmount    int64
	GasPrices     string
	GasAdjust     float64
	KeyringPass   string
}

func env(name, fallback string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return fallback
}

const (
	gas           = "0"
	ga            = "0.2"
	fallBackArmor = `
-----BEGIN TENDERMINT PRIVATE KEY-----
salt: BF94D84D7E0BFEF9AB735D9315AD271E
type: secp256k1
kdf: bcrypt

PECa11ktJ6mV4iTnhHGIL9nhjdXjplQDt5n+o5nddvnmS613AWbCL5FrC3WErdCR
vdsyKdlue2uLJizP46Ao3w6PKMBVYgIkKe97GjA=
=MZbT
-----END TENDERMINT PRIVATE KEY-----
`
	faucetAddr = "star1pdp388k2jj5zsxx67v02pxtttguf6r4jj79v00"
)

func NewConfiguration() (*Configuration, error) {
	gasPrices := env("GAS_PRICES", "10.0uvoi")
	ga := env("GAS_ADJUST", ga)
	gasAdjust, err := strconv.ParseFloat(ga, 64)
	if err != nil {
		return nil, errors.Wrap(err, "GAS_ADJUST")
	}

	sendStr := env("SEND_AMOUNT", "100")
	send, err := strconv.Atoi(sendStr)
	return &Configuration{
		TendermintRPC: env("TENDERMINT_RPC", "http://localhost:26657"),
		Port:          env("PORT", ":8080"),
		ChainID:       env("CHAIN_ID", "local"),
		CoinDenom:     env("COIN_DENOM", "tiov"),
		Armor:         env("ARMOR", fallBackArmor),
		Passphrase:    env("PASSPHRASE", "12345678"),
		Memo:          "sent by IOV with love",
		SendAmount:    int64(send),
		GasPrices:     gasPrices,
		GasAdjust:     gasAdjust,
	}, nil
}
