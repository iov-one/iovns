package executor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/utils"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
	"testing"
)

func TestDomain_Transfer(t *testing.T) {
	// defines test prereqs
	init := func() (k keeper.Keeper, ctx sdk.Context, ex *Domain) {
		k, ctx, _ = keeper.NewTestKeeper(t, false)
		domain := types.Domain{
			Name:       "test",
			Admin:      keeper.BobKey,
			ValidUntil: 1,
			Type:       types.OpenDomain,
			Broker:     nil,
		}
		acc1 := types.Account{
			Domain:       "test",
			Name:         utils.StrPtr("1"),
			Owner:        keeper.BobKey,
			ValidUntil:   1,
			Resources:    nil,
			Certificates: nil,
			Broker:       nil,
			MetadataURI:  "",
		}
		acc2 := types.Account{
			Domain:       "test",
			Name:         utils.StrPtr("2"),
			Owner:        keeper.BobKey,
			ValidUntil:   1,
			Resources:    nil,
			Certificates: nil,
			Broker:       nil,
			MetadataURI:  "",
		}
		// add account not owned
		acc3 := types.Account{
			Domain: "test",
			Name:   utils.StrPtr("not-owned"),
			Owner:  keeper.CharlieKey,
		}
		NewDomain(ctx, k, domain).Create()
		NewAccount(ctx, k, acc1).Create()
		NewAccount(ctx, k, acc2).Create()
		NewAccount(ctx, k, acc3).Create()
		ex = NewDomain(ctx, k, domain)
		return
	}
	t.Run("transfer owned", func(t *testing.T) {
		k, ctx, ex := init()
		ex.Transfer(types.TransferOwned, keeper.AliceKey)
		filter := k.AccountStore(ctx).Filter(&types.Account{
			Domain: "test",
		})
		for ; filter.Valid(); filter.Next() {
			acc := new(types.Account)
			filter.Read(acc)
			if !acc.Owner.Equals(keeper.AliceKey) && !acc.Owner.Equals(keeper.CharlieKey) {
				t.Fatal("owner mismatch")
			}
		}
	})
	t.Run("transfer-flush", func(t *testing.T) {
		k, ctx, ex := init()
		ex.Transfer(types.TransferFlush, keeper.AliceKey)
		filter := k.AccountStore(ctx).Filter(&types.Account{
			Domain: "test",
		})
		for ; filter.Valid(); filter.Next() {
			acc := new(types.Account)
			filter.Read(acc)
			// only empty account is expected
			if *acc.Name != types.EmptyAccountName {
				t.Fatalf("only empty account is expected to exist, got: %s", *acc.Name)
			}
		}
	})
	t.Run("transfer-reset-none", func(t *testing.T) {
		k, ctx, ex := init()
		ex.Transfer(types.TransferResetNone, keeper.AliceKey)
		filter := k.AccountStore(ctx).Filter(&types.Account{
			Domain: "test",
		})
		for ; filter.Valid(); filter.Next() {
			acc := new(types.Account)
			filter.Read(acc)
			switch *acc.Name {
			case types.EmptyAccountName:
				if !acc.Owner.Equals(keeper.AliceKey) {
					t.Fatal("owner mismatch")
				}
			case "1":
				if !acc.Owner.Equals(keeper.BobKey) {
					t.Fatal("owner mismatch")
				}
			case "2":
				if !acc.Owner.Equals(keeper.BobKey) {
					t.Fatal("owner mismatch")
				}
			case "not-owned":
				if !acc.Owner.Equals(keeper.CharlieKey) {
					t.Fatal("owner mismatch")
				}
			default:
				t.Fatalf("unexpected account found: %s", *acc.Name)
			}
		}
	})

}
