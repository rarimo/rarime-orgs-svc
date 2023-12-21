package pg

import (
	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"golang.org/x/net/context"
)

func (q UserQ) SelectCtx(ctx context.Context, selector data.UsersSelector) ([]data.User, error) {
	stmt := squirrel.Select("*").From("public.users")

	if selector.OrgID != nil {
		stmt = stmt.Where(squirrel.Eq{"org_id": selector.OrgID})
	}

	if selector.DID != nil {
		stmt = stmt.Where(squirrel.Eq{"did": selector.DID})
	}

	if selector.Role != nil {
		stmt = stmt.Where(squirrel.Eq{"role": selector.Role})
	}

	if selector.PageSize != 0 {
		stmt = stmt.Limit(selector.PageSize)
	}

	stmt = stmt.Offset(selector.PageCursor)

	if len(selector.Sort) == 0 {
		selector.Sort = pgdb.Sorts{"-time"}
	}

	stmt = selector.Sort.ApplyTo(stmt, map[string]string{
		"time": "created_at",
	})

	var users []data.User

	if err := q.db.SelectContext(ctx, &users, stmt); err != nil {
		return nil, errors.Wrap(err, "failed to select users")
	}

	return users, nil
}
