package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/iov-one/iovns/cmd/block-metrics/pkg"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
)

func main() {
	conf := pkg.Configuration{
		DBHost:           os.Getenv("POSTGRES_HOST"),
		DBName:           os.Getenv("POSTGRES_DB"),
		DBPass:           os.Getenv("POSTGRES_PASSWORD"),
		DBROPass:         os.Getenv("POSTGRES_RO_PASSWORD"),
		DBROUser:         os.Getenv("POSTGRES_RO_USER"),
		DBSSL:            os.Getenv("POSTGRES_SSL_ENABLE"),
		DBUser:           os.Getenv("POSTGRES_USER"),
		FeeDenom:         os.Getenv("FEE_DENOMINATION"),
		TendermintLcdUrl: os.Getenv("TENDERMINT_LCD_URL"),
		TendermintWsURI:  os.Getenv("TENDERMINT_WS_URI"),
	}

	if err := run(conf); err != nil {
		log.Fatal(err)
	}
}

func run(conf pkg.Configuration) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := pkg.EnsureDatabase(conf.DBUser, conf.DBPass, conf.DBHost, conf.DBName, conf.DBSSL); err != nil {
		return fmt.Errorf("ensure database: %s", err)
	}

	dbUri := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", conf.DBUser, conf.DBPass, conf.DBHost, conf.DBName, conf.DBSSL)
	db, err := sql.Open("postgres", dbUri)
	if err != nil {
		return fmt.Errorf("cannot connect to postgres: %s", err)
	}
	defer db.Close()

	if err := pkg.EnsureSchema(db, conf.DBName, conf.DBROUser, conf.DBROPass); err != nil {
		return fmt.Errorf("ensure schema: %s", err)
	}

	st := pkg.NewStore(db)

	tmc, err := pkg.DialTendermint(conf.TendermintWsURI)
	if err != nil {
		return errors.Wrap(err, "dial tendermint")
	}
	defer tmc.Close()

	inserted, err := pkg.Sync(ctx, tmc, st, conf.FeeDenom, conf.TendermintLcdUrl)
	if err != nil {
		return errors.Wrap(err, "sync")
	}

	fmt.Println("inserted:", inserted)

	return nil
}
