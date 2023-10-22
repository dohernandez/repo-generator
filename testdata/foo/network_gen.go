package foo

// Code generated by repo-generator v0.1.0. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strings"

	"github.com/dohernandez/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrNetworkScan is the error that indicates a Network scan failed.
	ErrNetworkScan = errors.New("scan")
	// ErrNetworkNotFound is the error that indicates a Network was not found.
	ErrNetworkNotFound = errors.New("not found")
	// ErrNetworkUpdate is the error that indicates a Network was not updated.
	ErrNetworkUpdate = errors.New("update")
	// ErrNetworkExists is returned when the Network already exists.
	ErrNetworkExists = errors.New("exists")
)

// NetworkRow is an interface for anything that can scan a Network, copying the columns from the matched
// row into the values pointed at by dest.
type NetworkRow interface {
	Scan(dest ...any) error
}

// NetworkRepo is a repository for the Network.
type NetworkRepo struct {
	// db is the database connection.
	db *sql.DB

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
	// colCreatedAt is the Network.CreatedAt column name. It can be used in a queries to specify the column.
	colCreatedAt string
	// colUpdatedAt is the Network.UpdatedAt column name. It can be used in a queries to specify the column.
	colUpdatedAt string
}

// NewNetworkRepo creates a new NetworkRepo.
func NewNetworkRepo(db *sql.DB, table string) *NetworkRepo {
	return &NetworkRepo{
		db:    db,
		table: table,

		colID:        "id",
		colToken:     "token",
		colURI:       "uri",
		colNumber:    "number",
		colTotal:     "total",
		colIP:        "ip",
		colCreatedAt: "created_at",
		colUpdatedAt: "updated_at",
	}
}

// Scan scans a Network from the given NetworkRow (sql.Row|sql.Rows).
func (repo *NetworkRepo) Scan(_ context.Context, s NetworkRow) (*Network, error) {
	var (
		m Network

		uRI    sql.NullString
		number sql.NullInt64
		total  int64
		iP     string
	)

	err := s.Scan(
		&m.ID,
		&m.Token,
		&uRI,
		&number,
		&total,
		&iP,
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

// Create creates a new Network and returns the persisted Network.
//
// The returned Network will contain the fields that were tag as "auto", which maybe were generated by the
// database..
func (repo *NetworkRepo) Create(ctx context.Context, m *Network) (*Network, error) {
	var (
		cols []string
		args []interface{}
	)

	if m.ID != "" {
		cols = append(cols, repo.colID)
		args = append(args, m.ID)
	}

	cols = append(cols, repo.colToken)
	args = append(args, m.Token)

	var uRI sql.NullString

	uRI.String = m.URI
	uRI.Valid = true

	cols = append(cols, repo.colURI)
	args = append(args, uRI)

	if m.Number != nil {
		var number sql.NullInt64

		number.Int64 = m.Number.Int64()
		number.Valid = true

		cols = append(cols, repo.colNumber)
		args = append(args, number)
	}

	cols = append(cols, repo.colTotal)
	args = append(args, m.Total.Int64())

	if m.IP != nil {
		cols = append(cols, repo.colIP)
		args = append(args, *m.IP)
	}

	if !m.CreatedAt.IsZero() {
		cols = append(cols, repo.colCreatedAt)
		args = append(args, m.CreatedAt)
	}

	if !m.UpdatedAt.IsZero() {
		cols = append(cols, repo.colUpdatedAt)
		args = append(args, m.UpdatedAt)
	}

	values := make([]string, len(cols))

	for i := range cols {
		values[i] = fmt.Sprintf("$%d", i+1)
	}

	qCols := strings.Join(cols, ", ")
	qValues := strings.Join(values, ", ")

	rCols := strings.Join([]string{
		repo.colID,
		repo.colToken,
		repo.colURI,
		repo.colNumber,
		repo.colTotal,
		repo.colIP,
		repo.colCreatedAt,
		repo.colUpdatedAt,
	}, ", ")

	sql := "INSERT INTO %s (%s) VALUES (%s) RETURNING %s"
	sql = fmt.Sprintf(sql, repo.table, qCols, qValues, rCols)

	return repo.Scan(ctx, repo.db.QueryRowContext(ctx, sql, args...))
}

// Insert inserts one or more Network records into the database.
//
// When using this method the Network fields that are tag as "auto" should be set as the other fields non tag as "auto".
// The same applies for those other fields that are tag as "omitempty".
func (repo *NetworkRepo) Insert(ctx context.Context, ms ...*Network) error {
	// Build values query.
	var (
		valuesQueryBuilder strings.Builder
		lms                = len(ms)
	)

	var cols []string

	cols = append(cols, repo.colID)
	cols = append(cols, repo.colToken)
	cols = append(cols, repo.colURI)
	cols = append(cols, repo.colNumber)
	cols = append(cols, repo.colTotal)
	cols = append(cols, repo.colIP)
	cols = append(cols, repo.colCreatedAt)
	cols = append(cols, repo.colUpdatedAt)

	lcols := len(cols)

	// Size is equal to the number of models (lms) multiplied by the number of columns (lcols).
	args := make([]interface{}, 0, lms*lcols)

	for idx := range ms {
		m := ms[idx]

		indexOffset := idx * lcols

		valuesQueryBuilder.WriteString(repo.valuesStatement(cols, indexOffset, idx != lms-1))

		args = append(args, m.ID)

		args = append(args, m.Token)

		var uRI sql.NullString

		uRI.String = m.URI
		uRI.Valid = true

		args = append(args, uRI)

		var number sql.NullInt64

		number.Int64 = m.Number.Int64()
		number.Valid = true

		args = append(args, number)

		args = append(args, m.Total.Int64())

		if m.IP != nil {
			args = append(args, *m.IP)
		}

		args = append(args, m.CreatedAt)

		args = append(args, m.UpdatedAt)
	}

	qCols := strings.Join(cols, ", ")

	sql := "INSERT INTO %s (%s) VALUES %s"
	sql = fmt.Sprintf(sql, repo.table, qCols, valuesQueryBuilder.String())

	_, err := repo.db.ExecContext(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return errors.WrapError(err, ErrNetworkExists)
		}

		return errors.Wrap(err, "exec context")
	}

	return nil
}

// valuesStatement returns a string with the values statement ($n) for the given columns,
// starting from the given offset.
func (repo *NetworkRepo) valuesStatement(cols []string, offset int, separator bool) string {
	var sep string

	if separator {
		sep = ","
	}

	values := make([]string, len(cols))
	for i := range cols {
		values[i] = fmt.Sprintf("$%d", offset+(i+1))
	}

	return fmt.Sprintf("(%s)%s", strings.Join(values, ", "), sep)
}

// Update updates a Network.
//
// skipZeroValues indicates whether to skip zero values from the update statement.
// In case of boolean fields, skipZeroValues is not applicable since false is the zero value of boolean and could be
// a potential update. Always set this type of fields.
//
// Returns the error ErrNetworkUpdate if the Network was not updated and database did not error,
// otherwise database error.
func (repo *NetworkRepo) Update(ctx context.Context, m *Network, skipZeroValues bool) error {
	var (
		sets   []string
		where  []string
		args   []interface{}
		offset = 1
	)

	where = append(where, fmt.Sprintf("%s = $%d", repo.colID, offset))
	args = append(args, m.ID)

	offset++

	where = append(where, fmt.Sprintf("%s = $%d", repo.colToken, offset))
	args = append(args, m.Token)

	offset++

	if skipZeroValues {
		if m.URI != "" {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colURI, offset))
			args = append(args, m.URI)

			offset++
		}

		if m.Number != nil {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colNumber, offset))
			args = append(args, m.Number.Int64())

			offset++
		}

		if m.Total.Int64() != 0 {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colTotal, offset))
			args = append(args, m.Total.Int64())

			offset++
		}

		if m.IP != nil {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colIP, offset))
			args = append(args, *m.IP)

			offset++
		}

		if !m.CreatedAt.IsZero() {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colCreatedAt, offset))
			args = append(args, m.CreatedAt)

			offset++
		}

		if !m.UpdatedAt.IsZero() {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colUpdatedAt, offset))
			args = append(args, m.UpdatedAt)

			offset++
		}
	} else {
		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colURI, offset))
		args = append(args, m.URI)

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colNumber, offset))
		args = append(args, m.Number.Int64())

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colTotal, offset))
		args = append(args, m.Total.Int64())

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colIP, offset))
		args = append(args, *m.IP)

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colCreatedAt, offset))
		args = append(args, m.CreatedAt)

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colUpdatedAt, offset))
		args = append(args, m.UpdatedAt)

		offset++

	}

	qSets := strings.Join(sets, ", ")
	qWhere := strings.Join(where, " AND ")

	sql := "UPDATE %s SET %s WHERE %s"
	sql = fmt.Sprintf(sql, repo.table, qSets, qWhere)

	res, err := repo.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNetworkUpdate
	}

	return nil
}

// Delete deletes a Network.
//
// The Network must have the fields that are tag as "key" set.
func (repo *NetworkRepo) Delete(ctx context.Context, m *Network) error {
	var (
		where []string
		args  []interface{}
	)

	where = append(where, fmt.Sprintf("%s = $%d", repo.colID, 1))
	args = append(args, m.ID)

	where = append(where, fmt.Sprintf("%s = $%d", repo.colToken, 2))
	args = append(args, m.Token)

	qWhere := strings.Join(where, " AND ")

	sql := "DELETE FROM %s WHERE %s"
	sql = fmt.Sprintf(sql, repo.table, qWhere)

	_, err := repo.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "exec context")
	}

	return nil
}
