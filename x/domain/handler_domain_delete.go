package domain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd/x/domain/keeper"
	"github.com/iov-one/iovnsd/x/domain/types"
)

func handleMsgDomainDelete(ctx sdk.Context, k keeper.Keeper, msg types.MsgDeleteDomain) (*sdk.Result, error) {
	panic("to implement")
}
