/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import (
	"encoding/json"
	"time"
)

type CredentialRequest struct {
	// The ID of the organization that the group belongs to
	CredentialSchema string `json:"credential_schema"`
	// JSON object containing the data to create a claim based on
	CredentialSubject json.RawMessage `json:"credential_subject"`
	// The time (UTC) when the claim expires
	Expiration time.Time `json:"expiration"`
	// Flag to include Merkle Tree proof or not
	MtProof bool `json:"mt_proof"`
	// Flag to include Signature proof or not
	SignatureProof bool `json:"signature_proof"`
	// The type of the claim
	Type string `json:"type"`
}
