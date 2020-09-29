package pkg

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	config2 "github.com/iov-one/iovns/app/config"
)

// ModuleCdc instantiates a new codec for the domain module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
	config := sdk.GetConfig()
	config.SetCoinType(config2.CoinType)
	config.SetFullFundraiserPath(config2.FullFundraiserPath)
	config.SetBech32PrefixForAccount(config2.Bech32PrefixAccAddr, config2.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(config2.Bech32PrefixValAddr, config2.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(config2.Bech32PrefixConsAddr, config2.Bech32PrefixConsPub)
	config.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	sdk.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
}
