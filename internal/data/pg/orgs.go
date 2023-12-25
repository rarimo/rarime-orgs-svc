package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (q OrganizationQ) New() data.OrganizationQ {
	return OrganizationQ{q.db}
}

func (q OrganizationQ) SelectCtx(ctx context.Context, selector data.OrgsSelector) ([]data.Organization, error) {
	stmt := squirrel.Select("*").From("public.organizations")

	if selector.UserDID != nil {
		stmt = squirrel.Select("o.*").
			From("public.organizations o").
			Join("public.users u ON o.id = u.org_id").
			Where(squirrel.Eq{"u.did": selector.UserDID})
	}

	if selector.Owner != nil {
		stmt = stmt.Where(squirrel.Eq{"owner": selector.Owner})
	}

	if selector.Status != nil {
		stmt = stmt.Where(squirrel.Eq{"status": selector.Status})
	}

	stmt = applyPagination(stmt, selector.Sort, selector.PageSize, selector.PageCursor)

	var orgs []data.Organization

	if err := q.db.SelectContext(ctx, &orgs, stmt); err != nil {
		return nil, errors.Wrap(err, "failed to select organizations")
	}

	return orgs, nil
}

func applyPagination(stmt squirrel.SelectBuilder, sorts pgdb.Sorts, size, cursor uint64) squirrel.SelectBuilder {
	if size != 0 {
		stmt = stmt.Limit(size)
	}

	stmt = stmt.Offset(cursor)

	if len(sorts) == 0 {
		sorts = pgdb.Sorts{"-time"}
	}

	stmt = sorts.ApplyTo(stmt, map[string]string{
		"time": "created_at",
	})

	return stmt
}
