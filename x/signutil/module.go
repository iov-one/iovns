package signutil

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

const ModuleName = "signutil"

// ModuleCdc instantiates a new codec for the domain module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgArbitrarySignature{}, fmt.Sprintf("%s/%s", ModuleName, "MsgArbitrarySignature"), nil)
}

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// nolint
// - - - FILL APP MODULE BASIC -- //
// AppModuleBasic implements the AppModuleBasic interface of the cosmos-sdk
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string                   { return ModuleName }
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) { RegisterCodec(cdc) }
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return []byte("{}")
}
func (AppModuleBasic) ValidateGenesis(_ json.RawMessage) (err error) {
	return nil
}

func (AppModuleBasic) RegisterRESTRoutes(_ context.CLIContext, _ *mux.Router) {

}
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return getTxCmd(ModuleName, cdc)
}
func (AppModuleBasic) GetQueryCmd(_ *codec.Codec) *cobra.Command {
	return nil
}

// - - FILL APP MODULE - -
type AppModule struct {
	AppModuleBasic
}

func NewAppModule() AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
	}
}
func (AppModule) Name() string                                       { return ModuleName }
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry)         {}
func (AppModule) Route() string                                      { return ModuleName }
func (AppModule) QuerierRoute() string                               { return ModuleName }
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
func (a AppModule) NewHandler() sdk.Handler {
	return func(_ sdk.Context, _ sdk.Msg) (*sdk.Result, error) {
		return nil, fmt.Errorf("invalid call to signutil module")
	}
}
func (a AppModule) NewQuerierHandler() sdk.Querier {
	return func(_ sdk.Context, _ []string, _ abci.RequestQuery) ([]byte, error) {
		return nil, fmt.Errorf("invalid query to signutil module")
	}
}

func (a AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return nil
}

func (a AppModule) ExportGenesis(_ sdk.Context) json.RawMessage {
	return []byte("{}")
}
