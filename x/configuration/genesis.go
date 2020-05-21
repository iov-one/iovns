package configuration

import (
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration/types"
	domain_types "github.com/iov-one/iovns/x/domain/types"
)

// GenesisState is used to unmarshal the genesis state
// when the app is initialized, and it is used to marshal
// the state when it needs to be exported.
type GenesisState struct {
	Config types.Config `json:"config"`
	Fees   *types.Fees  `json:"fees"`
}

// NewGenesisState is GenesisState constructor
func NewGenesisState(conf types.Config, fees *types.Fees) GenesisState {
	return GenesisState{
		Config: conf,
		Fees:   fees,
	}
}

// ValidateGenesis makes sure that the genesis state is valid
func ValidateGenesis(data GenesisState) error {
	conf := data.Config
	if conf.DomainRenew < 0 {
		return fmt.Errorf("empty domain renew")
	}
	if conf.Configurer == nil {
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
	if data.Fees == nil {
		return fmt.Errorf("empty fees")
	}
	if data.Fees.LevelFees == nil {
		return fmt.Errorf("empty length fees")
	}
	if data.Fees.DefaultFees == nil {
		return fmt.Errorf("empty default fees")
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
		ValidDomain:            "^(.*?)?",
		ValidName:              "^(.*?)?",
		ValidBlockchainID:      "^(.*?)?",
		ValidBlockchainAddress: "^(.*?)?",
		DomainRenew:            86400,
	}
	// set fees
	fees := types.NewFees()
	defFee := sdk.NewCoin("iov", sdk.NewInt(10))
	// add domain fees
	fees.UpsertDefaultFees(&domain_types.MsgRegisterDomain{}, defFee)
	fees.UpsertDefaultFees(&domain_types.MsgAddAccountCertificates{}, defFee)
	fees.UpsertDefaultFees(&domain_types.MsgDeleteAccountCertificate{}, defFee)
	fees.UpsertDefaultFees(&domain_types.MsgDeleteDomain{}, defFee)
	fees.UpsertDefaultFees(&domain_types.MsgDeleteAccount{}, defFee)
	fees.UpsertDefaultFees(&domain_types.MsgRegisterAccount{}, defFee)
	fees.UpsertDefaultFees(&domain_types.MsgRenewAccount{}, defFee)
	fees.UpsertDefaultFees(&domain_types.MsgRenewDomain{}, defFee)
	fees.UpsertDefaultFees(&domain_types.MsgReplaceAccountTargets{}, defFee)
	fees.UpsertDefaultFees(&domain_types.MsgTransferAccount{}, defFee)
	fees.UpsertDefaultFees(&domain_types.MsgTransferDomain{}, defFee)
	fees.UpsertDefaultFees(&domain_types.MsgSetAccountMetadata{}, defFee)
	// return genesis
	return GenesisState{
		Config: config,
		Fees:   fees,
	}
}

// InitGenesis sets the initial state of the configuration module
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetConfig(ctx, data.Config)
	k.SetFees(ctx, data.Fees)
}

// ExportGenesis saves the state of the configuration module
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Config: k.GetConfiguration(ctx),
		Fees:   k.GetFees(ctx),
	}
}
