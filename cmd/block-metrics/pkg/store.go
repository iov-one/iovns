package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/starname/types"

	"github.com/pkg/errors"
)

// dbTx is a database transaction used to batch inserts/updates.
// It Begin()s in BatchBegin() and is committed and reassigned in BatchCommit().
var dbTx *sql.Tx

// NewStore returns a store that provides an access to our database.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

type Store struct {
	db *sql.DB
}

type Block struct {
	Height  int64     `json:"height"`
	Hash    string    `json:"hash"`
	Time    time.Time `json:"time"`
	FeeFrac uint64    `json:"fee_frac"`
}

// convert an sdk.AccAddress to a string
func a2s(addr sdk.AccAddress) string {
	return sdk.AccAddress(addr).String()
}

// get accounts.id given a domain and a name
func getAccountID(ctx context.Context, domain string, name string) (int64, error) {
	var accountID int64
	err := dbTx.QueryRowContext(ctx, `
		SELECT id FROM accounts
		WHERE domain_id = (SELECT MAX(id) FROM domains  WHERE name = $1)
		AND          id = (SELECT MAX(id) FROM accounts WHERE name = $2)
	`, domain, name).Scan(&accountID)
	return accountID, castPgErr(err)
}

func (st *Store) RegisterDomain(ctx context.Context, msg *types.MsgRegisterDomain, height int64) (int64, error) {
	// create the domain...
	var id int64
	err := dbTx.QueryRowContext(ctx, `
		INSERT INTO domains (name, admin, type)
		VALUES ($1, $2, $3)
		RETURNING id
	`, msg.Name, a2s(msg.Admin), msg.DomainType).Scan(&id)
	if err != nil {
		return 0, castPgErr(err)
	}

	// ...and then the empty account
	return st.RegisterAccount(ctx, &types.MsgRegisterAccount{
		Domain:       msg.Name,
		Name:         "",
		Owner:        msg.Admin,
		Broker:       msg.Broker,
		FeePayerAddr: msg.FeePayerAddr,
	}, height)
}

func (st *Store) DeleteDomain(ctx context.Context, msg *types.MsgDeleteDomain, height int64) (int64, error) {
	accountID, err := st.DeleteAccount(ctx, &types.MsgDeleteAccount{
		Domain:       msg.Domain,
		Name:         "",
		Owner:        msg.Owner,
		FeePayerAddr: msg.FeePayerAddr,
	}, height)
	if err == nil {
		_, err = dbTx.ExecContext(ctx, `
			UPDATE domains
			SET deleted = (SELECT block_time FROM blocks WHERE block_height=$1)
			WHERE id = (SELECT MAX(id) FROM domains WHERE name = $2)
		`, height, msg.Domain)
	}
	return accountID, castPgErr(err)
}

func (st *Store) TransferDomain(ctx context.Context, msg *types.MsgTransferDomain, height int64) (int64, error) {
	// update the empty account...
	accountID, err := st.TransferAccount(ctx, &types.MsgTransferAccount{
		Domain:       msg.Domain,
		Name:         "",
		Owner:        msg.Owner,
		NewOwner:     msg.NewAdmin,
		FeePayerAddr: msg.FeePayerAddr,
		Reset:        msg.TransferFlag != types.TransferResetNone,
	}, height)
	if err == nil {
		// ...and then the domain...
		_, err = dbTx.ExecContext(ctx, `
			UPDATE domains
			SET admin = $1
			WHERE id = (SELECT MAX(id) FROM domains WHERE name = $2)
		`, a2s(msg.NewAdmin), msg.Domain)
		if err != nil {
			return accountID, castPgErr(err)
		}
		// ...and with the different transfer flags
		switch msg.TransferFlag {
		case types.TransferResetNone:
			// no-op
		case types.TransferFlush:
			_, err = dbTx.ExecContext(ctx, `
				UPDATE accounts
				SET deleted = (SELECT block_time FROM blocks WHERE block_height=$3)
				WHERE owner = $1
				AND domain_id = (SELECT MAX(id) FROM domains WHERE name = $2)
			`, a2s(msg.Owner), msg.Domain, height)
		case types.TransferOwned:
			_, err = dbTx.ExecContext(ctx, `
				UPDATE accounts
				SET owner = $3
				WHERE owner = $1
				AND domain_id = (SELECT MAX(id) FROM domains WHERE name = $2)
			`, a2s(msg.Owner), msg.Domain, a2s(msg.NewAdmin))
		}
	}
	return accountID, castPgErr(err)
}

func (st *Store) RenewDomain(ctx context.Context, msg *types.MsgRenewDomain, height int64) (int64, error) {
	// only valid_until needs to be updated and that's done in HandleLcdData()
	return getAccountID(ctx, msg.Domain, "")
}

func (st *Store) TransferAccount(ctx context.Context, msg *types.MsgTransferAccount, height int64) (int64, error) {
	accountID, err := getAccountID(ctx, msg.Domain, msg.Name)
	if err == nil {
		_, err = dbTx.ExecContext(ctx, `
			UPDATE accounts
			SET owner = $1
			WHERE id = $2
		`, a2s(msg.NewOwner), accountID)
	}
	return accountID, castPgErr(err)
}

func (st *Store) RegisterAccount(ctx context.Context, msg *types.MsgRegisterAccount, height int64) (int64, error) {
	var id int64
	err := dbTx.QueryRowContext(ctx, `
		INSERT INTO accounts (domain_id, name, owner)
		VALUES ((SELECT MAX(id) FROM domains WHERE name = $1), $2, $3)
		RETURNING id
	`, msg.Domain, msg.Name, a2s(msg.Owner)).Scan(&id)
	return id, castPgErr(err)
}

func (st *Store) ReplaceAccountResources(ctx context.Context, msg *types.MsgReplaceAccountResources, height int64) (int64, error) {
	accountID, err := getAccountID(ctx, msg.Domain, msg.Name)
	if err == nil {
		for _, r := range msg.NewResources {
			_, err = dbTx.ExecContext(ctx, `
				INSERT INTO resources (account_id, resource, uri)
				VALUES ($1, $2, $3)
				ON CONFLICT (id) DO UPDATE SET resource = EXCLUDED.resource, uri = EXCLUDED.uri
			`, accountID, r.Resource, r.URI)
			if err != nil {
				return accountID, castPgErr(err)
			}
		}
	}
	return accountID, castPgErr(err)
}

func (st *Store) ReplaceAccountMetadata(ctx context.Context, msg *types.MsgReplaceAccountMetadata, height int64) (int64, error) {
	accountID, err := getAccountID(ctx, msg.Domain, msg.Name)
	if err == nil {
		_, err = dbTx.ExecContext(ctx, `
			UPDATE accounts
			SET metadata = $1
			WHERE id = $2
		`, msg.NewMetadataURI, accountID)
	}
	return accountID, castPgErr(err)
}

func (st *Store) AddAccountCertificates(ctx context.Context, msg *types.MsgAddAccountCertificates, height int64) (int64, error) {
	accountID, err := getAccountID(ctx, msg.Domain, msg.Name)
	if err == nil {
		_, err = dbTx.ExecContext(ctx, `
			INSERT INTO certificates(account_id, certificate)
			VALUES ($1, $2)
		`, accountID, msg.NewCertificate)
	}
	return accountID, castPgErr(err)
}

func (st *Store) DeleteAccountCerts(ctx context.Context, msg *types.MsgDeleteAccountCertificate, height int64) (int64, error) {
	accountID, err := getAccountID(ctx, msg.Domain, msg.Name)
	if err == nil {
		_, err = dbTx.ExecContext(ctx, `
			DELETE FROM certificates
			WHERE account_id = $1
		`, accountID)
	}
	return accountID, castPgErr(err)
}

func (st *Store) DeleteAccount(ctx context.Context, msg *types.MsgDeleteAccount, height int64) (int64, error) {
	accountID, err := getAccountID(ctx, msg.Domain, msg.Name)
	if err == nil {
		_, err = dbTx.ExecContext(ctx, `
			UPDATE accounts
			SET deleted = (SELECT block_time FROM blocks WHERE block_height=$2)
			WHERE id = $1
		`, accountID, height)
	}
	return accountID, castPgErr(err)
}

func (st *Store) RenewAccount(ctx context.Context, msg *types.MsgRenewAccount, height int64) (int64, error) {
	// only valid_until needs to be updated and that's done in HandleLcdData()
	return getAccountID(ctx, msg.Domain, msg.Name)
}

func (st *Store) InsertBlock(ctx context.Context, b Block) error {
	_, err := dbTx.ExecContext(ctx, `
		INSERT INTO blocks (block_height, block_hash, block_time, fee_frac)
		VALUES ($1, $2, $3, $4)
	`, b.Height, b.Hash, b.Time.UTC(), b.FeeFrac)
	return wrapPgErr(err, "insert block")
}

func (st *Store) LatestBlock(ctx context.Context) (*Block, error) {
	blocks, err := st.LastNBlock(ctx, 1, 0)
	if err != nil {
		return nil, err
	}
	return blocks[0], nil
}

func (st *Store) LastNBlock(ctx context.Context, limit, after int) ([]*Block, error) {
	// max number of blocks that is allowed to retrieved is 100
	if limit > 100 {
		return nil, errors.Wrapf(ErrLimit, "limit exceeded")
	}

	var rows *sql.Rows
	var err error
	if after == 0 {
		rows, err = st.db.QueryContext(ctx, `
		SELECT block_height, block_hash, block_time, fee_frac
		FROM blocks
		ORDER BY block_height DESC
		LIMIT $1
	`, limit)
	} else {
		rows, err = st.db.QueryContext(ctx, `
		SELECT block_height, block_hash, block_time, fee_frac
		FROM blocks
		WHERE block_height < $1
		ORDER BY block_height DESC
		LIMIT $2
	`, after, limit)
	}
	defer rows.Close()

	if err != nil {
		err = castPgErr(err)
		if ErrNotFound.Is(err) {
			return nil, errors.Wrap(err, "no blocks")
		}
		return nil, errors.Wrap(castPgErr(err), "cannot select block")
	}

	var blocks []*Block

	for rows.Next() {
		var b Block
		err := rows.Scan(&b.Height, &b.Hash, &b.Time, &b.FeeFrac)
		if err != nil {
			err = castPgErr(err)
			if ErrNotFound.Is(err) {
				return nil, errors.Wrap(err, "no blocks")
			}
			return nil, errors.Wrap(castPgErr(err), "cannot select block")

		}
		// normalize it here, as not always stored like this in the db
		b.Time = b.Time.UTC()
		blocks = append(blocks, &b)
	}
	if len(blocks) == 0 {
		return nil, errors.Wrap(ErrNotFound, "no blocks")
	}
	return blocks, nil
}

func (st *Store) InsertGenesis(ctx context.Context, tmc *TendermintClient) error {
	gen, err := FetchGenesis(ctx, tmc)
	if err != nil {
		return errors.Wrapf(err, "genesis fetch failed")
	}
	// begin a batch insert
	if err = st.BatchBegin(ctx); err != nil {
		return errors.Wrap(err, "st.BatchBegin() failed")
	}
	defer st.BatchRollback()
	for _, domain := range gen.Domains {
		msg := types.MsgRegisterDomain{
			Name:         domain.Name,
			Admin:        domain.Admin,
			DomainType:   domain.Type,
			Broker:       domain.Broker,
			FeePayerAddr: domain.Admin,
		}
		if accountID, err := st.RegisterDomain(ctx, &msg, 0); err != nil {
			return errors.Wrapf(err, "cannot insert domain")
		} else {
			_, err = dbTx.ExecContext(ctx, `
				UPDATE domains
				SET valid_until = TO_TIMESTAMP(CAST($1 AS DECIMAL))
				WHERE id = (SELECT MAX(id) FROM domains WHERE name = $2)
			`, domain.ValidUntil, domain.Name)
			if err != nil {
				return errors.Wrapf(err, "failed to update valid_until on domain %s", domain.Name)
			}
			_, err = dbTx.ExecContext(ctx, `
				UPDATE accounts
				SET valid_until = TO_TIMESTAMP(CAST($1 AS DECIMAL))
				WHERE id = $2
			`, domain.ValidUntil, accountID)
			if err != nil {
				return errors.Wrapf(err, "failed to update valid_until on accountID %d", accountID)
			}
		}
	}
	for _, acc := range gen.Accounts {
		// skip the empty account because it was handled in st.RegisterDomain()
		if *acc.Name == "" {
			continue
		}
		msg := types.MsgRegisterAccount{
			Domain:    acc.Domain,
			Name:      *acc.Name,
			Owner:     acc.Owner,
			Resources: acc.Resources,
			Broker:    acc.Broker,
		}
		if accountID, err := st.RegisterAccount(ctx, &msg, 0); err != nil {
			return errors.Wrapf(err, "cannot insert domain")
		} else {
			_, err = dbTx.ExecContext(ctx, `
				UPDATE accounts
				SET valid_until = TO_TIMESTAMP(CAST($1 AS DECIMAL))
				WHERE id = $2
			`, acc.ValidUntil, accountID)
			if err != nil {
				return errors.Wrapf(err, "failed to update valid_until on accountID %d", accountID)
			}
		}
	}
	// commit the batch
	if err = st.BatchCommit(ctx); err != nil {
		return errors.Wrapf(err, "st.BatchCommit() failed")
	}
	return err
}

func (st *Store) BatchBegin(ctx context.Context) error {
	tx, err := st.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "cannot create database transaction")
	}
	dbTx = tx
	return nil
}

func (st *Store) BatchCommit(ctx context.Context) error {
	// commit before...
	if err := dbTx.Commit(); err != nil {
		return errors.Wrapf(err, "dbTx.Commit()")
	}
	// ...begining a new database transaction
	if err := st.BatchBegin(ctx); err != nil {
		return errors.Wrap(err, "cannot create transaction")
	}
	return nil
}

func (st *Store) BatchRollback() error {
	return dbTx.Rollback()
}

func find(needle string, haystack []sdk.Attribute) (string, error) {
	for _, pair := range haystack {
		if pair.Key == needle {
			return pair.Value, nil
		}
	}
	return "", errors.New(fmt.Sprintf("couldn't find %s in %s", needle, haystack))
}

func getPayment(attributes []sdk.Attribute, denom string) (string, int64, error) {
	payer, err := find("sender", attributes)
	if err != nil {
		return "", 0, err
	}
	denominated, err := find("amount", attributes)
	if err != nil {
		return "", 0, err
	}
	absolute := strings.Replace(denominated, denom, "", 1)
	amount, err := strconv.ParseInt(absolute, 10, 64)
	if err != nil {
		return "", 0, err
	}
	return payer, amount, nil
}

func updateDomainValidUntil(ctx context.Context, domain string, expires int64) error {
	_, err := dbTx.ExecContext(ctx, `
		UPDATE domains
		SET valid_until = TO_TIMESTAMP($1)
		WHERE id = (SELECT MAX(id) FROM domains WHERE name = $2)
	`, expires, domain)
	return castPgErr(err)
}

func updateAccountValidUntil(ctx context.Context, id int64, expires int64) error {
	_, err := dbTx.ExecContext(ctx, `
		UPDATE accounts
		SET valid_until = TO_TIMESTAMP($1)
		WHERE id = $2
	`, expires, id)
	return castPgErr(err)
}

func (st *Store) HandleLcdData(ctx context.Context, queries *[]*LcdRequestData, responses *[]*LcdResponseData, height int64, denom string) error {
	for i, query := range *queries {
		response := (*responses)[i]
		if *response.TxError != nil {
			return *response.TxError
		}
		if response.StarnameError != nil && *response.StarnameError != nil {
			return *response.StarnameError
		}
		events := response.TxResponse.Logs[0].Events
		event0 := events[0]
		event1 := events[1]
		if event0.Type != "message" {
			return errors.New(fmt.Sprintf("expected event type 'message' but got '%s'", event0.Type))
		}
		if event1.Type != "transfer" {
			return errors.New(fmt.Sprintf("expected event type 'transfer' but got '%s'", event1.Type))
		}
		action, err := find("action", event0.Attributes)
		if err != nil {
			return err
		}
		owner, err := find("owner", event0.Attributes)
		if err != nil {
			return err
		}
		payer, amount, err := getPayment(event1.Attributes, denom)
		if err != nil {
			return err
		}
		if owner == payer {
			payer = ""
		}
		broker, _ := find("broker", event0.Attributes) // not all actions have a broker, so ignore the err
		_, err = dbTx.ExecContext(ctx, `
			INSERT INTO product_fees (block, account_id, action, fee, payer, broker)
			VALUES ((SELECT block_height FROM blocks WHERE block_height=$1), $2, $3, $4, $5, $6)
		`, height, query.AccountID, action, amount, payer, broker)
		if err != nil {
			return castPgErr(err)
		}
		// update valid_until where appropriate
		if response.StarnameResponse != nil && response.StarnameResponse.Height != "" { // "" happens if a starname has been deleted
			account := response.StarnameResponse.Result.Account
			expires := account.ValidUntil
			if expires > 0 {
				switch action {
				case "register_domain", "renew_domain":
					if err := updateDomainValidUntil(ctx, account.Domain, expires); err != nil {
						return err
					}
				}
				if err := updateAccountValidUntil(ctx, query.AccountID, expires); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
