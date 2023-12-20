/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type GroupUser struct {
	Key
	Attributes    GroupUserAttributes    `json:"attributes"`
	Relationships GroupUserRelationships `json:"relationships"`
}
type GroupUserResponse struct {
	Data     GroupUser `json:"data"`
	Included Included  `json:"included"`
}

type GroupUserListResponse struct {
	Data     []GroupUser     `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *GroupUserListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *GroupUserListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustGroupUser - returns GroupUser from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustGroupUser(key Key) *GroupUser {
	var groupUser GroupUser
	if c.tryFindEntry(key, &groupUser) {
		return &groupUser
	}
	return nil
}
