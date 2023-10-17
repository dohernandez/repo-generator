package deps

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/dohernandez/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type Hash = common.Hash

var HexToHash = common.HexToHash

type ChainID int

const (
	EthereumChainID ChainID = 1

	PolygonChainID ChainID = 137
)

type Address = common.Address

var HexToAddress = common.HexToAddress

type AssetContractType string

const (
	// AssetTypeNative refers to the native asset within an EVM compatible network.
	AssetTypeNative AssetContractType = "native"

	// AssetTypeERC20 allows for the implementation of a standard API for tokens within
	// smart contracts. Provides basic functionality to transfer tokens, as well as allow
	// tokens to be approved, so they  can be spent by another on-chain third party.
	//
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-20.md
	// https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/ERC20.sol
	AssetTypeERC20 AssetContractType = "erc20"

	// AssetTypeERC179 allows for the implementation of a standard API for tokens within
	// smart contracts. A simpler version of the ERC20 standard token contract, removing
	// transfer approvals.
	//
	// https://github.com/ethereum/EIPs/issues/179
	AssetTypeERC179 AssetContractType = "erc179"
)

type ContractMetadataKey string

// Keys used for assets contract metadata.
const (
	// ContractNameMetadataKey is the key used to identify the metadata name field for an asset contract.
	// This is commonly present in almost all asset contract standards.
	ContractNameMetadataKey ContractMetadataKey = "name"

	// ContractSymbolMetadataKey is the key used to identify the metadata symbol field for an asset contract.
	// This is commonly present in almost all asset contract standards.
	ContractSymbolMetadataKey ContractMetadataKey = "symbol"

	// ContractDecimalsMetadataKey is the key used to identify the metadata decimals field for an asset contract.
	// This is present in 'erc179','erc20','erc223', and 'erc777' type asset contract standard.
	ContractDecimalsMetadataKey ContractMetadataKey = "decimals"

	// ContractBaseURIMetadataKey is the key used to identify the metadata baseURI field for an asset contract.
	// This is commonly present in the 'erc721' type asset contract standard.
	ContractBaseURIMetadataKey ContractMetadataKey = "baseURI"

	// ContractImplementationMetadataKey is the key used to identify the metadata implementation field for an asset contract.
	// This is present in 'proxy' type asset contract standard.
	ContractImplementationMetadataKey ContractMetadataKey = "implementation"
)

type ContractMetadata map[ContractMetadataKey]any

// Value implements the driver Valuer interface.
func (m *ContractMetadata) Value() (driver.Value, error) {
	eContractMetadata, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return eContractMetadata, nil
}

// Scan implements the Scanner interface.
func (m *ContractMetadata) Scan(src any) error {
	if src == nil {
		return nil
	}

	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("cannot scan")
	}

	if err := json.Unmarshal(data, m); err != nil {
		return err
	}

	return nil
}

// UnmarshalJSON parses the JSON-encoded data to ContractMetadata.
func (m *ContractMetadata) UnmarshalJSON(data []byte) error {
	var cm map[string]any

	err := json.Unmarshal(data, &cm)
	if err != nil {
		return err
	}

	if cm == nil {
		return nil
	}

	md := make(map[ContractMetadataKey]any)

	for k, v := range cm {
		key := ContractMetadataKey(k)

		if k == string(ContractDecimalsMetadataKey) {
			value, ok := v.(float64)
			if ok {
				md[key] = uint8(value)
			}

			continue
		}

		value, ok := v.(string)
		if ok {
			if k == string(ContractImplementationMetadataKey) {
				md[key] = HexToAddress(value)
			} else {
				md[key] = value
			}
		}
	}

	*m = md

	return nil
}

type State string

func (s State) String() string {
	return string(s)
}

func IsUUIDZero(u uuid.UUID) bool {
	return u == uuid.Nil
}
