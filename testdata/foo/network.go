package foo

import (
	"math/big"
	"time"

	"repo-generator/testdata/deps"
)

type Network struct {
	ID        string       `db:"id,key,auto"`
	URI       string       `db:"uri,nullable"`
	Token     string       `db:"token,key"`
	Number    *big.Int     `db:"number,nullable,omitempty" type:"int64" value:".Int64" scan:"big.NewInt"`
	Total     big.Int      `db:"total" type:"int64" value:".Int64" scan:"bigNewInt"`
	IP        *string      `db:"ip"`
	Address   deps.Address `db:"address,nullable,omitempty" type:"string" value:".String" scan:"stringToAddress"`
	CreatedAt time.Time    `db:"created_at,omitempty,auto"`
	UpdatedAt time.Time    `db:"updated_at,omitempty,auto"`
}

func bigNewInt(i int64) big.Int {
	bi := big.NewInt(i)

	return *bi
}

func stringToAddress(s string) deps.Address {
	return deps.HexToAddress(s)
}
