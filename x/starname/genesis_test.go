package starname

import (
	"github.com/iov-one/iovns/tutils"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/keeper/executor"
	"github.com/iov-one/iovns/x/starname/types"
	"testing"
)

func TestExportGenesis(t *testing.T) {

	k, ctx, _ := keeper.NewTestKeeper(t, true)
	executor.NewDomain(ctx, k, types.Domain{
		Name:       "test",
		Admin:      keeper.AliceKey,
		ValidUntil: 100,
		Type:       types.OpenDomain,
		Broker:     nil,
	}).Create()
	executor.NewAccount(ctx, k, types.Account{
		Domain:      "test",
		Name:        tutils.StrPtr("test"),
		Owner:       keeper.AliceKey,
		ValidUntil:  100,
		MetadataURI: "",
	}).Create()
	_ = ExportGenesis(ctx, k)
}
