package pkg

type Configuration struct {
	DBHost string
	DBUser string
	DBPass string
	DBName string
	DBSSL  string
	// Tendermint websocket URI
	TendermintWsURI string
	// Derivation path: "tiov" or "iov"
	Hrp   string
	Denom string
}
