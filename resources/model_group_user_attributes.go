/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type GroupUserAttributes struct {
	// The time (UTC) that the group user was created in RFC3339 format
	CreatedAt time.Time `json:"created_at"`
	// The ID of the group that the user belongs to
	GroupId string `json:"group_id"`
	// The role of the group user. `undefined` – The group user that was created but hasn't been verified yet by the organization owner or group admin or superadmin. `employee` – The group user that was verified by the organization owner or group admin or superadmin by verifying verification request. `admin` – The group user that was verified by the organization owner or group admin or superadmin as admin and has the ability to verify other group users. `superadmin` – The group user that was verified by the superadmin as superadmin and has the ability to verify other group users and has the highest permissions in the group.
	Role GroupUserRole `json:"role"`
	// The time (UTC) that the group user was updated in RFC3339 format
	UpdatedAt time.Time `json:"updated_at"`
	// The ID of the user that belongs to the group user
	UserId string `json:"user_id"`
}
