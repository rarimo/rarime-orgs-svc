package handlers

import (
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type groupByIDRequest struct {
	ID uuid.UUID

	IncludeGroupUsers bool `include:"group_users"`
}

func newGroupByIDRequest(r *http.Request) (*groupByIDRequest, error) {
	var request groupByIDRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode url")
	}

	request.ID, err = groupIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func GroupByID(w http.ResponseWriter, r *http.Request) {
	req, err := newGroupByIDRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	group, err := Storage(r).GroupQ().GroupByIDCtx(r.Context(), req.ID, false)
	if err != nil {
		panic(errors.Wrap(err, "failed to get group", logan.F{
			"group_id": req.ID,
		}))
	}
	if group == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	resGroup := populateGroup(*group)

	resp := resources.GroupResponse{
		Data:     resGroup,
		Included: resources.Included{},
	}

	if req.IncludeGroupUsers {
		groupUsers, err := Storage(r).GroupUserQ().SelectCtx(r.Context(), group.ID)
		if err != nil {
			panic(errors.Wrap(err, "failed to get group users", logan.F{
				"group_id": group.ID,
			}))
		}

		resGroup.Relationships = &resources.GroupRelationships{
			GroupUsers: resources.RelationCollection{
				Data: make([]resources.Key, len(groupUsers)),
			},
		}

		for i, groupUser := range groupUsers {
			resGroupUser := populateGroupUser(groupUser)
			resGroup.Relationships.GroupUsers.Data[i] = resGroupUser.Key
			resp.Included.Add(&resGroupUser)
		}

		resp.Data = resGroup
	}

	ape.Render(w, resp)
}
