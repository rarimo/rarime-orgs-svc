package data

import (
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type GroupsSelector struct {
	OrgID *uuid.UUID `json:"owner,omitempty"`

	PageCursor uint64     `json:"page_number,omitempty"`
	PageSize   uint64     `json:"page_size,omitempty"`
	Sort       pgdb.Sorts `json:"sort"`
}
