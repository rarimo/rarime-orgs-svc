/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type InvitationCreateEmail struct {
	Key
	Attributes InvitationCreateEmailAttributes `json:"attributes"`
}
type InvitationCreateEmailRequest struct {
	Data     InvitationCreateEmail `json:"data"`
	Included Included              `json:"included"`
}

type InvitationCreateEmailListRequest struct {
	Data     []InvitationCreateEmail `json:"data"`
	Included Included                `json:"included"`
	Links    *Links                  `json:"links"`
	Meta     json.RawMessage         `json:"meta,omitempty"`
}

func (r *InvitationCreateEmailListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *InvitationCreateEmailListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustInvitationCreateEmail - returns InvitationCreateEmail from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustInvitationCreateEmail(key Key) *InvitationCreateEmail {
	var invitationCreateEmail InvitationCreateEmail
	if c.tryFindEntry(key, &invitationCreateEmail) {
		return &invitationCreateEmail
	}
	return nil
}
