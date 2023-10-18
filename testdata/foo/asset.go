package foo

import (
	"time"

	"github.com/dohernandez/repo-generator/testdata/deps"
)

type Asset struct {
	// BlockHash is the hash of the block for when this asset contract was created.
	BlockHash deps.Hash `db:"block_hash" type:"string" value:".String" scan:"deps.HexToHash"`

	// ChainID is the identifier for the network on which this log was emitted.
	ChainID deps.ChainID `db:"chain_id,key"`

	// Address is the address for the asset contract in lower case.
	Address deps.Address `db:"address,key" type:"string" value:".String" scan:"deps.HexToAddress"`

	// Type is a slice of AssetContractType representing the ERC asset standard for this asset contract.
	//
	// Types include; 'native','proxy','erc179','erc20','erc223','erc721','erc777','erc1155'.
	Type []deps.AssetContractType `db:"types,arrayable" type:"string" value:"string" scan:"deps.AssetContractType"`

	// Name is the name of the asset.
	Name string `db:"name"`

	// Symbol is the symbol of the asset.
	Symbol string `db:"symbol"`

	// Metadata is an optional field which can contain JSON-encoded metadata related to the asset.
	//
	// {"implementation": string} for 'proxy' assets.
	// {"decimals": int} for 'native','erc179','erc20','erc223','erc777'.
	// {"baseURI": string} for 'erc721'.
	// {"tokens": {"{tokenId}":{"uri": string, "type": "nft", "name": string, "decimals": int, "description": string, "image": string, "properties": string}}} for 'erc1155'.
	Metadata deps.ContractMetadata `db:"metadata" nil:".IsEmpty"`

	// Immutable specifies whether the asset could be updated or not.
	// By default, this is set to false.
	Immutable bool `db:"immutable"`

	// CreatedAt is the timestamp at which this record was indexed into the EVM database.
	CreatedAt time.Time `db:"created_at"`

	// UpdatedAt is the timestamp at which this record was updated within EVM database (can contain manual updates).
	UpdatedAt time.Time `db:"updated_at"`
}
