package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/iov-one/iovns/x/domain/types"
	"github.com/spf13/cobra"
	"strconv"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	domainTxCmd := &cobra.Command{
		Use:                        storeKey,
		Short:                      fmt.Sprintf("%s transactions subcommands", storeKey),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	domainTxCmd.AddCommand(flags.PostCommands(
		// TODO: Add tx based commands
		getCmdRegisterDomain(cdc),
	)...)

	return domainTxCmd
}

func getCmdRegisterDomain(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:                        "register [domain-name] [has-superuser] [account-renew]",
		Short:                      "registers a domain",
		SuggestionsMinimumDistance: 2,
		Args:                       cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get superuser
			hasSuperUser, err := strconv.ParseBool(args[1])
			if err != nil {
				return fmt.Errorf("invalid superuser bool: %s", args[1])
			}
			// get account renew time
			accountRenew, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account renew: %s", args[2])
			}
			msg := types.MsgRegisterDomain{
				Name:         args[0],
				Admin:        cliCtx.GetFromAddress(),
				HasSuperuser: hasSuperUser,
				Broker:       nil,
				AccountRenew: accountRenew,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
}
