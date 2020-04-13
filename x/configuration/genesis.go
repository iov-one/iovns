package configuration

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration/types"
	"regexp"
)

// GenesisState is used to unmarshal the genesis state
// when the app is initialized, and it is used to marshal
// the state when it needs to be exported.
type GenesisState struct {
	Config types.Config `json:"config"`
}

// NewGenesisState is GenesisState constructor
func NewGenesisState(conf types.Config) GenesisState {
	return GenesisState{
		Config: conf,
	}
}

// ValidateGenesis makes sure that the genesis state is valid
func ValidateGenesis(data GenesisState) error {
	conf := data.Config
	if conf.DomainRenew < 0 {
		return fmt.Errorf("empty domain renew")
	}
	if conf.Owner == nil {
		return fmt.Errorf("empty owner")
	}
	if _, err := regexp.Compile(conf.ValidBlockchainAddress); err != nil {
		return fmt.Errorf("empty valid blockchain address regexp: %w", err)
	}
	if _, err := regexp.Compile(conf.ValidBlockchainID); err != nil {
		return fmt.Errorf("empty valid blockchain id regexp: %w", err)
	}
	if _, err := regexp.Compile(conf.ValidDomain); err != nil {
		return fmt.Errorf("empty valid domain regexp: %w", err)
	}
	if conf.ValidName == "" {
		return fmt.Errorf("empty valid name regexp")
	}
	return nil
}

// DefaultGenesisState returns the default genesis state
// TODO this needs to be updated, although it will be imported from iovns chain
func DefaultGenesisState() GenesisState {
	return GenesisState{Config: types.Config{
		Owner:                  sdk.AccAddress{},
		ValidDomain:            "^(.*?)?",
		ValidName:              "^(.*?)?",
		ValidBlockchainID:      "^(.*?)?",
		ValidBlockchainAddress: "^(.*?)?",
		DomainRenew:            0,
	}}
}

// InitGenesis sets the initial state of the configuration module
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetConfig(ctx, data.Config)
}

// ExportGenesis saves the state of the configuration module
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{Config: k.GetConfiguration(ctx)}
}
