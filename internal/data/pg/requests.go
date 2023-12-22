package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (q RequestQ) SelectCtx(ctx context.Context, selector data.RequestsSelector) ([]data.Request, error) {
	stmt := squirrel.Select("*").From("public.requests")

	if selector.OrgID != nil {
		stmt = stmt.Where(squirrel.Eq{"org_id": selector.OrgID})
	}
	if selector.GroupID != nil {
		stmt = stmt.Where(squirrel.Eq{"group_id": selector.GroupID})
	}
	if selector.UserDID != nil {
		stmt = stmt.Where(squirrel.Eq{"user_did": selector.UserDID})
	}
	if selector.Status != nil {
		stmt = stmt.Where(squirrel.Eq{"status": selector.Status})
	}

	stmt = applyPagination(stmt, selector.Sort, selector.PageSize, selector.PageCursor)

	var requests []data.Request

	if err := q.db.SelectContext(ctx, &requests, stmt); err != nil {
		return nil, errors.Wrap(err, "failed to select requests")
	}

	return requests, nil
}
