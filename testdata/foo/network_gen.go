package foo

// Code generated by repo-generator v0.1.0. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strings"

	"github.com/dohernandez/errors"
)

var (
	ErrNetworkScan     = errors.New("scan")
	ErrNetworkNotFound = errors.New("not found")
	ErrNetworkUpdate   = errors.New("update")
)

type NetworkScanner interface {
	Scan(dest ...any) error
}

type NetworkRepo struct {
	db *sql.DB

	table        string
	colID        string
	colToken     string
	colURI       string
	colNumber    string
	colTotal     string
	colIP        string
	colCreatedAt string
	colUpdatedAt string

	cols []string
}

func NewNetworkRepo(db *sql.DB, table string) *NetworkRepo {
	cols := []string{
		"id",
		"token",
		"uri",
		"number",
		"total",
		"ip",
		"created_at",
		"updated_at",
	}

	return &NetworkRepo{
		db:           db,
		table:        table,
		colID:        "id",
		colToken:     "token",
		colURI:       "uri",
		colNumber:    "number",
		colTotal:     "total",
		colIP:        "ip",
		colCreatedAt: "created_at",
		colUpdatedAt: "updated_at",

		cols: cols,
	}
}

func (repo *NetworkRepo) Scan(_ context.Context, s NetworkScanner) (*Network, error) {
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

	rCols := strings.Join(repo.cols, ", ")

	sql := "INSERT INTO %s (%s) VALUES (%s) RETURNING %s"
	sql = fmt.Sprintf(sql, repo.table, qCols, qValues, rCols)

	return repo.Scan(ctx, repo.db.QueryRowContext(ctx, sql, args...))
}

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

	for i := range ms {
		m := ms[i]

		indexOffset := i * lcols

		valuesQueryBuilder.WriteString(repo.valuesStatement(cols, indexOffset, i != lms-1))

		if m.ID != "" {
			args = append(args, m.ID)
		} else {
			args = append(args, nil)
		}

		args = append(args, m.Token)

		var uRI sql.NullString

		uRI.String = m.URI
		uRI.Valid = true

		args = append(args, uRI)

		var number sql.NullInt64

		if m.Number != nil {
			number.Int64 = m.Number.Int64()
			number.Valid = true
		}

		args = append(args, number)

		args = append(args, m.Total.Int64())

		if m.IP != nil {
			args = append(args, *m.IP)
		} else {
			args = append(args, nil)
		}

		if !m.CreatedAt.IsZero() {
			args = append(args, m.CreatedAt)
		} else {
			args = append(args, nil)
		}

		if !m.UpdatedAt.IsZero() {
			args = append(args, m.UpdatedAt)
		} else {
			args = append(args, nil)
		}

	}

	qCols := strings.Join(cols, ", ")

	sql := "INSERT INTO %s (%s) VALUES %s"
	sql = fmt.Sprintf(sql, repo.table, qCols, valuesQueryBuilder.String())

	_, err := repo.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "exec context")
	}

	return nil
}

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

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colURI, offset))
		args = append(args, m.URI)

		offset++

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

		where = append(where, fmt.Sprintf("%s = $%d", repo.colURI, offset))
		args = append(args, m.URI)

		offset++

		where = append(where, fmt.Sprintf("%s = $%d", repo.colNumber, offset))
		args = append(args, m.Number)

		offset++

		where = append(where, fmt.Sprintf("%s = $%d", repo.colTotal, offset))
		args = append(args, m.Total)

		offset++

		where = append(where, fmt.Sprintf("%s = $%d", repo.colIP, offset))
		args = append(args, m.IP)

		offset++

		where = append(where, fmt.Sprintf("%s = $%d", repo.colCreatedAt, offset))
		args = append(args, m.CreatedAt)

		offset++

		where = append(where, fmt.Sprintf("%s = $%d", repo.colUpdatedAt, offset))
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
