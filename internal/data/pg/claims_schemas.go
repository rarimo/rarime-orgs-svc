package pg

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (q ClaimsSchemaQ) SelectCtx(ctx context.Context, selector data.ClaimsSchemasSelector) ([]data.ClaimsSchema, error) {
	stmt := squirrel.Select("*").From("public.claims_schemas")

	if selector.ID != nil {
		stmt = stmt.Where(squirrel.Eq{"id": selector.ID})
	}
	if selector.ActionType != nil {
		stmt = stmt.Where(squirrel.Eq{"action_type": selector.ActionType})
	}
	if selector.SchemaType != nil {
		stmt = stmt.Where(squirrel.Eq{"schema_type": selector.SchemaType})
	}
	if selector.SchemaURL != nil {
		stmt = stmt.Where(squirrel.Eq{"schema_url": selector.SchemaURL})
	}

	stmt = applyPagination(stmt, selector.Sort, selector.PageLimit, selector.PageCursor)

	var requests []data.ClaimsSchema

	if err := q.db.SelectContext(ctx, &requests, stmt); err != nil {
		return nil, errors.Wrap(err, "failed to select requests")
	}

	return requests, nil
}
