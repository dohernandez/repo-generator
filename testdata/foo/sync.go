package foo

import (
	"math/big"
	"time"

	"github.com/google/uuid"

	"github.com/dohernandez/repo-generator/testdata/deps"
)

type Sync struct {
	ID               uuid.UUID    `db:"id,key,auto" nil:"deps.IsUUIDZero"`
	State            deps.State   `db:"state"`
	ChainID          deps.ChainID `db:"chain_id"`
	BlockNumber      *big.Int     `db:"block_number" type:"int64" value:".String" scan:"big.NewInt"`
	BlockHash        deps.Hash    `db:"block_hash,nullable" type:"string" value:".String" scan:"deps.HexToHash"`
	ParentHash       deps.Hash    `db:"parent_hash,nullable" type:"string" value:".String" scan:"deps.HexToHash"`
	BlockTimestamp   time.Time    `db:"block_timestamp,nullable"`
	BlockHeaderPath  string       `db:"block_header_path,nullable"`
	TransactionsPath string       `db:"transactions_path,nullable"`
	ReceiptsPath     string       `db:"receipts_path,nullable"`
	LogsPath         string       `db:"logs_path,nullable"`
	TracesPath       string       `db:"traces_path,nullable"`
	CreatedAt        time.Time    `db:"created_at"`
	UpdatedAt        time.Time    `db:"updated_at"`
}
