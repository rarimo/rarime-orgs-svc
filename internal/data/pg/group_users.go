package pg

import (
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"golang.org/x/net/context"
)

func (q GroupUserQ) SelectCtx(ctx context.Context, groupID uuid.UUID) ([]data.GroupUser, error) {
	stmt := squirrel.Select("*").
		From("public.group_users").
		Where(squirrel.Eq{"group_id": groupID})

	var users []data.GroupUser

	if err := q.db.SelectContext(ctx, &users, stmt); err != nil {
		return nil, errors.Wrap(err, "failed to select group users")
	}

	return users, nil
}
