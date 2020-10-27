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
CREATE TABLE IF NOT EXISTS blocks (
	block_height BIGINT NOT NULL PRIMARY KEY,
	block_hash TEXT NOT NULL UNIQUE,
	block_time TIMESTAMPTZ NOT NULL,
	fee_frac BIGINT NOT NULL
);

---

CREATE TABLE IF NOT EXISTS transactions (
	id BIGSERIAL PRIMARY KEY,
	transaction_hash TEXT NOT NULL UNIQUE,
	block_id BIGINT NOT NULL REFERENCES blocks(block_height),
	signatures BYTEA ARRAY,
	fee JSONB,
	memo text
);

---
CREATE TABLE IF NOT EXISTS domains (
	id BIGSERIAL PRIMARY KEY,
	name TEXT,
	admin TEXT NOT NULL,
	type TEXT NOT NULL,
	broker TEXT,
	fee_payer_addr TEXT,
	deleted_at TIMESTAMP
);

---
CREATE TABLE IF NOT EXISTS accounts (
	id BIGSERIAL PRIMARY KEY,
	domain_id BIGINT NOT NULL REFERENCES domains(id),
	domain TEXT NOT NULL,
	name TEXT,
	owner TEXT NOT NULL,
	registerer TEXT,
	broker TEXT,
	metadata_uri TEXT,
	fee_payer_addr TEXT,
	deleted_at TIMESTAMP
);

---
CREATE TABLE IF NOT EXISTS resources (
	id BIGSERIAL PRIMARY KEY,
	account_id BIGINT REFERENCES accounts(id),
	uri TEXT,
	resource TEXT
);

---
CREATE TABLE IF NOT EXISTS account_certificates (
	id BIGSERIAL PRIMARY KEY,
	account_id BIGINT REFERENCES accounts(id),
	certificate BYTEA,
	deleted_at TIMESTAMP
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
