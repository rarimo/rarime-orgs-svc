/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import (
	"encoding/json"
	"time"
)

type OrganizationAttributes struct {
	// The time (UTC) that the organization was created in RFC3339 format
	CreatedAt time.Time `json:"created_at"`
	// The DID of the organization, can be empty for the organizations with the status `unverified`
	Did *string `json:"did,omitempty"`
	// The domain of the organization
	Domain string `json:"domain"`
	// The number of claims issued by the organization
	IssuedClaimsCount string `json:"issued_claims_count"`
	// The number of members in the organization
	MembersCount string `json:"members_count"`
	// JSON object containing the metadata information of the organization
	Metadata json.RawMessage `json:"metadata"`
	// The status of the organization. `unverified` – The organization was created by the user and hasn't been verified yet. `verified` – The owner verified the organization's domain by adding code to the DNS record, verifying by the service, and creating the organization's DID issuer. As the result of the verification – the owner of the organization receives an \"owner role claim\".
	Status OrganizationStatus `json:"status"`
	// The time (UTC) that the organization was updated in RFC3339 format
	UpdatedAt time.Time `json:"updated_at"`
	// The base64 encoded verification code that was by service to verify the domain of the organization. Can be empty for the organizations with the status `unverified`
	VerificationCode *string `json:"verification_code,omitempty"`
}
