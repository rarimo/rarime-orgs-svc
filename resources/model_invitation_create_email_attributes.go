/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type InvitationCreateEmailAttributes struct {
	// The email address of the user that the invitation email is sent to
	Email string `json:"email"`
	// JSON object containing the rules of the group with the predefined values
	Rules *json.RawMessage `json:"rules,omitempty"`
}
