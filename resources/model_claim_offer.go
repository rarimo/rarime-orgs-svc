/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type ClaimOffer struct {
	Key
	Attributes ClaimOfferAttributes `json:"attributes"`
}
type ClaimOfferResponse struct {
	Data     ClaimOffer `json:"data"`
	Included Included   `json:"included"`
}

type ClaimOfferListResponse struct {
	Data     []ClaimOffer    `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *ClaimOfferListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *ClaimOfferListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustClaimOffer - returns ClaimOffer from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustClaimOffer(key Key) *ClaimOffer {
	var claimOffer ClaimOffer
	if c.tryFindEntry(key, &claimOffer) {
		return &claimOffer
	}
	return nil
}
