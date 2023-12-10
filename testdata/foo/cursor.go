package foo

import (
	"time"

	"github.com/google/uuid"
)

type Cursor struct {
	// ID is a unique identifier generated for each new Cursor, usually a UUID.
	ID uuid.UUID `db:"id,key"`

	// Name is a human-readable, unique identifier for each new Cursor. This should be descriptive and unique across all
	// cursors in an Inbox.
	Name string `db:"name"`

	// Position is a foreign-key which refers to the ID of the last event which was successfully processed by a
	// Consumer.
	Position uuid.UUID `db:"position,nullable" type:"string" value:".String" scan:"uuid.MustParse"`

	// Leader is a unique identifier of a Subscriber which is deemed the "leader" of the group. Only subscribers which
	// carry this identifier will establish event streams.
	//
	// This is not a mandatory field and is only set if leader election is enabled for a SubscriberGroup.
	Leader uuid.UUID `db:"leader,nullable" type:"string" value:".String" scan:"uuid.MustParse"`

	// LeaderElectedAt is a timestamp which indicates when the current leader was elected. This will be empty if a
	// leader has never been elected.
	LeaderElectedAt time.Time `db:"leader_elected_at,nullable" scan:".UTC"`
	CreatedAt       time.Time `db:"created_at" scan:".UTC"`
	UpdatedAt       time.Time `db:"updated_at" scan:".UTC"`
}
