package fee

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/fee/keeper"
	"github.com/iov-one/iovns/x/fee/types"
)

// GenesisState is used to unmarshal the genesis state
// when the app is initialized, and it is used to marshal
// the state when it needs to be exported.
type GenesisState struct {
	Fees *types.Fees `json:"fees"`
}

// NewGenesisState is GenesisState constructor
func NewGenesisState(fees *types.Fees) GenesisState {
	return GenesisState{
		Fees: fees,
	}
}

// ValidateGenesis makes sure that the genesis state is valid
func ValidateGenesis(data GenesisState) error {
	if err := data.Fees.Validate(); err != nil {
		return err
	}
	return nil
}

// DefaultGenesisState returns the default genesis state
// TODO this needs to be updated, although it will be imported from iovns chain
func DefaultGenesisState() GenesisState {
	// set fees
	// add domain module fees
	feeCoinDenom := "tiov" // set coin denom used for fees
	// generate new fees
	fees := types.NewFees()
	// set default fees
	fees.SetDefaults(feeCoinDenom)
	// return genesis
	return GenesisState{
		Fees: fees,
	}
}

// InitGenesis sets the initial state of the configuration module
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data GenesisState) {
	k.SetFees(ctx, data.Fees)
}

// ExportGenesis saves the state of the configuration module
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) GenesisState {
	return GenesisState{
		Fees: k.GetFees(ctx),
	}
}
