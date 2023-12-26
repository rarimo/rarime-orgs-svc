package data

import (
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type UsersSelector struct {
	OrgID *uuid.UUID `json:"owner,omitempty"`
	DID   *string    `json:"did,omitempty"`
	Role  *int       `json:"role,omitempty"`

	PageCursor uint64     `json:"page_cursor,omitempty"`
	PageLimit  uint64     `json:"page_limit,omitempty"`
	Sort       pgdb.Sorts `json:"sort"`
}
