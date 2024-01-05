package schemas

const (
	ActionDomainOwnership         = "DomainOwnership"
	ActionUserRole                = "UserRole"
	ActionEmployeeNameAndPosition = "EmployeeNameAndPosition"
	ActionEmployeeNickname        = "EmployeeNickname"
	ActionEmployeeWalletAddress   = "EmployeeWalletAddress"
	ActionEmployeeBirthdate       = "EmployeeBirthdate"

	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleUndefined = "undefined"

	SocialMediaEmail = "email"
)

type CreateCredentialRequest struct {
	CredentialSchema  string  `json:"credentialSchema"`
	CredentialSubject subject `json:"credentialSubject"`
	Type              string  `json:"type"`
}

// subject is a marker interface to mark a certain struct as a credential subject
// in order to be processed by the issuer. This provides type safety, compared to
// using interface{}.
type subject interface {
	f()
}

type DomainOwnership struct {
	subject    `json:"-"`
	IdentityID string `json:"id"`
	Domain     string `json:"domain"`
}

type UserRole struct {
	subject `json:"-"`
	Role    string `json:"role"`
}

type EmployeeNameAndPosition struct {
	subject      `json:"-"`
	EmployeeName string `json:"employeeName"`
	Position     string `json:"position"`
	Meta         string `json:"meta"`
}

type EmployeeNickname struct {
	subject             `json:"-"`
	EmployeeSocialMedia string `json:"employeeSocialMedia"`
	EmployeeNickname    string `json:"employeeNickname"`
	Meta                string `json:"meta"`
}

type EmployeeWalletAddress struct {
	subject `json:"-"`
	Address string `json:"address"`
	Meta    string `json:"meta"`
}

type EmployeeBirthdate struct {
	subject   `json:"-"`
	Birthdate int    `json:"birthdate"`
	Meta      string `json:"meta"`
}
