package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc instantiates a new codec for the domain module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers the sdk.Msg for the module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgRegisterDomain{}, fmt.Sprintf("%s/RegisterDomain", ModuleName), nil)
	cdc.RegisterConcrete(&MsgTransferDomain{}, fmt.Sprintf("%s/TransferDomain", ModuleName), nil)
	cdc.RegisterConcrete(&MsgTransferAccount{}, fmt.Sprintf("%s/TransferAccount", ModuleName), nil)
	cdc.RegisterConcrete(&MsgRenewAccount{}, fmt.Sprintf("%s/RenewAccount", ModuleName), nil)
	cdc.RegisterConcrete(&MsgAddAccountCertificates{}, fmt.Sprintf("%s/AddAccountCertificates", ModuleName), nil)
	cdc.RegisterConcrete(&MsgDeleteAccountCertificate{}, fmt.Sprintf("%s/DeleteAccountCertificates", ModuleName), nil)
	cdc.RegisterConcrete(&MsgDeleteAccount{}, fmt.Sprintf("%s/DeleteAccount", ModuleName), nil)
	cdc.RegisterConcrete(&MsgDeleteDomain{}, fmt.Sprintf("%s/DeleteDomain", ModuleName), nil)
	cdc.RegisterConcrete(&MsgFlushDomain{}, fmt.Sprintf("%s/FlusDomain", ModuleName), nil)
	cdc.RegisterConcrete(&MsgRegisterAccount{}, fmt.Sprintf("%s/RegisterAccount", ModuleName), nil)
	cdc.RegisterConcrete(&MsgRenewDomain{}, fmt.Sprintf("%s/RenewDomain", ModuleName), nil)
	cdc.RegisterConcrete(&MsgReplaceAccountTargets{}, fmt.Sprintf("%s/ReplaceAccountTargets", ModuleName), nil)
}
