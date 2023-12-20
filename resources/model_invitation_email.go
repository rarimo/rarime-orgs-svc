/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type InvitationEmail struct {
	Key
	Attributes    InvitationEmailAttributes    `json:"attributes"`
	Relationships InvitationEmailRelationships `json:"relationships"`
}
type InvitationEmailResponse struct {
	Data     InvitationEmail `json:"data"`
	Included Included        `json:"included"`
}

type InvitationEmailListResponse struct {
	Data     []InvitationEmail `json:"data"`
	Included Included          `json:"included"`
	Links    *Links            `json:"links"`
	Meta     json.RawMessage   `json:"meta,omitempty"`
}

func (r *InvitationEmailListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *InvitationEmailListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustInvitationEmail - returns InvitationEmail from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustInvitationEmail(key Key) *InvitationEmail {
	var invitationEmail InvitationEmail
	if c.tryFindEntry(key, &invitationEmail) {
		return &invitationEmail
	}
	return nil
}
