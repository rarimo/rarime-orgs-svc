/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type InvitationAcceptEmailAttributes struct {
	// The one-time password that the user must use to accept the invitation
	Otp string `json:"otp"`
	// The DID of the user that the request is associated with, can be empty if user hasn't accepted the invitation yet
	UserDid string `json:"user_did"`
}
