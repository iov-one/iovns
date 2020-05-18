package cli

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/types"
)

var (
	// CLI module error codes being with 4xx
	ErrCertificateNotProvided    = sdkerrors.Register(types.ModuleName, 400, "provide certificate")
	ErrCertificateProvideOnlyOne = sdkerrors.Register(types.ModuleName, 401, "provide either cert or cert-file")
)
