package resources

import "encoding/json"

type GroupUserRole int32

const (
	GroupUserRole_Undefined  GroupUserRole = 0
	GroupUserRole_Employee   GroupUserRole = 1
	GroupUserRole_Admin      GroupUserRole = 2
	GroupUserRole_Superadmin GroupUserRole = 3

	groupUserRole_Undefined_Str  = "undefined"
	groupUserRole_Employee_Str   = "employee"
	groupUserRole_Admin_Str      = "admin"
	groupUserRole_Superadmin_Str = "superadmin"
)

var groupUserRoleIntStr = map[GroupUserRole]string{
	GroupUserRole_Undefined:  groupUserRole_Undefined_Str,
	GroupUserRole_Employee:   groupUserRole_Employee_Str,
	GroupUserRole_Admin:      groupUserRole_Admin_Str,
	GroupUserRole_Superadmin: groupUserRole_Superadmin_Str,
}

var groupUserRoleStrInt = map[string]GroupUserRole{
	groupUserRole_Undefined_Str:  GroupUserRole_Undefined,
	groupUserRole_Employee_Str:   GroupUserRole_Employee,
	groupUserRole_Admin_Str:      GroupUserRole_Admin,
	groupUserRole_Superadmin_Str: GroupUserRole_Superadmin,
}

func GroupUserRoleFromString(s string) GroupUserRole {
	return groupUserRoleStrInt[s]
}

func (r GroupUserRole) String() string {
	return groupUserRoleIntStr[r]
}

func (r GroupUserRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(Flag{
		Name:  groupUserRoleIntStr[r],
		Value: int32(r),
	})
}

func (r *GroupUserRole) UnmarshalJSON(b []byte) error {
	var res Flag
	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}

	*r = GroupUserRole(res.Value)
	return nil
}

func (r GroupUserRole) Int() int {
	return int(r)
}
func (r GroupUserRole) Int16() int16 {
	return int16(r)
}
