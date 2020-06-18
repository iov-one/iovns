package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/pkg/index"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryHandlerFunc defines the query handler for this module
type QueryHandlerFunc func(ctx sdk.Context, path []string, query abci.RequestQuery, k Keeper) ([]byte, error)

// AvailableQueries returns the list of available queries in the module
func AvailableQueries() []iovns.QueryHandler {
	queries := []iovns.QueryHandler{
		&QueryAccountsInDomain{},
		&QueryResolveDomain{},
		&QueryResolveAccount{},
		&QueryAccountsWithOwner{},
		&QueryDomainsWithOwner{},
		&QueryTargetAccounts{},
	}
	return queries
}

// queryRouter defines a router for domain queries
type queryRouter map[string]QueryHandlerFunc

func buildRouter(queries []iovns.QueryHandler) queryRouter {
	// queryHandler extends the default query handler
	// to provide also an handler function required to
	// build a router
	type queryHandler interface {
		iovns.QueryHandler
		Handler() QueryHandlerFunc
	}
	// build router
	router := make(queryRouter, len(queries))
	for _, query := range queries {
		queryAndHandler, ok := query.(queryHandler)
		// if interface is not implemented then the query type formation is invalid
		if !ok {
			panic(fmt.Sprintf("invalid query type: %T", query))
		}
		router[queryAndHandler.QueryPath()] = queryAndHandler.Handler()
	}
	// return
	return router
}

// NewQuerier builds the query handler for the module
func NewQuerier(k Keeper) sdk.Querier {
	// get queries
	queries := AvailableQueries()
	router := buildRouter(queries)
	// return sdk.Querier
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		handler, ok := router[path[0]]
		// handler not found, query does not exist
		if !ok {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "%s", path[0])
		}
		// handler
		return handler(ctx, path, req, k)
	}
}

// QueryAccountsInDomain is the request model used to
// query accounts contained in a domain
type QueryAccountsInDomain struct {
	// Domain is the domain name
	Domain string `json:"domain" arg:"positional"`
	// ResultsPerPage is the results that each page should contain
	ResultsPerPage int `json:"results_per_page" arg:"positional"`
	// Offset is the page number
	Offset int `json:"offset" arg:"positional"`
}

// Use is a placeholder
func (q *QueryAccountsInDomain) Use() string {
	return "domain-accounts"
}

// Description is a placeholder
func (q *QueryAccountsInDomain) Description() string {
	return "returns all the accounts contained in a domain"
}

// Handler implements queryHandler
func (q *QueryAccountsInDomain) Handler() QueryHandlerFunc {
	return queryAccountsInDomainHandler
}

// Validate implements iovns.QueryHandler
func (q *QueryAccountsInDomain) Validate() error {
	if q.Domain == "" {
		return sdkerrors.Wrapf(types.ErrInvalidDomainName, "empty")
	}
	// if results per page is unset then use default
	if q.ResultsPerPage <= 0 {
		q.ResultsPerPage = 100
	}
	// if offset is zero then use default
	if q.Offset <= 0 {
		q.Offset = 1
	}
	return nil
}

// QueryPath implements iovns.QueryHandler
func (q *QueryAccountsInDomain) QueryPath() string {
	return "accountsInDomain"
}

// QueryAccountsInDomainResponse is the response model
// returned after a QueryAccountsInDomain query
type QueryAccountsInDomainResponse struct {
	// Accounts is a slice of the accounts found
	Accounts []types.Account `json:"accounts"`
}

// queryAccountsInDomainHandler returns all accounts in aliceAddr domain
func queryAccountsInDomainHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	query := new(QueryAccountsInDomain)
	err := iovns.DefaultQueryDecode(req.Data, query)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	// verify request
	if err = query.Validate(); err != nil {
		return nil, err
	}
	keys := make([][]byte, 0, query.ResultsPerPage)
	i := 0
	// calculate index range
	indexStart := query.ResultsPerPage*query.Offset - query.ResultsPerPage // this is the start
	indexEnd := indexStart + query.ResultsPerPage - 1                      // this is the end
	do := func(key []byte) bool {
		// check if our index is grater-equal than our start
		if i >= indexStart {
			keys = append(keys, key)
		}
		if i == indexEnd {
			return false
		}
		// increase index
		i++
		return true
	}
	// iterate keys
	k.GetAccountsInDomain(ctx, query.Domain, do)
	// get accounts
	accounts := make([]types.Account, 0, len(keys))
	for _, key := range keys {
		account, _ := k.GetAccount(ctx, query.Domain, accountKeyToString(key))
		accounts = append(accounts, account)
	}
	// return response
	respBytes, err := iovns.DefaultQueryEncode(QueryAccountsInDomainResponse{Accounts: accounts})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}

// QueryAccountsWithOwner queries all the accounts
// owned by a certain sdk.AccAddress
type QueryAccountsWithOwner struct {
	// Owner is the owner of the accounts
	Owner sdk.AccAddress `json:"owner"`
	// ResultsPerPage is the number of results returned in each page
	ResultsPerPage int `json:"results_per_page"`
	// Offset is the page number
	Offset int `json:"offset"`
}

// Use is a placeholder
func (q *QueryAccountsWithOwner) Use() string {
	return "owner-accounts"
}

// Description is a placeholder
func (q *QueryAccountsWithOwner) Description() string {
	return "gets all the accounts owned by a given address"
}

// Handler implements local queryHandler
func (q *QueryAccountsWithOwner) Handler() QueryHandlerFunc {
	return queryAccountsWithOwnerHandler
}

// QueryPath implements iovns.QueryHandler
func (q *QueryAccountsWithOwner) QueryPath() string {
	return "accountsWithOwner"
}

// Validate implements iovns.QueryHandler
func (q *QueryAccountsWithOwner) Validate() error {
	if q.Owner == nil {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "empty")
	}
	if q.ResultsPerPage == 0 {
		q.ResultsPerPage = 100
	}
	if q.Offset == 0 {
		q.Offset = 1
	}
	return nil
}

// QueryAccountsWithOwnerResponse is the response model
// returned by QueryAccountsWithOwner
type QueryAccountsWithOwnerResponse struct {
	// Accounts is a slice containing the accounts
	// returned by the query
	Accounts []types.Account `json:"accounts"`
}

// queryAccountsWithOwnerHandler gets all the accounts related to an account address
func queryAccountsWithOwnerHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	query := new(QueryAccountsWithOwner)
	err := iovns.DefaultQueryDecode(req.Data, query)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	// validate request
	if err = query.Validate(); err != nil {
		return nil, err
	}
	// generate expected keys
	keys := make([][]byte, 0, query.ResultsPerPage)
	i := 0
	// calculate index range
	indexStart := query.ResultsPerPage*query.Offset - query.ResultsPerPage // this is the start
	indexEnd := indexStart + query.ResultsPerPage - 1                      // this is the end
	do := func(key []byte) bool {
		// check if our index is grater-equal than our start
		if i >= indexStart {
			keys = append(keys, key)
		}
		if i == indexEnd {
			return false
		}
		// increase index
		i++
		return true
	}
	// iterate account keys
	err = k.iterAccountToOwner(ctx, query.Owner, do)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	// check if there are any keys
	if len(keys) == 0 {
		respBytes, err := iovns.DefaultQueryEncode(QueryAccountsWithOwnerResponse{})
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return respBytes, nil
	}
	// fill accounts
	accounts := make([]types.Account, 0, len(keys))
	for _, accKey := range keys {
		// get indexed account
		indexedAccount := types.Account{}
		err = index.Unpack(accKey, &indexedAccount)
		if err != nil {
			panic(err)
		}
		account, _ := k.GetAccount(ctx, indexedAccount.Domain, indexedAccount.Name)
		accounts = append(accounts, account)
	}
	// return response
	respBytes, err := iovns.DefaultQueryEncode(QueryAccountsWithOwnerResponse{Accounts: accounts})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}

// QueryDomainsWithOwner is the request model used
// to query domains owned by a sdk.AccAddress
type QueryDomainsWithOwner struct {
	// Owner is the address of the owner of the domains
	Owner sdk.AccAddress `json:"owner"`
	// ResultsPerPage is the number of results displayed in a page
	ResultsPerPage int `json:"results_per_page"`
	// Offset is the page number
	Offset int `json:"offset"`
}

// Use is a placeholder
func (q *QueryDomainsWithOwner) Use() string {
	return "owner-domains"
}

// Description is a placeholder
func (q *QueryDomainsWithOwner) Description() string {
	return "gets all the domains owned by the given address"
}

// Handler implements the local queryHandler
func (q *QueryDomainsWithOwner) Handler() QueryHandlerFunc {
	return queryDomainsWithOwnerHandler
}

// QueryPath implements iovns.QueryHandler
func (q *QueryDomainsWithOwner) QueryPath() string {
	return "domainsWithOwner"
}

// Validate implements iovns.QueryHandler
func (q *QueryDomainsWithOwner) Validate() error {
	if q.Owner == nil {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "empty")
	}
	if q.ResultsPerPage == 0 {
		q.ResultsPerPage = 100
	}
	if q.Offset == 0 {
		q.Offset = 1
	}
	return nil
}

// QueryDomainsWithOwnerResponse is the response
// returned by the QueryDomainsWithOwner query
type QueryDomainsWithOwnerResponse struct {
	// Domains is a slice of the domains
	// found by the query
	Domains []types.Domain
}

// queryDomainsWithOwnerHandler is the query handler used to get all the domains owned by an sdk.AccAddress
func queryDomainsWithOwnerHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	query := new(QueryDomainsWithOwner)
	err := iovns.DefaultQueryDecode(req.Data, query)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	// validate request
	if err = query.Validate(); err != nil {
		return nil, err
	}
	// get domain keys
	// generate expected keys
	keys := make([][]byte, 0, query.ResultsPerPage)
	i := 0
	// calculate i range
	indexStart := query.ResultsPerPage*query.Offset - query.ResultsPerPage // this is the start
	indexEnd := indexStart + query.ResultsPerPage - 1                      // this is the end
	do := func(key []byte) bool {
		// check if our i is grater-equal than our start
		if i >= indexStart {
			keys = append(keys, key)
		}
		if i == indexEnd {
			return false
		}
		// increase i
		i++
		return true
	}
	// fill domain keys
	err = k.iterDomainToOwner(ctx, query.Owner, do)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}
	// get domains
	domains := make([]types.Domain, 0, len(keys))
	for _, key := range keys {
		// prepare the domain index key unpacker
		indexedDomain := types.Domain{}
		err = index.Unpack(key, &indexedDomain)
		// if it's not found or unpacked then we need to panic
		if err != nil {
			panic("unexpected domain key not found: " + err.Error())
		}
		// get domain
		domain, _ := k.GetDomain(ctx, indexedDomain.Name)
		domains = append(domains, domain)
	}
	respBytes, err := iovns.DefaultQueryEncode(QueryDomainsWithOwnerResponse{Domains: domains})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}

// QueryResolveAccount is the query model
// used to resolve an account
type QueryResolveAccount struct {
	// Domain is the name of the domain of the account
	Domain string `json:"domain"`
	// Name is the name of the account
	Name string `json:"name"`
	// Starname is the representation of an account in domain*name format
	Starname string `json:"starname"`
}

// Use is a placeholder
func (q *QueryResolveAccount) Use() string {
	return "resolve-account"
}

// Description is a placeholder
func (q *QueryResolveAccount) Description() string {
	return "resolves the given account"
}

// Handler implements local queryHandler
func (q *QueryResolveAccount) Handler() QueryHandlerFunc {
	return queryResolveAccountHandler
}

// Validate implements iovns.QueryHandler
func (q *QueryResolveAccount) Validate() error {
	if q.Starname != "" && (q.Domain != "" || q.Name != "") {
		return types.ErrProvideStarnameOrDomainName
	}

	if q.Starname != "" {
		if !strings.Contains(q.Starname, string(iovns.Separator)) {
			return types.ErrStarnameNotContainSep
		}
		sname := strings.Split(q.Starname, string(iovns.Separator))
		if len(sname) != 2 {
			return types.ErrStarnameMultipleSeparator
		}
		q.Name = sname[0]
		q.Domain = sname[1]
		q.Starname = ""
	}

	if q.Domain == "" {
		return sdkerrors.Wrapf(types.ErrInvalidDomainName, "empty")
	}
	return nil
}

// QueryPath implements iovns.QueryHandler
func (q *QueryResolveAccount) QueryPath() string {
	return "resolve"
}

// QueryResolveAccountResponse is the response
// returned by the QueryResolveAccount query
type QueryResolveAccountResponse struct {
	// Account contains the resolved account
	Account types.Account `json:"account"`
}

// queryResolveAccountHandler is the query handler that takes care of resolving accounts
func queryResolveAccountHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	q := new(QueryResolveAccount)
	err := iovns.DefaultQueryDecode(req.Data, q)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	// validate
	if err = q.Validate(); err != nil {
		return nil, err
	}
	// do query
	account, exists := k.GetAccount(ctx, q.Domain, q.Name)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found in domain %s: %s", q.Domain, q.Name)
	}
	// return response
	respBytes, err := iovns.DefaultQueryEncode(QueryResolveAccountResponse{Account: account})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}

// QueryResolveDomain is the request model
// used to resolve a domain
type QueryResolveDomain struct {
	// Name is the domain name
	Name string `json:"name" arg:"positional"`
}

// Use is a placeholder
func (q *QueryResolveDomain) Use() string {
	return "resolve-domain"
}

// Description is a placeholder
func (q *QueryResolveDomain) Description() string {
	return "resolves a domain"
}

// Handler implements the local queryHandler
func (q *QueryResolveDomain) Handler() QueryHandlerFunc {
	return queryResolveDomainHandler
}

// QueryPath implements iovns.QueryHandler
func (q *QueryResolveDomain) QueryPath() string {
	return "domainInfo"
}

// Validate implements iovns.QueryHandler
func (q *QueryResolveDomain) Validate() error {
	if q.Name == "" {
		return sdkerrors.Wrapf(types.ErrInvalidDomainName, "empty")
	}
	return nil
}

// QueryResolveDomainResponse is response returned
// by the QueryResolveDomain query
type QueryResolveDomainResponse struct {
	// Domain contains the queried domain information
	Domain types.Domain `json:"domain"`
}

// queryResolveDomainHandler takes care of resolving domains
func queryResolveDomainHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	q := new(QueryResolveDomain)
	err := iovns.DefaultQueryDecode(req.Data, q)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if err = q.Validate(); err != nil {
		return nil, err
	}
	domain, ok := k.GetDomain(ctx, q.Name)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", q.Name)
	}
	// return response
	respBytes, err := iovns.DefaultQueryEncode(QueryResolveDomainResponse{Domain: domain})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}

//  QueryTargetAccounts gets accounts from a target
type QueryTargetAccounts struct {
	// Target is the blockchain target we want to query
	Target types.BlockchainAddress `json:"target"`
	// ResultsPerPage is the number of results displayed in a page
	ResultsPerPage int `json:"results_per_page"`
	// Offset is the page number
	Offset int `json:"offset"`
}

// QueryPath implements iovns.QueryHandler
func (q *QueryTargetAccounts) QueryPath() string {
	return "targetAccounts"
}

func (q *QueryTargetAccounts) Validate() error {
	if q.Target.Address == "" {
		return sdkerrors.Wrapf(types.ErrInvalidBlockchainTarget, "empty address")
	}
	if q.Target.ID == "" {
		return sdkerrors.Wrapf(types.ErrInvalidBlockchainTarget, "empty ID")
	}
	if q.ResultsPerPage == 0 {
		q.ResultsPerPage = 100
	}
	if q.Offset == 0 {
		q.Offset = 1
	}

	return nil
}

func (q *QueryTargetAccounts) Handler() QueryHandlerFunc {
	return queryTargetAccountsHandler
}

// QueryTargetAccountsResponse is the response returned by query target
type QueryTargetAccountsResponse struct {
	Accounts []types.Account `json:"accounts"`
}

func queryTargetAccountsHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	q := new(QueryTargetAccounts)
	// decode query
	err := iovns.DefaultQueryDecode(req.Data, q)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	// validate
	if err := q.Validate(); err != nil {
		return nil, err
	}
	// generate expected keys
	keys := make([][]byte, 0, q.ResultsPerPage)
	// calculate index range
	indexStart := q.ResultsPerPage*q.Offset - q.ResultsPerPage // start index
	indexEnd := indexStart + q.ResultsPerPage - 1              // index end
	i := 0
	do := func(key []byte) bool {
		// check if our index is grater-equal than our start
		if i >= indexStart {
			keys = append(keys, key)
		}
		if i == indexEnd {
			return false
		}
		// increase index
		i++
		return true
	}
	// fill keys
	err = k.iterateBlockchainTargetsAccounts(ctx, q.Target, do)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}
	// get accounts
	accounts := make([]types.Account, 0, len(keys))
	for _, key := range keys {
		// get account name
		account := types.Account{}
		index.MustUnpack(key, &account)
		account, ok := k.GetAccount(ctx, account.Domain, account.Name)
		if !ok {
			panic(fmt.Sprintf("account was not found for key %x", key))
		}
		accounts = append(accounts, account)
	}
	// return response
	b, err := iovns.DefaultQueryEncode(QueryTargetAccountsResponse{Accounts: accounts})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return b, nil
}
