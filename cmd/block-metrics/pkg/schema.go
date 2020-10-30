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
	created BIGINT REFERENCES blocks(block_height),
	updated BIGINT REFERENCES blocks(block_height),
	deleted BIGINT REFERENCES blocks(block_height)
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
	created BIGINT REFERENCES blocks(block_height),
	updated BIGINT REFERENCES blocks(block_height),
	deleted BIGINT REFERENCES blocks(block_height)
);

---
CREATE TABLE IF NOT EXISTS resources (
	id BIGSERIAL PRIMARY KEY,
	account_id BIGINT REFERENCES accounts(id),
	uri TEXT,
	resource TEXT,
	updated BIGINT REFERENCES blocks(block_height)
);

---
CREATE TABLE IF NOT EXISTS account_certificates (
	id BIGSERIAL PRIMARY KEY,
	account_id BIGINT REFERENCES accounts(id),
	certificate BYTEA,
	created BIGINT REFERENCES blocks(block_height),
	deleted BIGINT REFERENCES blocks(block_height)
);

---
DROP TYPE IF EXISTS action CASCADE;
CREATE TYPE action AS ENUM (
	'add_certificates_account',
	'delete_account',
	'delete_certificate_account',
	'delete_domain',
	'register_account',
	'register_domain',
	'renew_account',
	'renew_domain',
	'replace_account_resources',
	'set_account_metadata',
	'transfer_account',
	'transfer_domain'
);

---
CREATE TABLE IF NOT EXISTS product_fees (
	id BIGSERIAL PRIMARY KEY,
	block BIGINT REFERENCES blocks(block_height),
	account_id BIGINT REFERENCES accounts(id),
	action action,
	fee BIGINT,
	payer TEXT,
	broker TEXT
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
