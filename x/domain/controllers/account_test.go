package controllers

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"reflect"
	"testing"
)

func TestAccount_certNotExist(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		acc := &Account{
			account: &types.Account{
				Certificates: []types.Certificate{[]byte("test-cert")},
			},
		}
		err := acc.certNotExist([]byte("does not exist"), nil)
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("cert exists", func(t *testing.T) {
		acc := &Account{
			account: &types.Account{
				Certificates: []types.Certificate{[]byte("test-cert"), []byte("exists")},
			},
		}
		i := new(int)
		err := acc.certNotExist([]byte("exists"), i)
		if !errors.Is(err, types.ErrCertificateExists) {
			t.Fatalf("unexpected error: %s, wanted: %s", err, types.ErrCertificateExists)
		}
		if *i != 1 {
			t.Fatalf("unexpected index pointer: %d", *i)
		}
	})
}

func TestAccount_mustExist(t *testing.T) {
	type fields struct {
		name    string
		domain  string
		account *types.Account
		conf    *configuration.Config
		ctx     sdk.Context
		k       keeper.Keeper
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				name:    tt.fields.name,
				domain:  tt.fields.domain,
				account: tt.fields.account,
				conf:    tt.fields.conf,
				ctx:     tt.fields.ctx,
				k:       tt.fields.k,
			}
			if err := a.mustExist(); (err != nil) != tt.wantErr {
				t.Errorf("mustExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccount_mustNotExist(t *testing.T) {
	type fields struct {
		name    string
		domain  string
		account *types.Account
		conf    *configuration.Config
		ctx     sdk.Context
		k       keeper.Keeper
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				name:    tt.fields.name,
				domain:  tt.fields.domain,
				account: tt.fields.account,
				conf:    tt.fields.conf,
				ctx:     tt.fields.ctx,
				k:       tt.fields.k,
			}
			if err := a.mustNotExist(); (err != nil) != tt.wantErr {
				t.Errorf("mustNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccount_notExpired(t *testing.T) {
	type fields struct {
		name    string
		domain  string
		account *types.Account
		conf    *configuration.Config
		ctx     sdk.Context
		k       keeper.Keeper
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				name:    tt.fields.name,
				domain:  tt.fields.domain,
				account: tt.fields.account,
				conf:    tt.fields.conf,
				ctx:     tt.fields.ctx,
				k:       tt.fields.k,
			}
			if err := a.notExpired(); (err != nil) != tt.wantErr {
				t.Errorf("notExpired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccount_ownedBy(t *testing.T) {
	type fields struct {
		name    string
		domain  string
		account *types.Account
		conf    *configuration.Config
		ctx     sdk.Context
		k       keeper.Keeper
	}
	type args struct {
		addr sdk.AccAddress
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				name:    tt.fields.name,
				domain:  tt.fields.domain,
				account: tt.fields.account,
				conf:    tt.fields.conf,
				ctx:     tt.fields.ctx,
				k:       tt.fields.k,
			}
			if err := a.ownedBy(tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("ownedBy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccount_requireAccount(t *testing.T) {
	type fields struct {
		name    string
		domain  string
		account *types.Account
		conf    *configuration.Config
		ctx     sdk.Context
		k       keeper.Keeper
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				name:    tt.fields.name,
				domain:  tt.fields.domain,
				account: tt.fields.account,
				conf:    tt.fields.conf,
				ctx:     tt.fields.ctx,
				k:       tt.fields.k,
			}
			if err := a.requireAccount(); (err != nil) != tt.wantErr {
				t.Errorf("requireAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccount_requireConfiguration(t *testing.T) {
	type fields struct {
		name    string
		domain  string
		account *types.Account
		conf    *configuration.Config
		ctx     sdk.Context
		k       keeper.Keeper
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestAccount_validName(t *testing.T) {
	type fields struct {
		name    string
		domain  string
		account *types.Account
		conf    *configuration.Config
		ctx     sdk.Context
		k       keeper.Keeper
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				name:    tt.fields.name,
				domain:  tt.fields.domain,
				account: tt.fields.account,
				conf:    tt.fields.conf,
				ctx:     tt.fields.ctx,
				k:       tt.fields.k,
			}
			if err := a.validName(); (err != nil) != tt.wantErr {
				t.Errorf("validName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAccountController(t *testing.T) {
	type args struct {
		ctx    sdk.Context
		k      keeper.Keeper
		domain string
		name   string
	}
	tests := []struct {
		name string
		args args
		want *Account
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAccountController(tt.args.ctx, tt.args.k, tt.args.domain, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAccountController() = %v, want %v", got, tt.want)
			}
		})
	}
}
