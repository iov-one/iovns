package cli

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/starname/types"
)

var (
	// CLI module error codes being with 4xx
	// ErrCertificateNotProvided is returned by the CLI when certificates are not provided
	ErrCertificateNotProvided = sdkerrors.Register(types.ModuleName, 400, "provide certificate")
	// ErrCertificatedProvidedOnlyOne is returned when multiple certs + key value certs are provided via CLI
	ErrCertificateProvideOnlyOne = sdkerrors.Register(types.ModuleName, 401, "provide either cert or cert-file")
	// ErrInvalidCertificate is returned when the provided certificate is deemed to be invalid
	ErrInvalidCertificate = sdkerrors.Register(types.ModuleName, 402, "invalid certificate")
)
