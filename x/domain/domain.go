package domain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// Domain defines a domain
type Domain struct {
	// Name is the name of the domain
	Name string
	// Admin is the owner of the domain
	Admin sdk.AccAddress
	// ValidUntil is a unix timestamp that defines for how long the domain is valid
	ValidUntil int64
	// HasSuperuser checks if the domain is owned by a super user or not
	HasSuperuser bool
	// AccountRenew defines the duration of each created or renewed account
	// under the domain
	AccountRenew time.Duration
	// Broker TODO needs comment
	Broker sdk.AccAddress
}

func (d Domain) String() string {
	panic("implement plz")
}
