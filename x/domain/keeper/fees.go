package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/iov-one/iovns/x/domain/types"
	fee2 "github.com/iov-one/iovns/x/fee"
	types2 "github.com/iov-one/iovns/x/fee/types"
)

// CollectFees collects the fees of a msg and sends them
// to the distribution module to validators and stakers
func (k Keeper) CollectFees(ctx sdk.Context, msg types2.ProductMsg, domain types.Domain) error {
	moduleFees := k.FeeKeeper.GetFees(ctx)
	// create fee calculator
	Warning here!
	calculator := fee2.NewCalculator(ctx, k, moduleFees, domain)
	// get fee
	fee := calculator.GetFee(msg)
	// transfer fee to distribution
	return k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.FeePayer(), auth.FeeCollectorName, sdk.NewCoins(fee))
}
