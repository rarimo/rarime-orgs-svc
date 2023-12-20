/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type GroupCreate struct {
	Key
	Attributes GroupCreateAttributes `json:"attributes"`
}
type GroupCreateRequest struct {
	Data     GroupCreate `json:"data"`
	Included Included    `json:"included"`
}

type GroupCreateListRequest struct {
	Data     []GroupCreate   `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *GroupCreateListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *GroupCreateListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustGroupCreate - returns GroupCreate from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustGroupCreate(key Key) *GroupCreate {
	var groupCreate GroupCreate
	if c.tryFindEntry(key, &groupCreate) {
		return &groupCreate
	}
	return nil
}
