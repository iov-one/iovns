package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/iov-one/iovns/x/domain/controllers/fees"
	"github.com/iov-one/iovns/x/domain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CollectFees collects the fees of a msg and sends them
// to the distribution module to validators and stakers
func (k Keeper) CollectFees(ctx sdk.Context, msg sdk.Msg, addr sdk.AccAddress, domain types.Domain) error {
	moduleFees := k.ConfigurationKeeper.GetFees(ctx)
	// create fee calculator
	calculator := fees.NewController(ctx, k, moduleFees, domain)
	// get fee
	fee := calculator.GetFee(msg)
	// transfer fee to distribution
	return k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, addr, auth.FeeCollectorName, sdk.NewCoins(fee))
}
