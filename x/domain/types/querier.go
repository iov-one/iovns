package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// Query endpoints supported by the domain querier
const (
// TODO: Describe query parameters, update <action> with your query
// Query<Action>    = "<action>"
)

/*
Below you will be able how to set your own queries:


// QueryResList Queries Result Payload for a query
type QueryResList []string

// implement fmt.Stringer
func (n QueryResList) String() string {
	return strings.Join(n[:], "\n")
}

*/

type QueryResultDomain struct {
	Name         string         `json:"name"`
	Admin        sdk.AccAddress `json:"admin"`
	ValidUntil   int64          `json:"valid_until"`
	HasSuperuser bool           `json:"has_superuser"`
	AccountRenew time.Duration  `json:"account_renew"`
	Broker       sdk.AccAddress
}

func (r QueryResultDomain) String() string {
	return fmt.Sprintf("Domain<Name: %s; Admin: %s; ValidUntil: %s, HasSuperuser: %t, AccountRenew: %s, Broker: %s>",
		r.Name, r.Admin, time.Unix(r.ValidUntil, 0), r.HasSuperuser, r.AccountRenew, r.Broker)
}
