/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type ClaimSchema struct {
	Key
	Attributes ClaimSchemaAttributes `json:"attributes"`
}
type ClaimSchemaResponse struct {
	Data     ClaimSchema `json:"data"`
	Included Included    `json:"included"`
}

type ClaimSchemaListResponse struct {
	Data     []ClaimSchema   `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *ClaimSchemaListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *ClaimSchemaListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustClaimSchema - returns ClaimSchema from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustClaimSchema(key Key) *ClaimSchema {
	var claimSchema ClaimSchema
	if c.tryFindEntry(key, &claimSchema) {
		return &claimSchema
	}
	return nil
}
