package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
// var (
//	ErrInvalid = sdkerrors.Register(ModuleName, 1, "custom error message")
// )

var ErrInvalidDomainName = sdkerrors.Register(ModuleName, 1, "domain name provided is invalid")
var ErrDomainAlreadyExists = sdkerrors.Register(ModuleName, 2, "domain already exists")
var ErrUnauthorized = sdkerrors.Register(ModuleName, 3, "operation unauthorized")
var ErrInvalidRegisterDomainRequest = sdkerrors.Register(ModuleName, 4, "register domain request is not valid")
var ErrDomainDoesNotExist = sdkerrors.Register(ModuleName, 5, "domain does not exist")
