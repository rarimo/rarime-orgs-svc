package data

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

const MaxDigits = 6

type RequestsSelector struct {
	OrgID   *uuid.UUID `json:"org_id,omitempty"`
	GroupID *uuid.UUID `json:"group_id,omitempty"`

	UserDID *string `json:"user_did,omitempty"`
	Status  *int    `json:"status,omitempty"`

	PageCursor uint64     `json:"page_cursor,omitempty"`
	PageLimit  uint64     `json:"page_limit,omitempty"`
	Sort       pgdb.Sorts `json:"sort"`
}

func GenerateOTP() string {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(MaxDigits)))),
	)
	if err != nil {
		panic(errors.Wrap(err, "failed to generate uniform random value"))
	}
	return fmt.Sprintf("%0*d", MaxDigits, bi)
}
