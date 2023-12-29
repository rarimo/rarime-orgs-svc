package models

import (
	"github.com/rarimo/rarime-orgs-svc/internal/services/core/issuer"
	"time"
)

type DomainVerificationCredentialSubject struct {
	IdentityID string `json:"id"`
	Domain     string `json:"domain"`
}

// CreateClaimDomainVerificationRequest defines model for CreateCredentialRequest.
type CreateClaimDomainVerificationRequest struct {
	CredentialSchema  string                               `json:"credentialSchema"`
	CredentialSubject *DomainVerificationCredentialSubject `json:"credentialSubject"`
	Expiration        *time.Time                           `json:"expiration,omitempty"`
	Type              string                               `json:"type"`
}

func NewEmptyDomainVerificationCredentialSubject() *DomainVerificationCredentialSubject {
	return &DomainVerificationCredentialSubject{
		IdentityID: issuer.EmptyStringField,
		Domain:     issuer.EmptyStringField,
	}
}
