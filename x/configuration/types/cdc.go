package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgUpsertLevelFee{}, fmt.Sprintf("%s/MsgUpsertLevelFees", ModuleName), nil)
	cdc.RegisterConcrete(MsgUpsertDefaultFee{}, fmt.Sprintf("%s/MsgUpsertDefaultFees", ModuleName), nil)
	cdc.RegisterConcrete(MsgDeleteLevelFee{}, fmt.Sprintf("%s/MsgDeleteLevelFees", ModuleName), nil)
	cdc.RegisterConcrete(MsgUpdateConfig{}, fmt.Sprintf("%s/MsgUpdateConfig", ModuleName), nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
