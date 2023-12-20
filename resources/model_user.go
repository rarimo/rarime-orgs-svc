/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type User struct {
	Key
	Attributes UserAttributes `json:"attributes"`
}
type UserResponse struct {
	Data     User     `json:"data"`
	Included Included `json:"included"`
}

type UserListResponse struct {
	Data     []User          `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *UserListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *UserListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustUser - returns User from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustUser(key Key) *User {
	var user User
	if c.tryFindEntry(key, &user) {
		return &user
	}
	return nil
}
