/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Owner struct {
	Key
	Attributes OwnerAttributes `json:"attributes"`
}
type OwnerResponse struct {
	Data     Owner    `json:"data"`
	Included Included `json:"included"`
}

type OwnerListResponse struct {
	Data     []Owner         `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *OwnerListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *OwnerListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustOwner - returns Owner from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustOwner(key Key) *Owner {
	var owner Owner
	if c.tryFindEntry(key, &owner) {
		return &owner
	}
	return nil
}
