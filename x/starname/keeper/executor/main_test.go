package executor

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/mock"
	"github.com/iov-one/iovns/pkg/utils"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
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
var _, testAddrs = utils.GeneratePrivKeyAddressPairs(2)
var aliceKey sdk.AccAddress = testAddrs[0]
var bobKey sdk.AccAddress = testAddrs[1]
var testConfig = &configuration.Config{
	Configurer:           nil,
	DomainRenewalPeriod:  10 * time.Second,
	AccountRenewalPeriod: 20 * time.Second,
}

var testKeeper keeper.Keeper
var testAccount = types.Account{
	Domain:     "a-super-domain",
	Name:       utils.StrPtr("a-super-account"),
	Owner:      aliceKey,
	ValidUntil: 10000,
	Resources: []*types.Resource{
		{
			URI:      "a-super-uri",
			Resource: "a-super-res",
		},
	},
	Certificates: [][]byte{[]byte("a-random-cert")},
	Broker:       nil,
	MetadataURI:  "metadata",
}

var testDomain = types.Domain{
	Name:       "a-super-domain",
	Admin:      bobKey,
	ValidUntil: 100,
	Type:       types.ClosedDomain,
}

func newTest() error {
	mockConfig := mock.NewConfiguration(nil, testConfig)
	// gen test store
	testCdc = codec.New()
	mdb := db.NewMemDB()
	ms := store.NewCommitMultiStore(mdb)
	ms.MountStoreWithDB(testKey, sdk.StoreTypeIAVL, mdb)
	err := ms.LoadLatestVersion()
	if err != nil {
		return err
	}
	testCtx = sdk.NewContext(ms, tmtypes.Header{Time: time.Now()}, true, log.NewNopLogger())
	testKeeper = keeper.NewKeeper(testCdc, testKey, mockConfig, nil, nil)
	testKeeper.AccountStore(testCtx).Create(&testAccount)
	testKeeper.DomainStore(testCtx).Create(&testDomain)
	testKeeper.AccountStore(testCtx).Create(&types.Account{
		Domain:      testDomain.Name,
		Name:        utils.StrPtr(types.EmptyAccountName),
		Owner:       testDomain.Admin,
		ValidUntil:  testDomain.ValidUntil,
		MetadataURI: "",
	})
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
