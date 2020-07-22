package starname

import (
	"bytes"
	"errors"
	"github.com/iov-one/iovns/pkg/utils"
	"github.com/iov-one/iovns/x/starname/keeper/executor"
	"reflect"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
)

func Test_Close_handlerMsgAddAccountCertificates(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"does not respect account expiration": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					CertificateCountMax: 2,
					CertificateSizeMax:  100,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
					Type:       types.ClosedDomain,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: 0,
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.AliceKey,
					NewCertificate: nil,
				})
				if err != nil {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
		},
	}
	keeper.RunTests(t, cases)
}

func Test_Open_handlerMsgAddAccountCertificates(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					CertificateCountMax: 2,
					CertificateSizeMax:  100,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
					Type:       types.OpenDomain,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: 0,
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.AliceKey,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
		},
	}
	keeper.RunTests(t, cases)
}

func Test_Common_handlerMsgAddAccountCertificates(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"domain does not exist": {
			BeforeTest: nil,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "does not exist",
					Name:           "",
					Owner:          keeper.BobKey,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// add expired domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: 0,
					Admin:      keeper.BobKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "",
					Owner:          keeper.BobKey,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				//
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "does not exist",
					Owner:          nil,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"msg owner is not account owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.BobKey,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"admin cannot add cert": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.BobKey,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"certificate exists": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					CertificateCountMax: 2,
					CertificateSizeMax:  100,
					AccountGracePeriod:  1000 * time.Second,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:        keeper.AliceKey,
					Certificates: []types.Certificate{[]byte("test")},
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.AliceKey,
					NewCertificate: []byte("test"),
				})
				if !errors.Is(err, types.ErrCertificateExists) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrCertificateExists, err)
				}
			},
			AfterTest: nil,
		},
		"certificate size exceeded": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					CertificateCountMax: 2,
					CertificateSizeMax:  4,
					AccountGracePeriod:  1000 * time.Second,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:        keeper.AliceKey,
					Certificates: nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.AliceKey,
					NewCertificate: []byte("12345"),
				})
				if !errors.Is(err, types.ErrCertificateSizeExceeded) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrCertificateExists, err)
				}
				_, err = handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.AliceKey,
					NewCertificate: []byte("1234"),
				})
				if err != nil {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrCertificateExists, err)
				}
			},
			AfterTest: nil,
		},
		"certificate limit reached": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					CertificateCountMax: 2,
					CertificateSizeMax:  100,
					AccountGracePeriod:  1000 * time.Second,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:        keeper.AliceKey,
					Certificates: []types.Certificate{[]byte("1")},
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.AliceKey,
					NewCertificate: []byte("12345"),
				})
				if err != nil {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrCertificateExists, err)
				}
				_, err = handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.AliceKey,
					NewCertificate: []byte("1234"),
				})
				if !errors.Is(err, types.ErrCertificateLimitReached) {
					t.Fatalf("handlerMsgAddAccountCertificates() expected error: %s, got: %s", types.ErrCertificateExists, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					CertificateCountMax: 2,
					CertificateSizeMax:  4,
					AccountGracePeriod:  1000 * time.Second,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgAddAccountCertificates(ctx, k, &types.MsgAddAccountCertificates{
					Domain:         "test",
					Name:           "test",
					Owner:          keeper.AliceKey,
					NewCertificate: []byte("test"),
				})
				if err != nil {
					t.Fatalf("handlerMsgAddAccountCertificates() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				expected := []types.Certificate{[]byte("test")}
				account := new(types.Account)
				ok := k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("test")}).PrimaryKey(), account)
				if !ok {
					t.Fatal("account not found")
				}
				if !reflect.DeepEqual(account.Certificates, expected) {
					t.Fatalf("handlerMsgAddAccountCertificates: got: %#v, expected: %#v", account.Certificates, expected)
				}
			},
		},
	}
	keeper.RunTests(t, cases)
}

func Test_Closed_handlerMsgDeleteAccountCertificate(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"does not respect account valid until": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
					Type:       types.ClosedDomain,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					ValidUntil:   0,
					Owner:        keeper.AliceKey,
					Certificates: []types.Certificate{[]byte("test")},
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("test"),
					Owner:             keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccountCertificates() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				ok := k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("test")}).PrimaryKey(), account)
				if !ok {
					t.Fatal("account not found")
				}
				for _, cert := range account.Certificates {
					if bytes.Equal(cert, []byte("test")) {
						t.Fatalf("handlerMsgDeleteAccountCertificates() certificate not deleted")
					}
				}
				// success
			},
		},
	}

	keeper.RunTests(t, cases)
}

func Test_Open_handlerMsgDeleteAccountCertificate(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
					Type:       types.OpenDomain,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					ValidUntil:   0,
					Owner:        keeper.AliceKey,
					Certificates: []types.Certificate{[]byte("test")},
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("test"),
					Owner:             keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerMsgDeleteAccountCertificates() got error: %s", err)
				}
			},
		},
	}

	keeper.RunTests(t, cases)
}
func Test_Common_handlerMsgDeleteAccountCertificate(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "does not exist",
					DeleteCertificate: nil,
					Owner:             keeper.BobKey,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"msg signer is not account owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: nil,
					Owner:             keeper.BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"domain admin cannot delete cert": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      keeper.BobKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: nil,
					Owner:             keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"certificate does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("does not exist"),
					Owner:             keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrCertificateDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrCertificateDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      keeper.AliceKey,
				}).Create()
				// add mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:        keeper.AliceKey,
					Certificates: []types.Certificate{[]byte("test")},
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("test"),
					Owner:             keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccountCertificates() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// check if certificate is still present
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if !exists {
					t.Fatal("account not found")
				}
				for _, cert := range account.Certificates {
					if bytes.Equal(cert, []byte("test")) {
						t.Fatalf("handlerMsgDeleteAccountCertificates() certificate not deleted")
					}
				}
				// success
			},
		},
	}

	keeper.RunTests(t, cases)
}

func Test_Closed_handlerMsgDeleteAccount(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"domain admin can delete": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					Type:       types.ClosedDomain,
					ValidUntil: types.MaxValidUntil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  keeper.BobKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"domain expired": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					Type:       types.ClosedDomain,
					ValidUntil: 2,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  keeper.BobKey,
				}).Create()
			},
			TestBlockTime: 2,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
		},
		"account owner cannot delete": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					Type:       types.ClosedDomain,
					ValidUntil: types.MaxValidUntil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("unexpected error: %s", err)
				}
			},
		},
	}
	keeper.RunTests(t, cases)
}

func Test_Open_handlerMsgDeleteAccount(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"domain admin cannot can delete before grace period": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidResource:      keeper.RegexMatchNothing,
					ValidURI:           keeper.RegexMatchAll,
					AccountGracePeriod: 1000 * time.Second,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					Type:       types.OpenDomain,
					ValidUntil: types.MaxValidUntil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  keeper.BobKey,
				}).Create()
			},
			TestBlockTime: 3,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
		},
		"no domain valid until check": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidResource:       keeper.RegexMatchNothing,
					ValidURI:            keeper.RegexMatchAll,
					DomainRenewalPeriod: 10,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					Type:       types.OpenDomain,
					ValidUntil: 2,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  keeper.BobKey,
				}).Create()
			},
			TestBlockTime: 100,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"only account owner can delete before grace period": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidResource:      keeper.RegexMatchNothing,
					ValidURI:           keeper.RegexMatchAll,
					AccountGracePeriod: 10 * time.Second,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					Type:       types.OpenDomain,
					ValidUntil: types.MaxValidUntil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  keeper.AliceKey,
				}).Create()
			},
			TestBlockTime: 5,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// admin test
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("unexpected error: %v", err)
				}
				// anyone test
				_, err = handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.CharlieKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("unexpected error: %v", err)
				}
				// account owner test
				_, err = handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"domain admin can delete after grace": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidResource:      keeper.RegexMatchNothing,
					ValidURI:           keeper.RegexMatchAll,
					AccountGracePeriod: 10,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					Type:       types.OpenDomain,
					ValidUntil: types.MaxValidUntil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  keeper.AliceKey,
				}).Create()
			},
			TestBlockTime: 100,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// admin test
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"anyone can delete after grace": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidResource:      keeper.RegexMatchNothing,
					ValidURI:           keeper.RegexMatchAll,
					AccountGracePeriod: 10,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					Type:       types.OpenDomain,
					ValidUntil: types.MaxValidUntil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  keeper.AliceKey,
				}).Create()
			},
			TestBlockTime: 100,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// admin test
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.CharlieKey,
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
	}
	keeper.RunTests(t, cases)
}

func Test_Common_handlerMsgDeleteAccount(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "does not exist",
					Name:   "does not exist",
					Owner:  keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccount() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:  "test",
					Admin: keeper.BobKey,
				}).Create()

			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccount() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success domain owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:  "test",
					Admin: keeper.AliceKey,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  keeper.BobKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"success account owner": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:  "test",
					Admin: keeper.AliceKey,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  keeper.BobKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if exists {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
	}

	// run tests
	keeper.RunTests(t, cases)
}

func Test_ClosedDomain_handlerMsgRegisterAccount(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"only domain admin can register account": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidResource:        keeper.RegexMatchNothing,
					ValidURI:             keeper.RegexMatchAll,
					DomainRenewalPeriod:  10,
					AccountRenewalPeriod: 10,
				})
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Owner:      keeper.BobKey,
					Registerer: keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
				_, err = handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test2",
					Owner:      keeper.BobKey,
					Registerer: keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
			},
		},
		"account valid until is set to max": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       keeper.RegexMatchNothing, // don't match anything
					ValidURI:            keeper.RegexMatchAll,     // match all
					DomainRenewalPeriod: 10,
				})
				// add a closed domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Owner:      keeper.BobKey,
					Registerer: keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if !exists {
					t.Fatal("account test not found")
				}
				if account.ValidUntil != types.MaxValidUntil {
					t.Fatalf("unexpected account valid until %d", account.ValidUntil)
				}
			},
		},
		"account owner can be different than domain admin": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       keeper.RegexMatchNothing, // don't match anything
					ValidURI:            keeper.RegexMatchAll,     // match all
					DomainRenewalPeriod: 10,
				})
				// add a closed domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Registerer: keeper.BobKey,
					Owner:      keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
				_, err = handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test2",
					Registerer: keeper.BobKey,
					Owner:      keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
			},
		},
	}
	// run tests
	keeper.RunTests(t, testCases)
}

func Test_OpenDomain_handleMsgRegisterAccount(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"account valid until is now plus config account renew": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:        keeper.RegexMatchNothing, // don't match anything
					ValidURI:             keeper.RegexMatchAll,     // match all
					DomainRenewalPeriod:  10 * time.Second,
					AccountRenewalPeriod: 10 * time.Second,
				})
				// add a closed domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).Create()
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Owner:      keeper.BobKey,
					Registerer: keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if !exists {
					t.Fatal("account test not found")
				}
				expected := utils.TimeToSeconds(time.Unix(11, 0))
				if account.ValidUntil != expected {
					t.Fatalf("unexpected account valid until %d, expected %d", account.ValidUntil, expected)
				}
			},
		},
	}
	keeper.RunTests(t, testCases)
}

func Test_Common_handleMsgRegisterAccount(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"fail resource": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {

				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       keeper.RegexMatchNothing, // don't match anything
					ValidURI:            keeper.RegexMatchAll,     // match all
					DomainRenewalPeriod: 10,
				})
				// add a domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: 2,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
					Resources: []types.Resource{
						{
							URI:      "works",
							Resource: "won't work",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrInvalidResource) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrInvalidResource, err)
				}
			},
			AfterTest: nil,
		},
		"fail invalid uri": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidURI:            keeper.RegexMatchNothing, // don't match anything
					ValidResource:       keeper.RegexMatchAll,     // match all
					DomainRenewalPeriod: 10,
				})
				// add a domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: 2,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
					Resources: []types.Resource{
						{
							URI:      "invalid blockchain id",
							Resource: "valid blockchain address",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrInvalidResource) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrInvalidResource, err)
				}
			},
			AfterTest: nil,
		},
		"fail invalid account name": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       keeper.RegexMatchAll,     // match all
					ValidURI:            keeper.RegexMatchAll,     // match all
					ValidAccountName:    keeper.RegexMatchNothing, // match nothing
					DomainRenewalPeriod: 10,
				})
				// add a domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: 2,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "this won't match",
					Owner:  keeper.AliceKey,
					Resources: []types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrInvalidAccountName) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrInvalidAccountName, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain name does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in resources
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       keeper.RegexMatchAll, // match all
					ValidURI:            keeper.RegexMatchAll, // match all
					ValidAccountName:    keeper.RegexMatchAll, // match nothing
					DomainRenewalPeriod: 10,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "this does not exist",
					Name:   "works",
					Owner:  keeper.AliceKey,
					Resources: []types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"fail only owner of domain with superuser can register accounts": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in resources
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       keeper.RegexMatchAll, // match all
					ValidURI:            keeper.RegexMatchAll, // match all
					ValidAccountName:    keeper.RegexMatchAll, // match nothing
					DomainRenewalPeriod: 10,
				})
				// add a domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: 2,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey, // invalid owner
					Resources: []types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain has expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in resources
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       keeper.RegexMatchAll, // match all
					ValidURI:            keeper.RegexMatchAll, // match all
					ValidAccountName:    keeper.RegexMatchAll, // match nothing
					DomainRenewalPeriod: 10,
				})
				// add a domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: 0, // domain is expired
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
					Resources: []types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"fail account exists": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in resources
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       keeper.RegexMatchAll, // match all
					ValidURI:            keeper.RegexMatchAll, // match all
					ValidAccountName:    keeper.RegexMatchAll, // match nothing
					DomainRenewalPeriod: 10,
				})
				// add a domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
				// add an account that we are gonna try to overwrite
				executor.NewAccount(ctx, k, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("exists"),
					Owner:        keeper.AliceKey,
					ValidUntil:   0,
					Resources:    nil,
					Certificates: nil,
					Broker:       nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "exists",
					Owner:  keeper.BobKey,
					Resources: []types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrAccountExists) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrAccountExists, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set regexp match nothing in resources
				// get set config function
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       keeper.RegexMatchAll, // match all
					ValidURI:            keeper.RegexMatchAll, // match all
					ValidAccountName:    keeper.RegexMatchAll, // match nothing
					DomainRenewalPeriod: 10,
				})
				// add a domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Owner:      keeper.BobKey,
					Registerer: keeper.BobKey,
					Resources: []types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
					Broker: nil,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterAccount() got error: %s", err)
				}
			},
			AfterTest: nil, // TODO fill with matching data
		},
	}
	// run tests
	keeper.RunTests(t, testCases)
}

func Test_Closed_handlerMsgRenewAccount(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"account cannot be renewed since its max": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					AccountRenewalPeriod: 1,
					AccountGracePeriod:   5,
				})
				// set mock domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:  "test",
					Type:  types.ClosedDomain,
					Admin: keeper.BobKey,
				}).Create()
				// set mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Unix(1000, 0)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
					Domain: "test",
					Name:   "test",
				})
				if !errors.Is(err, types.ErrInvalidDomainType) {
					t.Fatalf("handlerMsgRenewAccount() want err: %s, got: %s", types.ErrInvalidDomainType, err)
				}
			},
		},
	}

	keeper.RunTests(t, cases)
}
func Test_Open_handlerMsgRenewAccount(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					AccountRenewalPeriod: 1,
					AccountGracePeriod:   5,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
					Domain: "does not exist",
					Name:   "",
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgRenewAccount() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					AccountRenewalPeriod: 1,
					AccountGracePeriod:   5,
				})
				// set mock domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:  "test",
					Type:  types.OpenDomain,
					Admin: keeper.BobKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
					Domain: "test",
					Name:   "does not exist",
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgRenewAccount() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success domain grace period not updated": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					AccountRenewalPeriod:   1 * time.Second,
					AccountRenewalCountMax: 200000,
					AccountGracePeriod:     5 * time.Second,
				})
				// set mock domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:  "test",
					Type:  types.OpenDomain,
					Admin: keeper.BobKey,
				}).Create()
				// set mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Unix(1, 0)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
					Domain: "test",
					Name:   "test",
				})
				if err != nil {
					t.Fatalf("handlerMsgRenewAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if !exists {
					t.Fatal("account not found")
				}
				want := ctx.BlockTime().Add(k.ConfigurationKeeper.GetConfiguration(ctx).AccountRenewalPeriod)
				if account.ValidUntil != utils.TimeToSeconds(want) {
					t.Fatalf("handlerMsgRenewAccount() want: %d, got: %d", want.Unix(), account.ValidUntil)
				}
			},
		},
		"success domain valid until updated": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					AccountRenewalPeriod:   1 * time.Second,
					AccountRenewalCountMax: 200000,
					AccountGracePeriod:     5 * time.Second,
					DomainGracePeriod:      2 * time.Second,
				})
				// set mock domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Type:       types.OpenDomain,
					Admin:      keeper.BobKey,
					ValidUntil: 2,
				}).Create()
				// set mock account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Unix(1, 0)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			TestBlockTime: 10,
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
					Domain: "test",
					Name:   "test",
				})
				if err != nil {
					t.Fatalf("handlerMsgRenewAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				domain := new(types.Domain)
				exists := k.DomainStore(ctx).Read((&types.Domain{Name: "test"}).PrimaryKey(), domain)
				if !exists {
					t.Fatal("domain not found")
				}
				account := new(types.Account)
				exists = k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if !exists {
					t.Fatal("account not found")
				}
				if domain.ValidUntil != account.ValidUntil {
					t.Fatalf("handlerMsgRenewAccount() want: %d, got: %d", domain.ValidUntil, account.ValidUntil)
				}
			},
		},
	}

	keeper.RunTests(t, cases)
}

func Test_Closed_handlerMsgReplaceAccountResources(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"fail does not respect account valid until": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      keeper.RegexMatchAll,
					ValidResource: keeper.RegexMatchAll,
					ResourcesMax:  5,
				})
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
					Type:       types.ClosedDomain,
				}).Create()
				// create account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: 0,
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountResources() got error: %s", err)
				}
			},
		},
	}

	keeper.RunTests(t, cases)
}

func Test_Open_handlerMsgReplaceAccountResources(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      keeper.RegexMatchAll,
					ValidResource: keeper.RegexMatchAll,
					ResourcesMax:  3,
				})
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
					Type:       types.OpenDomain,
				}).Create()
				// create account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: 0,
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
	}

	keeper.RunTests(t, cases)
}
func Test_Common_handlerMsgReplaceAccountResources(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"invalid blockchain resource": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      keeper.RegexMatchNothing,
					ValidResource: keeper.RegexMatchNothing,
				})
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				// create account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []types.Resource{
						{
							URI:      "invalid",
							Resource: "invalid",
						},
					},
					Owner: keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrInvalidResource) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrInvalidResource, err)
				}
			},
		},
		"resource limit exceeded": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      keeper.RegexMatchAll,
					ValidResource: keeper.RegexMatchAll,
					ResourcesMax:  2,
				})
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				// create account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
						{
							URI:      "valid1",
							Resource: "valid1",
						},
						{
							URI:      "valid2",
							Resource: "valid2",
						},
					},
					Owner: keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrResourceLimitExceeded) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrInvalidResource, err)
				}
				_, err = handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []types.Resource{
						{
							URI:      "invalid",
							Resource: "invalid",
						},
						{
							URI:      "invalid2",
							Resource: "invalid2",
						},
					},
					Owner: keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrInvalidResource, err)
				}
				_, err = handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []types.Resource{
						{
							URI:      "invalid",
							Resource: "invalid",
						},
					},
					Owner: keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrInvalidResource, err)
				}
			},
		},
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      keeper.RegexMatchAll,
					ValidResource: keeper.RegexMatchAll,
				})
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "does not exist",
					Name:   "",
					NewResources: []types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      keeper.RegexMatchAll,
					ValidResource: keeper.RegexMatchAll,
				})
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:  "test",
					Admin: keeper.BobKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					NewResources: []types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      keeper.RegexMatchAll,
					ValidResource: keeper.RegexMatchAll,
				})
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "does not exist",
					NewResources: []types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"signer is not owner of account": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      keeper.RegexMatchAll,
					ValidResource: keeper.RegexMatchAll,
				})
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				// create account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: keeper.BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// set config to match all
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      keeper.RegexMatchAll,
					ValidResource: keeper.RegexMatchAll,
					ResourcesMax:  5,
				})
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				// create account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountResources() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				expected := []types.Resource{{
					URI:      "valid",
					Resource: "valid",
				}}
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if !exists {
					t.Fatal("account not found")
				}
				if !reflect.DeepEqual(expected, account.Resources) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected: %+v, got %+v", expected, account.Resources)
				}
			},
		},
	}
	// run tests
	keeper.RunTests(t, cases)
}

func Test_Closed_handlerMsgReplaceAccountMetadata(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"account expiration not respected": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					MetadataSizeMax: 100,
				})
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
					Type:       types.ClosedDomain,
				}).Create()
				// create account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: 0,
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
		},
	}

	keeper.RunTests(t, cases)
}

func Test_Open_handlerMsgReplaceAccountMetadata(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					MetadataSizeMax: 100,
				})
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
					Type:       types.OpenDomain,
				}).Create()
				// create account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: 0,
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
	}

	keeper.RunTests(t, cases)
}
func Test_Common_handlerMsgReplaceAccountMetadata(t *testing.T) {
	cases := map[string]keeper.SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain: "does not exist",
					Name:   "",
					Owner:  keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:  "test",
					Admin: keeper.BobKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain:         "test",
					Name:           "",
					NewMetadataURI: "",
					Owner:          nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain: "test",
					Name:   "does not exist",
					Owner:  nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"signer is not owner of account": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				// create account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.CharlieKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"domain admin cannot replace metadata": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				// create account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  keeper.BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"metadata size exceeded": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					MetadataSizeMax: 2,
				})
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				// create account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain:         "test",
					Name:           "test",
					NewMetadataURI: "https://test.com",
					Owner:          keeper.AliceKey,
				})
				if !errors.Is(err, types.ErrMetadataSizeExceeded) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() got error: %s", err)
				}
				_, err = handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain:         "test",
					Name:           "test",
					NewMetadataURI: "12",
					Owner:          keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountMetadata() got error: %s", err)
				}
			},
		},
		"success": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				setConfig := keeper.GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					MetadataSizeMax: 100,
				})
				// create domain
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      keeper.BobKey,
				}).Create()
				// create account
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      keeper.AliceKey,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain:         "test",
					Name:           "test",
					NewMetadataURI: "https://test.com",
					Owner:          keeper.AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountMetadata() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				expected := "https://test.com"
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if !exists {
					t.Fatal("account not found")
				}
				if !reflect.DeepEqual(expected, account.MetadataURI) {
					t.Fatalf("handlerMsgSetMetadataURI expected: %+v, got %+v", expected, account.MetadataURI)
				}
			},
		},
	}
	// run tests
	keeper.RunTests(t, cases)
}

func Test_Closed_handlerAccountTransfer(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"only domain admin can transfer": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// domain owned by alice
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
				// account owned by bob
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					Owner:      keeper.BobKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// alice is domain owner and should transfer account owned by bob to alice
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.AliceKey,
					NewOwner: keeper.CharlieKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
				_, err = handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.BobKey,
					NewOwner: keeper.CharlieKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if !exists {
					panic("unexpected account deletion")
				}
				if !account.Owner.Equals(keeper.CharlieKey) {
					t.Fatalf("handlerAccounTransfer() expected new owner: %s, got: %s", keeper.CharlieKey, account.Owner)
				}
			},
		},
		"domain admin can reset account content": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// domain owned by alice
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
				// account owned by bob
				executor.NewAccount(ctx, k, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					Owner:        keeper.BobKey,
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					MetadataURI:  "lol",
					Certificates: []types.Certificate{[]byte("test")},
					Resources: []types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// alice is domain owner and should transfer account owned by bob to alice
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.AliceKey,
					NewOwner: keeper.CharlieKey,
					Reset:    true,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if !exists {
					panic("unexpected account deletion")
				}
				if account.Resources != nil {
					panic("resources not deleted")
				}
				if account.Certificates != nil {
					panic("certificates not deleted")
				}
				if account.MetadataURI != "" {
					panic("metadata not deleted")
				}
			},
		},
	}

	keeper.RunTests(t, testCases)
}

func Test_Open_handlerAccountTransfer(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"domain admin cannot transfer": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// domain owned by alice
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).Create()
				// account owned by bob
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					Owner:      keeper.BobKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.AliceKey,
					NewOwner: keeper.CharlieKey,
					Reset:    false,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}

				_, err = handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.BobKey,
					NewOwner: keeper.CharlieKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if !exists {
					panic("unexpected account deletion")
				}
				if !account.Owner.Equals(keeper.CharlieKey) {
					t.Fatalf("handlerAccounTransfer() expected new owner: %s, got: %s", keeper.CharlieKey, account.Owner)
				}
			},
		},
		"domain admin cannot reset account content": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// domain owned by alice
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).Create()
				// account owned by bob
				executor.NewAccount(ctx, k, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					Owner:        keeper.BobKey,
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					MetadataURI:  "lol",
					Certificates: []types.Certificate{[]byte("test")},
					Resources: []types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// alice is domain owner and should transfer account owned by bob to alice
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.AliceKey,
					NewOwner: keeper.CharlieKey,
					Reset:    true,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if !exists {
					t.Fatal("unexpected account deletion")
				}
				if account.Resources == nil {
					t.Fatal("resources deleted")
				}
				if account.Certificates == nil {
					t.Fatal("certificates deleted")
				}
				if account.MetadataURI == "" {
					t.Fatal("metadata not deleted")
				}
			},
		},
	}
	keeper.RunTests(t, testCases)
}

func Test_Common_handlerAccountTransfer(t *testing.T) {
	testCases := map[string]keeper.SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				// do nothing
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "does not exist",
					Name:     "does not exist",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain has expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "expired domain",
					Admin:      keeper.BobKey,
					ValidUntil: 0,
					Type:       types.OpenDomain,
					Broker:     nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "expired domain",
					Name:     "",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "this account does not exist",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account expired": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					Owner:        keeper.BobKey,
					ValidUntil:   0,
					Resources:    nil,
					Certificates: nil,
					Broker:       nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
		"if domain has super user only domain admin can transfer accounts": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					Owner:        keeper.BobKey,
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Resources:    nil,
					Certificates: nil,
					Broker:       nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.BobKey,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"if domain has no super user then only account owner can transfer accounts": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					Owner:        keeper.AliceKey,
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Resources:    nil,
					Certificates: nil,
					Broker:       nil,
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.BobKey,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success domain without superuser": {
			BeforeTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				executor.NewDomain(ctx, k, types.Domain{
					Name:       "test",
					Admin:      keeper.BobKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).Create()
				executor.NewAccount(ctx, k, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					Owner:      keeper.AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				}).Create()
			},
			Test: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    keeper.AliceKey,
					NewOwner: keeper.BobKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context, mocks *keeper.Mocks) {
				account := new(types.Account)
				exists := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account)
				if !exists {
					panic("unexpected account deletion")
				}
				if account.Resources != nil {
					t.Fatalf("handlerAccountTransfer() account resources were not deleted")
				}
				if account.Certificates != nil {
					t.Fatalf("handlerAccountTransfer() account certificates were not deleted")
				}
				if !account.Owner.Equals(keeper.BobKey) {
					t.Fatalf("handlerAccounTransfer() expected new owner: %s, got: %s", keeper.BobKey, account.Owner)
				}
			},
		},
	}
	keeper.RunTests(t, testCases)
}
