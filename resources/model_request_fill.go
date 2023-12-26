/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type RequestFill struct {
	Key
	Attributes RequestFillAttributes `json:"attributes"`
}
type RequestFillRequest struct {
	Data     RequestFill `json:"data"`
	Included Included    `json:"included"`
}

type RequestFillListRequest struct {
	Data     []RequestFill   `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *RequestFillListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *RequestFillListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustRequestFill - returns RequestFill from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustRequestFill(key Key) *RequestFill {
	var requestFill RequestFill
	if c.tryFindEntry(key, &requestFill) {
		return &requestFill
	}
	return nil
}
