package starname

import (
	"encoding/json"
	"testing"

	"github.com/iov-one/iovns/pkg/utils"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/keeper/executor"
	"github.com/iov-one/iovns/x/starname/types"
)

func TestExportGenesis(t *testing.T) {
	expected := `{"domains":[{"name":"test","admin":"cosmos1ze7y9qwdddejmy7jlw4cymqqlt2wh05ytm076d","broker":"","valid_until":100,"type":"open"}],"accounts":[{"domain":"test","name":"","owner":"cosmos1ze7y9qwdddejmy7jlw4cymqqlt2wh05ytm076d","broker":"","valid_until":100,"resources":null,"certificates":null,"metadata_uri":""},{"domain":"test","name":"test","owner":"cosmos1ze7y9qwdddejmy7jlw4cymqqlt2wh05ytm076d","broker":"","valid_until":100,"resources":null,"certificates":null,"metadata_uri":""}]}`
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
		Name:        utils.StrPtr("test"),
		Owner:       keeper.AliceKey,
		ValidUntil:  100,
		MetadataURI: "",
	}).Create()
	b, err := json.Marshal(ExportGenesis(ctx, k))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != expected {
		t.Fatalf("unexpected genesis state:\nGot: %s\nWanted: %s", b, expected)
	}
}
