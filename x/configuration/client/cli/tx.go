package cli

import (
	"bufio"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/iov-one/iovns/x/configuration/types"
	"github.com/spf13/cobra"
)

// GetTxCmd clubs together all the CLI tx commands
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	configTxCmd := &cobra.Command{
		Use:                        storeKey,
		Short:                      fmt.Sprintf("%s transactions subcommands", storeKey),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	configTxCmd.AddCommand(flags.PostCommands(
		getCmdUpdateConfig(cdc),
		getCmdUpsertDefaultFee(cdc),
		getCmdUpsertLevelFee(cdc),
		getCmdDeleteLevelFee(cdc),
	)...)
	return configTxCmd
}

func getCmdUpdateConfig(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-config",
		Short: "update domain configuration",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			configurerStr, err := cmd.Flags().GetString("configurer")
			if err != nil {
				return
			}
			configurer, err := sdk.AccAddressFromBech32(configurerStr)
			if err != nil {
				return
			}

			validDomainName, err := cmd.Flags().GetString("valid-domain-name")
			if err != nil {
				return err
			}

			validName, err := cmd.Flags().GetString("valid-account-name")
			if err != nil {
				return err
			}

			validBlockchainID, err := cmd.Flags().GetString("valid-blockchain-id")
			if err != nil {
				return err
			}

			validBlockchainAddress, err := cmd.Flags().GetString("valid-blockchain-address")
			if err != nil {
				return err
			}

			domainRenew, err := cmd.Flags().GetDuration("domain-renew")
			if err != nil {
				return err
			}
			domainRenewCountMax, err := cmd.Flags().GetUint32("domain-renew-count-max")
			if err != nil {
				return err
			}
			domainGracePeriod, err := cmd.Flags().GetDuration("domain-grace-period")
			if err != nil {
				return err
			}
			accountRenewPeriod, err := cmd.Flags().GetDuration("account-renew-period")
			if err != nil {
				return err
			}
			accountRenewCountMax, err := cmd.Flags().GetUint32("account-renew-count-max")
			if err != nil {
				return err
			}
			accountGracePeriod, err := cmd.Flags().GetDuration("account-grace-period")
			if err != nil {
				return err
			}
			blockchainTargetMax, err := cmd.Flags().GetUint32("blockchain-target-max")
			if err != nil {
				return err
			}
			certificateSizeMax, err := cmd.Flags().GetUint64("certificate-size-max")
			if err != nil {
				return err
			}
			certificateCountMax, err := cmd.Flags().GetUint32("certificate-count-max")
			if err != nil {
				return err
			}
			metadataSizeMax, err := cmd.Flags().GetUint64("metadata-size-max")
			if err != nil {
				return err
			}

			config := types.Config{
				Configurer:             configurer,
				ValidDomainName:        validDomainName,
				ValidAccountName:       validName,
				ValidBlockchainID:      validBlockchainID,
				ValidBlockchainAddress: validBlockchainAddress,
				DomainRenewalPeriod:    domainRenew,
				DomainRenewalCountMax:  domainRenewCountMax,
				DomainGracePeriod:      domainGracePeriod,
				AccountRenewalPeriod:   accountRenewPeriod,
				AccountRenewalCountMax: accountRenewCountMax,
				AccountGracePeriod:     accountGracePeriod,
				BlockchainTargetMax:    blockchainTargetMax,
				CertificateSizeMax:     certificateSizeMax,
				CertificateCountMax:    certificateCountMax,
				MetadataSizeMax:        metadataSizeMax,
			}
			if err := config.Validate(); err != nil {
				return err
			}
			// build msg
			msg := &types.MsgUpdateConfig{
				Configurer:       cliCtx.GetFromAddress(),
				NewConfiguration: config,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
	defaultDuration, _ := time.ParseDuration("1h")
	// add flags
	cmd.Flags().String("configurer", "", "configurer in bech32 format")
	cmd.Flags().String("valid-domain-name", "", "regexp that determines if domain name is valid or not")
	cmd.Flags().String("valid-account-name", "", "regexp that determines if account name is valid or not")
	cmd.Flags().String("valid-blockchain-id", "", "regexp that determines if blockchain id is valid or not")
	cmd.Flags().String("valid-blockchain-address", "", "regexp that determines if blockchain address is valid or not")

	cmd.Flags().Duration("domain-renew-period", defaultDuration, "domain renewal duration in seconds before expiration")
	cmd.Flags().Uint32("domain-renew-count-max", 10, "maximum number of applicable domain renewals")
	cmd.Flags().Duration("domain-grace-period", defaultDuration, "domain grace period duration in seconds")

	cmd.Flags().Duration("account-renew-period", defaultDuration, "domain renewal duration in seconds before expiration")
	cmd.Flags().Uint32("account-renew-count-max", 10, "maximum number of applicable account renewals")
	cmd.Flags().Duration("account-grace-period", defaultDuration, "account grace period duration in seconds")

	cmd.Flags().Uint32("blockchain-target-max", 15, "maximum number of blockchain targets could be saved under an account")
	cmd.Flags().Uint64("certificate-size-max", 10, "maximum size of a certificate that could be saved under an account")
	cmd.Flags().Uint32("certificate-count-max", 15, "maximum number of certificates that could be saved under an account")
	cmd.Flags().Uint64("metadata-size-max", 10, "maximum size of metadata that could be saved under an account")
	return cmd
}

func getCmdUpsertDefaultFee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-default-fee",
		Short: "upsert default fee configuration",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			configurerStr, err := cmd.Flags().GetString("configurer")
			if err != nil {
				return
			}
			configurer, err := sdk.AccAddressFromBech32(configurerStr)
			if err != nil {
				return
			}

			module, err := cmd.Flags().GetString("module")
			if err != nil {
				return err
			}

			msgType, err := cmd.Flags().GetString("msg-type")
			if err != nil {
				return err
			}

			feeStr, err := cmd.Flags().GetString("fee")
			fee, err := sdk.ParseCoin(feeStr)
			if err != nil {
				return err
			}

			// build msg
			msg := &types.MsgUpsertDefaultFee{
				Configurer: configurer,
				Module:     module,
				MsgType:    msgType,
				Fee:        fee,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
	// add flags
	cmd.Flags().String("configurer", "", "configurer in bech32 format")
	cmd.Flags().String("module", "", "what is this?")
	cmd.Flags().String("msg-type", "", "type of the message")
	cmd.Flags().String("fee", "10iov", "amount of the fee")
	return cmd
}

func getCmdUpsertLevelFee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-level-fee",
		Short: "upsert level fee configuration",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			configurerStr, err := cmd.Flags().GetString("configurer")
			if err != nil {
				return
			}
			configurer, err := sdk.AccAddressFromBech32(configurerStr)
			if err != nil {
				return
			}

			module, err := cmd.Flags().GetString("module")
			if err != nil {
				return err
			}

			msgType, err := cmd.Flags().GetString("msg-type")
			if err != nil {
				return err
			}

			feeStr, err := cmd.Flags().GetString("fee")
			fee, err := sdk.ParseCoin(feeStr)
			if err != nil {
				return err
			}

			// build msg
			msg := &types.MsgUpsertLevelFee{
				Configurer: configurer,
				Module:     module,
				MsgType:    msgType,
				Fee:        fee,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
	// add flags
	cmd.Flags().String("configurer", "", "configurer in bech32 format")
	cmd.Flags().String("module", "", "what is this?")
	cmd.Flags().String("msg-type", "", "type of the message")
	cmd.Flags().String("fee", "10iov", "amount of the fee")
	return cmd
}

func getCmdDeleteLevelFee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-level-fee",
		Short: "delete level fee configuration",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get flags
			configurerStr, err := cmd.Flags().GetString("configurer")
			if err != nil {
				return
			}
			configurer, err := sdk.AccAddressFromBech32(configurerStr)
			if err != nil {
				return
			}

			module, err := cmd.Flags().GetString("module")
			if err != nil {
				return err
			}

			msgType, err := cmd.Flags().GetString("msg-type")
			if err != nil {
				return err
			}

			// build msg
			msg := &types.MsgDeleteLevelFee{
				Configurer: configurer,
				Module:     module,
				MsgType:    msgType,
			}
			// check if valid
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			// broadcast request
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
	// add flags
	cmd.Flags().String("configurer", "", "configurer in bech32 format")
	cmd.Flags().String("module", "", "what is this?")
	cmd.Flags().String("msg-type", "", "type of the message")
	cmd.Flags().String("fee", "10iov", "amount of the fee")
	return cmd
}
