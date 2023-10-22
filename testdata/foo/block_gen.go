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
)

var (
	// ErrBlockScan is the error that indicates a Block scan failed.
	ErrBlockScan = errors.New("scan")
	// ErrBlockNotFound is the error that indicates a Block was not found.
	ErrBlockNotFound = errors.New("not found")
	// ErrBlockUpdate is the error that indicates a Block was not updated.
	ErrBlockUpdate = errors.New("update")
)

// BlockRow is an interface for anything that can scan a Block, copying the columns from the matched
// row into the values pointed at by dest.
type BlockRow interface {
	Scan(dest ...any) error
}

// BlockRepo is a repository for the Block.
type BlockRepo struct {
	// db is the database connection.
	db *sql.DB

	// table is the table name.
	table string

	// colID is the Block.ID column name. It can be used in a queries to specify the column.
	colID string
	// colChainID is the Block.ChainID column name. It can be used in a queries to specify the column.
	colChainID string
	// colHash is the Block.Hash column name. It can be used in a queries to specify the column.
	colHash string
	// colNumber is the Block.Number column name. It can be used in a queries to specify the column.
	colNumber string
	// colParentHash is the Block.ParentHash column name. It can be used in a queries to specify the column.
	colParentHash string
	// colBlockTimestamp is the Block.BlockTimestamp column name. It can be used in a queries to specify the column.
	colBlockTimestamp string
}

// NewBlockRepo creates a new BlockRepo.
func NewBlockRepo(db *sql.DB, table string) *BlockRepo {
	return &BlockRepo{
		db:    db,
		table: table,

		colID:             "id",
		colChainID:        "chain_id",
		colHash:           "hash",
		colNumber:         "number",
		colParentHash:     "parent_hash",
		colBlockTimestamp: "block_timestamp",
	}
}

// Scan scans a Block from the given BlockRow (sql.Row|sql.Rows).
func (repo *BlockRepo) Scan(_ context.Context, s BlockRow) (*Block, error) {
	var (
		m Block

		chainID        int
		hash           sql.NullString
		number         int64
		parentHash     sql.NullString
		blockTimestamp sql.NullTime
	)

	err := s.Scan(
		&m.ID,
		&chainID,
		&hash,
		&number,
		&parentHash,
		&blockTimestamp,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBlockNotFound
		}

		return nil, errors.WrapError(err, ErrBlockScan)
	}

	m.ChainID = deps.ChainID(chainID)

	if hash.Valid {
		m.Hash = deps.HexToHash(hash.String)
	}

	m.Number = big.NewInt(number)

	if parentHash.Valid {
		m.ParentHash = deps.HexToHash(parentHash.String)
	}

	if blockTimestamp.Valid {
		m.BlockTimestamp = blockTimestamp.Time.UTC()
	}

	return &m, nil
}

// ScanAll scans a slice of Block from the given sql.Rows.
func (repo *BlockRepo) ScanAll(ctx context.Context, rs *sql.Rows) ([]*Block, error) {
	var ms []*Block

	for rs.Next() {
		m, err := repo.Scan(ctx, rs)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	if len(ms) == 0 {
		return nil, ErrBlockNotFound
	}

	return ms, nil
}

// Create creates a new Block and returns the persisted Block.
//
// The returned Block will contain the fields that were tag as "auto", which maybe were generated by the
// database..
func (repo *BlockRepo) Create(ctx context.Context, m *Block) (*Block, error) {
	var (
		cols []string
		args []interface{}
	)

	if !deps.IsUUIDZero(m.ID) {
		cols = append(cols, repo.colID)
		args = append(args, m.ID)
	}

	cols = append(cols, repo.colChainID)
	args = append(args, int(m.ChainID))

	var hash sql.NullString

	hash.String = m.Hash.String()
	hash.Valid = true

	cols = append(cols, repo.colHash)
	args = append(args, hash)

	if m.Number != nil {
		cols = append(cols, repo.colNumber)
		args = append(args, m.Number.Int64())
	}

	var parentHash sql.NullString

	parentHash.String = m.ParentHash.String()
	parentHash.Valid = true

	cols = append(cols, repo.colParentHash)
	args = append(args, parentHash)

	if !m.BlockTimestamp.IsZero() {
		var blockTimestamp sql.NullTime

		blockTimestamp.Time = m.BlockTimestamp
		blockTimestamp.Valid = true

		cols = append(cols, repo.colBlockTimestamp)
		args = append(args, blockTimestamp)
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
		repo.colHash,
		repo.colNumber,
		repo.colParentHash,
		repo.colBlockTimestamp,
	}, ", ")

	sql := "INSERT INTO %s (%s) VALUES (%s) RETURNING %s"
	sql = fmt.Sprintf(sql, repo.table, qCols, qValues, rCols)

	return repo.Scan(ctx, repo.db.QueryRowContext(ctx, sql, args...))
}

// Insert inserts one or more Block records into the database.
//
// When using this method the Block fields that are tag as "auto" should be set as the other fields non tag as "auto".
// The same applies for those other fields that are tag as "omitempty".
func (repo *BlockRepo) Insert(ctx context.Context, ms ...*Block) error {
	// Build values query.
	var (
		valuesQueryBuilder strings.Builder
		lms                = len(ms)
	)

	var cols []string

	cols = append(cols, repo.colID)
	cols = append(cols, repo.colChainID)
	cols = append(cols, repo.colHash)
	cols = append(cols, repo.colNumber)
	cols = append(cols, repo.colParentHash)
	cols = append(cols, repo.colBlockTimestamp)

	lcols := len(cols)

	// Size is equal to the number of models (lms) multiplied by the number of columns (lcols).
	args := make([]interface{}, 0, lms*lcols)

	for idx := range ms {
		m := ms[idx]

		indexOffset := idx * lcols

		valuesQueryBuilder.WriteString(repo.valuesStatement(cols, indexOffset, idx != lms-1))

		args = append(args, m.ID)

		args = append(args, int(m.ChainID))

		var hash sql.NullString

		hash.String = m.Hash.String()
		hash.Valid = true

		args = append(args, hash)

		if m.Number != nil {
			args = append(args, m.Number.Int64())
		}

		var parentHash sql.NullString

		parentHash.String = m.ParentHash.String()
		parentHash.Valid = true

		args = append(args, parentHash)

		var blockTimestamp sql.NullTime

		blockTimestamp.Time = m.BlockTimestamp
		blockTimestamp.Valid = true

		args = append(args, blockTimestamp)
	}

	qCols := strings.Join(cols, ", ")

	sql := "INSERT INTO %s (%s) VALUES %s"
	sql = fmt.Sprintf(sql, repo.table, qCols, valuesQueryBuilder.String())

	_, err := repo.db.ExecContext(ctx, sql, args...)
	if err != nil {
		// TODO: Check if this is the error is duplicate and return sentinel error.
		return errors.Wrap(err, "exec context")
	}

	return nil
}

// valuesStatement returns a string with the values statement ($n) for the given columns,
// starting from the given offset.
func (repo *BlockRepo) valuesStatement(cols []string, offset int, separator bool) string {
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

// Update updates a Block.
//
// skipZeroValues indicates whether to skip zero values from the update statement.
// In case of boolean fields, skipZeroValues is not applicable since false is the zero value of boolean and could be
// a potential update. Always set this type of fields.
//
// Returns the error ErrBlockUpdate if the Block was not updated and database did not error,
// otherwise database error.
func (repo *BlockRepo) Update(ctx context.Context, m *Block, skipZeroValues bool) error {
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
		if int(m.ChainID) != 0 {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colChainID, offset))
			args = append(args, int(m.ChainID))

			offset++
		}

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colHash, offset))
		args = append(args, m.Hash.String())

		offset++

		if m.Number != nil {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colNumber, offset))
			args = append(args, m.Number.Int64())

			offset++
		}

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colParentHash, offset))
		args = append(args, m.ParentHash.String())

		offset++

		if !m.BlockTimestamp.IsZero() {
			sets = append(sets, fmt.Sprintf("%s = $%d", repo.colBlockTimestamp, offset))
			args = append(args, m.BlockTimestamp)

			offset++
		}
	} else {
		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colChainID, offset))
		args = append(args, int(m.ChainID))

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colHash, offset))
		args = append(args, m.Hash.String())

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colNumber, offset))
		args = append(args, m.Number.Int64())

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colParentHash, offset))
		args = append(args, m.ParentHash.String())

		offset++

		sets = append(sets, fmt.Sprintf("%s = $%d", repo.colBlockTimestamp, offset))
		args = append(args, m.BlockTimestamp)

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
		return ErrBlockUpdate
	}

	return nil
}

// Delete deletes a Block.
//
// The Block must have the fields that are tag as "key" set.
func (repo *BlockRepo) Delete(ctx context.Context, m *Block) error {
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
