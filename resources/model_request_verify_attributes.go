/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import (
	"encoding/json"
	"time"
)

type RequestVerifyAttributes struct {
	// The time (UTC) when claims will be issued
	ActivationDate time.Time `json:"activation_date"`
	// The approval status of the verification request. `true` – The verification request was approved. `false` – The verification request was rejected. In that case `metadata` and `role` fields could be empty.
	Approved bool `json:"approved"`
	// The time (UTC) when claims expire
	ExpirationDate time.Time `json:"expiration_date"`
	// JSON object containing the metadata (background, name, etc) information of the request
	Metadata json.RawMessage `json:"metadata"`
	// The role of the group user. 0 - `undefined` – The group user that was created but hasn't been verified yet by the organization owner or group admin or superadmin. 1 - `employee` – The group user that was verified by the organization owner or group admin or superadmin by verifying verification request. 2 - `admin` – The group user that was verified by the organization owner or group admin or superadmin as admin and has the ability to verify other group users. 3 - `superadmin` – The group user that was verified by the superadmin as superadmin and has the ability to verify other group users and has the highest permissions in the group.
	Role string `json:"role"`
}
