package pkg

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/iov-one/iovns/app"
)

// ModuleCdc instantiates a new codec for the domain module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
	config := sdk.GetConfig()
	config.SetCoinType(app.CoinType)
	config.SetFullFundraiserPath(app.FullFundraiserPath)
	config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	config.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	sdk.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
}
