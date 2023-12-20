/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import (
	"encoding/json"
	"time"
)

type GroupAttributes struct {
	// The time (UTC) that the group was created in RFC3339 format
	CreatedAt time.Time `json:"created_at"`
	// JSON object containing the metadata information of the group
	Metadata json.RawMessage `json:"metadata"`
	// The ID of the organization that the group belongs to
	OrgId string `json:"org_id"`
	// JSON object containing the rules of the group, which will be used to generate claims for the group members
	Rules json.RawMessage `json:"rules"`
}
