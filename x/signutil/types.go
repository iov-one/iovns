package signutil

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MsgSignText struct {
	Message string         `json:"message,omitempty"`
	Pairs   []Pair         `json:"pairs,omitempty"`
	Signer  sdk.AccAddress `json:"signer"`
}

func (m MsgSignText) Route() string {
	return ModuleName
}

func (m MsgSignText) Type() string {
	return "text_signature"
}

func (m MsgSignText) ValidateBasic() error {
	if len(m.Message) == 0 && len(m.Pairs) == 0 {
		return fmt.Errorf("empty msg and pairs")
	}
	if m.Signer.Empty() {
		return fmt.Errorf("missing signer")
	}
	return nil
}

// GetSignBytes implements sdk.Message
func (m MsgSignText) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Message
func (m MsgSignText) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.Signer} }
