package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type groupListRequest struct {
	OrgID uuid.UUID

	PageCursor uint64     `page:"cursor"`
	PageLimit  uint64     `page:"limit" default:"10"`
	Sorts      pgdb.Sorts `url:"sort" default:"-time"`

	IncludeGroupUsers bool `include:"group_users"`
}

func newGroupListRequest(r *http.Request) (*groupListRequest, error) {
	var req groupListRequest

	err := urlval.Decode(r.URL.Query(), &req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode url")
	}

	req.OrgID, err = orgIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func GroupList(w http.ResponseWriter, r *http.Request) {
	req, err := newGroupListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	groups, err := Storage(r).GroupQ().SelectCtx(r.Context(), data.GroupsSelector{
		OrgID:      &req.OrgID,
		PageCursor: req.PageCursor,
		PageLimit:  req.PageLimit,
		Sort:       req.Sorts,
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to get groups", logan.F{
			"org_id": req.OrgID,
		}))
	}

	resp := resources.GroupListResponse{
		Data:     make([]resources.Group, 0, len(groups)),
		Included: resources.Included{},
		Links: &resources.Links{
			Self: fmt.Sprintf("%s?%s", r.URL.Path, urlval.MustEncode(req)),
		},
	}

	if len(groups) == 0 {
		ape.Render(w, resp)
		return
	}

	for _, group := range groups {
		resGroup := populateGroup(group)

		if req.IncludeGroupUsers {
			groupUsers, err := Storage(r).GroupUserQ().SelectCtx(r.Context(), group.ID)
			if err != nil {
				panic(errors.Wrap(err, "failed to get group users", logan.F{
					"org_id":   req.OrgID,
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
		}

		resp.Data = append(resp.Data, resGroup)
	}

	cursor := req.PageCursor + req.PageLimit
	req.PageCursor = cursor
	resp.Links.Next = fmt.Sprintf("%s?%s", r.URL.Path, urlval.MustEncode(req))

	_ = resp.PutMeta(map[string]interface{}{
		"next_cursor": cursor,
	})

	ape.Render(w, resp)
}

func populateGroup(group data.Group) resources.Group {
	return resources.Group{
		Key: resources.Key{
			ID:   group.ID.String(),
			Type: resources.GROUPS,
		},
		Attributes: resources.GroupAttributes{
			OrgId:     group.OrgID.String(),
			Rules:     json.RawMessage(group.Rules),
			Metadata:  json.RawMessage(group.Metadata),
			CreatedAt: group.CreatedAt,
		},
	}
}

func populateGroupUser(groupUser data.GroupUser) resources.GroupUser {
	return resources.GroupUser{
		Key: resources.Key{
			ID:   groupUser.ID.String(),
			Type: resources.GROUP_USERS,
		},
		Attributes: resources.GroupUserAttributes{
			UserId:    groupUser.UserID.String(),
			GroupId:   groupUser.GroupID.String(),
			Role:      resources.GroupUserRole(groupUser.Role),
			CreatedAt: groupUser.CreatedAt,
			UpdatedAt: groupUser.UpdatedAt,
		},
	}
}
