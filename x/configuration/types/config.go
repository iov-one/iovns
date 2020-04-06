package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Config is the configuration of the network
type Config struct {
	// Owner is the configuration owner, the address allowed to register no super user domains
	Owner sdk.AccAddress
	// ValidDomain defines a regexp that determines if a domain name is valid or not
	ValidDomain string
	// ValidName defines a regexp that determines if an account name is valid or not
	ValidName string
	// ValidBlockchainID defines a regexp that determines if a blockchain id is valid or not
	ValidBlockchainID string
	// ValidBlockchainAddress determines a regexp for a valid blockchain address
	ValidBlockchainAddress string
	// DomainRenew defines the duration of the domain renewal period in seconds
	DomainRenew int64
}
