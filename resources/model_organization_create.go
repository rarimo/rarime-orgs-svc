/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type OrganizationCreate struct {
	Key
	Attributes OrganizationCreateAttributes `json:"attributes"`
}
type OrganizationCreateRequest struct {
	Data     OrganizationCreate `json:"data"`
	Included Included           `json:"included"`
}

type OrganizationCreateListRequest struct {
	Data     []OrganizationCreate `json:"data"`
	Included Included             `json:"included"`
	Links    *Links               `json:"links"`
	Meta     json.RawMessage      `json:"meta,omitempty"`
}

func (r *OrganizationCreateListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *OrganizationCreateListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustOrganizationCreate - returns OrganizationCreate from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustOrganizationCreate(key Key) *OrganizationCreate {
	var organizationCreate OrganizationCreate
	if c.tryFindEntry(key, &organizationCreate) {
		return &organizationCreate
	}
	return nil
}
