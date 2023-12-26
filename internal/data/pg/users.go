package pg

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
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

	stmt = applyPagination(stmt, selector.Sort, selector.PageLimit, selector.PageCursor)

	var users []data.User

	if err := q.db.SelectContext(ctx, &users, stmt); err != nil {
		return nil, errors.Wrap(err, "failed to select users")
	}

	return users, nil
}

func (q UserQ) UserByDidCtx(ctx context.Context, did string) (*data.User, error) {
	stmt := squirrel.Select("*").From("public.users").Where(squirrel.Eq{"did": did})

	var user data.User

	if err := q.db.GetContext(ctx, &user, stmt); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to select user")
	}

	return &user, nil
}
