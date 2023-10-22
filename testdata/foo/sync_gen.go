package foo

// Code generated by repo-generator v0.1.0. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strings"

	"github.com/dohernandez/errors"
	"github.com/dohernandez/repo-generator/testdata/deps"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrSyncScan is the error that indicates a Sync scan failed.
	ErrSyncScan = errors.New("scan")
	// ErrSyncNotFound is the error that indicates a Sync was not found.
	ErrSyncNotFound = errors.New("not found")
	// ErrSyncUpdate is the error that indicates a Sync was not updated.
	ErrSyncUpdate = errors.New("update")
	// ErrSyncExists is returned when the Sync already exists.
	ErrSyncExists = errors.New("exists")
)

// SyncRow is an interface for anything that can scan a Sync, copying the columns from the matched
// row into the values pointed at by dest.
type SyncRow interface {
	Scan(dest ...any) error
}

// SyncRepo is a repository for the Sync.
type SyncRepo struct {
	// db is the database connection.
	db *sql.DB

	// table is the table name.
	table string

	// colID is the Sync.ID column name. It can be used in a queries to specify the column.
	colID string
	// colState is the Sync.State column name. It can be used in a queries to specify the column.
	colState string
	// colChainID is the Sync.ChainID column name. It can be used in a queries to specify the column.
	colChainID string
	// colBlockNumber is the Sync.BlockNumber column name. It can be used in a queries to specify the column.
	colBlockNumber string
	// colBlockHash is the Sync.BlockHash column name. It can be used in a queries to specify the column.
	colBlockHash string
	// colParentHash is the Sync.ParentHash column name. It can be used in a queries to specify the column.
	colParentHash string
	// colBlockTimestamp is the Sync.BlockTimestamp column name. It can be used in a queries to specify the column.
	colBlockTimestamp string
	// colBlockHeaderPath is the Sync.BlockHeaderPath column name. It can be used in a queries to specify the column.
	colBlockHeaderPath string
	// colTransactionsPath is the Sync.TransactionsPath column name. It can be used in a queries to specify the column.
	colTransactionsPath string
	// colReceiptsPath is the Sync.ReceiptsPath column name. It can be used in a queries to specify the column.
	colReceiptsPath string
	// colLogsPath is the Sync.LogsPath column name. It can be used in a queries to specify the column.
	colLogsPath string
	// colTracesPath is the Sync.TracesPath column name. It can be used in a queries to specify the column.
	colTracesPath string
	// colCreatedAt is the Sync.CreatedAt column name. It can be used in a queries to specify the column.
	colCreatedAt string
	// colUpdatedAt is the Sync.UpdatedAt column name. It can be used in a queries to specify the column.
	colUpdatedAt string
}

// NewSyncRepo creates a new SyncRepo.
func NewSyncRepo(db *sql.DB, table string) *SyncRepo {
	return &SyncRepo{
		db:    db,
		table: table,

		colID:               "id",
		colState:            "state",
		colChainID:          "chain_id",
		colBlockNumber:      "block_number",
		colBlockHash:        "block_hash",
		colParentHash:       "parent_hash",
		colBlockTimestamp:   "block_timestamp",
		colBlockHeaderPath:  "block_header_path",
		colTransactionsPath: "transactions_path",
		colReceiptsPath:     "receipts_path",
		colLogsPath:         "logs_path",
		colTracesPath:       "traces_path",
		colCreatedAt:        "created_at",
		colUpdatedAt:        "updated_at",
	}
}

// Scan scans a Sync from the given SyncRow (sql.Row|sql.Rows).
func (repo *SyncRepo) Scan(_ context.Context, s SyncRow) (*Sync, error) {
	var (
		m Sync

		blockNumber      int64
		blockHash        sql.NullString
		parentHash       sql.NullString
		blockTimestamp   sql.NullTime
		blockHeaderPath  sql.NullString
		transactionsPath sql.NullString
		receiptsPath     sql.NullString
		logsPath         sql.NullString
		tracesPath       sql.NullString
	)

	err := s.Scan(
		&m.ID,
		&m.State,
		&m.ChainID,
		&blockNumber,
		&blockHash,
		&parentHash,
		&blockTimestamp,
		&blockHeaderPath,
		&transactionsPath,
		&receiptsPath,
		&logsPath,
		&tracesPath,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSyncNotFound
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, errors.WrapError(err, ErrSyncExists)
		}

		return nil, errors.WrapError(err, ErrSyncScan)
	}

	m.BlockNumber = big.NewInt(blockNumber)

	if blockHash.Valid {
		m.BlockHash = deps.HexToHash(blockHash.String)
	}

	if parentHash.Valid {
		m.ParentHash = deps.HexToHash(parentHash.String)
	}

	if blockTimestamp.Valid {
		m.BlockTimestamp = blockTimestamp.Time
	}

	if blockHeaderPath.Valid {
		m.BlockHeaderPath = blockHeaderPath.String
	}

	if transactionsPath.Valid {
		m.TransactionsPath = transactionsPath.String
	}

	if receiptsPath.Valid {
		m.ReceiptsPath = receiptsPath.String
	}

	if logsPath.Valid {
		m.LogsPath = logsPath.String
	}

	if tracesPath.Valid {
		m.TracesPath = tracesPath.String
	}

	return &m, nil
}

// ScanAll scans a slice of Sync from the given sql.Rows.
func (repo *SyncRepo) ScanAll(ctx context.Context, rs *sql.Rows) ([]*Sync, error) {
	var ms []*Sync

	for rs.Next() {
		m, err := repo.Scan(ctx, rs)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	if len(ms) == 0 {
		return nil, ErrSyncNotFound
	}

	return ms, nil
}

// Create creates a new Sync and returns the persisted Sync.
//
// The returned Sync will contain the fields that were tag as "auto", which maybe were generated by the
// database..
func (repo *SyncRepo) Create(ctx context.Context, m *Sync) (*Sync, error) {
	var (
		cols []string
		args []interface{}
	)

	if !deps.IsUUIDZero(m.ID) {
		cols = append(cols, repo.colID)
		args = append(args, m.ID)
	}

	cols = append(cols, repo.colState)
	args = append(args, m.State)

	cols = append(cols, repo.colChainID)
	args = append(args, m.ChainID)

	if m.BlockNumber != nil {
		cols = append(cols, repo.colBlockNumber)
		args = append(args, m.BlockNumber.String())
	}

	var blockHash sql.NullString

	blockHash.String = m.BlockHash.String()
	blockHash.Valid = true

	cols = append(cols, repo.colBlockHash)
	args = append(args, blockHash)

	var parentHash sql.NullString

	parentHash.String = m.ParentHash.String()
	parentHash.Valid = true

	cols = append(cols, repo.colParentHash)
	args = append(args, parentHash)

	var blockTimestamp sql.NullTime

	blockTimestamp.Time = m.BlockTimestamp
	blockTimestamp.Valid = true

	cols = append(cols, repo.colBlockTimestamp)
	args = append(args, blockTimestamp)

	var blockHeaderPath sql.NullString

	blockHeaderPath.String = m.BlockHeaderPath
	blockHeaderPath.Valid = true

	cols = append(cols, repo.colBlockHeaderPath)
	args = append(args, blockHeaderPath)

	var transactionsPath sql.NullString

	transactionsPath.String = m.TransactionsPath
	transactionsPath.Valid = true

	cols = append(cols, repo.colTransactionsPath)
	args = append(args, transactionsPath)

	var receiptsPath sql.NullString

	receiptsPath.String = m.ReceiptsPath
	receiptsPath.Valid = true

	cols = append(cols, repo.colReceiptsPath)
	args = append(args, receiptsPath)

	var logsPath sql.NullString

	logsPath.String = m.LogsPath
	logsPath.Valid = true

	cols = append(cols, repo.colLogsPath)
	args = append(args, logsPath)

	var tracesPath sql.NullString

	tracesPath.String = m.TracesPath
	tracesPath.Valid = true

	cols = append(cols, repo.colTracesPath)
	args = append(args, tracesPath)

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
		repo.colState,
		repo.colChainID,
		repo.colBlockNumber,
		repo.colBlockHash,
		repo.colParentHash,
		repo.colBlockTimestamp,
		repo.colBlockHeaderPath,
		repo.colTransactionsPath,
		repo.colReceiptsPath,
		repo.colLogsPath,
		repo.colTracesPath,
		repo.colCreatedAt,
		repo.colUpdatedAt,
	}, ", ")

	sql := "INSERT INTO %s (%s) VALUES (%s) RETURNING %s"
	sql = fmt.Sprintf(sql, repo.table, qCols, qValues, rCols)

	return repo.Scan(ctx, repo.db.QueryRowContext(ctx, sql, args...))
}

// Insert inserts one or more Sync records into the database.
//
// When using this method the Sync fields that are tag as "auto" should be set as the other fields non tag as "auto".
// The same applies for those other fields that are tag as "omitempty".
func (repo *SyncRepo) Insert(ctx context.Context, ms ...*Sync) error {
	// Build values query.
	var (
		valuesQueryBuilder strings.Builder
		lms                = len(ms)
	)

	var cols []string

	cols = append(cols, repo.colID)
	cols = append(cols, repo.colState)
	cols = append(cols, repo.colChainID)
	cols = append(cols, repo.colBlockNumber)
	cols = append(cols, repo.colBlockHash)
	cols = append(cols, repo.colParentHash)
	cols = append(cols, repo.colBlockTimestamp)
	cols = append(cols, repo.colBlockHeaderPath)
	cols = append(cols, repo.colTransactionsPath)
	cols = append(cols, repo.colReceiptsPath)
	cols = append(cols, repo.colLogsPath)
	cols = append(cols, repo.colTracesPath)
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

		args = append(args, m.State)

		args = append(args, m.ChainID)

		if m.BlockNumber != nil {
			args = append(args, m.BlockNumber.String())
		}

		var blockHash sql.NullString

		blockHash.String = m.BlockHash.String()
		blockHash.Valid = true

		args = append(args, blockHash)

		var parentHash sql.NullString

		parentHash.String = m.ParentHash.String()
		parentHash.Valid = true

		args = append(args, parentHash)

		var blockTimestamp sql.NullTime

		blockTimestamp.Time = m.BlockTimestamp
		blockTimestamp.Valid = true

		args = append(args, blockTimestamp)

		var blockHeaderPath sql.NullString

		blockHeaderPath.String = m.BlockHeaderPath
		blockHeaderPath.Valid = true

		args = append(args, blockHeaderPath)

		var transactionsPath sql.NullString

		transactionsPath.String = m.TransactionsPath
		transactionsPath.Valid = true

		args = append(args, transactionsPath)

		var receiptsPath sql.NullString

		receiptsPath.String = m.ReceiptsPath
		receiptsPath.Valid = true

		args = append(args, receiptsPath)

		var logsPath sql.NullString

		logsPath.String = m.LogsPath
		logsPath.Valid = true

		args = append(args, logsPath)

		var tracesPath sql.NullString

		tracesPath.String = m.TracesPath
		tracesPath.Valid = true

		args = append(args, tracesPath)

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
			return errors.WrapError(err, ErrSyncExists)
		}

		return errors.Wrap(err, "exec context")
	}

	return nil
}

// valuesStatement returns a string with the values statement ($n) for the given columns,
// starting from the given offset.
func (repo *SyncRepo) valuesStatement(cols []string, offset int, separator bool) string {
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

// Update updates a Sync.
//
// skipZeroValues indicates whether to skip zero values from the update statement.
// In case of boolean fields, skipZeroValues is not applicable since false is the zero value of boolean and could be
// a potential update. Always set this type of fields.
//
// Returns the error ErrSyncUpdate if the Sync was not updated and database did not error,
// otherwise database error.
func (repo *SyncRepo) Update(ctx context.Context, m *Sync, skipZeroValues bool) error {
	var (
		sets   []string
		where  []string
		args   []interface{}
		offset = 1
	)

	where = append(where, fmt.Sprintf("%s = $%d", repo.colID, offset))
	args = append(args, m.ID)

	offset++

	if skipZeroValues {
		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colState, offset))
		args = append(args, m.State)

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colChainID, offset))
		args = append(args, m.ChainID)

		offset++

		if m.BlockNumber != nil {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colBlockNumber, offset))
			args = append(args, m.BlockNumber.String())

			offset++
		}

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colBlockHash, offset))
		args = append(args, m.BlockHash.String())

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colParentHash, offset))
		args = append(args, m.ParentHash.String())

		offset++

		if !m.BlockTimestamp.IsZero() {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colBlockTimestamp, offset))
			args = append(args, m.BlockTimestamp)

			offset++
		}

		if m.BlockHeaderPath != "" {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colBlockHeaderPath, offset))
			args = append(args, m.BlockHeaderPath)

			offset++
		}

		if m.TransactionsPath != "" {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colTransactionsPath, offset))
			args = append(args, m.TransactionsPath)

			offset++
		}

		if m.ReceiptsPath != "" {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colReceiptsPath, offset))
			args = append(args, m.ReceiptsPath)

			offset++
		}

		if m.LogsPath != "" {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colLogsPath, offset))
			args = append(args, m.LogsPath)

			offset++
		}

		if m.TracesPath != "" {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colTracesPath, offset))
			args = append(args, m.TracesPath)

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
		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colState, offset))
		args = append(args, m.State)

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colChainID, offset))
		args = append(args, m.ChainID)

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colBlockNumber, offset))
		args = append(args, m.BlockNumber.String())

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colBlockHash, offset))
		args = append(args, m.BlockHash.String())

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colParentHash, offset))
		args = append(args, m.ParentHash.String())

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colBlockTimestamp, offset))
		args = append(args, m.BlockTimestamp)

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colBlockHeaderPath, offset))
		args = append(args, m.BlockHeaderPath)

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colTransactionsPath, offset))
		args = append(args, m.TransactionsPath)

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colReceiptsPath, offset))
		args = append(args, m.ReceiptsPath)

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colLogsPath, offset))
		args = append(args, m.LogsPath)

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colTracesPath, offset))
		args = append(args, m.TracesPath)

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
		return ErrSyncUpdate
	}

	return nil
}

// Delete deletes a Sync.
//
// The Sync must have the fields that are tag as "key" set.
func (repo *SyncRepo) Delete(ctx context.Context, m *Sync) error {
	var (
		where []string
		args  []interface{}
	)

	where = append(where, fmt.Sprintf("%s = $%d", repo.colID, 1))
	args = append(args, m.ID)

	qWhere := strings.Join(where, " AND ")

	sql := "DELETE FROM %s WHERE %s"
	sql = fmt.Sprintf(sql, repo.table, qWhere)

	_, err := repo.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "exec context")
	}

	return nil
}
