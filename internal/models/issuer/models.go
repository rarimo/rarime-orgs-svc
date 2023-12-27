package issuer

import (
	"github.com/pkg/errors"
	"github.com/rarimo/issuer/resources"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

type ClaimType string

func (c ClaimType) String() string {
	return string(c)
}

type DomainVerificationCredentialSubject struct {
	IdentityID string `json:"id"`
}

type IssueClaimResponse struct {
	Data IssueClaimResponseData `json:"data"`
}

func (i IssueClaimResponse) GetKey() resources.Key {
	return i.Data.Key
}

type IssueClaimResponseData struct {
	resources.Key
}

// CreateClaimDomainVerificationRequest defines model for CreateCredentialRequest.
type CreateClaimDomainVerificationRequest struct {
	CredentialSchema  string                               `json:"credentialSchema"`
	CredentialSubject *DomainVerificationCredentialSubject `json:"credentialSubject"`
	Type              string                               `json:"type"`
}

type UUIDResponse struct {
	Id string `json:"id"`
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
