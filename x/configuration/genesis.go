package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration/types"
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
	if err := conf.Validate(); err != nil {
		return err
	}
	return nil
}

// DefaultGenesisState returns the default genesis state
// TODO this needs to be updated, although it will be imported from iovns chain
func DefaultGenesisState() GenesisState {
	// get owner
	owner, err := sdk.AccAddressFromBech32("star1kxqay5tndu3w08ps5c27pkrksnqqts0zxeprzx")
	if err != nil {
		panic("invalid default owner provided")
	}
	// set default configs
	config := types.Config{
		Configurer:             owner,
		ValidDomainName:        "^[-_a-z0-9]{4,16}$",
		ValidAccountName:       "[-_\\.a-z0-9]{1,64}$",
		ValidBlockchainID:      "[-a-z0-9A-Z:]+$",
		ValidBlockchainAddress: "^[a-z0-9A-Z]+$",
		DomainRenewalPeriod:    300000000000,
		DomainRenewalCountMax:  2,
		DomainGracePeriod:      60000000000,
		AccountRenewalPeriod:   180000000000,
		AccountRenewalCountMax: 3,
		AccountGracePeriod:     60000000000,
		BlockchainTargetMax:    3,
		CertificateSizeMax:     10000,
		CertificateCountMax:    3,
		MetadataSizeMax:        86400,
	}
	return GenesisState{
		Config: config,
	}
}

// InitGenesis sets the initial state of the configuration module
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetConfig(ctx, data.Config)
}

// ExportGenesis saves the state of the configuration module
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Config: k.GetConfiguration(ctx),
	}
}
