package resources

import "encoding/json"

type UserRole int32

const (
	UserRole_Undefined  UserRole = 0
	UserRole_Owner      UserRole = 2
	UserRole_Superadmin UserRole = 3

	userRole_Undefined_Str  = "undefined"
	userRole_Owner_Str      = "owner"
	userRole_Superadmin_Str = "superadmin"
)

var userRoleIntStr = map[UserRole]string{
	UserRole_Undefined:  userRole_Undefined_Str,
	UserRole_Owner:      userRole_Owner_Str,
	UserRole_Superadmin: userRole_Superadmin_Str,
}

var userRoleStrInt = map[string]UserRole{
	userRole_Undefined_Str:  UserRole_Undefined,
	userRole_Owner_Str:      UserRole_Owner,
	userRole_Superadmin_Str: UserRole_Superadmin,
}

func UserRoleFromString(s string) UserRole {
	return userRoleStrInt[s]
}

func (r UserRole) String() string {
	return userRoleIntStr[r]
}

func (r UserRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(Flag{
		Name:  userRoleIntStr[r],
		Value: int32(r),
	})
}

func (r *UserRole) UnmarshalJSON(b []byte) error {
	var res Flag
	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}

	*r = UserRole(res.Value)
	return nil
}

func (r UserRole) Int() int {
	return int(r)
}
