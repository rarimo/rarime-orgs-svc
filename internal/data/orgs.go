package data

import (
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type OrgsSelector struct {
	Owner   *uuid.UUID `json:"owner,omitempty"`
	UserDID *string    `json:"user_did,omitempty"`
	Status  *int       `json:"status,omitempty"`

	PageCursor uint64     `json:"page_number,omitempty"`
	PageSize   uint64     `json:"page_size,omitempty"`
	Sort       pgdb.Sorts `json:"sort"`
}
