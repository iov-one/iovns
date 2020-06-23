package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/iov-one/iovns/x/domain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CollectFees collects the fees of a msg and sends them
// to the distribution module to validators and stakers
func (k Keeper) CollectFees(ctx sdk.Context, msg types.MsgWithFeePayer, fee sdk.Coin) error {
	// transfer fee to distribution
	return k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.FeePayer(), auth.FeeCollectorName, sdk.NewCoins(fee))
}
