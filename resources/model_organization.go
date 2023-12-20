/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Organization struct {
	Key
	Attributes    OrganizationAttributes    `json:"attributes"`
	Relationships OrganizationRelationships `json:"relationships"`
}
type OrganizationResponse struct {
	Data     Organization `json:"data"`
	Included Included     `json:"included"`
}

type OrganizationListResponse struct {
	Data     []Organization  `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *OrganizationListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *OrganizationListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustOrganization - returns Organization from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustOrganization(key Key) *Organization {
	var organization Organization
	if c.tryFindEntry(key, &organization) {
		return &organization
	}
	return nil
}
