/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type OrganizationCreateAttributes struct {
	// The domain of the organization
	Domain string `json:"domain"`
	// JSON object containing the metadata information of the organization
	Metadata json.RawMessage `json:"metadata"`
	// The DID of the owner of the organization
	OwnerDid string `json:"owner_did"`
}
