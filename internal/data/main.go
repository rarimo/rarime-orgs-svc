package data

//go:generate xo schema "postgres://orgs:orgs@localhost:15432/orgs-db?sslmode=disable" -o ./ --single=schema.xo.go --src templates
//go:generate xo schema "postgres://orgs:orgs@localhost:15432/orgs-db?sslmode=disable" -o pg --single=schema.xo.go --src=pg/templates --go-context=both
//go:generate goimports -w ./

type Storage interface {
	EmailInvitationQ() EmailInvitationQ
	GroupQ() GroupQ
	OrganizationQ() OrganizationQ
}

type EmailInvitationQ interface {
}

type GroupQ interface {
}

type GroupUserQ interface {
}

type OrganizationQ interface {
}

type RequestQ interface {
}

type UserQ interface {
}

type GorpMigrationQ interface {
}
