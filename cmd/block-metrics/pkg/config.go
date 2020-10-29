package pkg

type Configuration struct {
	DBHost string
	DBUser string
	DBPass string
	DBName string
	DBSSL  string
	// Tendermint websocket URI
	TendermintWsURI string
	// Denomination of the fee coin, eg uiov
	FeeDenom string
	// Tendermint light client daemon URL
	TendermintLcdUrl string
}
