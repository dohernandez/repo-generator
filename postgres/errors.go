package postgres

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsDbUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != pgerrcode.UniqueViolation {
		return false
	}

	if pgErr.TableName != "pg_database" {
		return false
	}

	if pgErr.ConstraintName != "pg_database_datname_index" {
		return false
	}

	return true
}
