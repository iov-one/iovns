package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Account defines an account that belongs to a domain
type Account struct {
	// Domain references the domain this account belongs to
	Domain string
	// Name is the name of the account
	Name string
	// Owner is the address that owns the account
	Owner sdk.AccAddress
	// ValidUntil defines a unix timestamp of the expiration of the account
	ValidUntil int64
	// Targets is the list of blockchain addresses this account belongs to
	Targets [][]byte
	// Certificates contains the list of certificates to identify the account owner
	Certificates [][]byte
	// Broker can be empty
	// it identifies an entity that facilitated the transaction of the account
	Broker sdk.AccAddress
}
