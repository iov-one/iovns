package pkg

import (
	"database/sql"
	"fmt"
	"strings"
)

func EnsureSchema(pg *sql.DB) error {
	tx, err := pg.Begin()
	if err != nil {
		return fmt.Errorf("transaction begin: %s", err)
	}

	for _, query := range strings.Split(schema, "\n---\n") {
		query = strings.TrimSpace(query)

		if _, err := tx.Exec(query); err != nil {
			return &QueryError{Query: query, Err: err}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit: %s", err)
	}

	_ = tx.Rollback()

	return nil
}

var schema = `
CREATE TABLE IF NOT EXISTS validators (
	id BIGSERIAL PRIMARY KEY,
	public_key BYTEA NOT NULL UNIQUE,
	address BYTEA NOT NULL UNIQUE,
	name TEXT,
	memo TEXT
);

---

CREATE TABLE IF NOT EXISTS blocks (
	block_height BIGINT NOT NULL PRIMARY KEY,
	block_hash TEXT NOT NULL UNIQUE,
	block_time TIMESTAMPTZ NOT NULL,
	proposer_id INT NOT NULL REFERENCES validators(id),
	messages TEXT[] NOT NULL,
	fee_frac BIGINT NOT NULL
);

---

CREATE TABLE IF NOT EXISTS block_participations (
	id BIGSERIAL PRIMARY KEY,
	validated BOOLEAN NOT NULL,
	block_id BIGINT NOT NULL REFERENCES blocks(block_height),
	validator_id INT NOT NULL REFERENCES validators(id),
	UNIQUE (block_id, validator_id)
);

---

CREATE TABLE IF NOT EXISTS transactions (
	id BIGSERIAL PRIMARY KEY,
	transaction_hash TEXT NOT NULL UNIQUE,
	block_id BIGINT NOT NULL REFERENCES blocks(block_height),
	signatures BYTEA[],
	fee JSONB,
	messages JSONB[],
	memo text
);

CREATE INDEX ON transactions (transaction_hash);
---

CREATE TABLE IF NOT EXISTS domains (
	id BIGSERIAL PRIMARY KEY,
	name TEXT,
	admin BYTEA,
	type TEXT,
	broker BYTEA,
	fee_payer_addr BYTEA
);

---
CREATE TABLE IF NOT EXISTS accounts (
	id BIGSERIAL PRIMARY KEY,
	domain TEXT,
	name TEXT,
	owner BYTEA,
	ValidUntil BIGINT,
	broker BYTEA,
	metadataURI TEXT
);

---
CREATE TABLE IF NOT EXISTS resources (
	id BIGSERIAL PRIMARY KEY,
	account_id INTEGER REFERENCES accounts(id),
	URI TEXT,
	resource TEXT
);

---
CREATE TABLE IF NOT EXISTS account_certificates (
	id BIGSERIAL PRIMARY KEY,
	account_id INTEGER REFERENCES accounts(id),
	certificate BYTEA
);
`

type QueryError struct {
	Query string
	Args  []interface{}
	Err   error
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("query error: %s\n%q", e.Err, e.Query)
}
