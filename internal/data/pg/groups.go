package pg

import (
	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"golang.org/x/net/context"
)

func (q GroupQ) SelectCtx(ctx context.Context, selector data.GroupsSelector) ([]data.Group, error) {
	stmt := squirrel.Select("*").From("public.groups")

	if selector.OrgID != nil {
		stmt = stmt.Where(squirrel.Eq{"org_id": selector.OrgID})
	}

	stmt = applyPagination(stmt, selector.Sort, selector.PageSize, selector.PageCursor)

	var groups []data.Group

	if err := q.db.SelectContext(ctx, &groups, stmt); err != nil {
		return nil, errors.Wrap(err, "failed to select groups")
	}

	return groups, nil
}
