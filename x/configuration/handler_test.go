package configuration

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/configuration/types"
)

func Test_HandleUpdateConfig(t *testing.T) {
	cases := map[string]SubTest{
		"only configurer can configure": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context) {
				conf := Config{
					Configurer: AliceKey,
				}
				k.SetConfig(ctx, conf)
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context) {
				msg := types.MsgUpdateConfig{
					Signer: CharlieKey,
				}
				_, err := handleUpdateConfig(ctx, msg, k)
				if !errors.Is(err, sdkerrors.ErrUnauthorized) {
					t.Fatalf("unexpected error: %s", err)
				}
				msg = types.MsgUpdateConfig{
					Signer: AliceKey,
					NewConfiguration: Config{
						Configurer: BobKey,
					},
				}
				_, err = handleUpdateConfig(ctx, msg, k)
				if err != nil {
					t.Fatalf("handlerMsgDeleteDomain() got error: %s", err)
				}
			},
		},
	}
	RunTests(t, cases)
}
