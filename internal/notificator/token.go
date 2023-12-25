package notificator

import (
	"encoding/base64"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type SendVerifyEmailParams struct {
	OrgID         string `json:"org_id"`
	GroupID       string `json:"group_id"`
	InviteEmailID string `json:"invite_email_id"`
	OTP           string `json:"otp"`
}

func (p SendVerifyEmailParams) Base64() string {
	params, err := json.Marshal(p)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal send verify email params", logan.F{
			"org_id":          p.OrgID,
			"group_id":        p.GroupID,
			"invite_email_id": p.InviteEmailID,
			"otp":             p.OTP,
		}))
	}
	return base64.URLEncoding.EncodeToString(params)
}
