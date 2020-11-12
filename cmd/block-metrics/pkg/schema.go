package pkg

import (
	"database/sql"
	"fmt"
	"strings"
)

func EnsureDatabase(user, password, host, database, ssl string) error {
	dbUri := fmt.Sprintf("postgres://%s:%s@%s/?sslmode=%s", user, password, host, ssl)
	db, err := sql.Open("postgres", dbUri)
	if err != nil {
		return err
	}
	defer db.Close()
	// ignore the error if the database already exists
	_, _ = db.Exec(fmt.Sprintf(`CREATE DATABASE %s`, database))
	return nil
}

func EnsureSchema(pg *sql.DB, database, rouser, ropassword string) error {
	// setup readonly on schema permissioned a la https://aws.amazon.com/blogs/database/managing-postgresql-users-and-roles/
	rows, err := pg.Query(`
		SELECT schema_name
		FROM information_schema.schemata
		WHERE schema_name = 'permissioned';
	`)
	if err != nil {
		return fmt.Errorf("type query: %s", err)
	}
	if !rows.Next() {
		sql := `
			CREATE SCHEMA permissioned;
			ALTER DATABASE __db__ SET search_path = "$user", permissioned, public
			REVOKE CREATE ON SCHEMA public FROM PUBLIC;
			REVOKE ALL ON DATABASE __db__ FROM PUBLIC;
			CREATE ROLE readonly;
			GRANT CONNECT ON DATABASE __db__ TO readonly;
			GRANT USAGE ON SCHEMA permissioned TO readonly;
			ALTER DEFAULT PRIVILEGES IN SCHEMA permissioned GRANT SELECT ON TABLES TO readonly;
			CREATE USER __rouser__ WITH PASSWORD '__ropassword__';
			GRANT readonly TO __rouser__;
		`
		replacements := make(map[string]string)
		replacements["__db__"] = database
		replacements["__rouser__"] = rouser
		replacements["__ropassword__"] = ropassword
		for key, value := range replacements {
			sql = strings.ReplaceAll(sql, key, value)
		}
		for _, query := range strings.Split(sql, "\n") {
			if _, err := pg.Exec(query); err != nil {
				return &QueryError{Query: query, Err: err}
			}
		}
	}

	// deal with the pesky TYPE 'action' that doesn't allow an IF NOT EXISTS clause
	rows, err = pg.Query(`
		SELECT pg_type.typname, pg_enum.enumlabel
		FROM pg_type
		JOIN pg_enum ON pg_enum.enumtypid = pg_type.oid
	`)
	if err != nil {
		return fmt.Errorf("type query: %s", err)
	}
	if !rows.Next() {
		_, err = pg.Exec(`
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
		`)
		if err != nil {
			return fmt.Errorf("type create: %s", err)
		}
	}

	// create tables, possibly
	for _, query := range strings.Split(schema, "\n---\n") {
		if _, err := pg.Exec(query); err != nil {
			return &QueryError{Query: query, Err: err}
		}
	}

	return nil
}

var schema = `
CREATE TABLE IF NOT EXISTS permissioned.blocks (
	block_height BIGINT NOT NULL PRIMARY KEY,
	block_hash TEXT NOT NULL UNIQUE,
	block_time TIMESTAMPTZ NOT NULL,
	fee_frac BIGINT NOT NULL
);

---
CREATE TABLE IF NOT EXISTS permissioned.transactions (
	id BIGSERIAL PRIMARY KEY,
	transaction_hash TEXT NOT NULL UNIQUE,
	block_id BIGINT NOT NULL REFERENCES permissioned.blocks(block_height),
	signatures BYTEA ARRAY,
	fee JSONB,
	memo text
);

---
CREATE TABLE IF NOT EXISTS permissioned.domains (
	id BIGSERIAL PRIMARY KEY,
	name TEXT,
	admin TEXT NOT NULL,
	type TEXT NOT NULL,
	valid_until TIMESTAMPTZ,
	deleted TIMESTAMPTZ
);

---
CREATE TABLE IF NOT EXISTS permissioned.accounts (
	id BIGSERIAL PRIMARY KEY,
	domain_id BIGINT NOT NULL REFERENCES permissioned.domains(id),
	name TEXT,
	owner TEXT NOT NULL,
	metadata TEXT,
	valid_until TIMESTAMPTZ,
	deleted TIMESTAMPTZ
);

---
CREATE TABLE IF NOT EXISTS permissioned.resources (
	id BIGSERIAL PRIMARY KEY,
	account_id BIGINT REFERENCES permissioned.accounts(id),
	uri TEXT,
	resource TEXT
);

---
CREATE TABLE IF NOT EXISTS permissioned.certificates (
	id BIGSERIAL PRIMARY KEY,
	account_id BIGINT REFERENCES permissioned.accounts(id),
	certificate BYTEA
);

---
CREATE TABLE IF NOT EXISTS permissioned.product_fees (
	id BIGSERIAL PRIMARY KEY,
	block BIGINT REFERENCES permissioned.blocks(block_height),
	account_id BIGINT REFERENCES permissioned.accounts(id),
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
