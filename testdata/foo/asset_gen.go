package foo

// Code generated by repo-generator. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dohernandez/errors"
	"github.com/dohernandez/repo-generator/testdata/deps"
	"github.com/lib/pq"
	"strings"
)

var (
	ErrAssetScan     = errors.New("scan")
	ErrAssetNotFound = errors.New("not found")
	ErrAssetUpdate   = errors.New("update")
)

type AssetScanner interface {
	Scan(dest ...any) error
}

type AssetRepo struct {
	db *sql.DB

	table        string
	colChainID   string
	colAddress   string
	colBlockHash string
	colType      string
	colName      string
	colSymbol    string
	colMetadata  string
	colImmutable string
	colCreatedAt string
	colUpdatedAt string

	cols []string
}

func NewAssetRepo(db *sql.DB, table string) *AssetRepo {
	cols := []string{
		"chain_id",
		"address",
		"block_hash",
		"types",
		"name",
		"symbol",
		"metadata",
		"immutable",
		"created_at",
		"updated_at",
	}

	return &AssetRepo{
		db:           db,
		table:        table,
		colChainID:   "chain_id",
		colAddress:   "address",
		colBlockHash: "block_hash",
		colType:      "types",
		colName:      "name",
		colSymbol:    "symbol",
		colMetadata:  "metadata",
		colImmutable: "immutable",
		colCreatedAt: "created_at",
		colUpdatedAt: "updated_at",

		cols: cols,
	}
}

func (repo *AssetRepo) Scan(_ context.Context, s AssetScanner) (*Asset, error) {
	var (
		m Asset

		address   string
		blockHash string
		typs      pq.StringArray
	)

	err := s.Scan(
		&m.ChainID,
		&address,
		&blockHash,
		&typs,
		&m.Name,
		&m.Symbol,
		&m.Metadata,
		&m.Immutable,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAssetNotFound
		}

		return nil, errors.WrapError(err, ErrAssetScan)
	}

	m.Address = deps.HexToAddress(address)
	m.BlockHash = deps.HexToHash(blockHash)

	for i := range typs {
		m.Type = append(m.Type, deps.AssetContractType(typs[i]))
	}

	return &m, nil
}

func (repo *AssetRepo) ScanAll(ctx context.Context, rs *sql.Rows) ([]*Asset, error) {
	var ms []*Asset

	for rs.Next() {
		m, err := repo.Scan(ctx, rs)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	if len(ms) == 0 {
		return nil, ErrAssetNotFound
	}

	return ms, nil
}

func (repo *AssetRepo) Create(ctx context.Context, m *Asset) (*Asset, error) {
	var (
		cols []string
		args []interface{}
	)

	cols = append(cols, repo.colChainID)
	args = append(args, m.ChainID)

	cols = append(cols, repo.colAddress)
	args = append(args, m.Address.String())

	cols = append(cols, repo.colBlockHash)
	args = append(args, m.BlockHash.String())

	cols = append(cols, repo.colType)

	typs := make([]string, len(m.Type))

	for i := range m.Type {
		typs[i] = string(m.Type[i])
	}

	args = append(args, typs)

	cols = append(cols, repo.colName)
	args = append(args, m.Name)

	cols = append(cols, repo.colSymbol)
	args = append(args, m.Symbol)

	cols = append(cols, repo.colMetadata)
	args = append(args, m.Metadata)

	cols = append(cols, repo.colImmutable)
	args = append(args, m.Immutable)

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

	sql := "INSERT INTO %s (%s) VALUES (%s)"
	sql = fmt.Sprintf(sql, repo.table, qCols, qValues)

	_, err := repo.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "exec context")
	}

	return m, nil
}

func (repo *AssetRepo) Insert(ctx context.Context, ms ...*Asset) error {
	// Build values query.
	var (
		valuesQueryBuilder strings.Builder
		lms                = len(ms)
	)

	var cols []string
	cols = append(cols, repo.colChainID)
	cols = append(cols, repo.colAddress)
	cols = append(cols, repo.colBlockHash)
	cols = append(cols, repo.colType)
	cols = append(cols, repo.colName)
	cols = append(cols, repo.colSymbol)
	cols = append(cols, repo.colMetadata)
	cols = append(cols, repo.colImmutable)
	cols = append(cols, repo.colCreatedAt)
	cols = append(cols, repo.colUpdatedAt)

	lcols := len(cols)

	// Size is equal to the number of models (lms) multiplied by the number of columns (lcols).
	args := make([]interface{}, 0, lms*lcols)

	for i := range ms {
		m := ms[i]

		indexOffset := i * lcols
		valuesQueryBuilder.WriteString(repo.valuesStatement(cols, indexOffset, i != lms-1))

		args = append(args, m.ChainID)

		args = append(args, m.Address.String())

		args = append(args, m.BlockHash.String())

		typs := make([]string, len(m.Type))

		for i := range m.Type {
			typs[i] = string(m.Type[i])
		}

		args = append(args, typs)

		args = append(args, m.Name)

		args = append(args, m.Symbol)

		args = append(args, m.Metadata)

		args = append(args, m.Immutable)

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

func (repo *AssetRepo) valuesStatement(cols []string, offset int, separator bool) string {
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

func (repo *AssetRepo) Update(ctx context.Context, m *Asset) error {
	var (
		sets   []string
		where  []string
		args   []interface{}
		offset = 1
	)
	where = append(where, fmt.Sprintf("%s = $%d", repo.colChainID, offset))
	args = append(args, m.ChainID)

	offset++
	where = append(where, fmt.Sprintf("%s = $%d", repo.colAddress, offset))
	args = append(args, m.Address)

	offset++

	if m.BlockHash.String() != "" {
		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colBlockHash, offset))
		args = append(args, m.BlockHash.String())

		offset++
	}

	if len(m.Type) > 0 {
		typs := make([]string, len(m.Type))

		for i := range m.Type {
			typs[i] = string(m.Type[i])
		}

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colType, offset))
		args = append(args, typs)

		offset++
	}

	if m.Name != "" {
		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colName, offset))
		args = append(args, m.Name)

		offset++
	}

	if m.Symbol != "" {
		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colSymbol, offset))
		args = append(args, m.Symbol)

		offset++
	}

	if !m.Metadata.IsEmpty() {
		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colMetadata, offset))
		args = append(args, m.Metadata)

		offset++
	}

	sets = append(sets, fmt.Sprintf("%s = $%d", repo.colImmutable, offset))
	args = append(args, m.Immutable)

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
		return ErrAssetUpdate
	}

	return nil
}
