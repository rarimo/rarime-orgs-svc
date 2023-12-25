/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ClaimOffer struct {
	Body     interface{} `json:"body"`
	From     string      `json:"from"`
	Id       string      `json:"id"`
	ThreadID string      `json:"threadID"`
	To       string      `json:"to"`
	Typ      string      `json:"typ"`
	Type     string      `json:"type"`
}

func (c ClaimOffer) GetKey() Key {
	return Key{c.Type, ResourceType(c.Id)}
}
