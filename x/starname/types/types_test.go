package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

func TestValidateDomainType(t *testing.T) {
	cases := map[string]struct {
		dType   DomainType
		wantErr *errors.Error
	}{
		"success open": {
			dType:   "open",
			wantErr: nil,
		},
		"success close": {
			dType:   "closed",
			wantErr: nil,
		},
		"fail one": {
			dType:   "sucuk doner",
			wantErr: ErrInvalidDomainType,
		},
		"fail two": {
			dType:   "",
			wantErr: ErrInvalidDomainType,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if err := ValidateDomainType(tc.dType); !tc.wantErr.Is(err) {
				t.Fatalf("unexpected error: %+v", err)
			}
		})
	}
}
