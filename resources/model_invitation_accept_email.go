/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type InvitationAcceptEmail struct {
	Key
	Attributes InvitationAcceptEmailAttributes `json:"attributes"`
}
type InvitationAcceptEmailRequest struct {
	Data     InvitationAcceptEmail `json:"data"`
	Included Included              `json:"included"`
}

type InvitationAcceptEmailListRequest struct {
	Data     []InvitationAcceptEmail `json:"data"`
	Included Included                `json:"included"`
	Links    *Links                  `json:"links"`
	Meta     json.RawMessage         `json:"meta,omitempty"`
}

func (r *InvitationAcceptEmailListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *InvitationAcceptEmailListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustInvitationAcceptEmail - returns InvitationAcceptEmail from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustInvitationAcceptEmail(key Key) *InvitationAcceptEmail {
	var invitationAcceptEmail InvitationAcceptEmail
	if c.tryFindEntry(key, &invitationAcceptEmail) {
		return &invitationAcceptEmail
	}
	return nil
}
