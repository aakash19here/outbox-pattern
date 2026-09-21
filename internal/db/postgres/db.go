package dbpostgres

import (
	"github.com/aakash19here/outbox-pattern/internal/config"
	"github.com/jmoiron/sqlx"
)

type DBPostgresOptions struct {
	DBname string
}

func getDBConnString(opts *DBPostgresOptions) string {
	if opts.DBname == "" {
		return config.GetDefaultConnString()
	}
	return config.GetConnString(opts.DBname)
}

func NewDBConn(opts *DBPostgresOptions) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", getDBConnString(opts))
	if err != nil {
		return nil, err
	}
	return db, nil
}
