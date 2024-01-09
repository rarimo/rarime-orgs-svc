package data

import (
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type ClaimsSchemasSelector struct {
	ID         *uuid.UUID `json:"id,omitempty"`
	ActionType *string    `json:"action_type,omitempty"`
	SchemaType *string    `json:"schema_type,omitempty"`
	SchemaURL  *string    `json:"schema_url,omitempty"`

	PageCursor uint64     `json:"page_cursor,omitempty"`
	PageLimit  uint64     `json:"page_limit,omitempty"`
	Sort       pgdb.Sorts `json:"sort"`
}
