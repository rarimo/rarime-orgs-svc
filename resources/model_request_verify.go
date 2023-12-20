/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type RequestVerify struct {
	Key
	Attributes RequestVerifyAttributes `json:"attributes"`
}
type RequestVerifyRequest struct {
	Data     RequestVerify `json:"data"`
	Included Included      `json:"included"`
}

type RequestVerifyListRequest struct {
	Data     []RequestVerify `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *RequestVerifyListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *RequestVerifyListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustRequestVerify - returns RequestVerify from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustRequestVerify(key Key) *RequestVerify {
	var requestVerify RequestVerify
	if c.tryFindEntry(key, &requestVerify) {
		return &requestVerify
	}
	return nil
}
