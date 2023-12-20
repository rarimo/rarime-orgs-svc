// Package data contains generated code for schema 'public'.
package data

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"database/sql/driver"
	"encoding/csv"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rarimo/xo/types/xo"

	"github.com/google/uuid"
)

// StringSlice is a slice of strings.
type StringSlice []string

// quoteEscapeRegex is the regex to match escaped characters in a string.
var quoteEscapeRegex = regexp.MustCompile(`([^\\]([\\]{2})*)\\"`)

// Scan satisfies the sql.Scanner interface for StringSlice.
func (ss *StringSlice) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		return errors.New("invalid StringSlice")
	}

	// change quote escapes for csv parser
	str := quoteEscapeRegex.ReplaceAllString(string(buf), `$1""`)
	str = strings.Replace(str, `\\`, `\`, -1)

	// remove braces
	str = str[1 : len(str)-1]

	// bail if only one
	if len(str) == 0 {
		*ss = StringSlice([]string{})
		return nil
	}

	// parse with csv reader
	cr := csv.NewReader(strings.NewReader(str))
	slice, err := cr.Read()
	if err != nil {
		fmt.Printf("exiting!: %v\n", err)
		return err
	}

	*ss = StringSlice(slice)

	return nil
}

// Value satisfies the driver.Valuer interface for StringSlice.
func (ss StringSlice) Value() (driver.Value, error) {
	v := make([]string, len(ss))
	for i, s := range ss {
		v[i] = `"` + strings.Replace(strings.Replace(s, `\`, `\\\`, -1), `"`, `\"`, -1) + `"`
	}
	return "{" + strings.Join(v, ",") + "}", nil
} // EmailInvitation represents a row from 'public.email_invitations'.
type EmailInvitation struct {
	ID        uuid.UUID `db:"id" json:"id" structs:"-"`                          // id
	ReqID     uuid.UUID `db:"req_id" json:"req_id" structs:"req_id"`             // req_id
	OrgID     uuid.UUID `db:"org_id" json:"org_id" structs:"org_id"`             // org_id
	GroupID   uuid.UUID `db:"group_id" json:"group_id" structs:"group_id"`       // group_id
	Email     string    `db:"email" json:"email" structs:"email"`                // email
	Otp       string    `db:"otp" json:"otp" structs:"otp"`                      // otp
	CreatedAt time.Time `db:"created_at" json:"created_at" structs:"created_at"` // created_at

}

// GorpMigration represents a row from 'public.gorp_migrations'.
type GorpMigration struct {
	ID        string       `db:"id" json:"id" structs:"-"`                          // id
	AppliedAt sql.NullTime `db:"applied_at" json:"applied_at" structs:"applied_at"` // applied_at

}

// Group represents a row from 'public.groups'.
type Group struct {
	ID        uuid.UUID `db:"id" json:"id" structs:"-"`                          // id
	OrgID     uuid.UUID `db:"org_id" json:"org_id" structs:"org_id"`             // org_id
	Metadata  xo.Jsonb  `db:"metadata" json:"metadata" structs:"metadata"`       // metadata
	Rules     xo.Jsonb  `db:"rules" json:"rules" structs:"rules"`                // rules
	CreatedAt time.Time `db:"created_at" json:"created_at" structs:"created_at"` // created_at

}

// GroupUser represents a row from 'public.group_users'.
type GroupUser struct {
	ID        uuid.UUID `db:"id" json:"id" structs:"-"`                          // id
	GroupID   uuid.UUID `db:"group_id" json:"group_id" structs:"group_id"`       // group_id
	UserID    uuid.UUID `db:"user_id" json:"user_id" structs:"user_id"`          // user_id
	Role      int16     `db:"role" json:"role" structs:"role"`                   // role
	CreatedAt time.Time `db:"created_at" json:"created_at" structs:"created_at"` // created_at
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" structs:"updated_at"` // updated_at

}

// Organization represents a row from 'public.organizations'.
type Organization struct {
	ID                uuid.UUID      `db:"id" json:"id" structs:"-"`                                                     // id
	Did               sql.NullString `db:"did" json:"did" structs:"did"`                                                 // did
	Owner             uuid.UUID      `db:"owner" json:"owner" structs:"owner"`                                           // owner
	Domain            string         `db:"domain" json:"domain" structs:"domain"`                                        // domain
	Metadata          xo.Jsonb       `db:"metadata" json:"metadata" structs:"metadata"`                                  // metadata
	Status            int16          `db:"status" json:"status" structs:"status"`                                        // status
	VerificationCode  sql.NullString `db:"verification_code" json:"verification_code" structs:"verification_code"`       // verification_code
	IssuedClaimsCount int64          `db:"issued_claims_count" json:"issued_claims_count" structs:"issued_claims_count"` // issued_claims_count
	MembersCount      int            `db:"members_count" json:"members_count" structs:"members_count"`                   // members_count
	CreatedAt         time.Time      `db:"created_at" json:"created_at" structs:"created_at"`                            // created_at
	UpdatedAt         time.Time      `db:"updated_at" json:"updated_at" structs:"updated_at"`                            // updated_at

}

// Request represents a row from 'public.requests'.
type Request struct {
	ID        uuid.UUID     `db:"id" json:"id" structs:"-"`                          // id
	OrgID     uuid.UUID     `db:"org_id" json:"org_id" structs:"org_id"`             // org_id
	GroupID   uuid.UUID     `db:"group_id" json:"group_id" structs:"group_id"`       // group_id
	UserID    uuid.NullUUID `db:"user_id" json:"user_id" structs:"user_id"`          // user_id
	Metadata  xo.Jsonb      `db:"metadata" json:"metadata" structs:"metadata"`       // metadata
	Status    int16         `db:"status" json:"status" structs:"status"`             // status
	CreatedAt time.Time     `db:"created_at" json:"created_at" structs:"created_at"` // created_at
	UpdatedAt time.Time     `db:"updated_at" json:"updated_at" structs:"updated_at"` // updated_at

}

// User represents a row from 'public.users'.
type User struct {
	ID        uuid.UUID `db:"id" json:"id" structs:"-"`                          // id
	Did       string    `db:"did" json:"did" structs:"did"`                      // did
	Role      int16     `db:"role" json:"role" structs:"role"`                   // role
	CreatedAt time.Time `db:"created_at" json:"created_at" structs:"created_at"` // created_at
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" structs:"updated_at"` // updated_at

}
