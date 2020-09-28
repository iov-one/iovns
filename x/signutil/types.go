package signutil

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MsgTextSignature struct {
	Message string         `json:"message,omitempty"`
	Pairs   []Pair         `json:"pairs,omitempty"`
	Signer  sdk.AccAddress `json:"signer"`
}

func (m MsgTextSignature) Route() string {
	return ModuleName
}

func (m MsgTextSignature) Type() string {
	return "text_signature"
}

func (m MsgTextSignature) ValidateBasic() error {
	if len(m.Message) == 0 && len(m.Pairs) == 0 {
		return fmt.Errorf("empty msg and pairs")
	}
	if m.Signer.Empty() {
		return fmt.Errorf("missing signer")
	}
	return nil
}

// GetSignBytes implements sdk.Message
func (m MsgTextSignature) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Message
func (m MsgTextSignature) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Signer} }
