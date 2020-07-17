package keeper

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/iov-one/iovns/mock"
	"github.com/iov-one/iovns/x/configuration"
	confCdc "github.com/iov-one/iovns/x/configuration/types"
	"github.com/iov-one/iovns/x/starname/types"
	"github.com/stretchr/testify/require"
	tmtypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

type DullMsg struct {
	signer sdk.AccAddress
}

func (m *DullMsg) FeePayer() sdk.AccAddress {
	return m.signer
}

func (m *DullMsg) Route() string {
	return "dull"
}

func (m *DullMsg) Type() string {
	return "dull_msg"
}

func (m *DullMsg) ValidateBasic() error {
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *DullMsg) GetSignBytes() []byte {
	return nil
}

// GetSigners implements sdk.Msg
func (m *DullMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.signer}
}

// NewTestCodec generates aliceAddr mock codec for keeper module
func NewTestCodec() *codec.Codec {
	// we should register this codec for all the modules
	// that are used and referenced by domain module
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	confCdc.RegisterCodec(cdc)
	return cdc
}

type Mocks struct {
	Supply *mock.SupplyKeeperMock
}

// NewTestKeeper generates aliceAddr keeper and aliceAddr context from it
func NewTestKeeper(t testing.TB, isCheckTx bool) (Keeper, sdk.Context, *Mocks) {
	cdc := NewTestCodec()
	// generate store
	mdb := db.NewMemDB()
	// generate multistore
	ms := store.NewCommitMultiStore(mdb)
	// generate store keys
	configurationStoreKey := sdk.NewKVStoreKey(configuration.StoreKey) // configuration module store key
	domainStoreKey := sdk.NewKVStoreKey(types.DomainStoreKey)          // domain module store key
	// generate sub store for each module referenced by the keeper
	ms.MountStoreWithDB(configurationStoreKey, sdk.StoreTypeIAVL, mdb) // mount configuration module
	ms.MountStoreWithDB(domainStoreKey, sdk.StoreTypeIAVL, mdb)        // mount domain module
	// test no errors
	require.Nil(t, ms.LoadLatestVersion())
	// create Mocks
	mocks := new(Mocks)
	// create mock supply keeper
	mocks.Supply = mock.NewSupplyKeeper()
	// create config keeper
	confKeeper := configuration.NewKeeper(cdc, configurationStoreKey, subspace.NewSubspace(cdc, nil, nil, "test"))
	// create context
	ctx := sdk.NewContext(ms, tmtypes.Header{Time: time.Now()}, isCheckTx, log.NewNopLogger())
	// create domain.Keeper
	return NewKeeper(cdc, domainStoreKey, confKeeper, mocks.Supply.Mock(), nil), ctx, mocks
}
