/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Request struct {
	Key
	Attributes    RequestAttributes     `json:"attributes"`
	Relationships *RequestRelationships `json:"relationships,omitempty"`
}
type RequestResponse struct {
	Data     Request  `json:"data"`
	Included Included `json:"included"`
}

type RequestListResponse struct {
	Data     []Request       `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *RequestListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *RequestListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustRequest - returns Request from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustRequest(key Key) *Request {
	var request Request
	if c.tryFindEntry(key, &request) {
		return &request
	}
	return nil
}
