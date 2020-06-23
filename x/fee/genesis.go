package fee

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/fatih/structs"
	"github.com/iov-one/iovns/x/fee/keeper"
	"github.com/iov-one/iovns/x/fee/types"
)

// GenesisState is used to unmarshal the genesis state
// when the app is initialized, and it is used to marshal
// the state when it needs to be exported.
type GenesisState struct {
	Fees *types.FeeConfiguration `json:"fee"`
}

// NewGenesisState is GenesisState constructor
func NewGenesisState(fees *types.FeeConfiguration) GenesisState {
	return GenesisState{
		Fees: fees,
	}
}

// ValidateGenesis makes sure that the genesis state is valid
func ValidateGenesis(data GenesisState) error {
	if err := data.Fees.Validate(); err != nil {
		return err
	}
	sds := structs.Map(data.Fees.FeeSeeds)
	for _, a := range types.ContractFeeSeeds {
		if _, ok := sds[a.ID]; !ok {
			return fmt.Errorf("contract fee seed %s does not exist in genesis file", a.ID)
		}
	}

	return nil
}

// DefaultGenesisState returns the default genesis state
func DefaultGenesisState() GenesisState {
	feeCoinDenom := "tiov" // set coin denom used for fees
	// generate new fees
	fees := types.NewFeeConfiguration()
	// set default fees
	fees.SetDefaults(feeCoinDenom)
	// return genesis
	return GenesisState{
		Fees: fees,
	}
}

// InitGenesis sets the initial state of the configuration module
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data GenesisState) {
	k.SetFeeConfiguration(ctx, *data.Fees)
}

// ExportGenesis saves the state of the configuration module
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) GenesisState {
	cfg := k.GetFeeConfiguration(ctx)
	return GenesisState{
		Fees: &cfg,
	}
}
