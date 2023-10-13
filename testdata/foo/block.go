package foo

import (
	"math/big"
	"time"

	"github.com/google/uuid"

	"github.com/dohernandez/repo-generator/testdata/deps"
)

type Block struct {
	ID             uuid.UUID    `db:"id,key"`
	ChainID        deps.ChainID `db:"chain_id"`
	Hash           deps.Hash    `db:"hash,nullable" type:"string" value:".String" scan:"deps.HexToHash"`
	Number         *big.Int     `db:"number" type:"int64" value:".Int64" scan:"big.NewInt"`
	ParentHash     deps.Hash    `db:"parent_hash,nullable" type:"string" value:".String" scan:"deps.HexToHash"`
	BlockTimestamp time.Time    `db:"block_timestamp,nullable,omitempty" scan:".UTC"`
}
