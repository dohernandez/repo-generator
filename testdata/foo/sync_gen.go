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
	ErrSyncUpdate   = errors.New("update")
)

type SyncScanner interface {
	Scan(dest ...any) error
}

type SyncRepo struct {
	db *sql.DB

	table               string
	colID               string
	colState            string
	colChainID          string
	colBlockNumber      string
	colBlockHash        string
	colParentHash       string
	colBlockTimestamp   string
	colBlockHeaderPath  string
	colTransactionsPath string
	colReceiptsPath     string
	colLogsPath         string
	colTracesPath       string
	colCreatedAt        string
	colUpdatedAt        string

	cols []string
}

func NewSyncRepo(db *sql.DB, table string) *SyncRepo {
	cols := []string{
		"id",
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

	return &SyncRepo{
		db:                  db,
		table:               table,
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

		cols: cols,
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

func (repo *SyncRepo) Update(ctx context.Context, m *Sync) error {
	var (
		sets   []string
		where  []string
		args   []interface{}
		offset = 1
	)
	where = append(where, fmt.Sprintf("%s = $%d", repo.colID, offset))
	args = append(args, m.ID)

	offset++

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
