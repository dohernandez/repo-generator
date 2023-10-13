package foo

// Code generated by repo-generator. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dohernandez/repo-generator/errors"
	"github.com/dohernandez/repo-generator/testdata/deps"
	"math/big"
	"strings"
	"time"
)

var (
	ErrTransferScan     = errors.New("scan")
	ErrTransferNotFound = errors.New("not found")
)

type TransferScanner interface {
	Scan(dest ...any) error
}

type TransferRepo struct {
	db    *sql.DB
	table string

	stateCols []string

	keyCols []string

	cols []string
}

func NewTransferRepo(db *sql.DB, table string) *TransferRepo {
	keyCols := []string{
		"id",
		"chain_id",
	}

	stateCols := []string{
		"block_hash",
		"block_timestamp",
		"transaction_hash",
		"method_id",
		"from_address",
		"to_address",
		"asset_contract",
		"value",
		"metadata",
		"trace_address",
		"created_at",
	}

	cols := append(keyCols, stateCols...)

	return &TransferRepo{
		db:    db,
		table: table,

		keyCols:   keyCols,
		stateCols: stateCols,
		cols:      cols,
	}
}

func (repo *TransferRepo) Scan(_ context.Context, s TransferScanner) (*Transfer, error) {
	var (
		m Transfer

		blockHash string

		transactionHash sql.NullString
		methodID        sql.NullString
		fromAddress     string
		toAddress       sql.NullString
		assetContract   string
		value           int64

		createdAt time.Time
	)

	err := s.Scan(
		&m.ID,
		&m.ChainID,
		&blockHash,
		&m.BlockTimestamp,
		&transactionHash,
		&methodID,
		&fromAddress,
		&toAddress,
		&assetContract,
		&value,
		&m.Metadata,
		&m.TraceAddress,
		&m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransferNotFound
		}

		return nil, errors.WrapWithError(err, ErrTransferScan)
	}

	m.BlockHash = deps.HexToHash(blockHash)

	if transactionHash.Valid {
		m.TransactionHash = toTransactionHash(transactionHash.String)
	}

	if methodID.Valid {
		tmp := methodID.String
		m.MethodID = &tmp
	}

	m.FromAddress = deps.HexToAddress(fromAddress)

	if toAddress.Valid {
		m.ToAddress = deps.HexToAddress(toAddress.String)
	}

	m.AssetContract = deps.HexToAddress(assetContract)

	m.Value = big.NewInt(value)

	m.CreatedAt = createdAt.UTC()

	return &m, nil
}

func (repo *TransferRepo) ScanAll(ctx context.Context, rs *sql.Rows) ([]*Transfer, error) {
	var ms []*Transfer

	for rs.Next() {
		m, err := repo.Scan(ctx, rs)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	if len(ms) == 0 {
		return nil, ErrTransferNotFound
	}

	return ms, nil
}

func (repo *TransferRepo) Create(ctx context.Context, m *Transfer) (*Transfer, error) {
	var (
		cols []string
		args []interface{}
	)

	cols = append(cols, "id")
	args = append(args, m.ID.String())

	cols = append(cols, "chain_id")
	args = append(args, m.ChainID)

	cols = append(cols, "block_hash")
	args = append(args, m.BlockHash.String())

	cols = append(cols, "block_timestamp")
	args = append(args, m.BlockTimestamp)

	if m.TransactionHash != nil {
		var transactionHash sql.NullString

		transactionHash.String = m.TransactionHash.String()
		transactionHash.Valid = true

		cols = append(cols, "transaction_hash")
		args = append(args, transactionHash)
	}

	if m.MethodID != nil {
		var methodID sql.NullString

		methodID.String = *m.MethodID
		methodID.Valid = true

		cols = append(cols, "method_id")
		args = append(args, methodID)
	}

	cols = append(cols, "from_address")
	args = append(args, m.FromAddress.String())

	var toAddress sql.NullString

	toAddress.String = m.ToAddress.String()
	toAddress.Valid = true

	cols = append(cols, "to_address")
	args = append(args, toAddress)

	cols = append(cols, "asset_contract")
	args = append(args, m.AssetContract.String())

	if m.Value != nil {
		cols = append(cols, "value")
		args = append(args, m.Value.Int64())
	}

	cols = append(cols, "metadata")
	args = append(args, m.Metadata)

	cols = append(cols, "trace_address")
	args = append(args, m.TraceAddress)

	if !m.CreatedAt.IsZero() {
		cols = append(cols, "created_at")
		args = append(args, m.CreatedAt)
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