package keeper

import (
	"log"

	"github.com/cosmos/cosmos-sdk/x/auth"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
)

// CollectFees collects the fees of a msg and sends them
// to the distribution module to validators and stakers
func (k Keeper) CollectFees(ctx sdk.Context, msg types.Feeable) error {
	var level int
	level = msg.CalculatedFee()
	feeConfig := k.ConfigurationKeeper.GetFees(ctx)
	fee, ok := feeConfig.CalculateLevelFees(msg, level)
	if !ok {
		// TODO we need to panic here
		log.Printf("WARNING unable to get expected fees for: %s/%s", types.ModuleName, msg.ID())
		return nil
	}
	// transfer fee to distribution
	return k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.FeePayer(), auth.FeeCollectorName, sdk.NewCoins(fee))
}
