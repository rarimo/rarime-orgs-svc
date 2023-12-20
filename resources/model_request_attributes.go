/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import (
	"encoding/json"
	"time"
)

type RequestAttributes struct {
	// The time (UTC) that the group was created in RFC3339 format
	CreatedAt time.Time `json:"created_at"`
	// The ID of the group
	GroupId string `json:"group_id"`
	// JSON object containing the metadata information of the request
	Metadata json.RawMessage `json:"metadata"`
	// The ID of the organization that the group belongs to
	OrgId string `json:"org_id"`
	// The request status. `created` – The request was created by the group admin or organization owner or superadmin and hasn't been filled yet by the user. `filled` – The user filled the request but it hasn't been approved or rejected yet by the group admin or organization owner or superadmin. `approved` – The request was approved by the group admin or organization owner or superadmin. `rejected` – The request was rejected by the group admin or organization owner or superadmin. `submitted` – The request becomes submitted when claims were issues for the all fields in attributes and was submitted to the chain by the issuer service.
	Status RequestStatus `json:"status"`
	// The time (UTC) that the organization was updated in RFC3339 format
	UpdatedAt time.Time `json:"updated_at"`
	// The DID of the user that the request is associated with, can be empty if user hasn't accepted the invitation yet
	UserDid *string `json:"user_did,omitempty"`
}
