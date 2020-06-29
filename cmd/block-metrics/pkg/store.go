package pkg

import (
	"context"
	"database/sql"
	"time"

	"github.com/iov-one/iovns/x/domain/types"

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

func (st *Store) RegisterDomain(ctx context.Context, msg *types.MsgRegisterDomain) (int64, error) {
	var id int64
	err := st.db.QueryRowContext(ctx, `
		INSERT INTO domains (name, admin, type, broker, fee_payer_addr)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, msg.Name, msg.Admin, msg.DomainType, msg.Broker, msg.FeePayerAddr).Scan(&id)
	return id, castPgErr(err)
}

func (st *Store) DeleteDomain(ctx context.Context, msg *types.MsgDeleteDomain) error {
	sqlStatement := `
		UPDATE domains SET deleted_at = now() 
		WHERE name = $1`
	_, err := st.db.ExecContext(ctx, sqlStatement, msg.Domain)
	return err
}

func (st *Store) TransferDomain(ctx context.Context, msg *types.MsgTransferDomain) error {
	sqlStatement := `
	UPDATE domains
	SET admin = $1
	WHERE name = $2`
	_, err := st.db.ExecContext(ctx, sqlStatement, msg.NewAdmin, msg.Domain)
	return err
}

func (st *Store) RegisterAccount(ctx context.Context, msg *types.MsgRegisterAccount) (int64, error) {
	var accountID int64

	sqlStatement := `SELECT id FROM domains WHERE name = $1;`
	row := st.db.QueryRowContext(ctx, sqlStatement, msg.Domain)
	var domainID int64
	if err := row.Scan(&domainID); err != nil {
		return domainID, castPgErr(err)
	}
	_, err := st.db.ExecContext(ctx, `
		INSERT INTO accounts (domain_id, domain, name, owner, registerer, broker, fee_payer_addr)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, domainID, msg.Domain, msg.Name, msg.Owner, msg.Registerer, msg.Broker, msg.FeePayerAddr)
	if err != nil {
		return accountID, castPgErr(err)
	}
	return accountID, nil
}

func (st *Store) ReplaceAccountResources(ctx context.Context, msg *types.MsgReplaceAccountResources) (int64, error) {
	var resourceID int64

	sqlStatement := `SELECT id FROM accounts WHERE domain = $1 and name = $2;`
	row := st.db.QueryRowContext(ctx, sqlStatement, msg.Domain, msg.Name)
	var accountID int64
	if err := row.Scan(&accountID); err != nil {
		return accountID, castPgErr(err)
	}

	tx, err := st.db.Begin()
	if err != nil {
		return accountID, castPgErr(err)
	}
	for _, r := range msg.NewResources {
		st := `INSERT INTO resources (account_id, resource, uri)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO UPDATE SET resource = EXCLUDED.resource, uri = EXCLUDED.uri;`
		stmt, err := tx.Prepare(st)
		if err != nil {
			tx.Rollback()
			return accountID, castPgErr(err)
		}
		_, err = stmt.ExecContext(ctx, accountID, r.Resource, r.URI)
		if err != nil {
			tx.Rollback()
			return accountID, castPgErr(err)
		}
		if err := stmt.Close(); err != nil {
			return accountID, castPgErr(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return accountID, castPgErr(err)
	}

	return resourceID, castPgErr(err)
}

func (st *Store) DeleteAccount(ctx context.Context, msg *types.MsgDeleteAccount) error {
	sqlStatement := `
		UPDATE accounts SET deleted_at = now() 
		WHERE domain = $1 AND name = $2`
	_, err := st.db.ExecContext(ctx, sqlStatement, msg.Domain, msg.Name)
	return err
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

	_ = tx.Rollback()

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
		if _, err := st.RegisterDomain(ctx, &msg); err != nil {
			return errors.Wrapf(err, "cannot insert domain")
		}
	}
	for _, acc := range gen.Accounts {
		msg := types.MsgRegisterAccount{
			Domain:    acc.Domain,
			Name:      acc.Name,
			Owner:     acc.Owner,
			Resources: acc.Resources,
			Broker:    acc.Broker,
		}
		if _, err := st.RegisterAccount(ctx, &msg); err != nil {
			return errors.Wrapf(err, "cannot insert domain")
		}
	}
	return err
}
