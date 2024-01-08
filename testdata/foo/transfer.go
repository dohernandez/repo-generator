package foo

import (
	"math/big"
	"time"

	"github.com/google/uuid"
	"repo-generator/testdata/deps"
)

type Transfer struct {
	ID uuid.UUID `db:"id,key,auto" value:".String"`
	// ChainID is the identifier for the network on which this log was emitted.
	ChainID deps.ChainID `db:"chain_id,key"`

	// BlockHash is the hash of the block including this transfer.
	BlockHash deps.Hash `db:"block_hash" type:"string" value:".String" scan:"deps.HexToHash"`

	// BlockTimestamp is the timestamp in seconds for when the block was collated.
	BlockTimestamp uint64 `db:"block_timestamp"`

	// TransactionHash is the hash of the transaction (tx.hash).
	TransactionHash *deps.Hash `db:"transaction_hash,nullable,omitempty" type:"string" value:".String" scan:"toTransactionHash"`

	// MethodID the function signature like `0x792004fe` taken from the transaction input,
	// i.e. first 10 char of the input if not empty.
	MethodID *string `db:"method_id,nullable,omitempty"`

	// FromAddress is the address of the sender in lower case.
	FromAddress deps.Address `db:"from_address" type:"string" value:".String" scan:"deps.HexToAddress"`

	// ToAddress is the address of the sender in lower case.
	// This will be set to an empty address when it's a contract creation transfer.
	ToAddress deps.Address `db:"to_address,nullable,omitempty" type:"string" value:".String" scan:"deps.HexToAddress"`

	// AssetContract is the address for the asset contract in lower case.
	AssetContract deps.Address `db:"asset_contract" type:"string" value:".String" scan:"deps.HexToAddress"`

	// Value is the value transferred in Wei for native transfers.
	// For non-native transfers this is referring to the value of
	// the underlying asset. For NFTs (non-native) this would be the token ID.
	Value *big.Int `db:"value" type:"int64" value:".Int64" scan:"big.NewInt"`

	//// Metadata is an optional field which can contain JSON-encoded metadata related to the transfer.
	Metadata []byte `db:"metadata"`
	//
	//// TraceAddress is the index for a given trace in the trace tree.
	//// Below is an example representation to understand how to read the position.
	////
	////      []
	////     /  \
	////   [0]  [1]
	////      /  |  \
	//// [1,0] [1,1] [1,2]
	////         |         \
	////      [1,1,0]    [1,2,0]
	////
	//// This will be null if the trace is of type "reward", "genesis" or "daofork".
	TraceAddress []int32 `db:"trace_address"`
	//
	//// CreatedAt is the timestamp at which this record was indexed into the EVM database.
	CreatedAt time.Time `db:"created_at,auto" scan:".UTC"`
}

func toTransactionHash(h string) *deps.Hash {
	tmp := deps.HexToHash(h)

	return &tmp
}
