package foo

// Code generated by repo-generator v0.1.0. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"repo-generator/testdata/deps"
	"strings"
	"time"

	"github.com/dohernandez/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrTransferScan is the error that indicates a Transfer scan failed.
	ErrTransferScan = errors.New("scan")
	// ErrTransferNotFound is the error that indicates a Transfer was not found.
	ErrTransferNotFound = errors.New("not found")
	// ErrTransferExists is returned when the Transfer already exists.
	ErrTransferExists = errors.New("exists")
)

// TransferRow is an interface for anything that can scan a Transfer, copying the columns from the matched
// row into the values pointed at by dest.
type TransferRow interface {
	Scan(dest ...any) error
}

// TransferSQLDB is an interface for anything that can execute the SQL statements needed to
// perform the Transfer operations.
type TransferSQLDB interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// TransferRepo is a repository for the Transfer.
type TransferRepo struct {
	// db is the database connection.
	db TransferSQLDB

	// table is the table name.
	table string

	// colID is the Transfer.ID column name. It can be used in a queries to specify the column.
	colID string
	// colChainID is the Transfer.ChainID column name. It can be used in a queries to specify the column.
	colChainID string
	// colBlockHash is the Transfer.BlockHash column name. It can be used in a queries to specify the column.
	colBlockHash string
	// colBlockTimestamp is the Transfer.BlockTimestamp column name. It can be used in a queries to specify the column.
	colBlockTimestamp string
	// colTransactionHash is the Transfer.TransactionHash column name. It can be used in a queries to specify the column.
	colTransactionHash string
	// colMethodID is the Transfer.MethodID column name. It can be used in a queries to specify the column.
	colMethodID string
	// colFromAddress is the Transfer.FromAddress column name. It can be used in a queries to specify the column.
	colFromAddress string
	// colToAddress is the Transfer.ToAddress column name. It can be used in a queries to specify the column.
	colToAddress string
	// colAssetContract is the Transfer.AssetContract column name. It can be used in a queries to specify the column.
	colAssetContract string
	// colValue is the Transfer.Value column name. It can be used in a queries to specify the column.
	colValue string
	// colMetadata is the Transfer.Metadata column name. It can be used in a queries to specify the column.
	colMetadata string
	// colTraceAddress is the Transfer.TraceAddress column name. It can be used in a queries to specify the column.
	colTraceAddress string
	// colCreatedAt is the Transfer.CreatedAt column name. It can be used in a queries to specify the column.
	colCreatedAt string
}

// NewTransferRepo creates a new TransferRepo.
func NewTransferRepo(db TransferSQLDB, table string) *TransferRepo {
	return &TransferRepo{
		db:    db,
		table: table,

		colID:              "id",
		colChainID:         "chain_id",
		colBlockHash:       "block_hash",
		colBlockTimestamp:  "block_timestamp",
		colTransactionHash: "transaction_hash",
		colMethodID:        "method_id",
		colFromAddress:     "from_address",
		colToAddress:       "to_address",
		colAssetContract:   "asset_contract",
		colValue:           "value",
		colMetadata:        "metadata",
		colTraceAddress:    "trace_address",
		colCreatedAt:       "created_at",
	}
}

// Table returns the table name.
func (repo *TransferRepo) Table() string {
	return repo.table
}

// Cols returns the represented cols of Transfer.
// Cols are returned in the order they are scanned.
func (repo *TransferRepo) Cols() []string {
	return []string{
		repo.colID,
		repo.colChainID,
		repo.colBlockHash,
		repo.colBlockTimestamp,
		repo.colTransactionHash,
		repo.colMethodID,
		repo.colFromAddress,
		repo.colToAddress,
		repo.colAssetContract,
		repo.colValue,
		repo.colMetadata,
		repo.colTraceAddress,
		repo.colCreatedAt,
	}
}

// Scan scans a Transfer from the given TransferRow (sql.Row|sql.Rows).
func (repo *TransferRepo) Scan(_ context.Context, s TransferRow) (*Transfer, error) {
	var (
		m Transfer

		blockHash string

		transactionHash sql.NullString
		methodId        sql.NullString
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
		&methodId,
		&fromAddress,
		&toAddress,
		&assetContract,
		&value,
		&m.Metadata,
		&m.TraceAddress,
		&createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransferNotFound
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, errors.WrapError(err, ErrTransferExists)
		}

		return nil, errors.WrapError(err, ErrTransferScan)
	}

	m.BlockHash = deps.HexToHash(blockHash)

	if transactionHash.Valid {
		m.TransactionHash = toTransactionHash(transactionHash.String)
	}

	if methodId.Valid {
		tmp := methodId.String
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

// ScanAll scans a slice of Transfer from the given sql.Rows.
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

// Create creates a new Transfer and returns the persisted Transfer.
//
// The returned Transfer will contain the fields that were tag as "auto", which maybe were generated by the
// database..
func (repo *TransferRepo) Create(ctx context.Context, m *Transfer) (*Transfer, error) {
	var (
		cols []string
		args []interface{}
	)

	// TODO: For the correct operation of the code, the nil method must be implemented for repo.colID.
	// Define the method using the tag 'nil'.
	cols = append(cols, repo.colID)
	args = append(args, m.ID.String())

	cols = append(cols, repo.colChainID)
	args = append(args, m.ChainID)

	cols = append(cols, repo.colBlockHash)
	args = append(args, m.BlockHash.String())

	cols = append(cols, repo.colBlockTimestamp)
	args = append(args, m.BlockTimestamp)

	if m.TransactionHash != nil {
		var transactionHash sql.NullString

		transactionHash.String = m.TransactionHash.String()
		transactionHash.Valid = true

		cols = append(cols, repo.colTransactionHash)
		args = append(args, transactionHash)
	}

	if m.MethodID != nil {
		var methodId sql.NullString

		methodId.String = *m.MethodID
		methodId.Valid = true

		cols = append(cols, repo.colMethodID)
		args = append(args, methodId)
	}

	cols = append(cols, repo.colFromAddress)
	args = append(args, m.FromAddress.String())

	var toAddress sql.NullString

	toAddress.String = m.ToAddress.String()
	toAddress.Valid = true

	cols = append(cols, repo.colToAddress)
	args = append(args, toAddress)

	cols = append(cols, repo.colAssetContract)
	args = append(args, m.AssetContract.String())

	if m.Value != nil {
		cols = append(cols, repo.colValue)
		args = append(args, m.Value.Int64())
	}

	cols = append(cols, repo.colMetadata)
	args = append(args, m.Metadata)

	cols = append(cols, repo.colTraceAddress)
	args = append(args, m.TraceAddress)

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
		repo.colChainID,
		repo.colBlockHash,
		repo.colBlockTimestamp,
		repo.colTransactionHash,
		repo.colMethodID,
		repo.colFromAddress,
		repo.colToAddress,
		repo.colAssetContract,
		repo.colValue,
		repo.colMetadata,
		repo.colTraceAddress,
		repo.colCreatedAt,
	}, ", ")

	sql := "INSERT INTO %s (%s) VALUES (%s) RETURNING %s"
	sql = fmt.Sprintf(sql, repo.table, qCols, qValues, rCols)

	return repo.Scan(ctx, repo.db.QueryRowContext(ctx, sql, args...))
}
