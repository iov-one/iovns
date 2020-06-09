package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
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
		getCmdUpdateFees(cdc),
	)...)
	return configTxCmd
}

var defaultDuration, _ = time.ParseDuration("1h")

const defaultRegex = "^(.*?)?"
const defaultNumber = 1

func getCmdUpdateFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-fees",
		Short: "update fees using a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			// get fees file
			feeFile, err := cmd.Flags().GetString("fees-file")
			if err != nil {
				return err
			}
			f, err := os.Open(feeFile)
			if err != nil {
				return fmt.Errorf("unable to open fee file: %s", err)
			}
			defer f.Close()
			newFees := new(types.Fees)
			err = json.NewDecoder(f).Decode(newFees)
			if err != nil {
				return err
			}
			msg := types.MsgUpdateFees{
				Fees:       newFees,
				Configurer: cliCtx.GetFromAddress(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid tx: %w", err)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String("fee-file", "fees.json", "fees file in json format")
	return cmd
}

func getCmdUpdateConfig(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-config",
		Short: "update domain configuration, provide the values you want to override in current configuration",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBuilder := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			rawCfg, _, err := cliCtx.QueryStore([]byte(types.ConfigKey), types.StoreKey)
			if err != nil {
				return err
			}
			var config types.Config
			cdc.MustUnmarshalBinaryBare(rawCfg, &config)

			// get flags
			configurerStr, err := cmd.Flags().GetString("configurer")
			if err != nil {
				return
			}
			if configurerStr != "" {
				configurer, err := sdk.AccAddressFromBech32(configurerStr)
				if err != nil {
					return err
				}
				config.Configurer = configurer
			}
			validDomainName, err := cmd.Flags().GetString("valid-domain-name")
			if err != nil {
				return err
			}
			if validDomainName != defaultRegex {
				config.ValidDomainName = validDomainName
			}
			validAccountName, err := cmd.Flags().GetString("valid-account-name")
			if err != nil {
				return err
			}
			if validAccountName != defaultRegex {
				config.ValidAccountName = validAccountName
			}
			validBlockchainID, err := cmd.Flags().GetString("valid-blockchain-id")
			if err != nil {
				return err
			}
			if validBlockchainID != defaultRegex {
				config.ValidBlockchainID = validBlockchainID
			}
			validBlockchainAddress, err := cmd.Flags().GetString("valid-blockchain-address")
			if err != nil {
				return err
			}
			if validBlockchainAddress != defaultRegex {
				config.ValidBlockchainAddress = validBlockchainAddress
			}
			domainRenew, err := cmd.Flags().GetDuration("domain-renew-period")
			if err != nil {
				return err
			}
			if domainRenew != defaultDuration {
				config.DomainRenewalPeriod = domainRenew
			}
			domainRenewCountMax, err := cmd.Flags().GetUint32("domain-renew-count-max")
			if err != nil {
				return err
			}
			if domainRenewCountMax != defaultNumber {
				config.DomainRenewalCountMax = domainRenewCountMax
			}
			domainGracePeriod, err := cmd.Flags().GetDuration("domain-grace-period")
			if err != nil {
				return err
			}
			if domainGracePeriod != defaultNumber {
				config.DomainGracePeriod = domainGracePeriod
			}
			accountRenewPeriod, err := cmd.Flags().GetDuration("account-renew-period")
			if err != nil {
				return err
			}
			if accountRenewPeriod != defaultNumber {
				config.AccountRenewalPeriod = accountRenewPeriod
			}
			accountRenewCountMax, err := cmd.Flags().GetUint32("account-renew-count-max")
			if err != nil {
				return err
			}
			if accountRenewCountMax != defaultNumber {
				config.AccountRenewalCountMax = accountRenewCountMax
			}
			accountGracePeriod, err := cmd.Flags().GetDuration("account-grace-period")
			if err != nil {
				return err
			}
			if accountGracePeriod != defaultDuration {
				config.AccountGracePeriod = accountGracePeriod
			}
			blockchainTargetMax, err := cmd.Flags().GetUint32("blockchain-target-max")
			if err != nil {
				return err
			}
			if blockchainTargetMax != defaultNumber {
				config.BlockchainTargetMax = blockchainTargetMax
			}
			certificateSizeMax, err := cmd.Flags().GetUint64("certificate-size-max")
			if err != nil {
				return err
			}
			if certificateSizeMax != defaultNumber {
				config.CertificateSizeMax = certificateSizeMax
			}
			certificateCountMax, err := cmd.Flags().GetUint32("certificate-count-max")
			if err != nil {
				return err
			}
			if certificateCountMax != defaultNumber {
				config.CertificateCountMax = certificateCountMax
			}
			metadataSizeMax, err := cmd.Flags().GetUint64("metadata-size-max")
			if err != nil {
				return err
			}
			if metadataSizeMax != defaultNumber {
				config.MetadataSizeMax = metadataSizeMax
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
	// add flags
	cmd.Flags().String("configurer", "", "configurer in bech32 format")
	cmd.Flags().String("valid-domain-name", defaultRegex, "regexp that determines if domain name is valid or not")
	cmd.Flags().String("valid-account-name", defaultRegex, "regexp that determines if account name is valid or not")
	cmd.Flags().String("valid-blockchain-id", defaultRegex, "regexp that determines if blockchain id is valid or not")
	cmd.Flags().String("valid-blockchain-address", defaultRegex, "regexp that determines if blockchain address is valid or not")

	cmd.Flags().Duration("domain-renew-period", defaultDuration, "domain renewal duration in seconds before expiration")
	cmd.Flags().Uint32("domain-renew-count-max", uint32(defaultNumber), "maximum number of applicable domain renewals")
	cmd.Flags().Duration("domain-grace-period", defaultDuration, "domain grace period duration in seconds")

	cmd.Flags().Duration("account-renew-period", defaultDuration, "domain renewal duration in seconds before expiration")
	cmd.Flags().Uint32("account-renew-count-max", uint32(defaultNumber), "maximum number of applicable account renewals")
	cmd.Flags().Duration("account-grace-period", defaultDuration, "account grace period duration in seconds")

	cmd.Flags().Uint32("blockchain-target-max", uint32(defaultNumber), "maximum number of blockchain targets could be saved under an account")
	cmd.Flags().Uint64("certificate-size-max", uint64(defaultNumber), "maximum size of a certificate that could be saved under an account")
	cmd.Flags().Uint32("certificate-count-max", uint32(defaultNumber), "maximum number of certificates that could be saved under an account")
	cmd.Flags().Uint64("metadata-size-max", uint64(defaultNumber), "maximum size of metadata that could be saved under an account")
	return cmd
}
