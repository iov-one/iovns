package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRegisterDomain{}, fmt.Sprintf("%s/RegisterDomain", ModuleName), nil)
}
