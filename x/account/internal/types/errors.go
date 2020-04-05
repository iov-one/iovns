package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
// var (
//	ErrInvalid = sdkerrors.Register(ModuleName, 1, "custom error message")
// )

var ErrAccountDoesNotExits = sdkerrors.Register(ModuleName, 1, "account does not exist")
var ErrInvalidDomain = sdkerrors.Register(ModuleName, 2, "invalid domain")
var ErrInvalidOwner = sdkerrors.Register(ModuleName, 3, "invalid owner")
