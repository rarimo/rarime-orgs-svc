package handlers

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

type orgUserListRequest struct {
	ID      uuid.UUID
	UserDID *string `filter:"user_did"`

	PageCursor uint64     `page:"cursor"`
	PageLimit  uint64     `page:"limit" default:"10"`
	Sorts      pgdb.Sorts `url:"sort" default:"-time"`

	IncludeOrg bool `include:"organization"`
}

func newOrgUserListRequest(r *http.Request) (*orgUserListRequest, error) {
	var req orgUserListRequest

	err := urlval.Decode(r.URL.Query(), &req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode url")
	}

	req.ID, err = orgIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	if req.UserDID != nil {
		err = validation.Validate(req.UserDID, ValidationDID)
		if err != nil {
			return nil, validation.Errors{
				"filter[user_did]": errors.Wrap(err, "invalid user_did filter value"),
			}
		}
	}

	return &req, nil
}

func OrgUserList(w http.ResponseWriter, r *http.Request) {
	// TODO: add auth but not sure for what
	req, err := newOrgUserListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	users, err := Storage(r).UserQ().SelectCtx(r.Context(), data.UsersSelector{
		OrgID:      &req.ID,
		DID:        req.UserDID,
		PageCursor: req.PageCursor,
		PageSize:   req.PageLimit,
		Sort:       req.Sorts,
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to get users", logan.F{
			"org_id": req.ID,
		}))
	}

	resp := resources.UserListResponse{
		Data:     make([]resources.User, 0, len(users)),
		Included: resources.Included{},
		Links: &resources.Links{
			Self: fmt.Sprintf("%s?%s", r.URL.Path, urlval.MustEncode(req)),
		},
	}

	if len(users) == 0 {
		ape.Render(w, resp)
		return
	}

	for _, user := range users {
		resp.Data = append(resp.Data, populateUser(user))
	}

	cursor := req.PageCursor + req.PageLimit
	req.PageCursor = cursor
	resp.Links.Next = fmt.Sprintf("%s?%s", r.URL.Path, urlval.MustEncode(req))

	_ = resp.PutMeta(map[string]interface{}{
		"next_cursor": cursor,
	})

	if req.IncludeOrg {
		org, err := Storage(r).OrganizationQ().OrganizationByIDCtx(r.Context(), req.ID, false)
		if err != nil {
			panic(errors.Wrap(err, "failed to get organization", logan.F{
				"org_id": req.ID,
			}))
		}
		if org == nil {
			panic(errors.From(errors.New("expected organization exist"), logan.F{
				"org_id": req.ID,
			}))
		}
		respOrg := populateOrg(*org)
		resp.Included.Add(&respOrg)
	}

	ape.Render(w, resp)
}
