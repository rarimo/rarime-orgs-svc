/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type GroupCreateAttributes struct {
	// JSON object containing the metadata information of the group
	Metadata json.RawMessage `json:"metadata"`
	// The ID of the organization that the group belongs to
	OrgId string `json:"org_id"`
	// JSON object containing the rules of the group, which will be used to generate claims for the group members
	Rules json.RawMessage `json:"rules"`
}
