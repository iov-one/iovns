package executor

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/tutils"
	tmtypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
	"os"
	"testing"
	"time"
)

var testCtx sdk.Context
var testKey = sdk.NewKVStoreKey("test")
var testCdc *codec.Codec

var aliceKey sdk.AccAddress
var bobKey sdk.AccAddress

func newTest() error {
	_, addr := tutils.GeneratePrivKeyAddressPairs(2)
	aliceKey = addr[0]
	bobKey = addr[1]
	testCdc = codec.New()
	mdb := db.NewMemDB()
	ms := store.NewCommitMultiStore(mdb)
	ms.MountStoreWithDB(testKey, sdk.StoreTypeIAVL, mdb)
	err := ms.LoadLatestVersion()
	if err != nil {
		return err
	}
	testCtx = sdk.NewContext(ms, tmtypes.Header{Time: time.Now()}, true, log.NewNopLogger())
	return nil
}

func TestMain(m *testing.M) {
	err := newTest()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}
