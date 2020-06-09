package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/fatih/structs"
)

func NewFees() *Fees {
	return &Fees{}
}

// Fees contains different type of fees
// to calculate coins to detract when
// processing different messages
type Fees struct {
	IovTokenPrice sdk.Dec

	DefaultFee sdk.Coin
	// account fees
	RegisterClosedAccount sdk.Coin
	RegisterOpenAccount   sdk.Coin
	TransferClosedAccount sdk.Coin
	TransferOpenAccount   sdk.Coin
	ReplaceAccountTargets sdk.Coin
	AddAccountCertificate sdk.Coin
	DelAccountCertificate sdk.Coin
	SetAccountMetadata    sdk.Coin
	// domain fees
	// Register domain
	RegisterDomain               sdk.Coin
	RegisterDomain1              sdk.Coin
	RegisterDomain2              sdk.Coin
	RegisterDomain3              sdk.Coin
	RegisterDomain4              sdk.Coin
	RegisterDomain5              sdk.Coin
	RegisterDomainDefault        sdk.Coin
	RegisterOpenDomainMultiplier sdk.Int
	// TransferDomain
	TransferDomainClosed sdk.Coin
	TransferDomainOpen   sdk.Coin
	// RenewDomain
	RenewOpenDomain sdk.Coin
}

func (f *Fees) Validate() error {
	if f == nil {
		return fmt.Errorf("fees is nil")
	}
	m := structs.New(f)
	for _, field := range m.Fields() {
		switch fee := field.Value().(type) {
		case sdk.Dec:
			if !fee.IsNil() {
				return fmt.Errorf("invalid dec: %s", field.Name())
			}
		case sdk.Coin:
			if fee.Amount.IsZero() {
				return fmt.Errorf("invalid dec coin: %s", field.Name())
			}
			if !fee.IsValid() {
				return fmt.Errorf("invalid dec coin: %s", field.Name())
			}
		case sdk.Int:
			if fee.IsZero() {
				return fmt.Errorf("invalid int multiplier: %s", field.Name())
			}
		default:
			panic(fmt.Sprintf("invalid type: %T", fee))
		}
	}
	return nil
}

// SetDefaults sets the default fees
func (f *Fees) SetDefaults(coin sdk.Coin) {
	if f == nil {
		panic("cannot set default fees for nil fees")
	}
	*f = Fees{
		IovTokenPrice:                sdk.NewDec(10),
		DefaultFee:                   coin,
		RegisterClosedAccount:        coin,
		RegisterOpenAccount:          coin,
		TransferClosedAccount:        coin,
		TransferOpenAccount:          coin,
		ReplaceAccountTargets:        coin,
		AddAccountCertificate:        coin,
		DelAccountCertificate:        coin,
		SetAccountMetadata:           coin,
		RegisterDomain:               coin,
		RegisterDomain1:              coin,
		RegisterDomain2:              coin,
		RegisterDomain3:              coin,
		RegisterDomain4:              coin,
		RegisterDomain5:              coin,
		RegisterDomainDefault:        coin,
		RegisterOpenDomainMultiplier: sdk.NewInt(2),
		TransferDomainClosed:         coin,
		TransferDomainOpen:           coin,
		RenewOpenDomain:              coin,
	}
}
