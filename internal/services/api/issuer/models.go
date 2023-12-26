package issuer

import (
	"time"

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

type IdentityProviderName string

func (ipn IdentityProviderName) String() string {
	return string(ipn)
}

type IdentityProvidersCredentialSubject struct {
	IdentityID       string               `json:"id"`
	Provider         IdentityProviderName `json:"provider"`
	IsNatural        int64                `json:"isNatural"`
	Address          string               `json:"address"`
	ProviderMetadata string               `json:"providerMetadata"`
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

// CreateCredentialRequest defines model for CreateCredentialRequest.
type CreateCredentialRequest struct {
	CredentialSchema  string                              `json:"credentialSchema"`
	CredentialSubject *IdentityProvidersCredentialSubject `json:"credentialSubject"`
	Expiration        *time.Time                          `json:"expiration,omitempty"`
	MtProof           *bool                               `json:"mtProof,omitempty"`
	SignatureProof    *bool                               `json:"signatureProof,omitempty"`
	Type              string                              `json:"type"`
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
