package issuer

import (
	"github.com/rarimo/issuer/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

const (
	EmptyStringField = "none"
)

type ClaimType string

func (c ClaimType) String() string {
	return string(c)
}

// UUIDResponse is updated to include the "IssueClaimResponse" method.
type UUIDResponse struct {
	Id string `json:"id"`
}

type IssueClaimResponseData struct {
	resources.Key
}

type IssueClaimResponse struct {
	Data IssueClaimResponseData `json:"data"`
}

func (i IssueClaimResponse) GetKey() resources.Key {
	return i.Data.Key
}

func (r UUIDResponse) IssueClaimResponse() IssueClaimResponse {
	return IssueClaimResponse{
		Data: IssueClaimResponseData{
			resources.Key{
				ID:   r.Id,
				Type: resources.CLAIM_ID,
			},
		},
	}
}
