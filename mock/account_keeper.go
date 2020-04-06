package mock

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd/x/account"
)

type AccountKeeper struct {
	m map[string]account.Account
}

func (a AccountKeeper) SetAccount(ctx sdk.Context, account account.Account) {
	a.m[a.genKey(account)] = account
}

func (a AccountKeeper) genKey(account account.Account) string {
	return fmt.Sprintf("%s*%s", account.Domain, account.Name)
}
