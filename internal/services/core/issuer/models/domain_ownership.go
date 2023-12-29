package models

import (
	"github.com/rarimo/rarime-orgs-svc/internal/services/core/issuer"
	"time"
)

type DomainOwnershipCredentialSubject struct {
	IdentityID string `json:"id"`
	Domain     string `json:"domain"`
}

// DomainOwnershipCredentialRequest defines model for CreateCredentialRequest.
type DomainOwnershipCredentialRequest struct {
	CredentialSchema  string                            `json:"credentialSchema"`
	CredentialSubject *DomainOwnershipCredentialSubject `json:"credentialSubject"`
	Expiration        *time.Time                        `json:"expiration,omitempty"`
	Type              string                            `json:"type"`
}

func NewEmptyDomainOwnershipCredentialSubject() *DomainOwnershipCredentialSubject {
	return &DomainOwnershipCredentialSubject{
		IdentityID: issuer.EmptyStringField,
		Domain:     issuer.EmptyStringField,
	}
}
