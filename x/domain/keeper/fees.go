package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/iov-one/iovns/x/domain/types"
	"log"
)

// CollectFees collects the fees of a msg and sends them
// to the distribution module to validators and stakers
func (k Keeper) CollectFees(ctx sdk.Context, msg sdk.Msg, addr sdk.AccAddress) error {
	var level int
	switch msg := msg.(type) {
	case *types.MsgTransferDomain:
		level = len(msg.Domain)
	case *types.MsgReplaceAccountTargets:
		level = len(msg.Domain)
	case *types.MsgAddAccountCertificates:
		level = len(msg.Domain)
	case *types.MsgDeleteAccountCertificate:
		level = len(msg.Domain)
	case *types.MsgDeleteAccount:
		level = len(msg.Domain)
	case *types.MsgDeleteDomain:
		level = len(msg.Domain)
	case *types.MsgFlushDomain:
		level = len(msg.Domain)
	case *types.MsgRegisterAccount:
		level = len(msg.Domain)
	case *types.MsgRegisterDomain:
		level = len(msg.Name)
	case *types.MsgRenewDomain:
		level = len(msg.Domain)
	case *types.MsgRenewAccount:
		level = len(msg.Domain)
	case *types.MsgTransferAccount:
		level = len(msg.Domain)
	}
	feeConfig := k.ConfigurationKeeper.GetFees(ctx)
	fee, ok := feeConfig.CalculateLevelFees(msg, level)
	if !ok {
		// TODO we need to panic here
		log.Printf("WARNING unable to get expected fees for: %s/%s", types.ModuleName, msg.Type())
		return nil
	}
	// transfer fee to distribution
	return k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, addr, distribution.ModuleName, sdk.NewCoins(fee))
}
