package resources

import "encoding/json"

type OrganizationStatus int32

const (
	OrganizationStatus_Unverified OrganizationStatus = 0
	OrganizationStatus_Verified   OrganizationStatus = 1

	organizationStatus_Unverified_Str = "unverified"
	organizationStatus_Verified_Str   = "verified"
)

var organizationStatusIntStr = map[OrganizationStatus]string{
	OrganizationStatus_Unverified: organizationStatus_Unverified_Str,
	OrganizationStatus_Verified:   organizationStatus_Verified_Str,
}

var organizationStatusStrInt = map[string]OrganizationStatus{
	organizationStatus_Unverified_Str: OrganizationStatus_Unverified,
	organizationStatus_Verified_Str:   OrganizationStatus_Verified,
}

func OrganizationStatusFromString(s string) OrganizationStatus {
	return organizationStatusStrInt[s]
}

func (s OrganizationStatus) String() string {
	return organizationStatusIntStr[s]
}

func (s OrganizationStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(Flag{
		Name:  organizationStatusIntStr[s],
		Value: int32(s),
	})
}

func (s *OrganizationStatus) UnmarshalJSON(b []byte) error {
	var res Flag
	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}

	*s = OrganizationStatus(res.Value)
	return nil
}

func (s OrganizationStatus) Int() int {
	return int(s)
}
func (s OrganizationStatus) Int16() int16 {
	return int16(s)
}
func (s OrganizationStatus) Intp() *int {
	res := int(s)
	return &res
}
