package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd builds the commands for queries in the domain module
func GetQueryCmd(moduleQueryPath string, cdc *codec.Codec) *cobra.Command {
	domainQueryCmd := &cobra.Command{
		Use:                        moduleQueryPath, // store key is same as module name
		Short:                      "querying commands for the domain module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
		Aliases:                    []string{"starname"},
	}
	domainQueryCmd.AddCommand(
		flags.GetCommands(
			getQueryResolveDomain(moduleQueryPath, cdc),
			getQueryResolveAccount(moduleQueryPath, cdc),
			getQueryDomainAccounts(moduleQueryPath, cdc),
			getQueryOwnerAccount(moduleQueryPath, cdc),
			getQueryOwnerDomain(moduleQueryPath, cdc),
			getQueryTargetAccounts(moduleQueryPath, cdc),
		)...,
	)
	return domainQueryCmd
}

func processQueryCmd(cdc *codec.Codec, path string, q interface{}, _ interface{}) (err error) {
	// get req byres
	b, err := iovns.DefaultQueryEncode(q)
	if err != nil {
		return
	}
	// get cli ctx
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	res, _, err := cliCtx.QueryWithData(path, b)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(cliCtx.Output, "%s\n\n", res)
	return err
}

func getQueryResolveDomain(modulePath string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain-info",
		Short: "resolve a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			// get query & validate
			q := keeper.QueryResolveDomain{Name: domain}
			if err = q.Validate(); err != nil {
				return err
			}
			// get query path
			path := fmt.Sprintf("custom/%s/%s", modulePath, q.QueryPath())
			return processQueryCmd(cdc, path, q, new(keeper.QueryResolveDomainResponse))
		},
	}
	// add flags
	cmd.Flags().String("domain", "", "the domain name you want to resolve")
	// return cmd
	return cmd
}

func getQueryDomainAccounts(modulePath string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain-accounts",
		Short: "get accounts in a domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			rpp, err := cmd.Flags().GetInt("rpp")
			if err != nil {
				return err
			}
			offset, err := cmd.Flags().GetInt("offset")
			if err != nil {
				return err
			}
			// get query & validate
			q := keeper.QueryAccountsInDomain{
				Domain:         domain,
				ResultsPerPage: rpp,
				Offset:         offset,
			}
			if err = q.Validate(); err != nil {
				return err
			}
			// get query path
			path := fmt.Sprintf("custom/%s/%s", modulePath, q.QueryPath())
			return processQueryCmd(cdc, path, q, new(keeper.QueryAccountsInDomainResponse))
		},
	}
	// add flags
	cmd.Flags().String("domain", "", "the domain name you want to resolve")
	cmd.Flags().Int("offset", 1, "the page offset")
	cmd.Flags().Int("rpp", 100, "results per page")
	// return cmd
	return cmd
}

func getQueryOwnerAccount(modulePath string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-accounts",
		Short: "get accounts owned by an address",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// get flags
			owner, err := cmd.Flags().GetString("owner")
			if err != nil {
				return err
			}
			// verify if address is correct
			accAddress, err := sdk.AccAddressFromBech32(owner)
			rpp, err := cmd.Flags().GetInt("rpp")
			if err != nil {
				return err
			}
			offset, err := cmd.Flags().GetInt("offset")
			if err != nil {
				return err
			}
			// get query & validate
			q := keeper.QueryAccountsWithOwner{
				Owner:          accAddress,
				ResultsPerPage: rpp,
				Offset:         offset,
			}
			if err = q.Validate(); err != nil {
				return err
			}
			// get query path
			path := fmt.Sprintf("custom/%s/%s", modulePath, q.QueryPath())
			return processQueryCmd(cdc, path, q, new(keeper.QueryAccountsWithOwnerResponse))
		},
	}
	// add flags
	cmd.Flags().String("owner", "", "the bech32 address of the owner you want to lookup")
	cmd.Flags().Int("offset", 1, "the page offset")
	cmd.Flags().Int("rpp", 100, "results per page")
	// return cmd
	return cmd
}

func getQueryOwnerDomain(modulePath string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-domains",
		Short: "get domains owned by an address",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// get flags
			owner, err := cmd.Flags().GetString("owner")
			if err != nil {
				return err
			}
			// verify if address is correct
			accAddress, err := sdk.AccAddressFromBech32(owner)
			rpp, err := cmd.Flags().GetInt("rpp")
			if err != nil {
				return err
			}
			offset, err := cmd.Flags().GetInt("offset")
			if err != nil {
				return err
			}
			// get query & validate
			q := keeper.QueryDomainsWithOwner{
				Owner:          accAddress,
				ResultsPerPage: rpp,
				Offset:         offset,
			}
			if err = q.Validate(); err != nil {
				return err
			}
			// get query path
			path := fmt.Sprintf("custom/%s/%s", modulePath, q.QueryPath())
			return processQueryCmd(cdc, path, q, new(keeper.QueryDomainsWithOwnerResponse))
		},
	}
	// add flags
	cmd.Flags().String("owner", "", "the bech32 address of the owner you want to lookup")
	cmd.Flags().Int("offset", 1, "the page offset")
	cmd.Flags().Int("rpp", 100, "results per page")
	// return cmd
	return cmd
}

func getQueryResolveAccount(modulePath string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "resolve an account, provide either starname or name/domain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// get flags
			domain, err := cmd.Flags().GetString("domain")
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}
			starname, err := cmd.Flags().GetString("starname")
			if err != nil {
				return err
			}
			// get query & validate
			q := keeper.QueryResolveAccount{
				Domain:   domain,
				Name:     name,
				Starname: starname,
			}
			if err = q.Validate(); err != nil {
				return err
			}
			// get query path
			path := fmt.Sprintf("custom/%s/%s", modulePath, q.QueryPath())
			return processQueryCmd(cdc, path, q, new(keeper.QueryResolveAccountResponse))
		},
	}
	// add flags
	cmd.Flags().String("starname", "", "the starname representation of the account")
	cmd.Flags().String("domain", "", "the domain name of the account")
	cmd.Flags().String("name", "", "the name of the account you want to resolve")
	// return cmd
	return cmd
}

func getQueryTargetAccounts(modulePath string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve-target",
		Short: "resolves a blockchain target into accounts",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// get flags
			id, err := cmd.Flags().GetString("blockchain-id")
			if err != nil {
				return err
			}
			addr, err := cmd.Flags().GetString("blockchain-address")
			if err != nil {
				return err
			}
			rpp, err := cmd.Flags().GetInt("rpp")
			if err != nil {
				return err
			}
			offset, err := cmd.Flags().GetInt("offset")
			if err != nil {
				return err
			}
			// get query & validate
			q := keeper.QueryTargetAccounts{
				Target: types.BlockchainAddress{
					ID:      id,
					Address: addr,
				},
				ResultsPerPage: rpp,
				Offset:         offset,
			}
			if err = q.Validate(); err != nil {
				return err
			}
			// get query path
			path := fmt.Sprintf("custom/%s/%s", modulePath, q.QueryPath())
			return processQueryCmd(cdc, path, q, new(keeper.QueryTargetAccountsResponse))
		},
	}
	// add flags
	cmd.Flags().String("blockchain-id", "", "the blockchain id")
	cmd.Flags().String("blockchain-address", "", "blockchain address")
	cmd.Flags().Int("offset", 1, "the page offset")
	cmd.Flags().Int("rpp", 100, "results per page")
	// return cmd
	return cmd
}
