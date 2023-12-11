package foo

// Code generated by repo-generator v0.1.0. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/dohernandez/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrNetworkScan is the error that indicates a Network scan failed.
	ErrNetworkScan = errors.New("scan")
	// ErrNetworkNotFound is the error that indicates a Network was not found.
	ErrNetworkNotFound = errors.New("not found")
	// ErrNetworkExists is returned when the Network already exists.
	ErrNetworkExists = errors.New("exists")
)

// NetworkRow is an interface for anything that can scan a Network, copying the columns from the matched
// row into the values pointed at by dest.
type NetworkRow interface {
	Scan(dest ...any) error
}

// NetworkSQLDB is an interface for anything that can execute the SQL statements needed to
// perform the Network operations.
type NetworkSQLDB interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error)
}

// NetworkRepo is a repository for the Network.
type NetworkRepo struct {
	// db is the database connection.
	db NetworkSQLDB

	// table is the table name.
	table string

	// colID is the Network.ID column name. It can be used in a queries to specify the column.
	colID string
	// colToken is the Network.Token column name. It can be used in a queries to specify the column.
	colToken string
	// colURI is the Network.URI column name. It can be used in a queries to specify the column.
	colURI string
	// colNumber is the Network.Number column name. It can be used in a queries to specify the column.
	colNumber string
	// colTotal is the Network.Total column name. It can be used in a queries to specify the column.
	colTotal string
	// colIP is the Network.IP column name. It can be used in a queries to specify the column.
	colIP string
	// colAddress is the Network.Address column name. It can be used in a queries to specify the column.
	colAddress string
	// colCreatedAt is the Network.CreatedAt column name. It can be used in a queries to specify the column.
	colCreatedAt string
	// colUpdatedAt is the Network.UpdatedAt column name. It can be used in a queries to specify the column.
	colUpdatedAt string
}

// NewNetworkRepo creates a new NetworkRepo.
func NewNetworkRepo(db NetworkSQLDB, table string) *NetworkRepo {
	return &NetworkRepo{
		db:    db,
		table: table,

		colID:        "id",
		colToken:     "token",
		colURI:       "uri",
		colNumber:    "number",
		colTotal:     "total",
		colIP:        "ip",
		colAddress:   "address",
		colCreatedAt: "created_at",
		colUpdatedAt: "updated_at",
	}
}

// Cols returns the represented cols of Network.
// Cols are returned in the order they are scanned.
func (repo *NetworkRepo) Cols() []string {
	return []string{
		repo.colID,
		repo.colToken,
		repo.colURI,
		repo.colNumber,
		repo.colTotal,
		repo.colIP,
		repo.colAddress,
		repo.colCreatedAt,
		repo.colUpdatedAt,
	}
}

// Scan scans a Network from the given NetworkRow (sql.Row|sql.Rows).
func (repo *NetworkRepo) Scan(_ context.Context, s NetworkRow) (*Network, error) {
	var (
		m Network

		uRI     sql.NullString
		number  sql.NullInt64
		total   int64
		iP      string
		address sql.NullString
	)

	err := s.Scan(
		&m.ID,
		&m.Token,
		&uRI,
		&number,
		&total,
		&iP,
		&address,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNetworkNotFound
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, errors.WrapError(err, ErrNetworkExists)
		}

		return nil, errors.WrapError(err, ErrNetworkScan)
	}

	if uRI.Valid {
		m.URI = uRI.String
	}

	if number.Valid {
		m.Number = big.NewInt(number.Int64)
	}

	m.Total = bigNewInt(total)

	tmp := iP
	m.IP = &tmp

	if address.Valid {
		m.Address = stringToAddress(address.String)
	}

	return &m, nil
}

// ScanAll scans a slice of Network from the given sql.Rows.
func (repo *NetworkRepo) ScanAll(ctx context.Context, rs *sql.Rows) ([]*Network, error) {
	var ms []*Network

	for rs.Next() {
		m, err := repo.Scan(ctx, rs)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	if len(ms) == 0 {
		return nil, ErrNetworkNotFound
	}

	return ms, nil
}
