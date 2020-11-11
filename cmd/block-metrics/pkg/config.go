package pkg

type Configuration struct {
	// database
	DBHost string
	DBName string
	DBSSL  string
	// Read-write user
	DBUser string
	DBPass string
	// Read-only user
	DBROUser string
	DBROPass string
	// Denomination of the fee coin, eg uiov
	FeeDenom string
	// Tendermint light client daemon URL
	TendermintLcdUrl string
	// Tendermint websocket URI
	TendermintWsURI string
}
