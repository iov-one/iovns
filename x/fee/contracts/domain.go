package contracts

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/fee/types"
)

// CONTRACT
// ContractFeeSeeds are seeds that will be used during genesis validation
var ContractFeeSeeds = []types.FeeSeed{
	{ID: "add_account_certificate", Amount: sdk.NewDec(0.5)},
	{ID: "del_account_certificate", Amount: sdk.NewDec(0.5)},
	{ID: "register_account_closed", Amount: sdk.NewDec(0.5)},
	{ID: "register_account_open", Amount: sdk.NewDec(10)},
	{ID: "register_domain_1", Amount: sdk.NewDec(10000)},
	{ID: "register_domain_2", Amount: sdk.NewDec(5000)},
	{ID: "register_domain_3", Amount: sdk.NewDec(2000)},
	{ID: "register_domain_4", Amount: sdk.NewDec(1000)},
	{ID: "register_domain_5", Amount: sdk.NewDec(500)},
	{ID: "register_domain_default", Amount: sdk.NewDec(250)},
	{ID: "register_domain_open", Amount: sdk.NewDec(12345)},
	{ID: "replace_account_targets", Amount: sdk.NewDec(10)},
	{ID: "set_account_metadata", Amount: sdk.NewDec(500)},
	{ID: "transfer_account_closed", Amount: sdk.NewDec(10)},
	{ID: "transfer_account_open", Amount: sdk.NewDec(10)},
	{ID: "transfer_domain_closed", Amount: sdk.NewDec(10)},
	{ID: "transfer_domain_open", Amount: sdk.NewDec(10)},
}
