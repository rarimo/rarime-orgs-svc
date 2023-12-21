/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type UserAttributes struct {
	// The time (UTC) that the user was created in RFC3339 format
	CreatedAt time.Time `json:"created_at"`
	// The DID of the user
	Did *string `json:"did,omitempty"`
	// The ID of the organization that the user belongs to
	OrgId string `json:"org_id"`
	// The global role of the user. `undefined` – Basic user role. `owner` – The owner of some organization. `superadmin` – The user with the highest privileges.
	Role UserRole `json:"role"`
	// The time (UTC) that the user was updated in RFC3339 format
	UpdatedAt time.Time `json:"updated_at"`
}
