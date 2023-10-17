package foo

// Code generated by repo-generator. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dohernandez/errors"
	"github.com/dohernandez/repo-generator/testdata/deps"
	"math/big"
	"strings"
)

var (
	ErrSyncScan     = errors.New("scan")
	ErrSyncNotFound = errors.New("not found")
)

type SyncScanner interface {
	Scan(dest ...any) error
}

type SyncRepo struct {
	db    *sql.DB
	table string

	stateCols []string

	keyCols []string

	cols []string
}

func NewSyncRepo(db *sql.DB, table string) *SyncRepo {
	keyCols := []string{
		"id",
	}

	stateCols := []string{
		"state",
		"chain_id",
		"block_number",
		"block_hash",
		"parent_hash",
		"block_timestamp",
		"block_header_path",
		"transactions_path",
		"receipts_path",
		"logs_path",
		"traces_path",
		"created_at",
		"updated_at",
	}

	cols := append(keyCols, stateCols...)

	return &SyncRepo{
		db:    db,
		table: table,

		keyCols:   keyCols,
		stateCols: stateCols,
		cols:      cols,
	}
}

func (repo *SyncRepo) Scan(_ context.Context, s SyncScanner) (*Sync, error) {
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

func (repo *SyncRepo) Create(ctx context.Context, m *Sync) (*Sync, error) {
	var (
		cols []string
		args []interface{}
	)

	if !deps.IsUUIDZero(m.ID) {
		cols = append(cols, "id")
		args = append(args, m.ID)
	}

	cols = append(cols, "state")
	args = append(args, m.State)

	cols = append(cols, "chain_id")
	args = append(args, m.ChainID)

	if m.BlockNumber != nil {
		cols = append(cols, "block_number")
		args = append(args, m.BlockNumber.String())
	}

	var blockHash sql.NullString

	blockHash.String = m.BlockHash.String()
	blockHash.Valid = true

	cols = append(cols, "block_hash")
	args = append(args, blockHash)

	var parentHash sql.NullString

	parentHash.String = m.ParentHash.String()
	parentHash.Valid = true

	cols = append(cols, "parent_hash")
	args = append(args, parentHash)

	var blockTimestamp sql.NullTime

	blockTimestamp.Time = m.BlockTimestamp
	blockTimestamp.Valid = true

	cols = append(cols, "block_timestamp")
	args = append(args, blockTimestamp)

	var blockHeaderPath sql.NullString

	blockHeaderPath.String = m.BlockHeaderPath
	blockHeaderPath.Valid = true

	cols = append(cols, "block_header_path")
	args = append(args, blockHeaderPath)

	var transactionsPath sql.NullString

	transactionsPath.String = m.TransactionsPath
	transactionsPath.Valid = true

	cols = append(cols, "transactions_path")
	args = append(args, transactionsPath)

	var receiptsPath sql.NullString

	receiptsPath.String = m.ReceiptsPath
	receiptsPath.Valid = true

	cols = append(cols, "receipts_path")
	args = append(args, receiptsPath)

	var logsPath sql.NullString

	logsPath.String = m.LogsPath
	logsPath.Valid = true

	cols = append(cols, "logs_path")
	args = append(args, logsPath)

	var tracesPath sql.NullString

	tracesPath.String = m.TracesPath
	tracesPath.Valid = true

	cols = append(cols, "traces_path")
	args = append(args, tracesPath)

	cols = append(cols, "created_at")
	args = append(args, m.CreatedAt)

	cols = append(cols, "updated_at")
	args = append(args, m.UpdatedAt)

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

func (repo *SyncRepo) Insert(ctx context.Context, ms ...*Sync) error {
	// Build values query.
	var (
		valuesQueryBuilder strings.Builder
		lms                = len(ms)
	)

	var cols []string
	cols = append(cols, "id")
	cols = append(cols, "state")
	cols = append(cols, "chain_id")
	cols = append(cols, "block_number")
	cols = append(cols, "block_hash")
	cols = append(cols, "parent_hash")
	cols = append(cols, "block_timestamp")
	cols = append(cols, "block_header_path")
	cols = append(cols, "transactions_path")
	cols = append(cols, "receipts_path")
	cols = append(cols, "logs_path")
	cols = append(cols, "traces_path")
	cols = append(cols, "created_at")
	cols = append(cols, "updated_at")

	lcols := len(cols)

	// Size is equal to the number of models (lms) multiplied by the number of columns (lcols).
	args := make([]interface{}, 0, lms*lcols)

	for i := range ms {
		m := ms[i]

		indexOffset := i * lcols
		valuesQueryBuilder.WriteString(repo.valuesStatement(cols, indexOffset, i != lms-1))

		if !deps.IsUUIDZero(m.ID) {
			args = append(args, m.ID)
		} else {
			args = append(args, nil)
		}

		args = append(args, m.State)

		args = append(args, m.ChainID)

		if m.BlockNumber != nil {
			args = append(args, m.BlockNumber.String())
		} else {
			args = append(args, nil)
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
		return errors.Wrap(err, "exec context")
	}

	return nil
}

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
