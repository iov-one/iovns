package executor

import (
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
	"reflect"
	"testing"
)

func TestAccount_AddCertificate(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	cert := []byte("a-cert")
	ex := NewAccount(testCtx, testKeeper, testAccount)
	ex.AddCertificate(cert)
	got := new(types.Account)
	testKeeper.AccountStore(testCtx).Read(testAccount.PrimaryKey(), got)
	if !reflect.DeepEqual(got.Certificates, append(testAccount.Certificates, cert)) {
		t.Fatal("unexpected result")
	}
}

func TestAccount_Create(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	acc := testAccount
	acc.Domain = "some-random-domain"
	ex := NewAccount(testCtx, testKeeper, acc)
	ex.Create()
	got := new(types.Account)
	testKeeper.AccountStore(testCtx).Read(acc.PrimaryKey(), got)
	if !reflect.DeepEqual(*got, acc) {
		t.Fatal("unexpected result")
	}
}

func TestAccount_DeleteCertificate(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	ex := NewAccount(testCtx, testKeeper, testAccount)
	ex.DeleteCertificate(0)
	got := new(types.Account)
	testKeeper.AccountStore(testCtx).Read(testAccount.PrimaryKey(), got)
	if len(got.Certificates) != 0 {
		t.Fatal("unexpected result")
	}
}

func TestAccount_Renew(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	NewAccount(testCtx, testKeeper, testAccount).Renew()
	newAcc := new(types.Account)
	ok := testKeeper.AccountStore(testCtx).Read(testAccount.PrimaryKey(), newAcc)
	if !ok {
		t.Fatal("account was deleted")
	}
	if newAcc.ValidUntil != testAccount.ValidUntil+int64(testConfig.AccountRenewalPeriod.Seconds()) {
		t.Fatal("time mismatch")
	}
}

func TestAccount_ReplaceResources(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	newRes := []*types.Resource{{
		URI:      "uri",
		Resource: "res",
	}}
	ex := NewAccount(testCtx, testKeeper, testAccount)
	ex.ReplaceResources(newRes)
	got := new(types.Account)
	testKeeper.AccountStore(testCtx).Read(testAccount.PrimaryKey(), got)
	if !reflect.DeepEqual(got.Resources, newRes) {
		t.Fatal("unexpected result")
	}

}

func TestAccount_State(t *testing.T) {

}

func TestAccount_Transfer(t *testing.T) {
	ex := NewAccount(testCtx, testKeeper, testAccount)
	t.Run("no-reset", func(t *testing.T) {
		testCtx, _ := testCtx.CacheContext()

		ex.Transfer(keeper.CharlieKey, false)
		got := new(types.Account)
		testKeeper.AccountStore(testCtx).Read(testAccount.PrimaryKey(), got)
		if !got.Owner.Equals(keeper.CharlieKey) {
			t.Fatal("unexpected owner")
		}
		if !reflect.DeepEqual(got.Resources, testAccount.Resources) {
			t.Fatal("unexpected resources")
		}
		if !reflect.DeepEqual(got.MetadataURI, testAccount.MetadataURI) {
			t.Fatal("unexpected metadata")
		}
		if !reflect.DeepEqual(got.Certificates, testAccount.Certificates) {
			t.Fatal("unexpected certs")
		}
	})
	t.Run("with-reset", func(t *testing.T) {
		testCtx, _ := testCtx.CacheContext()

		ex.Transfer(keeper.BobKey, true)
		got := new(types.Account)
		testKeeper.AccountStore(testCtx).Read(testAccount.PrimaryKey(), got)
		if !got.Owner.Equals(keeper.BobKey) {
			t.Fatal("owner mismatch")
		}
		if got.MetadataURI != "" || got.Resources != nil || got.Certificates != nil {
			t.Fatal("reset not performed")
		}
	})
}

func TestAccount_UpdateMetadata(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()

	newMeta := "a new meta"
	ex := NewAccount(testCtx, testKeeper, testAccount)
	ex.UpdateMetadata(newMeta)
	got := new(types.Account)
	testKeeper.AccountStore(testCtx).Read(testAccount.PrimaryKey(), got)
	if !reflect.DeepEqual(got.MetadataURI, newMeta) {
		t.Fatal("unexpected result")
	}
}

func TestAccount_Delete(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	ex := NewAccount(testCtx, testKeeper, testAccount)
	ex.Delete()
	got := new(types.Account)
	found := testKeeper.AccountStore(testCtx).Read(testAccount.PrimaryKey(), got)
	if found {
		t.Fatal("account was not deleted")
	}
}
