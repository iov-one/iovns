package pkg

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/starname/types"

	"github.com/pkg/errors"
)

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
func getAccountID(st *Store, ctx context.Context, domain string, name string) (int64, error) {
	var accountID int64
	err := st.db.QueryRowContext(ctx, `
		SELECT id FROM accounts
		WHERE domain_id = (SELECT MAX(id) FROM domains  WHERE name = $1)
		AND          id = (SELECT MAX(id) FROM accounts WHERE name = $2)
	`, domain, name).Scan(&accountID)

	return accountID, castPgErr(err)
}

func (st *Store) RegisterDomain(ctx context.Context, msg *types.MsgRegisterDomain, height int64) (int64, error) {
	tx, err := st.db.Begin()
	if err != nil {
		return 0, castPgErr(err)
	}

	// create the domain...
	var id int64
	err = st.db.QueryRowContext(ctx, `
		INSERT INTO domains (name, admin, type, broker, fee_payer_addr, created)
		VALUES ($1, $2, $3, $4, $5, (SELECT block_height FROM blocks WHERE block_height=$6))
		RETURNING id
	`, msg.Name, a2s(msg.Admin), msg.DomainType, a2s(msg.Broker), a2s(msg.FeePayerAddr), height).Scan(&id)
	if err != nil {
		return 0, castPgErr(err)
	}

	// ...and then the empty account
	msgEmptyAccount := types.MsgRegisterAccount{
		Domain:       msg.Name,
		Name:         "",
		Owner:        msg.Admin,
		Broker:       msg.Broker,
		FeePayerAddr: msg.FeePayerAddr,
	}
	accountID, err := st.RegisterAccount(ctx, &msgEmptyAccount, height)
	if err != nil {
		return accountID, err
	}

	err = tx.Commit()

	return accountID, castPgErr(err)
}

func (st *Store) DeleteDomain(ctx context.Context, msg *types.MsgDeleteDomain, height int64) (int64, error) {
	tx, err := st.db.Begin()
	if err != nil {
		return 0, castPgErr(err)
	}

	// delete the empty account...
	msgEmptyAccount := types.MsgDeleteAccount{
		Domain:       msg.Domain,
		Name:         "",
		Owner:        msg.Owner,
		FeePayerAddr: msg.FeePayerAddr,
	}
	accountID, err := st.DeleteAccount(ctx, &msgEmptyAccount, height)
	if err != nil {
		return accountID, err
	}

	// ...and then the domain
	_, err = st.db.ExecContext(ctx, `
		UPDATE domains
		SET deleted = (SELECT block_height FROM blocks WHERE block_height=$2)
		WHERE id = (SELECT MAX(id) FROM domains WHERE name = $1)
	`, msg.Domain, height)
	if err != nil {
		return 0, castPgErr(err)
	}

	err = tx.Commit()

	return accountID, castPgErr(err)
}

func (st *Store) TransferDomain(ctx context.Context, msg *types.MsgTransferDomain, height int64) (int64, error) {
	tx, err := st.db.Begin()
	if err != nil {
		return 0, castPgErr(err)
	}

	// update the empty account...
	msgEmptyAccount := types.MsgTransferAccount{
		Domain:       msg.Domain,
		Name:         "",
		Owner:        msg.Owner,
		NewOwner:     msg.NewAdmin,
		FeePayerAddr: msg.FeePayerAddr,
		Reset:        true, // TODO: deal with the different transfer flags
	}
	accountID, err := st.TransferAccount(ctx, &msgEmptyAccount, height)
	if err != nil {
		return accountID, err
	}

	// ...and then the domain
	_, err = st.db.ExecContext(ctx, `
		UPDATE domains
		SET admin = $1, updated = (SELECT block_height FROM blocks WHERE block_height=$3)
		WHERE id = (SELECT MAX(id) FROM domains WHERE name = $2)
	`, a2s(msg.NewAdmin), msg.Domain, height)
	if err != nil {
		return accountID, castPgErr(err)
	}

	err = tx.Commit()

	return accountID, castPgErr(err)
}

func (st *Store) TransferAccount(ctx context.Context, msg *types.MsgTransferAccount, height int64) (int64, error) {
	accountID, err := getAccountID(st, ctx, msg.Domain, msg.Name)
	if err != nil {
		return accountID, err
	}

	_, err = st.db.ExecContext(ctx, `
		UPDATE accounts
		SET owner = $1, updated = (SELECT block_height FROM blocks WHERE block_height=$3)
		WHERE id = $2
	`, a2s(msg.NewOwner), accountID, height)

	return accountID, castPgErr(err)
}

func (st *Store) RegisterAccount(ctx context.Context, msg *types.MsgRegisterAccount, height int64) (int64, error) {
	var id int64
	err := st.db.QueryRowContext(ctx, `
		INSERT INTO accounts (domain_id, domain, name, owner, registerer, broker, fee_payer_addr, created)
		VALUES ((SELECT MAX(id) FROM domains WHERE name = $1), $1, $2, $3, $4, $5, $6, (SELECT block_height FROM blocks WHERE block_height=$7))
		RETURNING id
	`, msg.Domain, msg.Name, a2s(msg.Owner), a2s(msg.Registerer), a2s(msg.Broker), a2s(msg.FeePayerAddr), height).Scan(&id)
	return id, castPgErr(err)
}

func (st *Store) ReplaceAccountResources(ctx context.Context, msg *types.MsgReplaceAccountResources, height int64) (int64, error) {
	tx, err := st.db.Begin()
	if err != nil {
		return 0, castPgErr(err)
	}

	accountID, err := getAccountID(st, ctx, msg.Domain, msg.Name)
	if err != nil {
		return accountID, err
	}

	for _, r := range msg.NewResources {
		st := `INSERT INTO resources (account_id, resource, uri, updated)
			VALUES ($1, $2, $3, (SELECT block_height FROM blocks WHERE block_height=$4))
			ON CONFLICT (id) DO UPDATE SET resource = EXCLUDED.resource, uri = EXCLUDED.uri;`
		stmt, err := tx.Prepare(st)
		if err != nil {
			tx.Rollback()
			return accountID, castPgErr(err)
		}
		_, err = stmt.ExecContext(ctx, accountID, r.Resource, r.URI, height)
		if err != nil {
			tx.Rollback()
			return accountID, castPgErr(err)
		}
		if err := stmt.Close(); err != nil {
			return accountID, castPgErr(err)
		}
	}

	err = tx.Commit()

	return accountID, castPgErr(err)
}

func (st *Store) ReplaceAccountMetadata(ctx context.Context, msg *types.MsgReplaceAccountMetadata, height int64) (int64, error) {
	tx, err := st.db.Begin()
	if err != nil {
		return 0, castPgErr(err)
	}

	accountID, err := getAccountID(st, ctx, msg.Domain, msg.Name)
	if err != nil {
		return accountID, err
	}

	_, err = st.db.ExecContext(ctx, `
		UPDATE accounts
		SET metadata_uri = $1, updated = (SELECT block_height FROM blocks WHERE block_height=$3)
		WHERE id = $2
	`, msg.NewMetadataURI, accountID, height)
	if err != nil {
		return accountID, castPgErr(err)
	}

	err = tx.Commit()

	return accountID, castPgErr(err)
}

func (st *Store) AddAccountCertificates(ctx context.Context, msg *types.MsgAddAccountCertificates, height int64) (int64, error) {
	tx, err := st.db.Begin()
	if err != nil {
		return 0, castPgErr(err)
	}

	accountID, err := getAccountID(st, ctx, msg.Domain, msg.Name)
	if err != nil {
		return accountID, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO account_certificates(account_id, certificate, created)
		VALUES ($1, $2, (SELECT block_height FROM blocks WHERE block_height=$3))
	`, accountID, msg.NewCertificate, height)
	if err != nil {
		return accountID, wrapPgErr(err, "insert block")
	}

	err = tx.Commit()

	return accountID, castPgErr(err)
}

func (st *Store) DeleteAccountCerts(ctx context.Context, msg *types.MsgDeleteAccountCertificate, height int64) (int64, error) {
	accountID, err := getAccountID(st, ctx, msg.Domain, msg.Name)
	if err != nil {
		return accountID, err
	}

	_, err = st.db.ExecContext(ctx, `
		UPDATE account_certificates
		SET deleted = (SELECT block_height FROM blocks WHERE block_height=$2)
		WHERE account_id = $1
	`, accountID, height)

	return accountID, castPgErr(err)
}

func (st *Store) DeleteAccount(ctx context.Context, msg *types.MsgDeleteAccount, height int64) (int64, error) {
	accountID, err := getAccountID(st, ctx, msg.Domain, msg.Name)
	if err != nil {
		return accountID, err
	}

	_, err = st.db.ExecContext(ctx, `
		UPDATE accounts
		SET deleted = (SELECT block_height FROM blocks WHERE block_height=$2)
		WHERE id = $1
	`, accountID, height)

	return accountID, castPgErr(err)
}

func (st *Store) InsertBlock(ctx context.Context, b Block) error {
	tx, err := st.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "cannot create transaction")
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO blocks (block_height, block_hash, block_time, fee_frac)
		VALUES ($1, $2, $3, $4)
	`, b.Height, b.Hash, b.Time.UTC(), b.FeeFrac)
	if err != nil {
		return wrapPgErr(err, "insert block")
	}

	err = tx.Commit()

	_ = tx.Rollback() // TODO: WTF?  Ask Orkun what this line was supposed to accomplish

	return wrapPgErr(err, "commit block tx")
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
	for _, domain := range gen.Domains {
		msg := types.MsgRegisterDomain{
			Name:         domain.Name,
			Admin:        domain.Admin,
			DomainType:   domain.Type,
			Broker:       domain.Broker,
			FeePayerAddr: domain.Admin,
		}
		if _, err := st.RegisterDomain(ctx, &msg, 0); err != nil {
			return errors.Wrapf(err, "cannot insert domain")
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
		if _, err := st.RegisterAccount(ctx, &msg, 0); err != nil {
			return errors.Wrapf(err, "cannot insert domain")
		}
	}
	return err
}
