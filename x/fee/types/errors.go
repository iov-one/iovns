package types

import sdk "github.com/cosmos/cosmos-sdk/types/errors"

var ErrDefaultFee = sdk.Register(ModuleName, 1, "default fee")
