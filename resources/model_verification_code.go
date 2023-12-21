/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type VerificationCode struct {
	Key
	Attributes VerificationCodeAttributes `json:"attributes"`
}
type VerificationCodeResponse struct {
	Data     VerificationCode `json:"data"`
	Included Included         `json:"included"`
}

type VerificationCodeListResponse struct {
	Data     []VerificationCode `json:"data"`
	Included Included           `json:"included"`
	Links    *Links             `json:"links"`
	Meta     json.RawMessage    `json:"meta,omitempty"`
}

func (r *VerificationCodeListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *VerificationCodeListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustVerificationCode - returns VerificationCode from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustVerificationCode(key Key) *VerificationCode {
	var verificationCode VerificationCode
	if c.tryFindEntry(key, &verificationCode) {
		return &verificationCode
	}
	return nil
}
