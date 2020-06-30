package executor

import (
	"github.com/iov-one/iovns/x/domain/types"
	"reflect"
	"testing"
)

func Test_accountStore(t *testing.T) {
	store := newAccountStore(testCtx, testKey, testCdc)
	account := types.Account{
		Domain:     "domain",
		Name:       "account",
		Owner:      aliceKey,
		ValidUntil: 0,
		Resources: []types.Resource{{
			URI:      "x",
			Resource: "y",
		}},
		Certificates: nil,
		Broker:       nil,
		MetadataURI:  "",
	}
	t.Run("create", func(t *testing.T) {
		store.create(account)
	})
	t.Run("read", func(t *testing.T) {
		got, exists := store.read(account.Domain, account.Name)
		if !exists {
			t.Fatalf("account not found")
		}
		if !reflect.DeepEqual(got, account) {
			t.Fatalf("expected: %+v, got: %+v", got, account)
		}
	})
	t.Run("update", func(t *testing.T) {
		changed := account
		changed.Owner = bobKey
		store.update(changed)
		got, exists := store.read(account.Domain, account.Name)
		if !exists {
			t.Fatalf("account not found")
		}
		if !reflect.DeepEqual(got, changed) {
			t.Fatalf("expected: %+v, got: %+v", changed, got)
		}
		account = changed
	})
	t.Run("delete", func(t *testing.T) {
		store.delete(account)
		if _, exists := store.read(account.Domain, account.Name); exists {
			t.Fatalf("account should not exist")
		}
	})
}
