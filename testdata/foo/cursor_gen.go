package foo

// Code generated by repo-generator v0.1.0. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dohernandez/errors"
	"github.com/dohernandez/repo-generator/testdata/deps"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrCursorScan is the error that indicates a Cursor scan failed.
	ErrCursorScan = errors.New("scan")
	// ErrCursorNotFound is the error that indicates a Cursor was not found.
	ErrCursorNotFound = errors.New("not found")
	// ErrCursorExists is returned when the Cursor already exists.
	ErrCursorExists = errors.New("exists")

	// ErrCursorUpdate is the error that indicates a Cursor was not updated.
	ErrCursorUpdate = errors.New("update")
)

// CursorRow is an interface for anything that can scan a Cursor, copying the columns from the matched
// row into the values pointed at by dest.
type CursorRow interface {
	Scan(dest ...any) error
}

// CursorCriteria is a criteria for the select Cursor(s).
//
// CursorCriteria is used to generate the where statement of the select. The order of the criteria items
// matter during generation of the where statement.
type CursorCriteria map[string]any

func (c CursorCriteria) toSql() (string, []interface{}) {
	var (
		where []string
		args  []interface{}
	)

	for k, v := range c {
		where = append(where, fmt.Sprintf("%s = $%d", k, len(where)+1))
		args = append(args, v)
	}

	return strings.Join(where, " AND "), args
}

// CursorRepo is a repository for the Cursor.
type CursorRepo struct {
	// db is the database connection.
	db *sql.DB

	// table is the table name.
	table string

	// colID is the Cursor.ID column name. It can be used in a queries to specify the column.
	colID string
	// colName is the Cursor.Name column name. It can be used in a queries to specify the column.
	colName string
	// colPosition is the Cursor.Position column name. It can be used in a queries to specify the column.
	colPosition string
	// colLeader is the Cursor.Leader column name. It can be used in a queries to specify the column.
	colLeader string
	// colLeaderElectedAt is the Cursor.LeaderElectedAt column name. It can be used in a queries to specify the column.
	colLeaderElectedAt string
	// colCreatedAt is the Cursor.CreatedAt column name. It can be used in a queries to specify the column.
	colCreatedAt string
	// colUpdatedAt is the Cursor.UpdatedAt column name. It can be used in a queries to specify the column.
	colUpdatedAt string
}

// NewCursorRepo creates a new CursorRepo.
func NewCursorRepo(db *sql.DB, table string) *CursorRepo {
	return &CursorRepo{
		db:    db,
		table: table,

		colID:              "id",
		colName:            "name",
		colPosition:        "position",
		colLeader:          "leader",
		colLeaderElectedAt: "leader_elected_at",
		colCreatedAt:       "created_at",
		colUpdatedAt:       "updated_at",
	}
}

// Scan scans a Cursor from the given CursorRow (sql.Row|sql.Rows).
func (repo *CursorRepo) Scan(_ context.Context, s CursorRow) (*Cursor, error) {
	var (
		m Cursor

		position        sql.NullString
		leader          sql.NullString
		leaderElectedAt sql.NullTime
		createdAt       time.Time
		updatedAt       time.Time
	)

	err := s.Scan(
		&m.ID,
		&m.Name,
		&position,
		&leader,
		&leaderElectedAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCursorNotFound
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, errors.WrapError(err, ErrCursorExists)
		}

		return nil, errors.WrapError(err, ErrCursorScan)
	}

	if position.Valid {
		m.Position = uuid.MustParse(position.String)
	}

	if leader.Valid {
		m.Leader = uuid.MustParse(leader.String)
	}

	if leaderElectedAt.Valid {
		m.LeaderElectedAt = leaderElectedAt.Time.UTC()
	}

	m.CreatedAt = createdAt.UTC()
	m.UpdatedAt = updatedAt.UTC()

	return &m, nil
}

// ScanAll scans a slice of Cursor from the given sql.Rows.
func (repo *CursorRepo) ScanAll(ctx context.Context, rs *sql.Rows) ([]*Cursor, error) {
	var ms []*Cursor

	for rs.Next() {
		m, err := repo.Scan(ctx, rs)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	if len(ms) == 0 {
		return nil, ErrCursorNotFound
	}

	return ms, nil
}

// Select selects a Cursor given CursorCriteria.
//
// Returns the error ErrCursorNotFound if the Cursor was not found and database did not error,
// otherwise database error.
func (repo *CursorRepo) Select(ctx context.Context, criteria CursorCriteria) (*Cursor, error) {
	const q = "SELECT %s FROM %s WHERE %s"

	where, args := criteria.toSql()

	cols := strings.Join([]string{
		repo.colID,
		repo.colName,
		repo.colPosition,
		repo.colLeader,
		repo.colLeaderElectedAt,
		repo.colCreatedAt,
		repo.colUpdatedAt}, ", ")

	query := fmt.Sprintf(q, cols, repo.table, where)

	return repo.Scan(ctx, repo.db.QueryRowContext(ctx, query, args...))
}

// Create creates a new Cursor and returns the persisted Cursor.
//
// The returned Cursor will contain the fields that were tag as "auto", which maybe were generated by the
// database..
func (repo *CursorRepo) Create(ctx context.Context, m *Cursor) (*Cursor, error) {
	var (
		cols []string
		args []interface{}
	)

	if !deps.IsUUIDZero(m.ID) {
		cols = append(cols, repo.colID)
		args = append(args, m.ID)
	}

	cols = append(cols, repo.colName)
	args = append(args, m.Name)

	var position sql.NullString

	position.String = m.Position.String()
	position.Valid = true

	cols = append(cols, repo.colPosition)
	args = append(args, position)

	var leader sql.NullString

	leader.String = m.Leader.String()
	leader.Valid = true

	cols = append(cols, repo.colLeader)
	args = append(args, leader)

	var leaderElectedAt sql.NullTime

	leaderElectedAt.Time = m.LeaderElectedAt
	leaderElectedAt.Valid = true

	cols = append(cols, repo.colLeaderElectedAt)
	args = append(args, leaderElectedAt)

	cols = append(cols, repo.colCreatedAt)
	args = append(args, m.CreatedAt)

	cols = append(cols, repo.colUpdatedAt)
	args = append(args, m.UpdatedAt)

	values := make([]string, len(cols))

	for i := range cols {
		values[i] = fmt.Sprintf("$%d", i+1)
	}

	qCols := strings.Join(cols, ", ")
	qValues := strings.Join(values, ", ")

	rCols := strings.Join([]string{
		repo.colID,
		repo.colName,
		repo.colPosition,
		repo.colLeader,
		repo.colLeaderElectedAt,
		repo.colCreatedAt,
		repo.colUpdatedAt,
	}, ", ")

	sql := "INSERT INTO %s (%s) VALUES (%s) RETURNING %s"
	sql = fmt.Sprintf(sql, repo.table, qCols, qValues, rCols)

	return repo.Scan(ctx, repo.db.QueryRowContext(ctx, sql, args...))
}

// UpdateCursorOptions is a type for specifying options when updating a Cursor.
type UpdateCursorOptions func(*updateCursorOptions)

type updateCursorOptions struct {
	// skipZeroValues indicates whether to skip zero values from the update statement.
	skipZeroValues bool
}

// SkipCursorZeroValues is an option for Update that indicates whether to skip zero values from the update statement.
func SkipCursorZeroValues() UpdateCursorOptions {
	return func(o *updateCursorOptions) {
		o.skipZeroValues = true
	}
}

// Update updates a Cursor.
//
// skipZeroValues indicates whether to skip zero values from the update statement.
// In case of boolean fields, skipZeroValues is not applicable since false is the zero value of boolean and could be
// a potential update. Always set this type of fields.
//
// Returns the error ErrCursorUpdate if the Cursor was not updated and database did not error,
// otherwise database error.
func (repo *CursorRepo) Update(ctx context.Context, m *Cursor, opts ...UpdateCursorOptions) error {
	var uOpts updateCursorOptions

	for _, opt := range opts {
		opt(&uOpts)
	}

	var (
		sets   []string
		where  []string
		args   []interface{}
		offset = 1
	)

	if deps.IsUUIDZero(m.ID) {
		return errors.Wrapf(ErrCursorUpdate, "field ID is required")
	}

	where = append(where, fmt.Sprintf("%s = $%d", repo.colID, offset))
	args = append(args, m.ID)

	offset++

	if uOpts.skipZeroValues {
		if m.Name != "" {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colName, offset))
			args = append(args, m.Name)

			offset++
		}

		if !deps.IsUUIDZero(m.Position) {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colPosition, offset))
			args = append(args, m.Position.String())

			offset++
		}

		if !deps.IsUUIDZero(m.Leader) {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colLeader, offset))
			args = append(args, m.Leader.String())

			offset++
		}

		if !m.LeaderElectedAt.IsZero() {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colLeaderElectedAt, offset))
			args = append(args, m.LeaderElectedAt)

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
		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colName, offset))
		args = append(args, m.Name)

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colPosition, offset))
		args = append(args, m.Position.String())

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colLeader, offset))
		args = append(args, m.Leader.String())

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colLeaderElectedAt, offset))
		args = append(args, m.LeaderElectedAt)

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
		return errors.WrapError(err, ErrCursorUpdate)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.WrapError(err, ErrCursorUpdate)
	}

	if rowsAffected == 0 {
		return ErrCursorUpdate
	}

	return nil
}
