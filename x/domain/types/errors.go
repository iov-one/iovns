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
var ErrAccountDoesNotExist = sdkerrors.Register(ModuleName, 6, "account does not exist")
var ErrAccountExpired = sdkerrors.Register(ModuleName, 7, "account has expired")
var ErrInvalidOwner = sdkerrors.Register(ModuleName, 8, "invalid owner")
var ErrInvalidAccountName = sdkerrors.Register(ModuleName, 9, "invalid account name")
var ErrInvalidBlockchainTarget = sdkerrors.Register(ModuleName, 10, "blockchain target provided is not valid")
var ErrDomainExpired = sdkerrors.Register(ModuleName, 11, "domain has expired")
var ErrAccountExists = sdkerrors.Register(ModuleName, 12, "account already exists")
var ErrInvalidRequest = sdkerrors.Register(ModuleName, 13, "malformed request")
var ErrCertificateExists = sdkerrors.Register(ModuleName, 14, "certificate already exists")
