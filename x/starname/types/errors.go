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

// ErrInvalidResource is returned when provided resource is not valid
var ErrInvalidResource = sdkerrors.Register(ModuleName, 10, "resource provided is not valid")

// ErrDomainExpired is returned when actions are performed on expired domains
var ErrDomainExpired = sdkerrors.Register(ModuleName, 11, "domain has expired")

// ErrDomainNotExpired is returned when actions are performed on not expired domains
var ErrDomainNotExpired = sdkerrors.Register(ModuleName, 12, "domain has not expired")

// ErrAccountExists is returned when a create action is done on an account that already exists
var ErrAccountExists = sdkerrors.Register(ModuleName, 13, "account already exists")

// ErrInvalidRequest is a general error that covers the uncommon cases of invalid request
var ErrInvalidRequest = sdkerrors.Register(ModuleName, 14, "malformed request")

// ErrCertificateExists is returned when a creation action is done on a certificate that already exists
var ErrCertificateExists = sdkerrors.Register(ModuleName, 15, "certificate already exists")

// ErrCertificateDoesNotExist is returned when an action is performed on a domain that already exists
var ErrCertificateDoesNotExist = sdkerrors.Register(ModuleName, 16, "certificate does not exist")

// ErrDomainGracePeriodNotFinished is returned when actions are performed on expired domains
var ErrDomainGracePeriodNotFinished = sdkerrors.Register(ModuleName, 17, "domain grace period has not finished")

// ErrInvalidDomainType is returned when domain type is invalid
var ErrInvalidDomainType = sdkerrors.Register(ModuleName, 18, "invalid domain type")

// ErrInvalidRegisterer is returned when the registerer address provided is not valid (empty, malformed, etc)
var ErrInvalidRegisterer = sdkerrors.Register(ModuleName, 19, "invalid registerer")

// ErrOpEmptyAcc is returned when an operation tried to be run on empty account
var ErrOpEmptyAcc = sdkerrors.Register(ModuleName, 20, "account name provided cannot be empty")

// ErrAccountGracePeriodNotFinished is returned when actions are performed on not expired domains
var ErrAccountGracePeriodNotFinished = sdkerrors.Register(ModuleName, 21, "account grace period has not finished")

// ErrResourceLimitExceeded is returned when resource limit is exceeded
var ErrResourceLimitExceeded = sdkerrors.Register(ModuleName, 22, "resource limit exceeded")

// ErrCertificateSizeExceeded is returned when certificate size exceeded
var ErrCertificateSizeExceeded = sdkerrors.Register(ModuleName, 23, "certificate size exceeded")

// ErrCertificateLimitReached is returned when certificate limit is exceeded
var ErrCertificateLimitReached = sdkerrors.Register(ModuleName, 24, "certificate limit reached")

// ErrMetadataSizeExceeded is returned when metadata size exceeded
var ErrMetadataSizeExceeded = sdkerrors.Register(ModuleName, 25, "metadata size exceeded")

// ErrClosedDomainAccExpire is returned when expiration related operation trying to be run on closed domain
var ErrClosedDomainAccExpire = sdkerrors.Register(ModuleName, 26, "accounts in closed domains do not expire")

// ErrMaxRenewExceeded is returned when max renew time exceeded
var ErrMaxRenewExceeded = sdkerrors.Register(ModuleName, 27, "max renew exceeded")

// ErrRenewalDeadlineExceeded is returned when the renewal deadline was surpassed
var ErrRenewalDeadlineExceeded = sdkerrors.Register(ModuleName, 31, "renewal deadline was exceeded")

// ----------- QUERY ----------

// ErrProvideStarnameOrDomainName is returned when both domain/name and starname provided
var ErrProvideStarnameOrDomainName = sdkerrors.Register(ModuleName, 28, "provide either starname or domain/name")

// ErrStarnameNotContainSep is returned when provided starname does not contain separator
var ErrStarnameNotContainSep = sdkerrors.Register(ModuleName, 29, "starname does not contain separator")

// ErrStarnameMultipleSeparator returned when provided starname contains more than one separator
var ErrStarnameMultipleSeparator = sdkerrors.Register(ModuleName, 30, "starname should contain single separator")
