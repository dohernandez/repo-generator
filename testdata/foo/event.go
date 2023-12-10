package foo

import (
	"time"

	"github.com/google/uuid"
)

// Event defines the core message structure which is transmitted via streams.
type Event struct {
	// ID is a unique identifier generated for each new Event, usually a UUID.
	ID uuid.UUID `db:"id,key"`

	// Topic is a partitioning key which is used to group similar events together.
	Topic string `db:"topic"`

	// Key is an optional extra partitioning key and reference to external data.
	Key string `db:"key"`

	// Sequence is a unique integer which is sequential for events being generated. It is used to ensure ordering and
	// assist in load balancing.
	Sequence uint64 `db:"sequence"`

	// Metadata is an optional field which can contain any serializable data.
	Metadata []byte `db:"metadata"`

	// CreatedAt is the timestamp of when the Event was created.
	CreatedAt time.Time `db:"created_at,auto"`
}
