package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
	"os"
	"reflect"
	"testing"
)

func genAddress() (sdk.AccAddress, sdk.AccAddress) {
	keyBase := keys.NewInMemory()
	addr1, _, err := keyBase.CreateMnemonic("alice", keys.English, "", keys.Secp256k1)
	if err != nil {
		fmt.Println("unable to generate mock addresses " + err.Error())
		os.Exit(1)
	}
	addr2, _, err := keyBase.CreateMnemonic("bob", keys.English, "", keys.Secp256k1)
	if err != nil {
		fmt.Println("unable to generate mock addresses " + err.Error())
		os.Exit(1)
	}
	return addr1.GetAddress(), addr2.GetAddress()
}

func Test_indexFunctionality(t *testing.T) {
	k, ctx := NewTestKeeper(t, true)
	k.CreateAccount(ctx, types.Account{
		Domain: "test",
		Name:   "1",
		Owner:  aliceAddr,
	})
	k.CreateAccount(ctx, types.Account{
		Domain: "test",
		Name:   "2",
		Owner:  bobAddr,
	})
	k.CreateAccount(ctx, types.Account{
		Domain: "test",
		Name:   "3",
		Owner:  aliceAddr,
	})
	accountKeys := k.iterAccountToOwner(ctx, aliceAddr)
	// expected two keys
	if len(accountKeys) != 2 {
		t.Fatalf("expected two keys, got: %d", len(accountKeys))
	}
	// transfer account
	acc, _ := k.GetAccount(ctx, "test", "1")
	k.TransferAccount(ctx, acc, bobAddr)
	// expected two keys for account bobAddr
	if len(k.iterAccountToOwner(ctx, bobAddr)) != 2 {
		t.Fatalf("expected two keys for %s, got: %d", bobAddr, len(k.iterAccountToOwner(ctx, bobAddr)))
	}
	// expect one key for aliceAddr
	if len(k.iterAccountToOwner(ctx, aliceAddr)) != 1 {
		t.Fatalf("expected two keys for %s, got: %d", bobAddr, len(k.iterAccountToOwner(ctx, aliceAddr)))
	}
	// delete account from bobAddr
	k.DeleteAccount(ctx, "test", "1") // belongs to bobAddr
	if len(k.iterAccountToOwner(ctx, bobAddr)) != 1 {
		t.Fatalf("expected two keys for %s, got: %d", bobAddr, len(k.iterAccountToOwner(ctx, bobAddr)))
	}

}

func TestKeeper_iterAccountToOwner(t *testing.T) {

}

func TestKeeper_iterDomainToOwner(t *testing.T) {

}

func TestKeeper_mapAccountToOwner(t *testing.T) {

}

func TestKeeper_mapDomainToOwner(t *testing.T) {
}

func TestKeeper_unmapAccountToOwner(t *testing.T) {

}

func TestKeeper_unmapDomainToOwner(t *testing.T) {

}

func Test_accAddrFromIndex(t *testing.T) {
	if !(aliceAddr.String() == accAddrFromIndex(indexAddr(aliceAddr)).String()) {
		t.Fatalf("mismatched addresses for: %s", aliceAddr.String())
	}
}

func Test_accountIndexStore(t *testing.T) {
	type args struct {
		store sdk.KVStore
	}
	tests := []struct {
		name string
		args args
		want sdk.KVStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := accountIndexStore(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("accountIndexStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_domainIndexStore(t *testing.T) {
	type args struct {
		store sdk.KVStore
	}
	tests := []struct {
		name string
		args args
		want sdk.KVStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := domainIndexStore(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("domainIndexStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOwnerToAccountKey(t *testing.T) {
	type args struct {
		owner   sdk.AccAddress
		domain  string
		account string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOwnerToAccountKey(tt.args.owner, tt.args.domain, tt.args.account); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOwnerToAccountKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOwnerToDomainKey(t *testing.T) {
	type args struct {
		owner  sdk.AccAddress
		domain string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOwnerToDomainKey(tt.args.owner, tt.args.domain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOwnerToDomainKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexAddr(t *testing.T) {
	type args struct {
		addr sdk.AccAddress
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := indexAddr(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("indexAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitOwnerToAccountKey(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name        string
		args        args
		wantAddr    sdk.AccAddress
		wantDomain  string
		wantAccount string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddr, gotDomain, gotAccount := splitOwnerToAccountKey(tt.args.key)
			if !reflect.DeepEqual(gotAddr, tt.wantAddr) {
				t.Errorf("splitOwnerToAccountKey() gotAddr = %v, want %v", gotAddr, tt.wantAddr)
			}
			if gotDomain != tt.wantDomain {
				t.Errorf("splitOwnerToAccountKey() gotDomain = %v, want %v", gotDomain, tt.wantDomain)
			}
			if gotAccount != tt.wantAccount {
				t.Errorf("splitOwnerToAccountKey() gotAccount = %v, want %v", gotAccount, tt.wantAccount)
			}
		})
	}
}

func Test_splitOwnerToDomainKey(t *testing.T) {
	type args struct {
		key []byte
	}
	tests := []struct {
		name       string
		args       args
		wantAddr   sdk.AccAddress
		wantDomain string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddr, gotDomain := splitOwnerToDomainKey(tt.args.key)
			if !reflect.DeepEqual(gotAddr, tt.wantAddr) {
				t.Errorf("splitOwnerToDomainKey() gotAddr = %v, want %v", gotAddr, tt.wantAddr)
			}
			if gotDomain != tt.wantDomain {
				t.Errorf("splitOwnerToDomainKey() gotDomain = %v, want %v", gotDomain, tt.wantDomain)
			}
		})
	}
}
