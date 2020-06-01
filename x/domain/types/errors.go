package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
// var (
//	ErrInvalid = sdkerrors.Register(ModuleName, 1, "custom error message")
// )

// ErrInvalidDomainName is returned when the domain name does not match the required standards
var ErrInvalidDomainName = sdkerrors.Register(ModuleName, 1, "domain name provided is invalid")

// ErrDomainAlreadyExists is returned when a create action is done on a domain that already exists
var ErrDomainAlreadyExists = sdkerrors.Register(ModuleName, 2, "domain already exists")

// ErrUnauthorized is returned when authentication process for an action fails
var ErrUnauthorized = sdkerrors.Register(ModuleName, 3, "operation unauthorized")

// ErrDomainDoesNotExist is returned when an action is performed on a domain that does not exist
var ErrDomainDoesNotExist = sdkerrors.Register(ModuleName, 5, "domain does not exist")

// ErrAccountDoesNotExist is returned when an action is performed on a domain that does not contain the specified account
var ErrAccountDoesNotExist = sdkerrors.Register(ModuleName, 6, "account does not exist")

// ErrAccountExpired is returned when actions are performed on expired accounts
var ErrAccountExpired = sdkerrors.Register(ModuleName, 7, "account has expired")

// ErrInvalidOwner is returned when the owner address provided is not valid (empty, malformed, etc)
var ErrInvalidOwner = sdkerrors.Register(ModuleName, 8, "invalid owner")

// ErrInvalidAccountName is returned when the account name does not match the required standards
var ErrInvalidAccountName = sdkerrors.Register(ModuleName, 9, "invalid account name")

// ErrInvalidBlockchainTarget is returned when provided blockchain target is not valid
var ErrInvalidBlockchainTarget = sdkerrors.Register(ModuleName, 10, "blockchain target provided is not valid")

// ErrDomainExpired is returned when actions are performed on expired domains
var ErrDomainExpired = sdkerrors.Register(ModuleName, 11, "domain has expired")

// ErrDomainExpired is returned when actions are performed on not expired domains
var ErrDomainNotExpired = sdkerrors.Register(ModuleName, 12, "domain has not expired")

// ErrAccountExists is returned when a create action is done on an account that already exists
var ErrAccountExists = sdkerrors.Register(ModuleName, 13, "account already exists")

// ErrInvalidRequest is a general error that covers the uncommon cases of invalid request
var ErrInvalidRequest = sdkerrors.Register(ModuleName, 14, "malformed request")

// ErrCertificateExists is returned when a creation action is done on a certificate that already exists
var ErrCertificateExists = sdkerrors.Register(ModuleName, 15, "certificate already exists")

// ErrCertificateDoesNotExist is returned when an action is performed on a domain that already exists
var ErrCertificateDoesNotExist = sdkerrors.Register(ModuleName, 16, "certificate does not exist")

// ErrGracePeriodNotFinished is returned when actions are performed on expired domains
var ErrGracePeriodNotFinished = sdkerrors.Register(ModuleName, 17, "domain grace period has not finished")

// ErrInvalidDomainType is returned when domain type is invalid
var ErrInvalidDomainType = sdkerrors.Register(ModuleName, 18, "invalid domain type")

// ErrInvalidRegisterer is returned when the registerer address provided is not valid (empty, malformed, etc)
var ErrInvalidRegisterer = sdkerrors.Register(ModuleName, 19, "invalid registerer")
