package foo

// Code generated by repo-generator v0.1.0. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/dohernandez/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrEventScan is the error that indicates a Event scan failed.
	ErrEventScan = errors.New("scan")
	// ErrEventNotFound is the error that indicates a Event was not found.
	ErrEventNotFound = errors.New("not found")
	// ErrEventExists is returned when the Event already exists.
	ErrEventExists = errors.New("exists")
)

// EventRow is an interface for anything that can scan a Event, copying the columns from the matched
// row into the values pointed at by dest.
type EventRow interface {
	Scan(dest ...any) error
}

// EventRepo is a repository for the Event.
type EventRepo struct {
	// db is the database connection.
	db *sql.DB

	// table is the table name.
	table string

	// colID is the Event.ID column name. It can be used in a queries to specify the column.
	colID string
	// colTopic is the Event.Topic column name. It can be used in a queries to specify the column.
	colTopic string
	// colKey is the Event.Key column name. It can be used in a queries to specify the column.
	colKey string
	// colSequence is the Event.Sequence column name. It can be used in a queries to specify the column.
	colSequence string
	// colMetadata is the Event.Metadata column name. It can be used in a queries to specify the column.
	colMetadata string
	// colCreatedAt is the Event.CreatedAt column name. It can be used in a queries to specify the column.
	colCreatedAt string
}

// NewEventRepo creates a new EventRepo.
func NewEventRepo(db *sql.DB, table string) *EventRepo {
	return &EventRepo{
		db:    db,
		table: table,

		colID:        "id",
		colTopic:     "topic",
		colKey:       "key",
		colSequence:  "sequence",
		colMetadata:  "metadata",
		colCreatedAt: "created_at",
	}
}

// Scan scans a Event from the given EventRow (sql.Row|sql.Rows).
func (repo *EventRepo) Scan(_ context.Context, s EventRow) (*Event, error) {
	var m Event

	err := s.Scan(
		&m.ID,
		&m.Topic,
		&m.Key,
		&m.Sequence,
		&m.Metadata,
		&m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventNotFound
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, errors.WrapError(err, ErrEventExists)
		}

		return nil, errors.WrapError(err, ErrEventScan)
	}

	return &m, nil
}

// ScanAll scans a slice of Event from the given sql.Rows.
func (repo *EventRepo) ScanAll(ctx context.Context, rs *sql.Rows) ([]*Event, error) {
	var ms []*Event

	for rs.Next() {
		m, err := repo.Scan(ctx, rs)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	if len(ms) == 0 {
		return nil, ErrEventNotFound
	}

	return ms, nil
}

// Create creates a new Event and returns the persisted Event.
//
// The returned Event will contain the fields that were tag as "auto", which maybe were generated by the
// database..
func (repo *EventRepo) Create(ctx context.Context, m *Event) (*Event, error) {
	var (
		cols []string
		args []interface{}
	)

	cols = append(cols, repo.colID)
	args = append(args, m.ID)

	cols = append(cols, repo.colTopic)
	args = append(args, m.Topic)

	cols = append(cols, repo.colKey)
	args = append(args, m.Key)

	cols = append(cols, repo.colSequence)
	args = append(args, m.Sequence)

	cols = append(cols, repo.colMetadata)
	args = append(args, m.Metadata)

	if !m.CreatedAt.IsZero() {
		cols = append(cols, repo.colCreatedAt)
		args = append(args, m.CreatedAt)
	}

	values := make([]string, len(cols))

	for i := range cols {
		values[i] = fmt.Sprintf("$%d", i+1)
	}

	qCols := strings.Join(cols, ", ")
	qValues := strings.Join(values, ", ")

	rCols := strings.Join([]string{
		repo.colID,
		repo.colTopic,
		repo.colKey,
		repo.colSequence,
		repo.colMetadata,
		repo.colCreatedAt,
	}, ", ")

	sql := "INSERT INTO %s (%s) VALUES (%s) RETURNING %s"
	sql = fmt.Sprintf(sql, repo.table, qCols, qValues, rCols)

	return repo.Scan(ctx, repo.db.QueryRowContext(ctx, sql, args...))
}
