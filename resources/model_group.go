/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Group struct {
	Key
	Attributes    GroupAttributes    `json:"attributes"`
	Relationships GroupRelationships `json:"relationships"`
}
type GroupResponse struct {
	Data     Group    `json:"data"`
	Included Included `json:"included"`
}

type GroupListResponse struct {
	Data     []Group         `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *GroupListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *GroupListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustGroup - returns Group from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustGroup(key Key) *Group {
	var group Group
	if c.tryFindEntry(key, &group) {
		return &group
	}
	return nil
}
