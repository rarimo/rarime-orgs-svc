/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type RequestVerifyAttributes struct {
	// The approval status of the verification request. `true` – The verification request was approved. `false` – The verification request was rejected. In that case `metadata` and `role` fields could be empty.
	Approved bool `json:"approved"`
	// JSON object containing the metadata information of the request
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	// The role of the group user. `undefined` – The group user that was created but hasn't been verified yet by the organization owner or group admin or superadmin. `employee` – The group user that was verified by the organization owner or group admin or superadmin by verifying verification request. `admin` – The group user that was verified by the organization owner or group admin or superadmin as admin and has the ability to verify other group users. `superadmin` – The group user that was verified by the superadmin as superadmin and has the ability to verify other group users and has the highest permissions in the group.
	Role *GroupUserRole `json:"role,omitempty"`
}
