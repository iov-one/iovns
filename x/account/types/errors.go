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
var ErrInvalidName = sdkerrors.Register(ModuleName, 4, "invalid account name")
var ErrInvalidBlockchainTarget = sdkerrors.Register(ModuleName, 5, "blockchain target provided is not valid")
var ErrUnauthorized = sdkerrors.Register(ModuleName, 6, "signer/s is/are not authorized to perform this action")
var ErrDomainExpired = sdkerrors.Register(ModuleName, 7, "domain has expired")
var ErrAccountExists = sdkerrors.Register(ModuleName, 8, "account already exists")
