package pkg

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/app"
)

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = app.MakeCodec()
	config := sdk.GetConfig()
	config.SetCoinType(app.CoinType)
	config.SetFullFundraiserPath(app.FullFundraiserPath)
	config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	config.Seal()
}
