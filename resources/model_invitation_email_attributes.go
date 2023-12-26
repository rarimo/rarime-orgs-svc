/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type InvitationEmailAttributes struct {
	// The time (UTC) that the email invitation was created in RFC3339 format
	CreatedAt time.Time `json:"created_at"`
	// The email address of the user that the invitation email is sent to
	Email string `json:"email"`
	// The ID of the group
	GroupId string `json:"group_id"`
	// The ID of the organization that the group belongs to
	OrgId string `json:"org_id"`
	// The ID of the request that the invitation email is associated with
	ReqId string `json:"req_id"`
}
