package healthcheck

import (
	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"github.com/tessig/flamingo-mysql/db"
)

type (
	// DB checks the database connection
	DB struct {
		db db.DB
	}
)

var _ healthcheck.Status = new(DB)

// Inject dependencies
func (d *DB) Inject(
	db db.DB,
) *DB {
	d.db = db

	return d
}

// Status of the database connection
func (d *DB) Status() (bool, string) {
	err := d.db.Connection().Ping()
	if err != nil {
		return false, err.Error()
	}

	return true, "DB connection alive"
}
