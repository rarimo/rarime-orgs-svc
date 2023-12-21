package data

import (
	"context"

	"github.com/google/uuid"
)

//go:generate xo schema "postgres://orgs:orgs@localhost:15432/orgs-db?sslmode=disable" -o ./ --single=schema.xo.go --src templates
//go:generate xo schema "postgres://orgs:orgs@localhost:15432/orgs-db?sslmode=disable" -o pg --single=schema.xo.go --src=pg/templates --go-context=both
//go:generate goimports -w ./

type Storage interface {
	Transaction(func() error) error
	EmailInvitationQ() EmailInvitationQ
	GroupQ() GroupQ
	GroupUserQ() GroupUserQ
	OrganizationQ() OrganizationQ
	RequestQ() RequestQ
	UserQ() UserQ
}

type EmailInvitationQ interface {
}

type GroupQ interface {
	InsertCtx(ctx context.Context, o *Group) error
	SelectCtx(ctx context.Context, selector GroupsSelector) ([]Group, error)
	GroupByIDCtx(ctx context.Context, id uuid.UUID, isForUpdate bool) (*Group, error)
}

type GroupUserQ interface {
	SelectCtx(ctx context.Context, groupID uuid.UUID) ([]GroupUser, error)
}

type OrganizationQ interface {
	InsertCtx(ctx context.Context, o *Organization) error
	UpdateCtx(ctx context.Context, o *Organization) error
	SelectCtx(ctx context.Context, selector OrgsSelector) ([]Organization, error)
	OrganizationByIDCtx(ctx context.Context, id uuid.UUID, isForUpdate bool) (*Organization, error)
}

type RequestQ interface {
}

type UserQ interface {
	InsertCtx(ctx context.Context, o *User) error
	SelectCtx(ctx context.Context, selector UsersSelector) ([]User, error)
	UserByIDCtx(ctx context.Context, id uuid.UUID, isForUpdate bool) (*User, error)
}

type GorpMigrationQ interface {
}
