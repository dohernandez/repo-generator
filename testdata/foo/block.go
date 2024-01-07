package foo

import (
	"math/big"
	"time"

	"github.com/google/uuid"
	"repo-generator/testdata/deps"
)

type Block struct {
	ID             uuid.UUID    `db:"id,key,auto" nil:"deps.IsUUIDZero"`
	ChainID        deps.ChainID `db:"chain_id" type:"int" value:"int" scan:"deps.ChainID"`
	Hash           deps.Hash    `db:"hash,nullable" type:"string" value:".String" scan:"deps.HexToHash"`
	Number         *big.Int     `db:"number" type:"int64" value:".Int64" scan:"big.NewInt"`
	ParentHash     deps.Hash    `db:"parent_hash,nullable" type:"string" value:".String" scan:"deps.HexToHash"`
	BlockTimestamp time.Time    `db:"block_timestamp,nullable,omitempty" scan:".UTC"`
}
